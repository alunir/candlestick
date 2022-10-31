package candlestick

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCandlestickChart(t *testing.T) {
	interval := time.Minute
	param := ChartParameter("time", interval, 5)
	chart := NewCandlestickChart(param)

	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC).Truncate(interval)
	chart.AddTrade(start.Add(time.Second), 5, 1)
	chart.AddTrade(start.Add(5*time.Second), 25, 1) // high
	chart.AddTrade(start.Add(59*time.Second), 3, 1) // low close
	chart.AddTrade(start.Add(interval), 25, 1)

	if c, ok := chart.GetLastCandle(); ok {
		if err := c.AssertOhlcv(t,
			start, start.Add(5*time.Second), start.Add(59*time.Second), start.Add(59*time.Second),
			5, 25, 3, 3, 3, 3); err != nil {
			t.Logf("test failed. %v", err)
			t.Fail()
		}
	}
}

func TestChartMarshalUnmashal(t *testing.T) {
	interval := time.Minute
	param := ChartParameter("time", interval, 5)
	chart := NewCandlestickChart(param)

	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC).Truncate(interval)
	chart.AddTrade(start.Add(time.Second), 5, 1)
	chart.AddTrade(start.Add(5*time.Second), 25, 1)
	chart.AddTrade(start.Add(59*time.Second), 3, 1)
	chart.AddTrade(start.Add(interval), 25, 1)

	d, err := chart.Marshal()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(d))

	chart2 := NewCandlestickChart(ChartParameter("time", interval, 5))
	if err := chart2.Unmarshal(d); err != nil {
		panic(err)
	}

	lc, _ := chart.GetLastCandle()
	lc2, _ := chart2.GetLastCandle()
	cc, _ := chart.GetCurrentCandle()
	cc2, _ := chart2.GetCurrentCandle()
	if !lc.Equal(lc2) {
		t.Errorf("last candle not equal. %v != %v", chart2, chart)
		t.Fail()
	}

	if !cc.Equal(cc2) {
		t.Errorf("current candle not equal. %v != %v", chart2, chart)
		t.Fail()
	}
}
