package user

import (
	"context"

	"github.com/bruguedes/gobid/internal/validator"
)

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req LoginUserRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	eval.CheckFieldError(validator.ValidateEmail(req.Email, validator.EmailRegex), "email", "must be a valid email address")
	eval.CheckFieldError(validator.NotBlank(req.Password), "password", "must be provided")

	return eval
}
