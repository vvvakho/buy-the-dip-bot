package api

import (
	"buy-the-dip-bot/database"
	"buy-the-dip-bot/internal/db"
	"buy-the-dip-bot/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
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

type RSI struct {
	RSI  float64
	Date time.Time
}

var avEnv = "ALPHA_VANTAGE_API_KEY"

func InitAlphaVantageClient() (*AVClient, error) {
	apiKey, err := utils.GetEnv(avEnv)
	if err != nil {
		log.Printf("Could not retrieve Alpha Vantage API Key: %s", err)
		return &AVClient{}, err
	}

	return &AVClient{ApiKey: apiKey}, nil
}

func (av *AVClient) FetchRSI(ticker string, date time.Time, queriesDB *db.Queries) (RSI, error) {
	rsiRow, err := database.CheckRSIinDB(ticker, date, queriesDB)
	if err != nil {
		if errors.Is(err, database.ErrRSINotFound) {
			rsi, err := requestRSI(ticker, queriesDB)
			if err != nil {
				return RSI{}, err
			}
			return RSI{RSI: rsi.RSI, Date: rsi.Date}, nil
		}
		return RSI{}, err
	}
	return RSI{RSI: rsiRow.Rsi, Date: rsiRow.Date}, nil
}

func requestRSI(ticker string, queriesDB *db.Queries) (RSI, error) {
	//url := fmt.Sprintf("https://www.alphavantage.co/query?function=RSI&symbol=%s&interval=weekly&time_period=10&series_type=open&apikey=%s", ticker, c.ApiKey)

	//log.Print("Initiating Get request")
	//resp, err := http.Get(url)
	//if err != nil {
	//	return 0, fmt.Errorf("failed to fetch RSI: %v", err)
	//}
	//defer resp.Body.Close()

	fileName := fmt.Sprintf("%s_rsi.json", ticker)
	//outFile, err := os.Create(fileName)
	//if err != nil {
	//	return 0, fmt.Errorf("failed to create file: %v", err)
	//}
	//defer outFile.Close()

	//_, err = io.Copy(outFile, resp.Body)
	//if err != nil {
	//	return 0, fmt.Errorf("failed to write response to file: %v", err)
	//}

	//log.Printf("Repsonse successfully written to %s", fileName)

	file, err := os.Open(fileName)
	if err != nil {
		return RSI{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var data AVDataRSI
	//if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
	//	return 0, fmt.Errorf("failed to decode RSI data: %v", err)
	//}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return RSI{}, fmt.Errorf("failed to decode file data: %v", err)
	}

	todaysDate := time.Now()
	for i := 0; i < 10; i++ {
		dateKey := todaysDate.Format("2006-01-02")
		if rsiData, ok := data.TechnicalAnalysisRSI[dateKey]; ok {
			rsiFloat, err := strconv.ParseFloat(rsiData.RSI, 64)
			if err != nil {
				return RSI{}, err
			}

			if err := database.AddRSI(ticker, rsiFloat, data.MetaData.LastRefreshed, queriesDB); err != nil {
				log.Printf("Error saving RSI record to database: %v", err)
			}
			return RSI{RSI: rsiFloat, Date: todaysDate}, nil
		}

		todaysDate = todaysDate.AddDate(0, 0, -1)
	}

	return RSI{}, fmt.Errorf("no RSI data available for %s", ticker)
}
