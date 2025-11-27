package logutils

import (
	"time"
)

// AppendTimeRFC3339 appends a time.Time to a byte slice in RFC3339 format.
func AppendTimeRFC3339(bytes []byte, timestamp time.Time) []byte {
	year, month, day := timestamp.Date()
	hour, minutes, sec := timestamp.Clock()
	bytes = append4Digits(bytes, year)
	bytes = append(bytes, '-')
	bytes = append2Digits(bytes, int(month))
	bytes = append(bytes, '-')
	bytes = append2Digits(bytes, day)
	bytes = append(bytes, 'T')
	bytes = append2Digits(bytes, hour)
	bytes = append(bytes, ':')
	bytes = append2Digits(bytes, minutes)
	bytes = append(bytes, ':')
	bytes = append2Digits(bytes, sec)
	bytes = append(bytes, 'Z')
	return bytes
}

func append2Digits(bytes []byte, value int) []byte {
	bytes = append(bytes, '0'+byte(value/10)) //nolint:mnd
	bytes = append(bytes, '0'+byte(value%10)) //nolint:mnd
	return bytes
}

func append4Digits(bytes []byte, value int) []byte {
	bytes = append(bytes, '0'+byte((value/1000)%10)) //nolint:mnd
	bytes = append(bytes, '0'+byte((value/100)%10))  //nolint:mnd
	bytes = append(bytes, '0'+byte((value/10)%10))   //nolint:mnd
	bytes = append(bytes, '0'+byte(value%10))        //nolint:mnd
	return bytes
}
