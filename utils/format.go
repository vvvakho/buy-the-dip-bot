package utils

import (
	"fmt"
	"time"
)

func FormatMessage(date time.Time, ticker string, rsi float64) string {
	formattedDate := date.Format("January 2, 2006")

	formattedMessage := fmt.Sprintf(
		"%s\n\n%s ($%s)\nRSI: %.2f",
		formattedDate,
		"S&P500",
		ticker,
		rsi,
	)

	return formattedMessage
}
