// This file contains types that are used in the repository layer.
package repository

import "errors"

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

var ErrUserNotFound = errors.New("user not found")
var ErrPhoneNumberAlreadyExists = errors.New("phone number already exists")

type User struct {
	ID                   string
	FullName             string
	Password             string
	PhoneNumber          string
	SuccessfulLoginCount int
}

type UpdateUser struct {
	ID          string
	FullName    *string
	PhoneNumber *string
}
