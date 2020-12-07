package count

import (
	"fmt"
	"testing"
	"time"

	"github.com/alunir/candlestick/candle"
)

func TestCountCandles(t *testing.T) {
	candleNum := 4
	var chart = &CountChart{
		Chart: candle.Chart{
			CandleNum:  candleNum,
			Candles:    make([]*candle.Candle, 0, candleNum),
			TimeSeries: map[time.Time]*candle.Candle{},
			Clock:      make(chan *candle.Candle),
		},
		Chunk:   2,
		Buysell: candle.ALL,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	chart.AddTrade(start, 5, 1)
	chart.AddTrade(start.Add(5*time.Second), 25, 1)
	chart.AddTrade(start.Add(25*time.Second), 3, 1)
	var c1 = chart.Candles[0]

	chart.AddTrade(start.Add(60*time.Second), 12, 5)
	chart.AddTrade(start.Add(70*time.Second), 13, 2)
	var c2 = chart.Candles[1]

	// Intentionally empty data series included here, to test flat candles
	chart.AddTrade(start.Add(240*time.Second), 15, 6)
	var c3 = chart.Candles[2]

	if !(c1.Count == 2 && c1.Volume == 2 && c1.Open == 5 && c1.Close == 25 &&
		c1.High == 25 && c1.Low == 5) {
		t.Logf("Got wrong c1 val: %v", c1)
		t.Fail()
	}

	if !(c2.Count == 2 && c2.Volume == 6 && c2.Open == 3 && c2.Close == 12 &&
		c2.High == 12 && c2.Low == 3) {
		t.Logf("Got wrong c2 val: %v", c2)
		t.Fail()
	}

	if !(c3.Count == 2 && c3.Volume == 8 && c3.Open == 13 && c3.Close == 15 &&
		c3.High == 15 && c3.Low == 13) {
		t.Logf("Got wrong c3 val: %v", c3)
		t.Fail()
	}

	if len(chart.Candles) != 3 {
		t.Logf("Got wrong len: %v", len(chart.Candles))
		t.Fail()
	}

	// Candles should be like a queue
	fmt.Printf("Got cap: %v len: %v\n", cap(chart.Candles), len(chart.Candles))
	chart.AddTrade(start.Add(300*time.Second), 10, 2)
	chart.AddTrade(start.Add(310*time.Second), 3, 6)
	chart.AddTrade(start.Add(370*time.Second), 54, 36)

	var c4 = chart.Candles[3]
	if !(c4.Count == 1 && c4.Volume == 36 && c4.Open == 54 && c4.Close == 54 &&
		c4.High == 54 && c4.Low == 54) {
		t.Logf("Got wrong c4 val: %v", c4)
		t.Fail()
	}

	var c5 = chart.Candles[0]
	if !(c5.Count == 2 && c5.Volume == 6 && c5.Open == 3 && c5.Close == 12 &&
		c5.High == 12 && c5.Low == 3) {
		t.Logf("Got wrong c5 val: %v", c5)
		t.Fail()
	}
	fmt.Printf("Got cap: %v len: %v\n", cap(chart.Candles), len(chart.Candles))
	if len(chart.Candles) != 4 {
		t.Logf("Got wrong len: %v", len(chart.Candles))
		t.Fail()
	}
}
