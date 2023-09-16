package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	saltCost = 12
)

func (r *Repository) InsertUser(ctx context.Context, user User) (userID int64, err error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), saltCost)
	if err != nil {
		return userID, err
	}

	user.Password = string(hashedPasswordBytes)

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

func (r *Repository) GetUsers(ctx context.Context, request UserFilter) (users []User, err error) {
	query, params, err := buildQueryGetUser(request)
	if err != nil {
		return []User{}, err
	}

	rows, err := r.Db.QueryContext(ctx, query, params...)
	if err != nil {
		return []User{}, err
	}

	defer rows.Close()
	for rows.Next() {
		user := User{}

		if err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.PhoneNumber,
			&user.Password,
			&user.CreatedTime,
			&user.UpdatedTime,
		); err != nil {
			return []User{}, err
		}

		users = append(users, user)
	}

	return users, nil
}

func buildQueryGetUser(in UserFilter) (string, []interface{}, error) {
	var (
		query  string = querySelectUsers
		params []interface{}
		offset int = 0
	)

	if in.PhoneNumber != "" {
		query += fmt.Sprintf(whereUserPhoneNumber, offset+1)
		params = append(
			params,
			in.PhoneNumber,
		)
		offset++
	}

	if in.UserID != 0 {
		query += fmt.Sprintf(whereUserID, offset+1)
		params = append(
			params,
			strconv.Itoa(int(in.UserID)),
		)
		offset++
	}

	return query, params, nil
}
