package errorrate

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

const (
	testingRedundancyFactor = 100
)

func TestHandler(t *testing.T) {
	rand.Seed(42)
	handler := NewHandler()

	for i := 0; i <= errorProbabilityInertness*testingRedundancyFactor; i++ {
		if handler.IsExceeded() {
			t.Errorf("Just created handler returned true on IsExceeded()")
		}
	}

	handler.ConsiderEvent(true)
	exceededCount := 0
	for i := 0; i < errorProbabilityInertness*testingRedundancyFactor; i++ {
		if handler.IsExceeded() {
			exceededCount++
		}
	}
	if exceededCount > errorProbabilityInertness*testingRedundancyFactor*0.99 {
		t.Errorf("Too many exceeded: %v > %v*0.99", exceededCount, errorProbabilityInertness*testingRedundancyFactor)
	}
	if exceededCount == 0 {
		t.Errorf("We have only one event and it's an error, but never got a true on IsExceeded()")
	}
	var wg sync.WaitGroup
	for i := 0; i < errorProbabilityInertness*testingRedundancyFactor; i++ {
		go func() {
			wg.Add(1)
			runtime.Gosched()
			handler.ConsiderEvent(true)
			wg.Done()
		}()
	}
	wg.Wait()
	exceededCount = 0
	for i := 0; i < errorProbabilityInertness*testingRedundancyFactor; i++ {
		if handler.IsExceeded() {
			exceededCount++
		}
	}
	if exceededCount < errorProbabilityInertness*testingRedundancyFactor*0.95 {
		t.Errorf("Not enough exceeded: %v < %v*0.95", exceededCount, errorProbabilityInertness*testingRedundancyFactor)
	}

	if handler.GetErrorProbability() < 0.95 {
		t.Errorf("probability (%v) < 0.95", handler.GetErrorProbability())
	}

	b, err := handler.MarshalJSON()
	if err != nil {
		t.Errorf(`Cannot MarshalJSON(): %v`, err)
	}

	handler.SetErrorProbability(0)
	if handler.GetErrorProbability() != 0 {
		t.Errorf("probability (%v) != 0", handler.GetErrorProbability())
	}

	err = handler.UnmarshalJSON(b)
	if err != nil {
		t.Errorf(`Cannot UnmarshalJSON(): %v`, err)
	}
	if handler.GetErrorProbability() < 0.95 {
		t.Errorf("probability (%v) < 0.95", handler.GetErrorProbability())
	}

	for i := 0; i <= errorProbabilityInertness*testingRedundancyFactor; i++ {
		go func() {
			wg.Add(1)
			runtime.Gosched()
			handler.ConsiderEvent(false)
			wg.Done()
		}()
	}
	wg.Wait()

	if handler.IsExceeded() {
		t.Errorf("We have a long history without any error, but got a true on IsExceeded()")
	}

	if handler.GetErrorProbability() > 0.05 {
		t.Errorf("probability (%v) > 0.05", handler.GetErrorProbability())
	}
}
