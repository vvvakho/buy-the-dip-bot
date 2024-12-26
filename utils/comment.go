package utils

func CommentRSI(rsi float64) string {
	switch {
	case rsi >= 70:
		return "RSI is over 70. The stock is likely overbought, and a reversal might be imminent."
	case rsi > 60:
		return "RSI is above 60. The stock is showing strong bullish momentum."
	case rsi > 40 && rsi <= 60:
		return "RSI is between 40 and 60. The stock is in a neutral range, indicating consolidation or balance between buyers and sellers."
	case rsi <= 40 && rsi > 30:
		return "RSI is below 40. The stock is showing weak momentum and might be approaching oversold conditions."
	case rsi <= 30:
		return "RSI is below 30. The stock is likely oversold, and a potential rebound could be near."
	default:
		return "RSI value is out of the expected range."
	}
}
