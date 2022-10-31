package count

import (
	"fmt"
	"testing"
	"time"

	"github.com/alunir/candlestick/candle"
)

func TestCountCandles(t *testing.T) {
	candleNum := 3
	var chart = &CountChart{
		Chart:   candle.NewChart(candleNum),
		Chunk:   2,
		Buysell: candle.ALL,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	chart.AddTrade(start, 5, 1)                     // open0 low0
	chart.AddTrade(start.Add(5*time.Second), 25, 1) // high0 close0
	chart.AddTrade(start.Add(25*time.Second), 3, 1) // open1 low1

	chart.AddTrade(start.Add(60*time.Second), 12, 5)  // high1 close1
	chart.AddTrade(start.Add(70*time.Second), 13, 2)  // open2 low2
	chart.AddTrade(start.Add(240*time.Second), 15, 6) // high2 close2

	// Intentionally empty data series included here, to test flat candles
	chart.AddTrade(start.Add(300*time.Second), 10, 2)

	if err := chart.Candles[0].AssertOhlcv(t,
		start, start.Add(5*time.Second), start, start.Add(5*time.Second),
		5, 25, 5, 25, 2, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.Candles[1].AssertOhlcv(t,
		start.Add(25*time.Second), start.Add(60*time.Second), start.Add(25*time.Second), start.Add(60*time.Second),
		3, 12, 3, 12, 6, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.Candles[2].AssertOhlcv(t,
		start.Add(70*time.Second), start.Add(240*time.Second), start.Add(70*time.Second), start.Add(240*time.Second),
		13, 15, 13, 15, 8, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	// Candles should be like a queue
	fmt.Printf("Got cap: %v len: %v\n", cap(chart.Candles), len(chart.Candles))
	chart.AddTrade(start.Add(310*time.Second), 3, 6)
	chart.AddTrade(start.Add(370*time.Second), 54, 36)

	// {2009-11-10 23:30:30 +0000 UTC 3 12 3 12 6 63 2 1}
	// {2009-11-10 23:31:15 +0000 UTC 13 15 13 15 8 116 2 1}
	// {2009-11-10 23:35:05 +0000 UTC 10 10 3 3 8 38 2 1}
	// for _, c := range chart.Candles {
	// 	fmt.Printf("%v\n", c)
	// }

	if err := chart.Candles[0].AssertOhlcv(t,
		start.Add(25*time.Second), start.Add(60*time.Second), start.Add(25*time.Second), start.Add(60*time.Second),
		3, 12, 3, 12, 6, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.Candles[2].AssertOhlcv(t,
		start.Add(300*time.Second), start.Add(300*time.Second), start.Add(310*time.Second), start.Add(310*time.Second),
		10, 10, 3, 3, 8, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if len(chart.Candles) != candleNum {
		t.Logf("Candles are not fulfilled. Size is %v", len(chart.Candles))
		t.Fail()
	}
}
