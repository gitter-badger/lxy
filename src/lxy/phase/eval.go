package phase

import (
)

// EvalPhasing evaluates the quality of a phasing solution relative to a known
// correct phasing.
func EvalPhasing(phasing []bool, key []bool) (float64, error) {

	matches := 0.0
	comparisons := 0.0

	for i, _ := range phasing {
		for j, _ := range phasing {
			comparisons += 1
			pmatch := (phasing[i] == phasing[j])
			kmatch := (key[i] == key[j])
			if (pmatch && kmatch) || (!pmatch && !kmatch) {
				matches += 1
			}
		}
	}

	return float64(matches/comparisons), nil

}

