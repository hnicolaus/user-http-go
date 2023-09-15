package repository

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	saltCost = 12
)

func (r *Repository) InsertUser(ctx context.Context, user User) (userID int64, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), saltCost)
	if err != nil {
		return userID, err
	}
	user.Password = string(hashedPassword)

	query, params := buildQueryInsertUsers([]User{user})

	err = r.Db.QueryRowContext(ctx, query, params...).Scan(&userID)

	return
}

func buildQueryInsertUsers(in []User) (string, []interface{}) {
	var (
		query  string = queryInsertUsers
		params []interface{}
		offset int = 0
	)

	for _, row := range in {
		query += fmt.Sprintf(
			valuesInsertUsersF,
			offset+1, offset+2, offset+3, offset+4,
		)

		params = append(
			params,
			row.FullName,
			row.PhoneNumber,
			row.Password,
		)

		// created_time is DB internal timestamp for new row creation
		params = append(params, time.Now())

		offset = offset + 4
	}

	// trim the last comma
	query = fmt.Sprintln(query[0:len(query)-1], returnLastInsertedUserID)

	return query, params
}
