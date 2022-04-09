package models

import "strings"

const (
	/*ErrNotFound is returned when a resource cannot be found in the database*/
	ErrNotFound modelError = "models: Resource not found"
	/*ErrInvalidEmailOrPassword is returned when invalid password is typed in*/
	ErrInvalidEmailOrPassword modelError = "models: Invalid email or password. Please try again"
	ErrInvalidToken           modelError = "models: The token is invalid or empty"

	ErrEmailRequired     modelError = "models: Email address is required"
	ErrEmailInvalid      modelError = "models: Email address is invalid"
	ErrEmailIsRegistered modelError = "models: Email is already registered"

	ErrPasswordIsShort    modelError = "models: Password must be at least eight characters long"
	ErrPasswordIsRequired modelError = "models: Password is required"

	ErrTitleRequired modelError = "models: Title is required"

	ErrTokenInvalid modelError = "models: token provided is not valid"

	/*ErrInvalidID is returned when invalid ID is provided to a method like Delete*/
	ErrInvalidID privateError = "models: ID provided was invalid"

	ErrRememberIsShort        privateError = "models: Token must be at least 32 bytes"
	ErrRememberHashIsRequired privateError = "models: Remember hash is required"

	ErrUserIDRequired privateError = "models: User ID is required "
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
