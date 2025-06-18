package product

import (
	"context"
	"time"

	"github.com/bruguedes/gobid/internal/validator"
	"github.com/google/uuid"
)

type CreateProductRequest struct {
	SellerID    uuid.UUID `json:"seller_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	AuctionEnd  time.Time `json:"auction_end"`
}

func (req CreateProductRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckFieldError(validator.NotBlank(req.ProductName), "product_name", "this field must be provided")
	eval.CheckFieldError(validator.MinChar(req.ProductName, 3) && validator.MaxChar(req.ProductName, 100),
		"product_name",
		"this field must be between 3 and 100 characters")

	eval.CheckFieldError(validator.NotBlank(req.Description), "description", "this field must be provided")
	eval.CheckFieldError(validator.MinChar(req.Description, 10) && validator.MaxChar(req.Description, 255),
		"description",
		"this field must be between 10 and 255 characters")

	eval.CheckFieldError(validator.Price(req.BasePrice), "base_price", "this field must be a positive number")

	eval.CheckFieldError(validator.AuctionEnd(req.AuctionEnd), "auction_end", "this field must be at least two hours duration")

	return eval
}
