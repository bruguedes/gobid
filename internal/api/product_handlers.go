package api

import (
	"net/http"

	"github.com/bruguedes/gobid/internal/jsonutils"
	"github.com/bruguedes/gobid/internal/usecase/product"
	"github.com/google/uuid"
)

func (api *API) HandlerCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[product.CreateProductRequest](r)

	if err != nil {
		if problems != nil {
			jsonutils.EncodeJSON(w, r, http.StatusBadRequest, problems)
			return
		}
		jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	userID, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserID").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]string{"error": "must be logged in to create a product"})
		return
	}

	data.SellerID = userID

	id, err := api.ProductService.CreateProduct(r.Context(), data)

	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{"error": "failed to create product"})
		return
	}
	jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{
		"message":    "product created successfully",
		"product_id": id,
	})
}
