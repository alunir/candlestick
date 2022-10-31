package volume

import (
	"math"
	"time"

	c "github.com/alunir/candlestick/candle"
	"github.com/shopspring/decimal"
)

type VolumeChart struct {
	*c.Chart
	Chunk   decimal.Decimal
	Buysell c.BuySellType
}

func NewVolumeChart(chunk float64, buysell c.BuySellType, candle_num int) *VolumeChart {
	return &VolumeChart{
		Chart:   c.NewChart(candle_num),
		Chunk:   decimal.NewFromFloat(chunk),
		Buysell: buysell,
	}
}

func (chart VolumeChart) AddTrade(ti time.Time, value, volume float64) {
	chart.Lock()
	defer chart.Unlock()

	if chart.Buysell == c.ALL {
		volume = math.Abs(volume)
	}
	chart.addTradeToVolumeCandle(ti, decimal.NewFromFloat(value), decimal.NewFromFloat(volume))
}

func (chart VolumeChart) addTradeToVolumeCandle(ti time.Time, value, volume decimal.Decimal) {
	currentStack := decimal.Zero
	if chart.CurrentCandle != nil {
		currentStack = chart.CurrentCandle.Stack
	}
	diffStack := volume.Abs()
	newStack := currentStack.Add(diffStack)
	chunknum, remain := newStack.QuoRem(chart.Chunk, 0)

	if chunknum.IsZero() && remain.IsZero() {
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
			chart.CurrentCandle.AddCandleWithBuySell(chart.Buysell, ti, value, chart.Chunk.Sub(currentStack), chart.Chunk.Sub(currentStack))
			chart.SetLastCandle(*chart.CurrentCandle)
			chunknum = chunknum.Sub(decimal.NewFromInt(1))
		}
		for i := 0; i < int(chunknum.IntPart()); i++ {
			chart.AddCandle(c.NewCandleWithBuySell(chart.Buysell, ti, value, chart.Chunk, chart.Chunk))
			chart.SetLastCandle(*chart.CurrentCandle)
		}
		chart.CurrentCandle = nil
	} else {
		// <--------------chunk-------------------><--------------chunk------------------->
		// |------------------|-------------------|---------------------------------------|
		// <--Current.stack-->|<------------diffStack------------------>
		if chart.CurrentCandle != nil {
			chart.CurrentCandle.AddCandleWithBuySell(chart.Buysell, ti, value, chart.Chunk.Sub(currentStack), chart.Chunk.Sub(currentStack))
			chart.SetLastCandle(*chart.CurrentCandle)
			chunknum = chunknum.Sub(decimal.NewFromInt(1))
		}
		for i := 0; i < int(chunknum.IntPart()); i++ {
			chart.AddCandle(c.NewCandleWithBuySell(chart.Buysell, ti, value, chart.Chunk, chart.Chunk))
			chart.SetLastCandle(*chart.CurrentCandle)
		}
		chart.AddCandle(c.NewCandleWithBuySell(chart.Buysell, ti, value, remain, remain))
	}
}

func (chart VolumeChart) AddLv2DataCallback(ti time.Time, askPrices []float64, askSizes []float64, bidPrices []float64, bidSizes []float64) {
}

func (chart VolumeChart) GetChartInfo() map[string]interface{} {
	return make(map[string]interface{})
}
