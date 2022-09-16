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
		Chart:   candle.NewChart(candleNum),
		Chunk:   2,
		Buysell: candle.ALL,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	chart.AddTrade(start, 5, 1)
	chart.AddTrade(start.Add(5*time.Second), 25, 1)
	chart.AddTrade(start.Add(25*time.Second), 3, 1)

	chart.AddTrade(start.Add(60*time.Second), 12, 5)
	chart.AddTrade(start.Add(70*time.Second), 13, 2)

	// Intentionally empty data series included here, to test flat candles
	chart.AddTrade(start.Add(240*time.Second), 15, 6)

	if err := chart.Candles[0].AssertOhlcv(t, start, 5, 25, 5, 25, 2, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.Candles[1].AssertOhlcv(t, start.Add(25*time.Second), 3, 12, 3, 12, 6, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.Candles[2].AssertOhlcv(t, start.Add(70*time.Second), 13, 15, 13, 15, 8, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	// Candles should be like a queue
	fmt.Printf("Got cap: %v len: %v\n", cap(chart.Candles), len(chart.Candles))
	chart.AddTrade(start.Add(300*time.Second), 10, 2)
	chart.AddTrade(start.Add(310*time.Second), 3, 6)
	chart.AddTrade(start.Add(370*time.Second), 54, 36)

	if err := chart.Candles[3].AssertOhlcv(t, start.Add(370*time.Second), 54, 54, 54, 54, 36, 1); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.Candles[0].AssertOhlcv(t, start.Add(25*time.Second), 3, 12, 3, 12, 6, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if len(chart.Candles) != candleNum {
		t.Logf("Candles are not fulfilled. Size is %v", len(chart.Candles))
		t.Fail()
	}
}
