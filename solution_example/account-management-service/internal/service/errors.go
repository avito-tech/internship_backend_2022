package service

import "fmt"

var (
	ErrCannotSignToken  = fmt.Errorf("cannot sign token")
	ErrCannotParseToken = fmt.Errorf("cannot parse token")

	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrCannotCreateUser  = fmt.Errorf("cannot create user")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrCannotGetUser     = fmt.Errorf("cannot get user")

	ErrAccountAlreadyExists = fmt.Errorf("account already exists")
	ErrCannotCreateAccount  = fmt.Errorf("cannot create account")
	ErrAccountNotFound      = fmt.Errorf("account not found")
	ErrCannotGetAccount     = fmt.Errorf("cannot get account")

	ErrCannotCreateReservation = fmt.Errorf("cannot create reservation")
)
