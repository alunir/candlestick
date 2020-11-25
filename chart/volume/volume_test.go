package volume

import (
	"fmt"
	"testing"
	"time"

	"github.com/alunir/candlestick/candle"
)

func TestVolumeCandles(t *testing.T) {
	candleNum := 4
	var chart = &VolumeChart{
		Chart: candle.Chart{
			CandleNum:  candleNum,
			Candles:    make([]*candle.Candle, 0, candleNum),
			TimeSeries: map[time.Time]*candle.Candle{},
			Clock:      make(chan *candle.Candle),
		},
		Chunk:   4.0,
		Buysell: candle.ALL,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	chart.AddTrade(start, 5, 1)
	chart.AddTrade(start.Add(5*time.Second), 25, 1)
	chart.AddTrade(start.Add(25*time.Second), 3, 1)
	chart.AddTrade(start.Add(60*time.Second), 12, 5)
	var c1 = chart.Candles[0]
	var c2 = chart.Candles[1]

	chart.AddTrade(start.Add(70*time.Second), 13, 2)

	// Intentionally empty data series included here, to test flat candles
	chart.AddTrade(start.Add(240*time.Second), 15, 5)
	var c3 = chart.Candles[2]
	var c4 = chart.Candles[3]

	if !(c1.Count == 4 && c1.Open == 5 && c1.Close == 12 &&
		c1.High == 25 && c1.Low == 3 && c1.Volume == 4) {
		t.Logf("Got wrong c1 val: %v", c1)
		t.Fail()
	}

	if !(c2.Count == 1 && c2.Open == 12 && c2.Close == 12 &&
		c2.High == 12 && c2.Low == 12 && c2.Volume == 4) {
		t.Logf("Got wrong c2 val: %v", c2)
		t.Fail()
	}

	if !(c3.Count == 2 && c3.Open == 13 && c3.Close == 15 &&
		c3.High == 15 && c3.Low == 13 && c3.Volume == 4) {
		t.Logf("Got wrong c3 val: %v", c3)
		t.Fail()
	}

	if !(c4.Count == 1 && c4.Open == 15 && c4.Close == 15 &&
		c4.High == 15 && c4.Low == 15 && c4.Volume == 3) {
		t.Logf("Got wrong c4 val: %v", c4)
		t.Fail()
	}

	if !(chart.LastCandle.Count == 2 && chart.LastCandle.Open == 13 && chart.LastCandle.Close == 15 &&
		chart.LastCandle.High == 15 && chart.LastCandle.Low == 13 && chart.LastCandle.Volume == 4) {
		t.Logf("Got wrong chart.LastCandle val: %v", chart.LastCandle)
		t.Fail()
	}

	if !(chart.CurrentCandle.Count == 1 && chart.CurrentCandle.Open == 15 && chart.CurrentCandle.Close == 15 &&
		chart.CurrentCandle.High == 15 && chart.CurrentCandle.Low == 15 && chart.CurrentCandle.Volume == 3) {
		t.Logf("Got wrong chart.CurrentCandle val: %v", chart.CurrentCandle)
		t.Fail()
	}

	if len(chart.Candles) != 4 {
		t.Logf("Got wrong len: %v", len(chart.Candles))
		t.Fail()
	}

	// Candles should be like a queue
	fmt.Printf("Got cap: %v len: %v\n", cap(chart.Candles), len(chart.Candles))
	chart.AddTrade(start.Add(300*time.Second), 10, 2)
	chart.AddTrade(start.Add(310*time.Second), 3, 6)
	chart.AddTrade(start.Add(370*time.Second), 54, 36)

	var c5 = chart.Candles[0]
	if !(c5.Count == 1 && c5.Open == 54 && c5.Close == 54 &&
		c5.High == 54 && c5.Low == 54 && c5.Volume == 4) {
		t.Logf("Got wrong c5 val: %v", c5)
		t.Fail()
	}
	fmt.Printf("Got cap: %v len: %v\n", cap(chart.Candles), len(chart.Candles))
	if len(chart.Candles) != 4 {
		t.Logf("Got wrong len: %v", len(chart.Candles))
		t.Fail()
	}

	if !(chart.LastCandle.Count == 1 && chart.LastCandle.Open == 54 && chart.LastCandle.Close == 54 &&
		chart.LastCandle.High == 54 && chart.LastCandle.Low == 54 && chart.LastCandle.Volume == 4) {
		t.Logf("Got wrong chart.LastCandle val: %v", chart.LastCandle)
		t.Fail()
	}

	if !(chart.CurrentCandle.Count == 1 && chart.CurrentCandle.Open == 54 && chart.CurrentCandle.Close == 54 &&
		chart.CurrentCandle.High == 54 && chart.CurrentCandle.Low == 54 && chart.CurrentCandle.Volume == 3) {
		t.Logf("Got wrong chart.CurrentCandle val: %v", chart.CurrentCandle)
		t.Fail()
	}
}
