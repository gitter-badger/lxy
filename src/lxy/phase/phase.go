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

/*
e.g.
gb build lxy; lxy phase infer --links data/test.var.links --output data/test.var.out
*/

func Phase(links *util.Links) {

	if (*links).Size() <= 0 {
		glog.Fatal("error: link object empty, cannot phase without link data")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	m := ga.NewMultiMutator()
	//inv := new(HaplotypeInversionMutator)
	msw := new(ga.GASwitchMutator)
	//m.Add(inv)
	m.Add(msw)

	param := ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector:    ga.NewGATournamentSelector(0.7, 5),
		Breeder:     new(ga.GA2PointBreeder),
		Mutator:     m,
		PMutate:     0.7,
		PBreed:      0.7}

	//gao := ga.NewGA(param)
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



