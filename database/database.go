package database

import (
	"buy-the-dip-bot/internal/db"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var ErrAlreadySubscribed = errors.New("subscribtion is already active")

func AddSubscription(userID int64, ticker string, queriesDB *db.Queries) error {
	subRow, err := GetSubscriptionRow(userID, ticker, queriesDB)
	if err != nil {
		if err == sql.ErrNoRows {
			addSubParams := db.AddSubscriptionParams{
				SubID:  uuid.New(),
				UserID: userID,
				Ticker: ticker,
			}
			if err := queriesDB.AddSubscription(context.Background(), addSubParams); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if subRow.Active.Valid && subRow.Active.Bool {
		return ErrAlreadySubscribed
	} else {
		reacSubParams := db.ReactivateSubscriptionParams{
			UserID: userID,
			Ticker: ticker,
		}
		if err := queriesDB.ReactivateSubscription(context.Background(), reacSubParams); err != nil {
			return err
		}
	}

	return nil
}

func Unsubscribe(userID int64, ticker string, queriesDB *db.Queries) error {
	err := queriesDB.Unsubscribe(context.Background(), db.UnsubscribeParams{UserID: userID, Ticker: ticker})
	if err != nil {
		return err
	}
	return nil
}

func GetSubscriptionRow(userID int64, ticker string, queriesDB *db.Queries) (db.UserSubscription, error) {
	getSubParams := db.GetSubscriptionParams{
		UserID: userID,
		Ticker: ticker,
	}
	subRow, err := queriesDB.GetSubscription(context.Background(), getSubParams)
	if err != nil {
		return db.UserSubscription{}, err
	}
	return subRow, nil
}

func GetAllUserIDs(queriesDB *db.Queries) ([]int64, error) {
	// TODO: Change to proper user DB
	users, err := queriesDB.TempGetAllUsers(context.Background())
	if err != nil {
		return []int64{}, err
	}
	return users, nil
}
