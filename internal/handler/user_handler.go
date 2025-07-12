package handler

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/service"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
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

func (u *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := util.GetUserID(ctx)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	userProfile, err := u.userService.GetUserProfile(ctx, userID)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	payload := util.NewJSONPayload()
	payload.Add("data", util.JSONPayload{"user": userProfile})

	WriteJSONHTTP(w, payload, http.StatusOK, nil)
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
