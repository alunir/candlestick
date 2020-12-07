package volume

import (
	"math"
	"time"

	c "github.com/alunir/candlestick/candle"
)

type VolumeChart struct {
	c.Chart
	Chunk   float64
	Buysell c.BuySellType
}

func (chart *VolumeChart) AddTrade(ti time.Time, value float64, volume float64) {
	if chart.Buysell == c.ALL {
		volume = math.Abs(volume)
	}
	chart.addTradeToVolumeCandle(ti, value, volume)
}

func (chart *VolumeChart) addTradeToVolumeCandle(ti time.Time, value float64, volume float64) {
	currentStack := 0.0
	if chart.CurrentCandle != nil {
		currentStack = chart.CurrentCandle.Stack
	}
	diffStack := math.Abs(volume)
	newStack := currentStack + diffStack

	remain := math.Mod(newStack, chart.Chunk)
	chunknum := (newStack - remain) / chart.Chunk

	if chunknum == 0 && remain == 0.0 {
		// <--------------chunk------------------->
		// |------------------|-------------------|
		// <--Current.stack-->|<----diffStack----->
		if chart.CurrentCandle != nil {
			chart.CurrentCandle.AddCandleWithBuySell(chart.Buysell, value, volume, diffStack)
			chart.SetLastCandle(chart.CurrentCandle)
			chart.CurrentCandle = nil
		} else {
			panic("No need to addTradeToCandle")
		}
	} else if chunknum == 0.0 {
		// <--------------chunk------------------->
		// |------------------|-------------------|
		// <--Current.stack-->|<--diffStack-->
		if chart.CurrentCandle != nil {
			chart.CurrentCandle.AddCandleWithBuySell(chart.Buysell, value, volume, diffStack)
			chart.CurrentCandleNew = false
		} else {
			chart.AddCandle(c.NewCandleWithBuySell(chart.Buysell, ti, value, volume, diffStack))
		}
	} else if remain == 0.0 {
		// <--------------chunk-------------------><--------------chunk------------------->
		// |------------------|-------------------|---------------------------------------|
		// <--Current.stack-->|<--------------------diffStack----------------------------->
		if chart.CurrentCandle != nil {
			chart.CurrentCandle.AddCandleWithBuySell(chart.Buysell, value, chart.Chunk-currentStack, chart.Chunk-currentStack)
			chart.SetLastCandle(chart.CurrentCandle)
			chunknum -= 1.0
		}
		for i := 0; i < int(chunknum); i++ {
			chart.AddCandle(c.NewCandleWithBuySell(chart.Buysell, ti, value, chart.Chunk, chart.Chunk))
			chart.SetLastCandle(chart.CurrentCandle)
		}
		chart.CurrentCandle = nil
	} else {
		// <--------------chunk-------------------><--------------chunk------------------->
		// |------------------|-------------------|---------------------------------------|
		// <--Current.stack-->|<------------diffStack------------------>
		if chart.CurrentCandle != nil {
			chart.CurrentCandle.AddCandleWithBuySell(chart.Buysell, value, chart.Chunk-currentStack, chart.Chunk-currentStack)
			chart.SetLastCandle(chart.CurrentCandle)
			chunknum -= 1.0
		}
		for i := 0; i < int(chunknum); i++ {
			chart.AddCandle(c.NewCandleWithBuySell(chart.Buysell, ti, value, chart.Chunk, chart.Chunk))
			chart.SetLastCandle(chart.CurrentCandle)
		}
		chart.AddCandle(c.NewCandleWithBuySell(chart.Buysell, ti, value, remain, remain))
	}
	return
}

func (chart *VolumeChart) AddLv2DataCallback(ti time.Time, askPrices []float64, askSizes []float64, bidPrices []float64, bidSizes []float64) {
}
