package scaff

// EvalScaffolding evaluates the quality of a scaffolding solution relative to a known
// correct scaffolding.
func EvalScaffolding(scaff []string, key []string) (float64, float64, error) {

	// Cache the order of each element in the key
	keyOrder := map[string]int{}
	for i, v := range key {
		keyOrder[v] = i
	}

	comparisons := 0.0
	correct := 0.0
	neighbors := 0.0
	neighborsCorrect := 0.0

	// For each triplet in the inferred solution, check whether 
	// that triplet is in the same order in the key.
	for i, v1 := range scaff {
		for j, v2 := range scaff {
			for k, v3 := range scaff {
				if (i < j) && (k > j) {
					comparisons += 1
					if (keyOrder[v1] < keyOrder[v2]) && (keyOrder[v2] < keyOrder[v3]) {
						correct += 1
					}
					if (j == i + 1) && (k == j + 1) {
						neighbors += 1
						if (keyOrder[v1] < keyOrder[v2]) && (keyOrder[v2] < keyOrder[v3]) {
							neighborsCorrect += 1
						}						
					}
				}
			}
		}
	}

	return float64(correct/comparisons), float64(neighborsCorrect/neighbors), nil

}


