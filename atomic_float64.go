package errorrate

import (
	"math"
	"sync/atomic"
)

type atomicFloat64 uint64

func (f *atomicFloat64) Get() float64 {
	return math.Float64frombits(atomic.LoadUint64((*uint64)(f)))
}

func (f *atomicFloat64) Set(n float64) {
	atomic.StoreUint64((*uint64)(f), math.Float64bits(n))
}
