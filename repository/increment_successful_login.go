package repository

import (
	"context"
	"errors"
)

func (r *Repository) IncrementSuccessfulLoginCount(ctx context.Context, userID int64) error {
	result, err := r.Db.ExecContext(ctx, queryIncrementSuccessfulLoginCount, userID)
	if err != nil {
		return err
	}

	// Check the affected rows count
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// No rows updated means user does not exist
	if affectedRows == 0 {
		return errors.New("user not found")
	}

	return nil
}
