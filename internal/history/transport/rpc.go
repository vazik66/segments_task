package history

import (
	"avito-segment/internal/history"
	"net/http"
)

type RPCHandler struct {
	service *history.HistoryService
}

func NewHandler(service *history.HistoryService) *RPCHandler {
	return &RPCHandler{
		service: service,
	}
}

// @Summary Create CSV Report UserSegments History
// @Description Returns link for csv report about users UserSegments history
// @Description Method name: history.CreateReport
// @Tags history
// @Accept json
// @Produce json
// @Param Request body pkg.JsonRPCRequest{params=history.CreateReportParams} true "Params"
// @Success 200 {object} pkg.JsonRPCSuccessResponse{result=history.ReportUrlResponse}
// @Failure 400 {object} pkg.JsonRPCErrorResponse
// @Router /jsonrpc/CreateReport [post]
func (h *RPCHandler) CreateReport(r *http.Request, args *history.CreateReportParams, reply *history.ReportUrlResponse) error {
    url, err := h.service.CreateReport(r.Context(), args.UserID, args.Date)   
    if err != nil {
        return err
    }
    reply.Url = url
	return nil
}
