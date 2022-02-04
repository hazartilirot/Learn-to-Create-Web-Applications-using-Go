package models

import "strings"

const (
	/*ErrNotFound is returned when a resource cannot be found in the database*/
	ErrNotFound modelError = "models: Resource not found"
	/*ErrInvalidID is returned when invalid ID is provided to a method like Delete*/
	ErrInvalidID modelError = "models: ID must be other than 0"
	/*ErrInvalidEmailOrPassword is returned when invalid password is typed in*/
	ErrInvalidEmailOrPassword modelError = "models: Invalid email or password. Please try again"
	ErrInvalidToken           modelError = "models: The token is invalid or empty"

	ErrEmailRequired     modelError = "models: Email address is required"
	ErrEmailInvalid      modelError = "models: Email address is invalid"
	ErrEmailIsRegistered modelError = "modules: Email is already registered"

	ErrPasswordIsShort    modelError = "modules: Password must be at least eight characters long"
	ErrPasswordIsRequired modelError = "modules: Password is required"

	ErrRememberIsShort        modelError = "modules: Token must be at least 32 bytes"
	ErrRememberHashIsRequired modelError = "modules: Remember hash is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "modules: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
