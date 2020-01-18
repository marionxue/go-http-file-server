package serverHandler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type jsonItem struct {
	IsDir   bool      `json:"isDir"`
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

type jsonResponseData struct {
	Item     *jsonItem   `json:"item"`
	SubItems []*jsonItem `json:"subItems"`
}

func getJsonItem(info os.FileInfo) *jsonItem {
	return &jsonItem{
		IsDir:   info.IsDir(),
		Name:    info.Name(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}
}

func getJsonData(data *responseData) *jsonResponseData {
	var item *jsonItem
	var subItems []*jsonItem

	if data.Item != nil {
		item = getJsonItem(data.Item)

		if data.Item.IsDir() {
			subItems = make([]*jsonItem, len(data.SubItems))
			for i, dataItem := range data.SubItems {
				subItems[i] = getJsonItem(dataItem.Info)
			}
		}
	}

	return &jsonResponseData{
		Item:     item,
		SubItems: subItems,
	}
}

func (h *handler) json(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Cache-Control", "public, max-age=0")

	if data.hasInternalError {
		w.WriteHeader(http.StatusInternalServerError)
	} else if data.hasNotFoundError {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if needResponseBody(r.Method) {
		jsonData := getJsonData(data)
		encoder := json.NewEncoder(w)
		err := encoder.Encode(jsonData)
		h.errHandler.LogError(err)
	}
}