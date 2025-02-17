package user

import (
	"github.com/szucik/trade-helper/apperrors"
	"regexp"
	"time"

	"github.com/szucik/trade-helper/portfolio"
)

type User struct {
	Username  string
	Email     string
	Password  string
	TokenHash string
	Created   time.Time
	Updated   time.Time
}

type AuthCredentials struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	Username  string                `json:"username" validate:"required"`
	Email     string                `json:"email" validate:"required"`
	Portfolio []portfolio.Portfolio `json:"portfolios" validate:"required"`
	Created   time.Time             `json:"created"`
}

func (u User) NewAggregate() (Aggregate, error) {
	switch {
	case !isLengthValid(u.Username, 2):
		return Aggregate{}, apperrors.Error(
			"User name is to short",
			"UserParamsValidation",
			400,
		)

	case !isEmailValid(u.Email):
		return Aggregate{}, apperrors.Error(
			"Invalid user email",
			"UserParamsValidation",
			400,
		)

	case len(u.Password) < 8:
		return Aggregate{},
			apperrors.Error(
				"Password is to short, it should be longer than 8 characters",
				"UserParamsValidation",
				400,
			)
	}

	return Aggregate{
		user: u,
	}, nil
}

func isLengthValid(value string, length int) bool {
	if len(value) < length {
		return false
	}
	return true
}

func isEmailValid(email string) bool {
	var emailRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"" +
		"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")" +
		"@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|" +
		"1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:" +
		"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")
	return isLengthValid(email, 6) && emailRegex.MatchString(email)
}
