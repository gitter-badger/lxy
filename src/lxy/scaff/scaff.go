package scaff

import (
	"fmt"
	"math/rand"
	"time"
	"os"
	"bufio"
	"os/exec"

	"github.com/thoj/go-galib"
	log "github.com/Sirupsen/logrus"

	util "sequtil"
)

var scores int

/*
e.g. lxy scaff infer --links data/test/GM.1mbp.X.links --output data/test/scaff.real.longrun.out --key data/test/testkey.txt --viz data/test/GM.1mbp.X.png
*/


func WriteScaffolding(scaffolding []string, path string) error {

	out, err := os.Create(path)
	if err != nil {
		fmt.Printf("Couldn't open output file (%s) for writing: %s\n", path, err)
	}
	defer out.Close()

	for _, v := range scaffolding {
		//fmt.Println(v)
		out.WriteString(v + "\n")
	}

	return err

}

func ReadScaffolding(path string) []string {

	in, err := os.Open(path)
	if err != nil {
		fmt.Printf("Couldn't open input file (%s) for reading: %s\n", path, err)
	}
	defer in.Close()

	scaff := []string{}
	s := bufio.NewScanner(in)
	for s.Scan() {
		scaff = append(scaff, s.Text())
	}

	return scaff

}

/*
func ScaffoldNew(links *util.Links) []int {

	iterMax := 10
	breedProb := 0.6
	mutateProb := 0.6
	numProc := 7
	//genomeSize := 10
	popSize := 5
	selectRate := 0.5
	verbosity := 3

	ga := galib.New(popSize, breedProb, mutateProb, selectRate, links)

	// Run the GA for a specified number of iterations
	e := ga.Run(numProc, iterMax, verbosity)
	if e != nil {
		fmt.Println(e)
	}

	//fmt.Println(ga.Best().Genes)

	return ga.Best().Genes

}
*/

func Scaffold(links *util.Links, outPath string) []string {

	rand.Seed(time.Now().UTC().UnixNano())

	m := ga.NewMultiMutator()
	msh := new(ga.GAShiftMutator)
	msw := new(ga.GASwitchMutator)
	inv := new(GAInvertMutator)
	m.Add(msh)
	m.Add(msw)
	m.Add(inv)

	param := ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector:    ga.NewGATournamentSelector(0.7, 5),
		Breeder:     new(ga.GA2PointBreeder),
		Mutator:     m,
		PMutate:     0.6,
		PBreed:      0.2}

	gao := ga.NewGAParallel(param, 7)
	init := (*links).IntIDs()
	genome := NewOrderedIntGenome(init, score)

	(*genome).data = links

	gao.Init(40, genome)

	numiter := 5000
	ct := 0

	for {

		ct += 1
		gao.Optimize(1)
		best := gao.Best().(*GAOrderedIntGenome)
		fmt.Println("best:", best.Score())
		fmt.Printf("Doing iteration %d (of %d)\n", ct, numiter)
		if ct >= numiter {
			break
		}
	
	}

	best := gao.Best().(*GAOrderedIntGenome)
	scaffolding, _ := (*links).Decode(best.Gene)
	err := WriteScaffolding(scaffolding, outPath)
	if err != nil {
		fmt.Printf("Error writing scaffolding: ", err)
	}
	fmt.Println(scaffolding)

	fmt.Printf("Calls to score = %d\n", scores)
	fmt.Printf("%s\n", m.Stats())

	return scaffolding

}


func VisualizeScaffolding(scaffPath, keyPath, outPath string) {

	// DEV
	lxyScriptsDir := "/Users/cb/code/src/github.com/cb01/lxy/scripts"

	cmd := exec.Command("python", lxyScriptsDir + "/scaffplot.py", "--inferred", scaffPath, "--actual", keyPath, "--outpath", outPath)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Visualizing contig order dotplot...")
	err = cmd.Wait()
	if err != nil {
		log.Fatal("Command finished with error: %v", err)
	}

}



