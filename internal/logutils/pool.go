package logutils

import "sync"

const (
	// The magic number 36 is the result of adding:
	// date size in RFC3339 that it's 25 bytes
	// 2 x separator and two spaces around them, that it's 6 bytes
	// level in string format that can be 4 or 5 bytes, so we choose 5 to avoid the calc.
	magicNumber = 36

	extraSmallPoolSize = 1*1024 + magicNumber
	smallPoolSize      = 4*1024 + magicNumber
	mediumPoolSize     = 16*1024 + magicNumber
	largePoolSize      = 32*1024 + magicNumber
	extraLargePoolSize = 64*1024 + magicNumber
)

// PoolID represents the size of the pool.
type PoolID int

const (
	extraSmallPoolID PoolID = iota
	smallPoolID
	mediumPoolID
	largePoolID
	extraLargePoolID
)

// ResponsivePool is a map of pools that are responsive to the size of the data.
type ResponsivePool map[PoolID]*sync.Pool

// BytesPools is the pool of byte slices.
//
//nolint:gochecknoglobals
var BytesPools = ResponsivePool{
	extraSmallPoolID: {
		New: func() any {
			b := make([]byte, 0, extraSmallPoolSize)
			return &b
		},
	},
	smallPoolID: {
		New: func() any {
			b := make([]byte, 0, smallPoolSize)
			return &b
		},
	},
	mediumPoolID: {
		New: func() any {
			b := make([]byte, 0, mediumPoolSize)
			return &b
		},
	},
	largePoolID: {
		New: func() any {
			b := make([]byte, 0, largePoolSize)
			return &b
		},
	},
	extraLargePoolID: {
		New: func() any {
			b := make([]byte, 0, extraLargePoolSize)
			return &b
		},
	},
}

// SimplePool is a pool of fixed length byte slice matching 64KB.
//
//nolint:gochecknoglobals
var SimplePool = &sync.Pool{
	New: func() any {
		b := make([]byte, 0, extraLargePoolSize)
		return &b
	},
}

// GetPool returns the pool that matches the size of the data.
func (p ResponsivePool) GetPool(size int) *sync.Pool {
	switch {
	case size <= extraSmallPoolSize:
		return p[extraSmallPoolID]
	case size <= smallPoolSize:
		return p[smallPoolID]
	case size <= mediumPoolSize:
		return p[mediumPoolID]
	case size <= largePoolSize:
		return p[largePoolID]
	default:
		return p[extraLargePoolID]
	}
}

// PutPool returns to the given pool the given byte slice only if it's not too big.
func PutPool(pool *sync.Pool, bytesPtr *[]byte) {
	// We avoid returning huge byte arrays to the pool to avoid increasing memory usage.
	if len(*bytesPtr) > extraLargePoolSize {
		return
	}
	pool.Put(bytesPtr)
}
