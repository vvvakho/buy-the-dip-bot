package database

import (
	"buy-the-dip-bot/internal/db"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

var ErrAlreadySubscribed = errors.New("subscribtion is already active")
var ErrRSINotFound = errors.New("RSI for specified ticker and date not found")

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

func CheckRSIinDB(ticker string, date time.Time, queriesDB *db.Queries) (db.Rsi, error) {
	dateCheck := date
	for i := 0; i < 10; i++ {
		rsiRow, err := queriesDB.CheckRSI(context.Background(), db.CheckRSIParams{Ticker: ticker, Column2: dateCheck})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// No record found, move to the previous day
				dateCheck = dateCheck.AddDate(0, 0, -1)
				continue
			}
			return db.Rsi{}, err
		}
		log.Printf("RSI record found in DB")
		return db.Rsi{Rsi: rsiRow.Rsi, Date: rsiRow.Date}, nil
	}
	log.Printf("RSI not found in DB")
	return db.Rsi{}, ErrRSINotFound

}

func AddRSI(ticker string, rsi float64, dateKey string, queriesDB *db.Queries) error {
	parsedDate, err := time.Parse("2006-01-02", dateKey)
	if err != nil {
		return err
	}
	_, err = queriesDB.AddRSI(context.Background(), db.AddRSIParams{
		RsiID:  uuid.New(),
		Ticker: ticker,
		Rsi:    rsi,
		Date:   parsedDate,
	})

	if err != nil {
		return err
	}
	return nil
}
