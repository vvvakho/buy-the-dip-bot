// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: ManageSubscriptions.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const addSubscription = `-- name: AddSubscription :exec
INSERT INTO user_subscriptions (sub_id, user_id, ticker)
VALUES ($1, $2, $3)
`

type AddSubscriptionParams struct {
	SubID  uuid.UUID
	UserID int64
	Ticker string
}

func (q *Queries) AddSubscription(ctx context.Context, arg AddSubscriptionParams) error {
	_, err := q.db.ExecContext(ctx, addSubscription, arg.SubID, arg.UserID, arg.Ticker)
	return err
}

const getSubscription = `-- name: GetSubscription :one
SELECT sub_id, active FROM user_subscriptions
WHERE user_id = $1 AND ticker = $2
LIMIT 1
`

type GetSubscriptionParams struct {
	UserID int64
	Ticker string
}

type GetSubscriptionRow struct {
	SubID  uuid.UUID
	Active sql.NullBool
}

func (q *Queries) GetSubscription(ctx context.Context, arg GetSubscriptionParams) (GetSubscriptionRow, error) {
	row := q.db.QueryRowContext(ctx, getSubscription, arg.UserID, arg.Ticker)
	var i GetSubscriptionRow
	err := row.Scan(&i.SubID, &i.Active)
	return i, err
}

const getSubscriptionsByTicker = `-- name: GetSubscriptionsByTicker :many
SELECT sub_id, user_id FROM user_subscriptions
WHERE ticker = $1 AND active = TRUE
`

type GetSubscriptionsByTickerRow struct {
	SubID  uuid.UUID
	UserID int64
}

func (q *Queries) GetSubscriptionsByTicker(ctx context.Context, ticker string) ([]GetSubscriptionsByTickerRow, error) {
	rows, err := q.db.QueryContext(ctx, getSubscriptionsByTicker, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSubscriptionsByTickerRow
	for rows.Next() {
		var i GetSubscriptionsByTickerRow
		if err := rows.Scan(&i.SubID, &i.UserID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const reactivateSubscription = `-- name: ReactivateSubscription :exec
UPDATE user_subscriptions
SET active = true, date_subscribed = NOW()
WHERE user_id = $1 AND ticker = $2
`

type ReactivateSubscriptionParams struct {
	UserID int64
	Ticker string
}

func (q *Queries) ReactivateSubscription(ctx context.Context, arg ReactivateSubscriptionParams) error {
	_, err := q.db.ExecContext(ctx, reactivateSubscription, arg.UserID, arg.Ticker)
	return err
}
