package candlestick

import (
	"strings"

	c "github.com/alunir/candlestick/candle"
	"go.uber.org/zap/zapcore"
)

type ChartParameters struct {
	mode       c.ChartMode
	CandleNum  int
	Resolution string
}

type ParameterOption func(*ChartParameters)

func Resolution(resol string) ParameterOption {
	return func(op *ChartParameters) {
		op.Resolution = strings.Replace(resol, "\"", "", -1)
	}
}

func ChartParameter(model string, resolution string, candleNum int, ops ...ParameterOption) *ChartParameters {
	params := ChartParameters{
		mode:       c.ChartMode{Value: strings.ToUpper(model)},
		Resolution: resolution,
		CandleNum:  candleNum,
	}
	for _, option := range ops {
		option(&params)
	}
	return &params
}

func (cp ChartParameters) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("model", cp.mode.String())
	enc.AddInt("candle_num", cp.CandleNum)
	enc.AddString("resolution", cp.Resolution)
	return nil
}
