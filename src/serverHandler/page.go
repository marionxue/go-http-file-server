package serverHandler

import (
	tplutil "../tpl/util"
	"html/template"
	"net/http"
)

const TypeDir = template.HTML("dir")
const TypeFile = template.HTML("file")

func updateSubItemsHtml(data *responseData) {
	length := len(data.SubItems)

	data.SubItemsHtml = make([]*itemHtml, length)

	for i, info := range data.SubItems {
		name := info.Name()
		displayName := tplutil.FormatFilename(name)

		var typ template.HTML
		var url string
		var readableSize template.HTML

		if info.IsDir() {
			typ = TypeDir
			url = data.SubItemPrefix + name + "/"
			displayName += "/"
		} else {
			typ = TypeFile
			url = data.SubItemPrefix + name
			readableSize = tplutil.FormatSize(info.Size())
		}

		data.SubItemsHtml[i] = &itemHtml{
			Type:        typ,
			Url:         url,
			DisplayName: displayName,
			DisplaySize: readableSize,
			DisplayTime: tplutil.FormatTime(info.ModTime()),
		}
	}
}

func (h *handler) page(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Cache-Control", "public, max-age=0")

	w.WriteHeader(data.Status)

	if needResponseBody(r.Method) {
		updateSubItemsHtml(data)
		err := h.template.Execute(w, data)
		h.errHandler.LogError(err)
	}
}
