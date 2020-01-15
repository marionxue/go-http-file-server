package serverHandler

import (
	"../util"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type pathEntry struct {
	Name string
	Path string
}

type responseData struct {
	rawReqPath     string
	handlerReqPath string

	hasNotFoundError bool
	hasInternalError bool

	IsRoot        bool
	Path          string
	Paths         []*pathEntry
	RootRelPath   string
	File          *os.File
	Item          os.FileInfo
	ItemName      string
	SubItems      []os.FileInfo
	SubItemPrefix string
	CanUpload     bool
	CanArchive    bool
	CanCors       bool
	NeedAuth      bool
	Errors        []error
}

func isSlash(c rune) bool {
	return c == '/'
}

func getPathEntries(path string, tailSlash bool) []*pathEntry {
	paths := []string{"/"}
	paths = append(paths, strings.FieldsFunc(path, isSlash)...)

	displayPathsCount := len(paths)

	pathsCount := displayPathsCount
	if !tailSlash {
		pathsCount--
	}

	pathEntries := make([]*pathEntry, displayPathsCount)
	for i := 0; i < displayPathsCount; i++ {
		var rPath string
		switch {
		case i < pathsCount-1:
			rPath = strings.Repeat("../", pathsCount-1-i)
		case i == pathsCount-1:
			rPath = "./"
		default:
			rPath = "./" + strings.Join(paths[pathsCount:], "/") + "/"
		}

		pathEntries[i] = &pathEntry{
			Name: paths[i],
			Path: rPath,
		}
	}

	return pathEntries
}

func stat(reqFsPath string, visitFs bool) (file *os.File, item os.FileInfo, err error) {
	if !visitFs {
		return
	}

	file, err = os.Open(reqFsPath)
	if err != nil {
		return
	}

	item, err = file.Stat()
	if err != nil {
		return
	}

	return
}

func readdir(file *os.File, item os.FileInfo) (subItems []os.FileInfo, errs []error) {
	if file == nil || item == nil || !item.IsDir() {
		return
	}

	var err error
	subItems, err = file.Readdir(0)
	if err != nil {
		errs = append(errs, err)
		return
	}

	return
}

func (h *handler) mergeAlias(rawRequestPath string, subItems *[]os.FileInfo) []error {
	errs := []error{}

	for _, alias := range h.aliases {
		aliasUrlPath := alias.urlPath
		aliasFsPath := alias.fsPath

		if len(aliasUrlPath) <= len(rawRequestPath) {
			continue
		}

		suffixIndex := len(rawRequestPath)
		aliasPrefix := aliasUrlPath[:suffixIndex]
		aliasSuffix := aliasUrlPath[suffixIndex:]

		if aliasPrefix != rawRequestPath {
			continue
		}

		if len(aliasPrefix) != 1 && aliasSuffix[0] != '/' {
			// aliasUrlPath doesn't cover the whole directory name
			// e.g:
			// rawReqPath == "/abc/def/ghi"
			// aliasPrefix == "/abc/de"
			continue
		}
		if aliasSuffix[0] == '/' {
			aliasSuffix = aliasSuffix[1:]
		}

		slashIndex := strings.Index(aliasSuffix, "/")
		var nextName string
		if slashIndex >= 0 {
			nextName = aliasSuffix[:slashIndex]
		} else {
			nextName = aliasSuffix
		}

		var aliasSubItem os.FileInfo
		if path.Dir(aliasUrlPath) == rawRequestPath { // reached second deepest path of alias
			var err error
			aliasSubItem, err = os.Stat(aliasFsPath)
			if err == nil {
				aliasSubItem = newRenamedFileInfo(nextName, aliasSubItem)
			} else {
				errs = append(errs, err)
				aliasSubItem = newFakeFileInfo(nextName, true)
			}
		} else {
			aliasSubItem = newFakeFileInfo(nextName, true)
		}

		replaced := false
		for i, subItem := range *subItems {
			if subItem.Name() == nextName {
				(*subItems)[i] = aliasSubItem
				replaced = true
				break
			}
		}

		if !replaced {
			*subItems = append(*subItems, aliasSubItem)
		}
	}

	return errs
}

func getSubItemPrefix(requestPath string, tailSlash bool) (subItemPrefix string) {
	if tailSlash {
		subItemPrefix = "./"
	} else {
		subItemPrefix = "./" + path.Base(requestPath) + "/"
	}
	return
}

func sortSubItems(subItems []os.FileInfo) {
	sort.Slice(
		subItems,
		func(prevIndex, nextIndex int) bool {
			prevItem := subItems[prevIndex]
			nextItem := subItems[nextIndex]

			prevIsDir := prevItem.IsDir()
			nextIsDir := nextItem.IsDir()

			if prevIsDir != nextIsDir {
				return prevIsDir
			}

			return util.CompareNumInStr(prevItem.Name(), nextItem.Name())
		},
	)
}

func getItemName(item os.FileInfo, r *http.Request) (itemName string) {
	if item != nil {
		itemName = item.Name()
	}
	if len(itemName) == 0 || itemName == "." {
		itemName = strings.Replace(r.Host, ":", "_", -1)
	}
	return
}

func (h *handler) getResponseData(r *http.Request) (data *responseData) {
	requestUri := r.URL.Path
	tailSlash := requestUri[len(requestUri)-1] == '/'

	rawReqPath := util.CleanUrlPath(requestUri)
	reqPath := util.CleanUrlPath(rawReqPath[len(h.urlPrefix):]) // strip url prefix path
	errs := []error{}
	notFound := false
	internalError := false

	isRoot := rawReqPath == "/"

	pathEntries := getPathEntries(rawReqPath, tailSlash)
	var rootRelPath string
	if len(pathEntries) > 0 {
		rootRelPath = pathEntries[0].Path
	} else {
		rootRelPath = "./"
	}

	reqFsPath, _absErr := filepath.Abs(h.root + reqPath)
	if _absErr != nil {
		reqFsPath = path.Clean(h.root + reqPath)
	}

	file, item, _statErr := stat(reqFsPath, !h.emptyRoot)
	if _statErr != nil {
		errs = append(errs, _statErr)
		notFound = os.IsNotExist(_statErr)
		internalError = internalError || !notFound
	}

	itemName := getItemName(item, r)

	subItems, _readdirErrs := readdir(file, item)
	errs = append(errs, _readdirErrs...)
	internalError = internalError || len(_readdirErrs) > 0

	_mergeErrs := h.mergeAlias(rawReqPath, &subItems)
	errs = append(errs, _mergeErrs...)
	internalError = internalError || len(_mergeErrs) > 0

	subItems = h.FilterItems(subItems)
	sortSubItems(subItems)

	subItemPrefix := getSubItemPrefix(reqPath, tailSlash)

	canUpload := h.getCanUpload(item, rawReqPath, reqFsPath)
	canArchive := h.getCanArchive(subItems, rawReqPath, reqFsPath)
	canCors := h.getCanCors(rawReqPath, reqFsPath)
	needAuth := h.getNeedAuth(rawReqPath, reqFsPath)

	data = &responseData{
		rawReqPath:     rawReqPath,
		handlerReqPath: reqPath,

		hasNotFoundError: notFound,
		hasInternalError: internalError,

		IsRoot:        isRoot,
		Path:          rawReqPath,
		Paths:         pathEntries,
		RootRelPath:   rootRelPath,
		File:          file,
		Item:          item,
		ItemName:      itemName,
		SubItems:      subItems,
		SubItemPrefix: subItemPrefix,

		CanUpload:  canUpload,
		CanArchive: canArchive,
		CanCors:    canCors,
		NeedAuth:   needAuth,

		Errors: errs,
	}

	return
}
