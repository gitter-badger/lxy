package scaff


func score(g *GAOrderedIntGenome) float64 {

	//fmt.Println(*g.data)

	var total float64
	
	for i, c := range g.Gene {
	
		/*for j:=1; j < (*g.data).Size(); j++ {
			if (i + j) < (*g.data).Size() {
				//fmt.Printf("adding to total %d, %d, got %d\n", c, g.Gene[i+1], (*g.data).Get(c, g.Gene[i+1]))
				total += float64(1/j)*(*g.data).Get(c, g.Gene[i+1])
			}
			if (i - j) > 0 {
				total += float64(1.0/j)*(*g.data).Get(c, g.Gene[i-1])
			}
		}*/
		

		/*
		if (i + 1) < (*g.data).Size() {
			total += (*g.data).Get(c, g.Gene[i+1])
		}
		if (i - 1) > 0 {
			total += (*g.data).Get(c, g.Gene[i-1])
		}
		if (i + 2) < (*g.data).Size() {
			total += 0.5*(*g.data).Get(c, g.Gene[i+2])
		}
		if (i - 2) > 0 {
			total += 0.5*(*g.data).Get(c, g.Gene[i-2])
		}
		if (i + 3) < (*g.data).Size() {
			total += 0.33*(*g.data).Get(c, g.Gene[i+3])
		}
		if (i - 3) > 0 {
			total += 0.33*(*g.data).Get(c, g.Gene[i-3])
		}*/


		if (i + 1) < (*g.data).Size() {
			total += (*g.data).Get(c, g.Gene[i+1])
		}
		if (i - 1) > 0 {
			total += (*g.data).Get(c, g.Gene[i-1])
		}
		if (i + 2) < (*g.data).Size() {
			total += 0.5*(*g.data).Get(c, g.Gene[i+2])
		}
		if (i - 2) > 0 {
			total += 0.5*(*g.data).Get(c, g.Gene[i-2])
		}
		if (i + 3) < (*g.data).Size() {
			total += 0.33*(*g.data).Get(c, g.Gene[i+3])
		}
		if (i - 3) > 0 {
			total += 0.33*(*g.data).Get(c, g.Gene[i-3])
		}

		if (i + 5) < (*g.data).Size() {
			total += 0.2*(*g.data).Get(c, g.Gene[i+5])
		}
		if (i - 5) > 0 {
			total += 0.2*(*g.data).Get(c, g.Gene[i-5])
		}

		if (i + 11) < (*g.data).Size() {
			total += 0.1*(*g.data).Get(c, g.Gene[i+11])
		}
		if (i - 11) > 0 {
			total += 0.1*(*g.data).Get(c, g.Gene[i-11])
		}

		if (i + 20) < (*g.data).Size() {
			total += 0.05*(*g.data).Get(c, g.Gene[i+20])
		}
		if (i - 20) > 0 {
			total += 0.05*(*g.data).Get(c, g.Gene[i-20])
		}

	}
	
	scores++

	return float64(-total)

}

func distScore(d, c1, c2 int) int {

	actualDist := (c1 - c2)
	if actualDist < 0 { actualDist = -1*actualDist }

	score := d - actualDist
	if score < 0 {return -1*score}
	return score

}



