package phase

import (
	"testing"
)

func TestReadPhasing(t *testing.T) {

}

func TestWritePhasing(t *testing.T) {

}

func TestPhase(t *testing.T) {
	// construct toy problem and test that it gets the right answer
	
}

func TestEvalPhasing(t *testing.T) {

	phasing := []bool{true, false, true}
	key := []bool{true, false, true}

	if s, _ := EvalPhasing(phasing, key); s != 1.0 {
		t.Errorf("phase/eval: failed to correctly score correct phasing as correct")
	}

}