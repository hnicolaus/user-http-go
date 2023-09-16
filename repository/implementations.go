package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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

func (r *Repository) IncrementSuccessfulLoginCount(ctx context.Context, userID int64) error {
	rows, err := r.Db.QueryContext(ctx, queryIncrementSuccessfulLoginCount, userID)
	if err != nil {
		return err
	}

	defer rows.Close()
	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, in User) (int64, error) {
	var (
		query     string
		setFields []string
		params    []interface{}
		offset    int = 0
	)

	if in.PhoneNumber != "" {
		setFields = append(setFields, fmt.Sprintf(setUserPhoneNumberF, offset+1))
		params = append(
			params,
			in.PhoneNumber,
		)
		offset++
	}

	if in.FullName != "" {
		setFields = append(setFields, fmt.Sprintf(setUserFullNameF, offset+1))
		params = append(
			params,
			in.FullName,
		)
		offset++
	}

	setFields = append(setFields, fmt.Sprintf(setUserUpdatedTimeF, offset+1))
	params = append(
		params,
		time.Now(),
	)
	offset++

	query = fmt.Sprintf(queryUpdateUserF, strings.Join(setFields, ","))

	query += fmt.Sprintf(whereUserID, offset+1)
	params = append(
		params,
		in.ID,
	)
	offset++

	result, err := r.Db.ExecContext(ctx, query, params...)
	if err != nil {
		return 0, err
	}

	// Check the affected rows count
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}
