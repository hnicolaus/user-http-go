// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	InsertUser(ctx context.Context, user User) (userID int64, err error)
	GetUsers(ctx context.Context, request UserFilter) (users []User, err error)
}
