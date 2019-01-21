package errorrate

import (
	"math"
	"strconv"
	"sync/atomic"
)

type atomicFloat64 uint64

func (f *atomicFloat64) Get() float64 {
	return math.Float64frombits(atomic.LoadUint64((*uint64)(f)))
}

func (f *atomicFloat64) Set(n float64) {
	atomic.StoreUint64((*uint64)(f), math.Float64bits(n))
}

func (f *atomicFloat64) UnmarshalJSON(data []byte) error {
	newProbability, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}
	/*if newProbability < 0 || newProbability > 1 {
		// TODO: ...
	}*/
	f.Set(newProbability)
	return nil
}

func (f *atomicFloat64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(f.Get(), 'f', -1, 64)), nil
}
