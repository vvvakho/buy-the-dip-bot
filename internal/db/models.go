// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"

	"github.com/google/uuid"
)

type UserSubscription struct {
	SubID          uuid.UUID
	UserID         int64
	Ticker         string
	DateSubscribed sql.NullTime
	LastNotified   sql.NullTime
	Active         sql.NullBool
}