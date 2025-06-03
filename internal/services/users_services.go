package services

import (
	"context"
	"errors"

	"github.com/bruguedes/gobid/internal/store/pgstore"
	"github.com/bruguedes/gobid/internal/usecase/user"
	"github.com/bruguedes/gobid/internal/validator"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
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

func (us *UserService) AuthenticateUser(ctx context.Context, data user.LoginUserRequest) (uuid.UUID, error) {
	user, err := us.queries.GetUserByEmail(ctx, data.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, validator.ErrInvalidCredentials
		}
		return uuid.UUID{}, err
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(data.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return uuid.UUID{}, validator.ErrInvalidCredentials
		}
		return uuid.UUID{}, err
	}

	return user.ID, nil

}
