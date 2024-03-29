package amount

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alunir/candlestick/candle"
	"github.com/shopspring/decimal"
)

func TestAmountCandles(t *testing.T) {
	candleNum := 6
	var chart = &AmountChart{
		Chart:   candle.NewChart(candleNum),
		Chunk:   decimal.NewFromInt(30),
		Buysell: candle.ALL,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	chart.AddTrade(start, 5, 1)                     // low0
	chart.AddTrade(start.Add(5*time.Second), 25, 1) // high0 close0
	chart.AddTrade(start.Add(25*time.Second), 5, 1) // open1 low1

	chart.AddTrade(start.Add(60*time.Second), 10, 5) // high1 close1 open2 low2
	chart.AddTrade(start.Add(70*time.Second), 20, 1) // high2 close2 open3 high3

	// Intentionally empty data series included here, to test flat candles
	chart.AddTrade(start.Add(240*time.Second), 15, 6) // low3 close3 open4 low4 high4 close4

	if err := chart.Candles[0].AssertOhlcv(t,
		start, start.Add(5*time.Second), start, start.Add(5*time.Second),
		5, 25, 5, 25, 2, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.Candles[1].AssertOhlcv(t,
		start.Add(25*time.Second), start.Add(60*time.Second), start.Add(25*time.Second), start.Add(60*time.Second),
		5, 10, 5, 10, 3.5, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.Candles[2].AssertOhlcv(t,
		start.Add(60*time.Second), start.Add(70*time.Second), start.Add(60*time.Second), start.Add(70*time.Second),
		10, 20, 10, 20, 2.75, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}
	if err := chart.Candles[3].AssertOhlcv(t,
		start.Add(70*time.Second), start.Add(70*time.Second), start.Add(240*time.Second), start.Add(240*time.Second),
		20, 20, 15, 15, 1.75, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}
	if err := chart.LastCandle.AssertOhlcv(t,
		start.Add(240*time.Second), start.Add(240*time.Second), start.Add(240*time.Second), start.Add(240*time.Second),
		15, 15, 15, 15, 2, 1); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.CurrentCandle.AssertOhlcv(t,
		start.Add(240*time.Second), start.Add(240*time.Second), start.Add(240*time.Second), start.Add(240*time.Second),
		15, 15, 15, 15, 1, 1); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if len(chart.Candles) != candleNum {
		t.Logf("Candles are not fulfilled. Size is %v", len(chart.Candles))
		t.Fail()
	}

	// Candles should be like a queue
	fmt.Printf("Got cap: %v len: %v\n", cap(chart.Candles), len(chart.Candles))
	chart.AddTrade(start.Add(300*time.Second), 10, 2)
	chart.AddTrade(start.Add(310*time.Second), 3, 6)
	chart.AddTrade(start.Add(370*time.Second), 50, 36) // open0 high0 low0 close0

	if err := chart.Candles[0].AssertOhlcv(t,
		start.Add(370*time.Second), start.Add(370*time.Second), start.Add(370*time.Second), start.Add(370*time.Second),
		50, 50, 50, 50, 0.6, 1); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}
}

func TestBuySellAmountCandles(t *testing.T) {
	t.Skip()
	candleNum := 4
	var b_chart = &AmountChart{
		Chart:   candle.NewChart(candleNum),
		Chunk:   decimal.NewFromInt(30),
		Buysell: candle.BUY,
	}
	var s_chart = &AmountChart{
		Chart:   candle.NewChart(candleNum),
		Chunk:   decimal.NewFromInt(30),
		Buysell: candle.SELL,
	}

	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	b_chart.AddTrade(start, 5, 1)
	s_chart.AddTrade(start, 5, 1)

	b_chart.AddTrade(start.Add(5*time.Second), 25, 1)
	s_chart.AddTrade(start.Add(5*time.Second), 25, 1)

	b_chart.AddTrade(start.Add(25*time.Second), 3, 1)
	s_chart.AddTrade(start.Add(25*time.Second), 3, 1)

	b_chart.AddTrade(start.Add(60*time.Second), 12, -5)
	s_chart.AddTrade(start.Add(60*time.Second), 12, -5)

	b_chart.AddTrade(start.Add(70*time.Second), 13, 2)
	s_chart.AddTrade(start.Add(70*time.Second), 13, 2)

	b_chart.AddTrade(start.Add(240*time.Second), 15, -5)
	s_chart.AddTrade(start.Add(240*time.Second), 15, -5)

	// assert
	if len(b_chart.Candles) != 4 {
		t.Logf("Got wrong len: %v", len(b_chart.Candles))
		t.Fail()
	}
	if len(s_chart.Candles) != 4 {
		t.Logf("Got wrong len: %v", len(s_chart.Candles))
		t.Fail()
	}

	// assertCandle(t, b_chart.Candles[0], 3, 2.316666666666667, 12, 15, 15, 12, start.Add(60*time.Second))
	// assertCandle(t, s_chart.Candles[0], 1, 0, 12, 12, 12, 12, start.Add(60*time.Second))

	// assertCandle(t, b_chart.Candles[1], 1, 2, 15, 15, 15, 15, start.Add(240*time.Second))
	// assertCandle(t, s_chart.Candles[1], 1, 0, 15, 15, 15, 15, start.Add(240*time.Second))

	// assertCandle(t, b_chart.Candles[2], 1, 2, 15, 15, 15, 15, start.Add(240*time.Second))
	// assertCandle(t, s_chart.Candles[2], 1, 0, 15, 15, 15, 15, start.Add(240*time.Second))

	// assertCandle(t, b_chart.Candles[3], 1, 0.9333333333333333, 15, 15, 15, 15, start.Add(240*time.Second))
	// assertCandle(t, s_chart.Candles[3], 1, 0, 15, 15, 15, 15, start.Add(240*time.Second))
}

func TestBuyAmountOnlyCandles(t *testing.T) {
	t.Skip()
	candleNum := 4
	var b_chart = &AmountChart{
		Chart:   candle.NewChart(candleNum),
		Chunk:   decimal.NewFromInt(30),
		Buysell: candle.BUY,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	b_chart.AddTrade(start, 5, 1)
	b_chart.AddTrade(start.Add(5*time.Second), 25, 1)
	b_chart.AddTrade(start.Add(25*time.Second), 3, 1)
	b_chart.AddTrade(start.Add(60*time.Second), 12, -8) // 8 * 12 > 30
	b_chart.AddTrade(start.Add(70*time.Second), 13, 3)  // 3 * 13 = 39
	b_chart.AddTrade(start.Add(240*time.Second), 15, -5)

	// assert
	if len(b_chart.Candles) != 4 {
		t.Logf("Got wrong len: %v", len(b_chart.Candles))
		t.Fail()
	}

	// count, volume, open, close, high, low
	// assertCandle(t, b_chart.Candles[0], 2, 2.184615384615385, 13, 15, 15, 13, start.Add(70*time.Second))
	// assertCandle(t, b_chart.Candles[1], 1, 2, 15, 15, 15, 15, start.Add(240*time.Second))
	// assertCandle(t, b_chart.Candles[2], 1, 2, 15, 15, 15, 15, start.Add(240*time.Second))
	// assertCandle(t, b_chart.Candles[3], 1, 0.2, 15, 15, 15, 15, start.Add(240*time.Second))
}

func TestAMOUNTGetLastCandleClock(t *testing.T) {
	candleNum := 4
	var chart = &AmountChart{
		Chart:   candle.NewChart(candleNum),
		Chunk:   decimal.NewFromInt(30),
		Buysell: candle.ALL,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context) {
		for {
			select {
			case lastCandle := <-chart.GetLastCandleUpdate():
				fmt.Printf("lastCandle is updated: %v\n", lastCandle)
			}
		}
	}(ctx)

	chart.AddTrade(start, 5, 1)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("[TestAMOUNTGetLastCandleClock] 1 lastCandle should be updated next.")
	chart.AddTrade(start.Add(5*time.Second), 25, 1)
	chart.AddTrade(start.Add(25*time.Second), 3, 1)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("[TestAMOUNTGetLastCandleClock] 2 lastCandle should be updated next.")
	chart.AddTrade(start.Add(60*time.Second), 12, 5)
	chart.AddTrade(start.Add(70*time.Second), 13, 2)
}

func TestBUYAMOUNTGetLastCandleClock(t *testing.T) {
	candleNum := 4
	var chart = &AmountChart{
		Chart:   candle.NewChart(candleNum),
		Chunk:   decimal.NewFromInt(30),
		Buysell: candle.BUY,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context) {
		for {
			select {
			case lastCandle := <-chart.GetLastCandleUpdate():
				fmt.Printf("lastCandle is updated: %v\n", lastCandle)
			}
		}
	}(ctx)

	chart.AddTrade(start, 5, 1)
	time.Sleep(100 * time.Millisecond)
	chart.AddTrade(start.Add(5*time.Second), 25, 1)
	fmt.Println("[TestBUYAMOUNTGetLastCandleClock] 1 lastCandle should be updated as follows")
	chart.AddTrade(start.Add(25*time.Second), 3, 1)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("[TestBUYAMOUNTGetLastCandleClock] 2 lastCandle should be updated as follows")
	chart.AddTrade(start.Add(60*time.Second), 12, 5)
	chart.AddTrade(start.Add(70*time.Second), 13, 2)
}
