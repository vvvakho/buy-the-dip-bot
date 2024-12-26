package api

import (
	"buy-the-dip-bot/database"
	"buy-the-dip-bot/internal/db"
	"buy-the-dip-bot/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

type AVDataDaily struct {
	MetaData struct {
		Information   string `json:"1. Information"`
		Symbol        string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
		OutputSize    string `json:"4. Output Size"`
		TimeZone      string `json:"5. Time Zone"`
	} `json:"Meta Data"`
	TimeSeriesDaily map[string]DailyJson `json:"Time Series (Daily)"`
}

type DailyJson struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type Daily struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume uint32
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
			rsi, err := av.requestRSI(ticker, queriesDB)
			if err != nil {
				return RSI{}, err
			}
			return RSI{RSI: rsi.RSI, Date: rsi.Date}, nil
		}
		return RSI{}, err
	}
	return RSI{RSI: rsiRow.Rsi, Date: rsiRow.Date}, nil
}

func (av *AVClient) requestRSI(ticker string, queriesDB *db.Queries) (RSI, error) {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=RSI&symbol=%s&interval=weekly&time_period=10&series_type=open&apikey=%s", ticker, av.ApiKey)

	log.Print("Initiating Get request")
	resp, err := http.Get(url)
	if err != nil {
		return RSI{}, fmt.Errorf("failed to fetch RSI: %v", err)
	}
	defer resp.Body.Close()

	var data AVDataRSI
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return RSI{}, fmt.Errorf("failed to decode RSI data: %v", err)
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

func (av *AVClient) FetchDaily(ticker string, date time.Time, queriesDB *db.Queries) (Daily, error) {

	log.Printf("fetching daily data for ticker %s, date %v\n", ticker, date)
	daily, err := av.requestDaily(ticker, date, queriesDB)
	if err != nil {
		return Daily{}, fmt.Errorf("error fetching daily data: %v", err)
	}

	return daily, nil
}

func (av *AVClient) requestDaily(ticker string, date time.Time, queriesDB *db.Queries) (Daily, error) {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s", ticker, av.ApiKey)

	resp, err := http.Get(url)
	if err != nil {
		return Daily{}, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	data := AVDataDaily{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Daily{}, fmt.Errorf("error decoding response: %v", err)
	}

	log.Println(data)

	latestDate := date
	for i := 0; i < 7; i++ {
		dateString := latestDate.Format("2006-01-02")
		todaysData, ok := data.TimeSeriesDaily[dateString]
		if !ok {
			latestDate = latestDate.AddDate(0, 0, -1)
		} else {
			openFloat, _ := strconv.ParseFloat(todaysData.Open, 64)
			highFloat, _ := strconv.ParseFloat(todaysData.High, 64)
			lowFloat, _ := strconv.ParseFloat(todaysData.Low, 64)
			closeFloat, _ := strconv.ParseFloat(todaysData.Close, 64)
			volUint, _ := strconv.ParseUint(todaysData.Volume, 10, 32)

			dataFormatted := Daily{
				Open:   openFloat,
				High:   highFloat,
				Low:    lowFloat,
				Close:  closeFloat,
				Volume: uint32(volUint),
			}
			return dataFormatted, nil
		}
	}

	return Daily{}, errors.New("error retrieving daily data")

}
