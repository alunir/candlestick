# candlestick
Provides sophisticated and easy to handle candlestick chart data type in Go.

## How to use
 ```golang
 import (
    "github.com/alunir/candlestick"
 )
 // Create 10 candles of volume chart
 // which triggered on 1e7 volume traded.
 param := candlestick.ChartParameter("volume", 1e7, 10)
 chart := candlestick.NewCandlestickChart(param)
 // Traded 0.01 quantity at price 100.0
 chart.AddTrade(time.Now(), 100.0, 0.01)
 candles := chart.GetCandles()
 ```

## Run test
```
go test ./...
```
