package phase

func score(g *GAFixedBitstringGenome) float64 {

	scores++
	total := 0.0
	for i, c := range g.Gene {

		if (i + 1) < (*g.data).Size() {
			c2 := g.Gene[i+1]
			if c != c2 {
				total -= (*g.data).Get(i, (i+1))
			} else {
				total += (*g.data).Get(i, (i+1))				
			}
		}
		if (i - 1) > 0 {
			c2 := g.Gene[i-1]
			if c != c2 {
				total -= (*g.data).Get(i, (i-1))
			} else {
				total += (*g.data).Get(i, (i-1))				
			}
		}

		/*
		if (i + 2) < (*g.data).Size() {
			c2 := g.Gene[i+2]
			if c != c2 {
				total -= 0.5*(*g.data).Get(i, (i+2))
			} else {
				total += 0.5*(*g.data).Get(i, (i+2))				
			}
		}
		if (i - 2) > 0 {
			c2 := g.Gene[i-2]
			if c != c2 {
				total -= 0.5*(*g.data).Get(i, (i-2))
			} else {
				total += 0.5*(*g.data).Get(i, (i-2))				
			}
		}*/

		/*
		if (i + 3) < (*g.data).Size() {
			c2 := g.Gene[i+3]
			if c != c2 {
				total -= 0.33*(*g.data).Get(i, (i+3))
			} else {
				total += 0.33*(*g.data).Get(i, (i+3))				
			}
		}
		if (i - 3) > 0 {
			c2 := g.Gene[i-3]
			if c != c2 {
				total -= 0.33*(*g.data).Get(i, (i-3))
			} else {
				total += 0.33*(*g.data).Get(i, (i-3))				
			}
		}*/

	}

	return float64(-total)
}



