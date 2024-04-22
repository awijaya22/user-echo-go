// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetUserByID(ctx context.Context, userID string) (user User, err error)
	CreateUser(ctx context.Context, user User) (ID string, err error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user User, err error)
	UpdateSuccessfulLoginCount(ctx context.Context, userID string) (err error)
	UpdateUser(ctx context.Context, updateReq UpdateUser) (err error)
}
