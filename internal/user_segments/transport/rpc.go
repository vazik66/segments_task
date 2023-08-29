package usersegments

import (
	usersegments "avito-segment/internal/user_segments"
	"avito-segment/pkg"
	"log"
	"net/http"
)

type RPCHandler struct {
	service *usersegments.UserSegmentsService
}

func NewHandler(service *usersegments.UserSegmentsService) *RPCHandler {
	return &RPCHandler{
		service: service,
	}
}

// @Summary Add Add Segments To User
// @Description Adds list of segments to user
// @Description Ttl - seconds to live
// @Description Method name: usersegments.AddSegmentsToUser
// @Tags usersegments
// @Accept json
// @Produce json
// @Param Request body pkg.JsonRPCRequest{params=usersegments.AddSegmentsToUserParams} true "Desc"
// @Success 200 {object} pkg.JsonRPCSuccessResponse{result=pkg.EmptyResponse}
// @Failure 400 {object} pkg.JsonRPCErrorResponse
// @Router /jsonrpc/AddSegmentsToUser [post]
func (h *RPCHandler) AddSegmentsToUser(r *http.Request, args *usersegments.AddSegmentsToUserParams, reply *pkg.EmptyResponse) error {
	err := h.service.AddToUser(r.Context(), args)
	if err != nil {
		return err
	}
	return nil
}

// @Summary Remove Segments From User
// @Description Removes list of segments from user
// @Description Method name: usersegments.RemoveSegmentsFromUser
// @Tags usersegments
// @Accept json
// @Produce json
// @Param Request body pkg.JsonRPCRequest{params=usersegments.RemoveSegmentsFromUserParams} true "Desc"
// @Success 200 {object} pkg.JsonRPCSuccessResponse{result=pkg.EmptyResponse}
// @Failure 400 {object} pkg.JsonRPCErrorResponse
// @Router /jsonrpc/RemoveSegmentsFromUser [post]
func (h *RPCHandler) RemoveSegmentsFromUser(r *http.Request, args *usersegments.RemoveSegmentsFromUserParams, reply *pkg.EmptyResponse) error {
	err := h.service.RemoveFromUser(r.Context(), args)
	if err != nil {
		return err
	}
	return nil
}

// @Summary Get Segment By User
// @Description Returns list of segments associated with user
// @Description Method name: usersegments.GetByUserID
// @Tags usersegments
// @Accept json
// @Produce json
// @Param Request body pkg.JsonRPCRequest{params=usersegments.GetUserSegmentsParams} true "Desc"
// @Success 200 {object} pkg.JsonRPCSuccessResponse{result=[]usersegments.UserSegment}
// @Failure 400 {object} pkg.Error
// @Router /jsonrpc/GetSegmentsByUser [post]
func (h *RPCHandler) GetByUserID(r *http.Request, args *usersegments.GetUserSegmentsParams, reply *[]usersegments.UserSegment) error {
	userSegments, err := h.service.GetByUser(r.Context(), args)
	if err != nil {
		return err
	}
    log.Printf("User segments: %+v", userSegments)
	*reply = *userSegments
	return nil
}
