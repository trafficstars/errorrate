package errorrate

// relaxation-exponent curve that represents errors rate (as a probability values — in [0:1))

import (
	"bytes"
	"encoding/json"
	"math"
	"math/rand"
)

const (
	errorProbabilityInertness        = 1000     // In events amount. It defines how much event is required to essentially change errorProbability value.
	errorProbabilityThreshold        = 0.67     // On which level of error probability it's required to bounce (IsExceeded() -> true).
	errorCounterTestRandomPassFactor = 1.0 / 16 // If we reached the threshold then we still need to pass events sometimes (to be able to down the probability value back down in future). The more this value the more new events are passed even the threshold is reached.
)

type handler struct {
	errorProbability atomicFloat64
}

// Handler represents an interface of relaxation-exponent-curve
// error-rate handler.
// It could be used to limit the error rate of any processor
type Handler interface {
	// ConsiderEvent adds a new recent event result to the history
	ConsiderEvent(isError bool)

	// GetErrorProbability returns the probability of an error on the next try
	GetErrorProbability() float64

	// SetErrorProbability sets the probability of an error on the next try
	SetErrorProbability(float64)

	// IsExceeded checks if the error rate is exceeded and we cannot process a new event
	IsExceeded() bool

	// UnmarshalJSON is the custom JSON unmarshaler
	UnmarshalJSON(data []byte) error

	// MarshalJSON is the custom JSON Marshaler
	MarshalJSON() ([]byte, error)
}

// NewHandler creates a new ready-to-use Handler
func NewHandler() *handler {
	h := &handler{}
	h.SetErrorProbability(errorProbabilityThreshold) // We start from the threshold to be more sensitive on the start. If there will be no errors in the start then there will be no bounces.

	return h
}

// ConsiderEvent adds a new recent event result to the history
func (h *handler) ConsiderEvent(isError bool) {
	currentErrorValue := float64(0)
	if isError {
		currentErrorValue = 1
	}
	h.errorProbability.Set((h.errorProbability.Get()*errorProbabilityInertness + currentErrorValue) / (errorProbabilityInertness + 1))
}

// SetErrorProbability sets the probability of an error on the next try
func (h *handler) SetErrorProbability(newErrorProbability float64) {
	h.errorProbability.Set(newErrorProbability)
}

// GetErrorProbability returns the probability of an error on the next try
func (h *handler) GetErrorProbability() float64 {
	return h.errorProbability.Get()
}

// IsExceeded checks if the error rate is exceeded and we cannot process a new event
func (h *handler) IsExceeded() bool {
	// It's required to get current error probability (based on history of considered before events) and do not pass if it exceeded the threshold.
	// But it will be impossible to lower down the error probability value after it reach the threshold, so sometimes we still randomly pass events.
	return h.GetErrorProbability()*math.Pow(rand.Float64(), errorCounterTestRandomPassFactor) > errorProbabilityThreshold
}

func (h *handler) UnmarshalJSON(data []byte) error {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for k, v := range raw {
		switch k {
		case `error_probability`:
			if err := h.errorProbability.UnmarshalJSON(v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *handler) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')

	{
		// writting the "h.errorProbability"

		buf.WriteString(`"error_probability":`)
		b, err := h.errorProbability.MarshalJSON()
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}
