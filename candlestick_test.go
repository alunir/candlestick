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
	chart.AddTrade(start.Add(5*time.Second), 25, 1)
	chart.AddTrade(start.Add(59*time.Second), 3, 1)
	chart.AddTrade(start.Add(interval), 25, 1)

	if err := chart.GetLastCandle().AssertOhlcv(t, start.Truncate(interval), 5, 25, 3, 3, 3, 3); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}
}

func TestMarshalUnmashal(t *testing.T) {
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

	if !chart.GetLastCandle().Equal(chart.GetLastCandle()) {
		t.Errorf("last candle not equal. %v != %v", chart2, chart)
		t.Fail()
	}

	if !chart.GetCurrentCandle().Equal(chart.GetCurrentCandle()) {
		t.Errorf("current candle not equal. %v != %v", chart2, chart)
		t.Fail()
	}
}
