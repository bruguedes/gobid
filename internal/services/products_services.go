package services

import (
	"context"

	"github.com/bruguedes/gobid/internal/store/pgstore"
	"github.com/bruguedes/gobid/internal/usecase/product"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewProductService(pool *pgxpool.Pool) ProductService {
	return ProductService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, data product.CreateProductRequest) (uuid.UUID, error) {
	params := pgstore.CreateProductParams{
		SellerID:    data.SellerID,
		ProductName: data.ProductName,
		Description: data.Description,
		BasePrice:   data.BasePrice,
		AuctionEnd:  data.AuctionEnd,
	}

	newProductID, err := ps.queries.CreateProduct(ctx, params)
	if err != nil {
		return uuid.UUID{}, err
	}
	return newProductID, nil
}
