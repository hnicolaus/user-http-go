// This file contains types that are used in the repository layer.
package repository

import "time"

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type User struct {
	ID          int64     `db:"id"`
	FullName    string    `db:"full_name"`
	PhoneNumber string    `db:"phone_number"`
	Password    string    `db:"password"`
	CreatedTime time.Time `db:"created_time"`
	UpdatedTime time.Time `db:"updated_time"`
}

type UserFilter struct {
	PhoneNumber string `db:"phone_number"`
	Password    string `db:"password"`
}
