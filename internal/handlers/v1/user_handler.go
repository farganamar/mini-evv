package handler

import (
	"encoding/json"
	"net/http"

	"github.com/farganamar/evv-service/helpers/failure"
	"github.com/farganamar/evv-service/internal/model/v1/user/dto"
	"github.com/farganamar/evv-service/transport/http/response"
)

// UserLogin godoc
// @Summary      User Login
// @Description  User Login
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "Login Request"
// @Success 200 {object} response.Base
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 403 {object} response.Base
// @Failure 404 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/evv/user/login [post]
func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	ctx := r.Context()
	var userLoginRequest dto.LoginRequest

	if err := decoder.Decode(&userLoginRequest); err != nil {
		response.WithError(w, failure.BadRequest(err))
	}

	if err := userLoginRequest.Validate(); err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userLoginResponse, err := h.UserServiceV1.Login(ctx, userLoginRequest)
	if err != nil {
		code := failure.GetCode(err)
		response.WithJSON(w, code, nil, err.Error())
		return
	}

	response.WithJSON(w, http.StatusOK, userLoginResponse, "OK")
}
