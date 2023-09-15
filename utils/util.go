package utils

import (
	"github.com/lib/pq"
)

// IsUniqueConstraintViolation determines if an error is a Postgres UNIQUE constraint error.
func IsUniqueConstraintViolation(err error) bool {
	// http://godoc.org/github.com/lib/pq#Error
	switch e := err.(type) {
	case *pq.Error:
		if e.Code == "23505" {
			return true
		}
	}

	return false
}
