package logutils

import (
	"log/slog"
	"strconv"
)

func AppendValue(b []byte, v slog.Value) []byte {
	switch v.Kind() {
	case slog.KindString:
		return append(b, v.String()...)
	case slog.KindInt64:
		return strconv.AppendInt(b, v.Int64(), 10)
	case slog.KindUint64:
		return strconv.AppendUint(b, v.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.AppendFloat(b, v.Float64(), 'f', -1, 64)
	case slog.KindBool:
		return strconv.AppendBool(b, v.Bool())
	default:
		return append(b, v.String()...)
	}
}

func AppendSeparator(b []byte, sep byte) []byte {
	if sep == ' ' {
		return append(b, sep)
	}
	return append(b, ' ', sep, ' ')
}
