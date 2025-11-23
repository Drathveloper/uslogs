package logutils

import (
	"time"
)

func AppendTimeRFC3339(b []byte, t time.Time) []byte {
	year, month, day := t.Date()
	hour, minutes, sec := t.Clock()
	b = append4Digits(b, year)
	b = append(b, '-')
	b = append2Digits(b, int(month))
	b = append(b, '-')
	b = append2Digits(b, day)
	b = append(b, 'T')
	b = append2Digits(b, hour)
	b = append(b, ':')
	b = append2Digits(b, minutes)
	b = append(b, ':')
	b = append2Digits(b, sec)
	b = append(b, 'Z')
	return b
}

func append2Digits(b []byte, v int) []byte {
	b = append(b, '0'+byte(v/10))
	b = append(b, '0'+byte(v%10))
	return b
}

func append4Digits(b []byte, v int) []byte {
	b = append(b, '0'+byte((v/1000)%10))
	b = append(b, '0'+byte((v/100)%10))
	b = append(b, '0'+byte((v/10)%10))
	b = append(b, '0'+byte(v%10))
	return b
}
