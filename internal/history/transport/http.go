package history

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPHandler struct {
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

//	@Summary		Get Report File
//	@Description	Returns report file
//	@Tags			history
//	@Accept			json
//	@Produce		octet-stream
//	@Param			filename	path	string	true	"filename"
//	@Success		200			{file}	file	"ok"
//	@Failure		404
//	@Router			/files/{filename} [get]
func (h *HTTPHandler) GetReportFile(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	w.Header().Add("content-type", "application/octet-stream")
	http.ServeFile(w, r, fmt.Sprintf("%s/%s", "files", filename))
}
