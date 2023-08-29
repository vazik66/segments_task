package user

import (
	"avito-segment/internal/user"
	"avito-segment/pkg"
	"net/http"
)

type RPCHandler struct {
	service *user.UserService
}

func NewHandler(service *user.UserService) *RPCHandler {
	return &RPCHandler{
		service: service,
	}
}

// @Summary Create user
// @Description Creates user and returns it
// @Description Method name: users.Create
// @Tags users
// @Accept json
// @Produce json
// @Param Request body pkg.JsonRPCRequest{params=array} true "Desc"
// @Success 200 {object} pkg.JsonRPCSuccessResponse{result=user.User}
// @Failure 400 {object} pkg.JsonRPCErrorResponse
// @Router /jsonrpc [post]
func (h *RPCHandler) Create(r *http.Request, args *pkg.EmptyArgs, reply *user.User) error {
	user, err := h.service.Create(r.Context())
	if err != nil {
		return err
	}
    *reply = *user
	return nil
}

type UserIDArg struct {
	UserID uint `json:"userId" example:"1"`
}

// @Summary Get user by ID
// @Description Returns user with given id
// @Description Method name: users.GetByID
// @Tags users
// @Accept json
// @Produce json
// @Param Request body pkg.JsonRPCRequest{params=user.UserIDArg} true "User ID"
// @Success 200 {object} pkg.JsonRPCSuccessResponse{result=user.User}
// @Failure 400 {object} pkg.JsonRPCErrorResponse
// @Router /jsonrpc/GetUserByID [post]
func (h *RPCHandler) GetByID(r *http.Request, args *UserIDArg, reply *user.User) error {
	user, err := h.service.GetByID(r.Context(), args.UserID)
	if err != nil {
		return err
	}
	reply.ID = user.ID
	return nil
}

// @Summary Delete User
// @Description Deletes User by id
// @Description Method name: users.Delete
// @Tags users
// @Accept json
// @Produce json
// @Param Request body pkg.JsonRPCRequest{params=UserIDArg} true "Desc"
// @Success 200 {array} pkg.JsonRPCSuccessResponse{result=pkg.EmptyResponse}
// @Failure 400 {object} pkg.JsonRPCErrorResponse
// @Router /jsonrpc/DelteUser [post]
func (h *RPCHandler) Delete(r *http.Request, args *UserIDArg, reply *pkg.EmptyResponse) error {
	err := h.service.Delete(r.Context(), args.UserID)
	if err != nil {
		return err
	}
	return nil
}

// @Summary List Users
// @Description Returns list of all users
// @Description Method name: users.List
// @Tags users
// @Accept json
// @Produce json
// @Param Request body pkg.JsonRPCRequest{params=pkg.EmptyArgs} true "Desc"
// @Success 200 {array} pkg.JsonRPCSuccessResponse{result=[]user.User}
// @Failure 400 {object} pkg.JsonRPCErrorResponse
// @Router /jsonrpc/ListUsers [post]
func (h *RPCHandler) List(r *http.Request, args *pkg.EmptyArgs, reply *[]user.User) error {
	users, err := h.service.List(r.Context())
	if err != nil {
		return err
	}
	*reply = *users
	return nil
}
