package errorrate

import (
	"math/rand"
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
	for i := 0; i < errorProbabilityInertness*testingRedundancyFactor; i++ {
		handler.ConsiderEvent(true)
	}
	exceededCount = 0
	for i := 0; i < errorProbabilityInertness*testingRedundancyFactor; i++ {
		if handler.IsExceeded() {
			exceededCount++
		}
	}
	if exceededCount < errorProbabilityInertness*testingRedundancyFactor*0.95 {
		t.Errorf("Not enough exceeded: %v < %v*0.95", exceededCount, errorProbabilityInertness*testingRedundancyFactor)
	}

	for i := 0; i <= errorProbabilityInertness*testingRedundancyFactor; i++ {
		handler.ConsiderEvent(false)
	}

	if handler.IsExceeded() {
		t.Errorf("We have a long history without any error, but got a true on IsExceeded()")
	}
}
