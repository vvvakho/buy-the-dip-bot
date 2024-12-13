package api

import (
	"buy-the-dip-bot/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type AVClient struct {
	ApiKey string
}

type AVDataRSI struct {
	MetaData struct {
		Symbol        string `json:"1: Symbol"`
		Indicator     string `json:"2: Indicator"`
		LastRefreshed string `json:"3: Last Refreshed"`
		Interval      string `json:"4: Interval"`
		TimePeriod    int    `json:"5: Time Period"`
		SeriesType    string `json:"6: Series Type"`
		TimeZone      string `json:"7: Time Zone"`
	} `json:"Meta Data"`
	TechnicalAnalysisRSI map[string]struct {
		RSI string `json:"RSI"`
	} `json:"Technical Analysis: RSI"`
}

var avEnv = "ALPHA_VANTAGE_API_KEY"
var TodaysRSI float64

func InitAlphaVantageClient() (*AVClient, error) {
	apiKey, err := utils.GetEnv(avEnv)
	if err != nil {
		log.Printf("Could not retrieve Alpha Vantage API Key: %s", err)
		return &AVClient{}, err
	}

	return &AVClient{ApiKey: apiKey}, nil
}

func (c *AVClient) FetchRSI(ticker string) (float64, error) {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=RSI&symbol=%s&interval=weekly&time_period=10&series_type=open&apikey=%s", ticker, c.ApiKey)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch RSI: %v", err)
	}
	defer resp.Body.Close()

	var data AVDataRSI
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("failed to decode RSI data: %v", err)
	}

	for date, rsiData := range data.TechnicalAnalysisRSI {
		rsi, err := strconv.ParseFloat(rsiData.RSI, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid RSI value for date %s: %v", date, err)
		}
		TodaysRSI = rsi
		return rsi, nil
	}

	return 0, fmt.Errorf("no RSI data available for %s", ticker)
}
