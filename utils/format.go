package utils

import (
	"fmt"
	"time"
)

func FormatMessage(date time.Time, ticker string, rsi float64, comment string) string {
	formattedDate := date.Format("January 2, 2006")

	formattedMessage := fmt.Sprintf(
		"%s ($%s)\n%s\n\nRSI: %.2f\n\n%s",
		"S&P500",
		ticker,
		formattedDate,
		rsi,
		comment,
	)

	return formattedMessage
}
