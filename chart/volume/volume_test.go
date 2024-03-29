package volume

import (
	"fmt"
	"testing"
	"time"

	"github.com/alunir/candlestick/candle"
)

func TestVolumeCandles(t *testing.T) {
	candleNum := 4
	var chart = NewVolumeChart(4.0, candle.ALL, candleNum)
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	chart.AddTrade(start, 5, 1)                      // open0
	chart.AddTrade(start.Add(5*time.Second), 25, 1)  // high0
	chart.AddTrade(start.Add(25*time.Second), 3, 1)  // low0
	chart.AddTrade(start.Add(60*time.Second), 12, 5) // close0 ohlc1

	chart.AddTrade(start.Add(70*time.Second), 13, 2) // open2 low2

	// Intentionally empty data series included here, to test flat candles
	chart.AddTrade(start.Add(240*time.Second), 15, 7) // high2 close2 ohlc3 ohlc4 ohlc-current

	if len(chart.Candles) != candleNum {
		t.Logf("Candles are not fulfilled. Size is %v", len(chart.Candles))
		t.Fail()
	}
	if err := chart.Candles[0].AssertOhlcv(t,
		start, start.Add(5*time.Second), start.Add(25*time.Second), start.Add(60*time.Second),
		5, 25, 3, 12, 4, 4); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}
	if err := chart.Candles[1].AssertOhlcv(t,
		start.Add(60*time.Second), start.Add(60*time.Second), start.Add(60*time.Second), start.Add(60*time.Second),
		12, 12, 12, 12, 4, 1); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}
	if err := chart.Candles[2].AssertOhlcv(t,
		start.Add(70*time.Second), start.Add(240*time.Second), start.Add(70*time.Second), start.Add(240*time.Second),
		13, 15, 13, 15, 4, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}
	if err := chart.Candles[3].AssertOhlcv(t,
		start.Add(240*time.Second), start.Add(240*time.Second), start.Add(240*time.Second), start.Add(240*time.Second),
		15, 15, 15, 15, 4, 1); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}
	if err := chart.LastCandle.AssertOhlcv(t,
		start.Add(240*time.Second), start.Add(240*time.Second), start.Add(240*time.Second), start.Add(240*time.Second),
		15, 15, 15, 15, 4, 1); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}
	if err := chart.CurrentCandle.AssertOhlcv(t,
		start.Add(240*time.Second), start.Add(240*time.Second), start.Add(240*time.Second), start.Add(240*time.Second),
		15, 15, 15, 15, 1, 1); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	// Candles should be like a queue
	fmt.Printf("Got cap: %v len: %v\n", cap(chart.Candles), len(chart.Candles))
	chart.AddTrade(start.Add(300*time.Second), 10, 2)
	chart.AddTrade(start.Add(310*time.Second), 3, 6)
	chart.AddTrade(start.Add(370*time.Second), 54, 36) // ohlc0 ohlc-current

	if err := chart.Candles[0].AssertOhlcv(t,
		start.Add(370*time.Second), start.Add(370*time.Second), start.Add(370*time.Second), start.Add(370*time.Second),
		54, 54, 54, 54, 4, 1); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	fmt.Printf("Got cap: %v len: %v\n", cap(chart.Candles), len(chart.Candles))
	if len(chart.Candles) != candleNum {
		t.Logf("Candles are not fulfilled. Size is %v", len(chart.Candles))
		t.Fail()
	}

	if err := chart.LastCandle.AssertOhlcv(t,
		start.Add(370*time.Second), start.Add(370*time.Second), start.Add(370*time.Second), start.Add(370*time.Second),
		54, 54, 54, 54, 4, 1); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.CurrentCandle.AssertOhlcv(t,
		start.Add(370*time.Second), start.Add(370*time.Second), start.Add(370*time.Second), start.Add(370*time.Second),
		54, 54, 54, 54, 1, 1); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}
}
