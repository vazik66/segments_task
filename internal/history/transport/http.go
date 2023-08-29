package history

import (
	"net/http"
)

type HTTPHandler struct {}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

// @Summary Get Report File
// @Description Returns report file
// @Tags history
// @Accept json
// @Produce octet-stream
// @Param filename path string true "filename"
// @Success 200 {file} file "ok"
// @Failure 404
// @Router /files/ [get]
func (h *HTTPHandler) GetReportFile(w http.ResponseWriter, r *http.Request) error {
    // just for swagger
    // handled by http.FileServer
	return nil
}

