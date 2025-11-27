package logutils

import (
	"log/slog"
	"strconv"
)

const (
	numBase = 10

	floatPrecision = -1
	floatSize      = 64
)

// AppendValue appends a slog value to a byte slice.
func AppendValue(bytes []byte, value slog.Value) []byte {
	//nolint:exhaustive
	switch value.Kind() {
	case slog.KindString:
		return append(bytes, value.String()...)
	case slog.KindInt64:
		return strconv.AppendInt(bytes, value.Int64(), numBase)
	case slog.KindUint64:
		return strconv.AppendUint(bytes, value.Uint64(), numBase)
	case slog.KindFloat64:
		return strconv.AppendFloat(bytes, value.Float64(), 'f', floatPrecision, floatSize)
	case slog.KindBool:
		return strconv.AppendBool(bytes, value.Bool())
	default:
		return append(bytes, value.String()...)
	}
}

// AppendSeparator appends a separator to a byte slice.
func AppendSeparator(bytes []byte, sep byte) []byte {
	if sep == ' ' {
		return append(bytes, sep)
	}
	return append(bytes, ' ', sep, ' ')
}
