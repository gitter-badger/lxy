package scaff

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	
	util "sequtil"
)

/*
lxy scaff infer --links data/test/GM.1mbp.X.links --output data/test/scaff.real.longrun.out --key data/test/testkey.txt --viz data/test/GM.1mbp.X.png

*/

func ScaffoldCommand() cli.Command {
    return cli.Command{
		Name:  "scaff",
		Usage: "Order and orient contigs using Hi-C contact frequency data.",
		Subcommands: []cli.Command{
	        cli.Command{
				Name:  "infer",
				Usage: "Infer a scaffolding from link data, e.g. lxy scaff infer --links data/test/GM.1mbp.links --output data/test/testscaffolding.txt --subset X --key data/test/testkey.txt --viz data/test/testorderfig.png",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name: "debug",
						Usage: "Whether to print detailed debugging information.",
					},
				  	cli.StringFlag{
				  		Name: "links", 
				  		Value: "", 
				  		Usage: "Path to the Hi-C links file.",
				  	},
				  	cli.StringFlag{
				  		Name: "subset", 
				  		Value: "", 
				  		Usage: "Tag on basis of which to subset links.",
				  	},
				  	cli.StringFlag{
				  		Name: "output", 
				  		Value: "", 
				  		Usage: "Output path for inferred phasing.",
				  	},
				  	cli.StringFlag{
				  		Name: "key", 
				  		Value: "", 
				  		Usage: "Path to scaffolding key file.",
				  	},
				  	cli.StringFlag{
				  		Name: "viz", 
				  		Value: "", 
				  		Usage: "Path to output the ordering viz, requires --key.",
				  	},
				  	cli.IntFlag{
				  		Name: "iterations", 
				  		Value: 1000, 
				  		Usage: "Number of iterations or GA 'generations' to perform.",
				  	},
				  	cli.IntFlag{
				  		Name: "popSize", 
				  		Value: 50, 
				  		Usage: "Population size for genetic algorithm.",
				  	},
				  	cli.Float64Flag{
				  		Name: "breedProb", 
				  		Value: 0.8, 
				  		Usage: "Probability of breeding between two members of population.",
				  	},	
				  	cli.Float64Flag{
				  		Name: "mutProb", 
				  		Value: 0.8, 
				  		Usage: "Probability of a mutation occurring.",
				  	},		  	  	
				},
				Action: scaffoldInferCommand,
    		},
	        cli.Command{
				Name:  "eval",
				Usage: "Evaluate a scaffolding using a key ordering.",
				Flags: []cli.Flag{
				  	cli.StringFlag{
				  		Name: "scaffolding", 
				  		Value: "", 
				  		Usage: "Path to scaffolding file.",
				  	},
				  	cli.StringFlag{
				  		Name: "key", 
				  		Value: "", 
				  		Usage: "Path to scaffolding key file.",
				  	},
				  	cli.StringFlag{
				  		Name: "output", 
				  		Value: "", 
				  		Usage: "Output path for evaluation stats file.",
				  	},
				},
				Action: evalScaffoldingCommand,
    		},
	        cli.Command{
				Name:  "prep",
				Usage: "Generate a links file from a set of aligned Hi-C reads.",
				Flags: []cli.Flag{
				  	cli.StringFlag{
				  		Name: "sam", 
				  		Value: "",
				  		Usage: "Path to sam file.",
				  	},
				  	cli.StringFlag{
				  		Name: "output", 
				  		Value: "", 
				  		Usage: "Output path for links file.",
				  	},
				},
				Action: prepScaffoldingCommand,
    		},
    	},
    }
}

func scaffoldInferCommand(c *cli.Context) {

	if (len(c.String("links")) == 0) {
		fmt.Printf("error: must provide a path to a links file with --links\n")
		return
	}

	if _, err := os.Stat(c.String("links")); os.IsNotExist(err) {
		fmt.Printf("error: the specified links file does not exist: %s\n", c.String("links"))
		return
	}

	/*

	if (len(c.String("assembly")) == 0) {
		fmt.Printf("error: must provide a path to a assembly file with --assembly\n")
		return
	}

	if _, err := os.Stat(c.String("assembly")); os.IsNotExist(err) {
		fmt.Printf("error: the specified assembly file does not exist: %s\n", c.String("assembly"))
		return
	}*/

	links, _ := util.LoadLinks(c.String("links"))
	//fmt.Println(links)
	if len(c.String("subset")) > 0 {
		links.Subset(c.String("subset"))
	}
	
	//fmt.Println(links.IntIDs())
	//fmt.Println(links.StringIDs())
	//fmt.Println(links.Decode(links.IntIDs()))

	scaffolding := Scaffold(&links, c.String("output"))
	
	/*scaffolding := ScaffoldNew(&links)
	scaffolding, _ := links.Decode(scaffolding)
	err := WriteScaffolding(scaffolding, c.String("output"))
	if err != nil {
		fmt.Printf("Error writing scaffolding: ", err)
	}*/

	//fmt.Println(scaffolding)
	//fmt.Println(scaffDecoded)

	if len(c.String("key")) > 0 {
		key := ReadScaffolding(c.String("key"))
		//fmt.Println(key)
		score, nscore, _ := EvalScaffolding(scaffolding, key)
		fmt.Printf("Evaluated scaffolding with score %f and neighbor score %f\n", score, nscore)
		if len(c.String("viz")) > 0 {
			VisualizeScaffolding(c.String("output"), c.String("key"), c.String("viz"))
		}
	}

	/*if len(c.String("key")) > 0 {
		key := ReadScaffolding(c.String("key"))
		//fmt.Println(key)
		score, nscore, _ := EvalScaffolding(scaffolding, key)
		fmt.Printf("Evaluated scaffolding with score %f and neighbor score %f\n", score, nscore)
		if len(c.String("viz")) > 0 {
			VisualizeScaffolding(c.String("output"), c.String("key"), c.String("viz"))
		}
	}*/

}

func evalScaffoldingCommand(c *cli.Context) {
	
	scaff := ReadScaffolding(c.String("scaffolding"))
	key := ReadScaffolding(c.String("key"))
	score, _, _ := EvalScaffolding(scaff, key)
	fmt.Printf("Evaluated scaffolding with score: %f\n", score)

}


func prepScaffoldingCommand(c *cli.Context) {

	util.ScaffoldLinksFromSam(c.String("sam"), c.String("output"))

}
