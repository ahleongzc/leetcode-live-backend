package handler

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(
	userService service.UserService,
) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := ReadJSONHTTPReq(w, r, request)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	err = u.userService.RegisterNewUser(ctx, request.Email, request.Password)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	WriteJSONHTTP(w, nil, http.StatusCreated, nil)
}
