package api

type MarketDataProvider interface {
	FetchRSI(symbol string) (float64, error)
	FetchPrice(symbol string) (float64, error)
}
