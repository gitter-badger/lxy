package phase

import (
	"fmt"
	ga "github.com/thoj/go-galib"
	"math/rand"
	"time"
	"os"
	"bufio"

	util "sequtil"

	"github.com/golang/glog"
)

var scores int

// Phase infers a haplotype phasing from a set of variant links
func Phase(links *util.Links) {

	if (*links).Size() <= 0 {
		glog.Fatal("error: link object empty, cannot phase without link data")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	m := ga.NewMultiMutator()
	inv := new(GAInvertMutator)
	msw := new(ga.GASwitchMutator)
	m.Add(inv)
	m.Add(msw)

	param := ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector:    ga.NewGATournamentSelector(0.7, 5),
		Breeder:     new(ga.GA2PointBreeder),
		Mutator:     m,
		PMutate:     0.7,
		PBreed:      0.7}

	gao := ga.NewGAParallel(param, 7)

	genome := NewFixedBitstringGenome(make([]bool, (*links).Size()), score)

	(*genome).data = links

	gao.Init(10, genome)

	numiter := 10000
	ct := 0

	for {

		ct += 1

		fmt.Printf("Doing iteration %d (of %d)\n", ct, numiter)

		gao.Optimize(1)
		best := gao.Best().(*GAFixedBitstringGenome)
		fmt.Println("best:", best.Score())
		
		if ct >= numiter {
			break
		}

	}

	glog.Infof("Finished optimization")
	fmt.Printf("Calls to score = %d\n", scores)
	best := gao.Best().(*GAFixedBitstringGenome)
	fmt.Println(best)

}


// readPhasing reads a phasing solution from disk and returns an array of booleans corresponding
// to the phase of variants.
//
// readPhasing will return an error if a filesystem path is provided which does not exist on
// the system. Furthermore, since readPhasing is currently designed only for diploids, an 
// error will be returned if a value is found in the phasing other than 0 or 1.
func readPhasing(path string) ([]bool, error) {
	
	pf, err := os.Open(path)
	if err != nil {
		fmt.Errorf("Couldn't open input file with path %s\n", path)
	}
	defer pf.Close()

	s := bufio.NewScanner(pf)
	phasingMap := map[int]string{}
	ct := 0
	for s.Scan() {
		val := s.Text()
		if (val != "0") && (val != "1") {
			return []bool{}, fmt.Errorf("Tried to read phasing file with more than two haplotypes, not currently supported.\n")
		}
		phasingMap[ct] = val
		ct += 1
	}

	phasing := make([]bool, ct)
	for k, v := range phasingMap{
		phasing[k] = (v == "1")
	}

	return phasing, nil

}


func writePhasing(phasing []bool) {

}


// score determine the quality score of a haplotype phasing as
// represented in a GAFixedBitstringGenome object.
func score(g *GAFixedBitstringGenome) float64 {

	scores++
	total := 0.0
	for i, c := range g.Gene {

		if (i + 1) < (*g.data).Size() {
			c2 := g.Gene[i+1]
			if c != c2 {
				val, _ := (*g.data).Get(i, (i+1))
				total -= val
			} else {
				val, _ := (*g.data).Get(i, (i+1))
				total += val
			}
		}
		if (i - 1) > 0 {
			c2 := g.Gene[i-1]
			if c != c2 {
				val, _ := (*g.data).Get(i, (i-1))
				total -= val
			} else {
				val, _ := (*g.data).Get(i, (i-1))
				total += val
			}
		}

	}

	return float64(-total)
}

// EvalPhasing evaluates the quality of a phasing solution relative to a known
// correct phasing.
func EvalPhasing(phasing, key []bool) (float64, error) {

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


