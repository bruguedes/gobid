package api

import (
	"errors"
	"net/http"

	"github.com/bruguedes/gobid/internal/jsonutils"
	"github.com/bruguedes/gobid/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *API) HandleSubscribeToAuction(w http.ResponseWriter, r *http.Request) {
	rawProductID := chi.URLParam(r, "product_id")

	productID, err := uuid.Parse(rawProductID)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusBadRequest, map[string]any{"error": "invalid product ID"})
		return
	}

	_, err = api.ProductService.GetProductByID(r.Context(), productID)
	if err != nil {
		if errors.Is(err, validator.ErrProductNotFound) {
			jsonutils.EncodeJSON(w, r, http.StatusNotFound, map[string]any{"error": "no product with given id"})
		}
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{"error": "internal server error"})
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), "authenticatedUserID").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{"error": "unexpected error occurred"})
		return

	}

	conn, err := api.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{"error": "failed to upgrade connection"})
		return
	}

}
