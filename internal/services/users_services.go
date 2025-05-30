package services

import (
	"context"
	"errors"

	"github.com/bruguedes/gobid/internal/store/pgstore"
	"github.com/bruguedes/gobid/internal/usecase/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserOrEmailAlreadyExists = errors.New("user or email already exists")

type UserService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewUserService(pool *pgxpool.Pool) UserService {
	return UserService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (us *UserService) CreateUser(ctx context.Context, data user.CreateUserRequest) (uuid.UUID, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 12)
	if err != nil {
		return uuid.UUID{}, err
	}

	params := pgstore.CreateUserParams{
		UserName:     data.UserName,
		Email:        data.Email,
		PasswordHash: passwordHash,
		Bio:          data.Bio,
	}

	newUser, err := us.queries.CreateUser(ctx, params)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.UUID{}, ErrUserOrEmailAlreadyExists

		}
		return uuid.UUID{}, err

	}
	return newUser.ID, nil
}
