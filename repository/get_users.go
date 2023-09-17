package repository

import (
	"context"
	"fmt"
	"strconv"
)

func (r *Repository) GetUsers(ctx context.Context, request UserFilter) (users []User, err error) {
	query, params, err := buildQueryGetUsers(request)
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

func buildQueryGetUsers(in UserFilter) (string, []interface{}, error) {
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
