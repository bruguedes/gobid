package user

import (
	"context"

	"github.com/bruguedes/gobid/internal/validator"
)

type CreateUserRequest struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (req CreateUserRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckFieldError(validator.NotBlank(req.UserName), "user_name", "must be provided")
	eval.CheckFieldError(validator.MinChar(req.UserName, 3) && validator.MaxChar(req.UserName, 50),
		"user_name",
		"must be between 3 and 50 characters")

	eval.CheckFieldError(validator.NotBlank(req.Email), "email", "must be provided")
	eval.CheckFieldError(validator.ValidateEmail(req.Email, validator.EmailRegex), "email", "must be a valid email address")

	eval.CheckFieldError(validator.NotBlank(req.Password), "password", "must be provided")
	eval.CheckFieldError(validator.MinChar(req.Password, 8), "password", "must be at least 8 characters")

	eval.CheckFieldError(validator.NotBlank(req.Bio), "bio", "must be provided")
	eval.CheckFieldError(validator.MinChar(req.Bio, 10) && validator.MaxChar(req.Bio, 255),
		"bio",
		"must be between 10 and 255 characters")

	return eval
}
