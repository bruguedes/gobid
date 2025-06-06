package api

import (
	"errors"
	"net/http"

	"github.com/bruguedes/gobid/internal/jsonutils"
	"github.com/bruguedes/gobid/internal/services"
	"github.com/bruguedes/gobid/internal/usecase/user"
	"github.com/bruguedes/gobid/internal/validator"
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
			jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "user or email already exists",
			})
			return
		}

	}

	jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]uuid.UUID{
		"user_id": userID,
	})
}

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {

	data, problems, err := jsonutils.DecodeValidJson[user.LoginUserRequest](r)

	if err != nil {
		if problems != nil {
			jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
			return
		}

		if errors.Is(err, validator.ErrInvalidEmail) || errors.Is(err, validator.ErrNotBlank) {
			jsonutils.EncodeJSON(w, r, http.StatusBadRequest, map[string]any{
				"error": "invalid email or password",
			})
			return
		}

		jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, err.Error())

		return

	}

	userID, err := api.UserService.AuthenticateUser(r.Context(), data)

	if err != nil {
		if errors.Is(err, validator.ErrInvalidCredentials) {
			jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]any{
				"error": "invalid email or password",
			})
			return
		}

		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	err = api.Sessions.RenewToken(r.Context()) //Renova o token da sessão antes de adicionar o usuário
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal error",
		})
		return
	}

	api.Sessions.Put(r.Context(), "AuthenticatedUserID", userID) // Armazena o ID do usuário autenticado na sessão
	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]any{
		"message": "user authenticated successfully",
	})

}

func (api *API) HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal error",
		})
		return
	}

	api.Sessions.Remove(r.Context(), "AuthenticatedUserID") // Remove o ID do usuário autenticado da sessão
	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]any{
		"message": "user logged out successfully",
	})
}
