package logutils_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/Drathveloper/uslogs/internal/logutils"
)

var testPatterns = []logutils.MaskPattern{
	logutils.NewMaskPattern("password=", '*', '&', ' ', '\n'),
	logutils.NewMaskPattern("token=", '*', '&', ' ', '\n'),
	logutils.NewMaskPattern("secret=", '*', '&', ' ', '\n'),
	logutils.NewMaskPattern("apikey=", '*', '&', ' ', '\n'),
}

var testDictionary = []string{
	"password=",
	"token=",
	"secret=",
	"apikey=",
}

func generateInput(size int) []byte {
	base := []byte("user=john&password=SuperSecret123&token=ABC123XYZ&apikey=987654321&location=EU ")
	return bytes.Repeat(base, size/len(base)+1)[:size]
}

func BenchmarkMask_1KB(b *testing.B) {
	m := logutils.NewMasker(testDictionary...)
	in := generateInput(1024)
	buf := make([]byte, len(in))
	copy(buf, in)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = m.Mask(buf, testPatterns)
	}
}

func BenchmarkMaskRegex_1KB(b *testing.B) {
	in := generateInput(1 * 1024)
	buf := make([]byte, len(in))
	regexMaskPatterns := []string{
		"password=([^&]*)&",
		"token=([^&]*)&",
		"secret=([^&]*)&",
		"apikey=([^&]*)&",
	}

	regexes := make([]*regexp.Regexp, len(regexMaskPatterns))
	for i, p := range regexMaskPatterns {
		regexes[i] = regexp.MustCompile(p)
	}
	copy(buf, in)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, re := range regexes {
			_ = re.ReplaceAllFunc(buf, func(m []byte) []byte {
				for x := range m {
					m[x] = '*'
				}
				return m
			})
		}
	}
}

func BenchmarkMask_8KB(b *testing.B) {
	m := logutils.NewMasker(testDictionary...)
	in := generateInput(8 * 1024)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, len(in))
		copy(buf, in)
		_ = m.Mask(buf, testPatterns)
	}
}

func BenchmarkMask_4_64KB(b *testing.B) {
	m := logutils.NewMasker(testDictionary...)
	in := generateInput(64 * 1024)
	buf := make([]byte, len(in))
	copy(buf, in)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = m.Mask(buf, testPatterns)
	}
}

func BenchmarkMaskRegex_4_64KB(b *testing.B) {
	in := generateInput(64 * 1024)
	buf := make([]byte, len(in))
	regexMaskPatterns := []string{
		"(password=)([^&]*)&",
		"(token=)([^&]*)&",
		"(secret)=([^&]*)&",
		"(apikey)=([^&]*)&",
	}

	regexes := make([]*regexp.Regexp, len(regexMaskPatterns))
	for i, p := range regexMaskPatterns {
		regexes[i] = regexp.MustCompile(p)
	}
	copy(buf, in)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, re := range regexes {
			_ = re.ReplaceAllFunc(buf, func(m []byte) []byte {
				idx := re.FindSubmatchIndex(m)
				if idx == nil || len(idx) < 4 {
					return m
				}

				startValue := idx[2]
				endValue := idx[3]

				for i := startValue; i < endValue; i++ {
					m[i] = '*'
				}

				return m
			})
		}
	}
}

func BenchmarkMask_16_64KB(b *testing.B) {
	m := logutils.NewMasker(testDictionary...)
	in := generateInput(64 * 1024)
	buf := make([]byte, len(in))
	copy(buf, in)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = m.Mask(buf, testPatterns)
	}
}

func BenchmarkMaskRegex_16_64KB(b *testing.B) {
	in := generateInput(64 * 1024)
	buf := make([]byte, len(in))
	regexMaskPatterns := []string{
		"(password=)([^&]*)&",
		"(token=)([^&]*)&",
		"(secret)=([^&]*)&",
		"(apikey)=([^&]*)&",
		"(apikey)=([^&]*)&",
	}

	regexes := make([]*regexp.Regexp, len(regexMaskPatterns))
	for i, p := range regexMaskPatterns {
		regexes[i] = regexp.MustCompile(p)
	}
	copy(buf, in)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, re := range regexes {
			_ = re.ReplaceAllFunc(buf, func(m []byte) []byte {
				idx := re.FindSubmatchIndex(m)
				if idx == nil || len(idx) < 4 {
					return m
				}

				startValue := idx[2]
				endValue := idx[3]

				for i := startValue; i < endValue; i++ {
					m[i] = '*'
				}

				return m
			})
		}
	}
}
