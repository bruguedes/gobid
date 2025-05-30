package api

import (
	"errors"
	"net/http"

	"github.com/bruguedes/gobid/internal/jsonutils"
	"github.com/bruguedes/gobid/internal/services"
	"github.com/bruguedes/gobid/internal/usecase/user"
	"github.com/google/uuid"
)

func (api *API) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[user.CreateUserRequest](r)
	if err != nil {
		if problems != nil {
			jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
			return
		}
		jsonutils.EncodeJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := api.UserService.CreateUser(r.Context(), data)

	if err != nil {
		if errors.Is(err, services.ErrUserOrEmailAlreadyExists) {
			jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, map[string]string{
				"error": "User or email already exists",
			})
			return
		}

	}

	jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]uuid.UUID{
		"user_id": userID,
	})
}

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Handle user login
}

func (api *API) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Handle user logout
}
