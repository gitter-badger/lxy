package phase

import (
	ga "github.com/thoj/go-galib"
	"math/rand"
)

type GAInvertMutator struct{}

func (m GAInvertMutator) Mutate(a ga.GAGenome) ga.GAGenome {
	n := a.Copy()
	p1 := rand.Intn(a.Len())
	p2 := rand.Intn(a.Len())
	if p1 > p2 {
		p1, p2 = p2, p1
	}

	n.Invert(p1, p2)

	return n
}
func (m GAInvertMutator) String() string { return "GAInvertMutator" }





