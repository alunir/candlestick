# candlestick
Provides sophisticated and easy to handle candlestick chart data type in Go.

## how to use

 1. import

 ```
 import (
    "github.com/alunir/candlestick"
 )
 ```

 1. set ChartParameters
 ```
 param := candlestickChartParameters{"volume", "1e7", 10}
 ```


 1. new Chart
 ```
 chart := candlestick.NewCandlestickChart(param)
 ```

 1. add Trade
 ```
 chart.AddTrade(time.Now(), 100.0, 0.01)
 ```

 1. Candles
 ```
 candles := chart.GetCandles()
 ```
