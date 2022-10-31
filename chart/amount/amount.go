package amount

import (
	"math"
	"time"

	c "github.com/alunir/candlestick/candle"
	"github.com/shopspring/decimal"
)

type AmountChart struct {
	*c.Chart
	Chunk   decimal.Decimal
	Buysell c.BuySellType
}

func NewAmountChart(chunk float64, buysell c.BuySellType, candle_num int) *AmountChart {
	return &AmountChart{
		Chart:   c.NewChart(candle_num),
		Chunk:   decimal.NewFromFloat(chunk),
		Buysell: buysell,
	}
}

func (chart AmountChart) AddTrade(ti time.Time, value float64, volume float64) {
	chart.Lock()
	defer chart.Unlock()

	if chart.Buysell == c.ALL {
		volume = math.Abs(volume)
	}
	chart.addTradeToAmountCandle(ti, decimal.NewFromFloat(value), decimal.NewFromFloat(volume))
}

func (chart AmountChart) addTradeToAmountCandle(ti time.Time, value, volume decimal.Decimal) {
	currentStack := decimal.Zero
	if chart.CurrentCandle != nil {
		currentStack = chart.CurrentCandle.Stack
	}
	diffStack := value.Mul(volume).Abs()
	newStack := currentStack.Add(diffStack)
	chunknum, remain := newStack.QuoRem(chart.Chunk, 0)

	if chunknum.Equal(decimal.NewFromInt(1)) && remain.IsZero() {
		// <--------------chunk------------------->
		// |------------------|-------------------|
		// <--Current.stack-->|<----diffStack----->
		if chart.CurrentCandle != nil {
			chart.CurrentCandle.AddCandleWithBuySell(chart.Buysell, ti, value, volume, diffStack)
			chart.SetLastCandle(*chart.CurrentCandle)
			chart.CurrentCandle = nil
		} else {
			panic("No need to addTradeToCandle")
		}
	} else if chunknum.IsZero() {
		// <--------------chunk------------------->
		// |------------------|-------------------|
		// <--Current.stack-->|<--diffStack-->
		if chart.CurrentCandle != nil {
			chart.CurrentCandle.AddCandleWithBuySell(chart.Buysell, ti, value, volume, diffStack)
			chart.CurrentCandleNew = false
		} else {
			chart.AddCandle(c.NewCandleWithBuySell(chart.Buysell, ti, value, volume, diffStack))
		}
	} else if remain.IsZero() {
		// <--------------chunk-------------------><--------------chunk------------------->
		// |------------------|-------------------|---------------------------------------|
		// <--Current.stack-->|<--------------------diffStack----------------------------->
		if chart.CurrentCandle != nil {
			chart.CurrentCandle.AddCandleWithBuySell(chart.Buysell, ti, value, chart.Chunk.Sub(currentStack).Div(value), chart.Chunk.Sub(currentStack))
			chart.SetLastCandle(*chart.CurrentCandle)
			chunknum = chunknum.Sub(decimal.NewFromInt(1))
		}
		for i := 0; i < int(chunknum.IntPart()); i++ {
			chart.AddCandle(c.NewCandleWithBuySell(chart.Buysell, ti, value, chart.Chunk.Div(value), chart.Chunk))
			chart.SetLastCandle(*chart.CurrentCandle)
		}
		chart.CurrentCandle = nil
	} else {
		// <--------------chunk-------------------><--------------chunk------------------->
		// |------------------|-------------------|---------------------------------------|
		// <--Current.stack-->|<------------diffStack------------------>
		if chart.CurrentCandle != nil {
			chart.CurrentCandle.AddCandleWithBuySell(chart.Buysell, ti, value, chart.Chunk.Sub(currentStack).Div(value), chart.Chunk.Sub(currentStack))
			chart.SetLastCandle(*chart.CurrentCandle)
			chunknum = chunknum.Sub(decimal.NewFromInt(1))
		}
		for i := 0; i < int(chunknum.IntPart()); i++ {
			chart.AddCandle(c.NewCandleWithBuySell(chart.Buysell, ti, value, chart.Chunk.Div(value), chart.Chunk))
			chart.SetLastCandle(*chart.CurrentCandle)
		}
		chart.AddCandle(c.NewCandleWithBuySell(chart.Buysell, ti, value, remain.Div(value), remain))
	}
}

func (chart AmountChart) AddLv2DataCallback(ti time.Time, askPrices []float64, askSizes []float64, bidPrices []float64, bidSizes []float64) {
}

func (chart AmountChart) GetChartInfo() map[string]interface{} {
	return make(map[string]interface{})
}
