package candlestick

import (
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
