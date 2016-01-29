package scaff

import (
	ga "github.com/thoj/go-galib"
	//"fmt"
	"math/rand"
)

/*
type InversionMutator struct{}

// Performs an order inversion between two randomly sampled points.
func (m InversionMutator) Mutate(a ga.GAGenome) ga.GAGenome {
	
	start := rand.Intn(a.Len())
	end := rand.Intn(a.Len())
	if start > end {
		start, end = end, start
	}
	n := a.Copy()
	n.Invert(start, end)

	return n
}

func (m InversionMutator) String() string { 
	return "InversionMutator" 
}
*/

type GAInvertMutator struct{}

func (m GAInvertMutator) Mutate(a ga.GAGenome) ga.GAGenome {
	n := a.Copy()
	p1 := rand.Intn(a.Len())
	p2 := rand.Intn(a.Len())
	if p1 > p2 {
		p1, p2 = p2, p1
	}

	n.Invert(p1, p2)

	/*
	// Until you reach the center
	for {
		if p1 >= p2{
			break
		}
		n.Switch(p1, p2)
		p1 += 1
		p2 -= 1
	}*/

	return n
}
func (m GAInvertMutator) String() string { return "GAInvertMutator" }





