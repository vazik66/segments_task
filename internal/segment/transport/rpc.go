package segment

import (
	"avito-segment/internal/segment"
	"avito-segment/pkg"
	"net/http"
)

type RPCHandler struct {
	service *segment.SegmentService
}

func NewHandler(service *segment.SegmentService) *RPCHandler {
	return &RPCHandler{
		service: service,
	}
}

// @Summary Create Segment
// @Description Creates segment and returns it
// @Description Method name: segments.Create
// @Tags segments
// @Accept json
// @Product json
// @Param Request body pkg.JsonRPCRequest{params=segment.CreateSegmentParams} true "Segment"
// @Success 200 {object} pkg.JsonRPCSuccessResponse{result=segment.Segment}
// @Failure 400 {object} pkg.JsonRPCErrorResponse
// @Router /jsonrpc/CreateSegment [post]
func (h *RPCHandler) Create(r *http.Request, args *segment.CreateSegmentParams, reply *segment.Segment) error {
	newSegment, err := h.service.Create(r.Context(), args)
	if err != nil {
        return err
	}
    *reply = *newSegment
	return nil
}

// @Summary Get Segment By Slug
// @Description Returns segment by slug
// @Description Method name: segments.GetBySlug
// @Tags segments
// @Accept json
// @Product json
// @Param Request body pkg.JsonRPCRequest{params=string} true "Slug"
// @Success 200 {object} pkg.JsonRPCSuccessResponse{result=segment.Segment}
// @Failure 400 {object} pkg.JsonRPCErrorResponse
// @Router /jsonrpc/GetBySlug [post]
func (h *RPCHandler) GetBySlug(r *http.Request, args *string, reply *segment.Segment) error {
	segment, err := h.service.GetBySlug(r.Context(), args)
	if err != nil {
		return err
	}
	*reply = *segment
	return nil
}

// @Summary List Segments
// @Description Returns list of segments
// @Description Method name: segments.List
// @Tags segments
// @Accept json
// @Product json
// @Param Request body pkg.JsonRPCRequest{params=pkg.EmptyArgs} true "Empty"
// @Success 200 {object} pkg.JsonRPCSuccessResponse{result=[]segment.Segment}
// @Failure 400 {object} pkg.JsonRPCErrorResponse
// @Router /jsonrpc/ListSegments [post]
func (h *RPCHandler) List(r *http.Request, args *pkg.EmptyArgs, reply *[]segment.Segment) error {
	segments, err := h.service.List(r.Context())
	if err != nil {
		return err
	}
	*reply = *segments
	return nil
}

// @Summary Delete Segment
// @Description Deletes segment by slug
// @Description Method name: segments.Delete
// @Tags segments
// @Accept json
// @Product json
// @Param Request body pkg.JsonRPCRequest{params=string} true "Slug"
// @Success 200 {object} pkg.JsonRPCSuccessResponse{result=pkg.EmptyResponse}
// @Failure 400 {object} pkg.JsonRPCErrorResponse
// @Router /jsonrpc/DeleteSegment [post]
func (h *RPCHandler) Delete(r *http.Request, args *string, reply *pkg.EmptyResponse) error {
	err := h.service.Delete(r.Context(), args)
	if err != nil {
		return err
	}
	return nil
}

