package services

import (
	"context"
	"errors"

	"github.com/bruguedes/gobid/internal/store/pgstore"
	"github.com/bruguedes/gobid/internal/validator"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BidsService struct {
	poll    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewBidsService(poll *pgxpool.Pool) BidsService {
	return BidsService{
		poll:    poll,
		queries: pgstore.New(poll),
	}
}

func (bs *BidsService) PlaceBid(ctx context.Context, product_id, bidder_id uuid.UUID, amount float64) (pgstore.Bid, error) {
	params := pgstore.CreateBidParams{
		ProductID: product_id,
		BidderID:  bidder_id,
		BidAmount: amount,
	}

	product, err := bs.queries.GetProductByID(ctx, product_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
		return pgstore.Bid{}, err
	}

	highestBid, err := bs.queries.GetHighestBidByProductId(ctx, product_id)
	if err != nil {
		// caso não haja lances, devemos continuar a execução
		if !errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}

	if product.BasePrice >= amount || highestBid.BidAmount >= amount {
		return pgstore.Bid{}, validator.ErrBidAmountTooLow
	}

	newHighestBid, err := bs.queries.CreateBid(ctx, params)

	if err != nil {
		return pgstore.Bid{}, err
	}
	return newHighestBid, nil

}
