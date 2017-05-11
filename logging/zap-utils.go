package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// FieldSet is an alias for a slice of zapcore Fields. Functions that aggregate fields can return a FieldSet to allow
// further aggregation.
type FieldSet []zapcore.Field

// AppendFieldSet appends a FieldSet to the current FieldSet.
func (f FieldSet) AppendFieldSet(fields FieldSet) FieldSet {
	return append([]zapcore.Field(f), fields...)
}

// Append appends fields to the current FieldSet.
func (f FieldSet) Append(fields ...zapcore.Field) FieldSet {
	return append([]zapcore.Field(f), fields...)
}

func newProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        timeKey,
		LevelKey:       levelKey,
		NameKey:        nameKey,
		CallerKey:      callerKey,
		MessageKey:     messageKey,
		StacktraceKey:  stacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newProductionConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    newProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func newDevelopmentConfig() zap.Config {
	return zap.NewDevelopmentConfig()
}
