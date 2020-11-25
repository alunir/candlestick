package candle

type BuySellType int

const (
	ALL BuySellType = iota
	BUY
	SELL
)

type ChartMode struct{ Value string }

func (c ChartMode) String() string {
	if c.Value == "" {
		return "UNDEFINED"
	}
	return c.Value
}

var (
	TIME        = ChartMode{"TIME"}
	VOLUME      = ChartMode{"VOLUME"}
	BUY_VOLUME  = ChartMode{"BUY_VOLUME"}
	SELL_VOLUME = ChartMode{"SELL_VOLUME"}
	AMOUNT      = ChartMode{"AMOUNT"}
	BUY_AMOUNT  = ChartMode{"BUY_AMOUNT"}
	SELL_AMOUNT = ChartMode{"SELL_AMOUNT"}
	COUNT       = ChartMode{"COUNT"}
	RATIO       = ChartMode{"RATIO"}
	BUY_PRICE   = ChartMode{"BUY_PRICE"}
	SELL_PRICE  = ChartMode{"SELL_PRICE"}
)
