package phase

import (
	"fmt"
	"os"
	"github.com/codegangsta/cli"
	util "sequtil"

	"github.com/golang/glog"
)

/*

pulling NA12878 out of 1kgenome data:
cut -f1,2,3,4,5,6,7,8,1762 data/GM12878/1kgenome/ALL.chr1.phase3_shapeit2_mvncall_integrated_v5a.20130502.genotypes.vcf | head -n 300

*/

func PhaseCommand() cli.Command {
    return cli.Command{
		Name:  "phase",
		Usage: "Commands to infer and evaluate haplotype phasing solutions",
		Subcommands: []cli.Command{
	        cli.Command{
				Name:  "infer",
				Usage: "Given a bipartite variant graph, determine the phase of the two partitions.",
				Flags: []cli.Flag{
				  	cli.StringFlag{
				  		Name: "links", 
				  		Value: "", 
				  		Usage: "Path to the Hi-C links file.",
				  	},
				  	cli.StringFlag{
				  		Name: "output", 
				  		Value: "", 
				  		Usage: "Output path for inferred phasing.",
				  	},
				},
				Action: phaseInferCommand,
    		},
	        cli.Command{
				Name:  "eval",
				Usage: "Evaluate a haplotype phasing solution.",
				Flags: []cli.Flag{
				  	cli.StringFlag{
				  		Name: "phasing", 
				  		Value: "", 
				  		Usage: "Path to phasing file.",
				  	},
				  	cli.StringFlag{
				  		Name: "key", 
				  		Value: "", 
				  		Usage: "Path to phasing key file.",
				  	},
				  	cli.StringFlag{
				  		Name: "output", 
				  		Value: "", 
				  		Usage: "Output path for evaluation stats file.",
				  	},
				},
				Action: evalPhasingCommand,
    		},
    		cli.Command{
				Name:  "prep",
				Usage: "Generate a variant links object from a set of aligned Hi-C reads.",
				Subcommands: []cli.Command{
					cli.Command{
						Name:  "varlinks",
						Usage: "Generate a variant links object from a set of aligned Hi-C reads.",
						Flags: []cli.Flag{
						  	cli.StringFlag{
						  		Name: "sam", 
						  		Value: "", 
						  		Usage: "Path to sam file.",
						  	},
						  	cli.StringFlag{
						  		Name: "vcf", 
						  		Value: "", 
						  		Usage: "Path to vcf file giving variants.",
						  	},
						  	cli.StringFlag{
						  		Name: "output", 
						  		Value: "", 
						  		Usage: "Output path for links file.",
						  	},
						},
						Action: prepPhasingCommand,
					},
				},
    		},
    	},
    }

}

func phaseInferCommand(c *cli.Context) {

	if (len(c.String("links")) == 0) {
		glog.Errorf("error: must provide a path to a links file with --links.")
		return
	}

	if _, err := os.Stat(c.String("links")); os.IsNotExist(err) {
		glog.Errorf("error: the specified links file does not exist: %s\n", c.String("links"))
		return
	}

	links, e2 := util.LoadLinks(c.String("links"))
	if e2 != nil {
		glog.Errorf("Error loading links: %s\n", e2)
	}

	Phase(&links)

}

func evalPhasingCommand(c *cli.Context) {

	if (len(c.String("phasing")) == 0) {
		glog.Errorf("error: must provide a path to a phasing file with --phasing\n")
		return
	}

	if _, err := os.Stat(c.String("phasing")); os.IsNotExist(err) {
		glog.Errorf("error: the specified phasing file does not exist: %s\n", c.String("phasing"))
		return
	}

	if (len(c.String("key")) == 0) {
		glog.Errorf("error: must provide a path to a key file with --key\n")
		return
	}

	if _, err := os.Stat(c.String("key")); os.IsNotExist(err) {
		glog.Errorf("error: the specified key file does not exist: %s\n", c.String("key"))
		return
	}

	phasing, ep := readPhasing(c.String("phasing"))
	if ep != nil {
		glog.Errorf("error: %s", ep)
		return
	}

	key, ek := readPhasing(c.String("key"))
	if ek != nil {
		glog.Errorf("error: %s\n", ek)
		return
	}

	score, e := EvalPhasing(phasing, key)
	if e != nil {
		glog.Errorf("error evaluating phasing: %s\n", e)
	}

	fmt.Printf("Evaluated phasing with score: %f\n", score)

	if (len(c.String("output"))!=0) {

	}

}

func prepPhasingCommand(c *cli.Context) {

	if (len(c.String("vcf")) == 0) {
		glog.Errorf("error: must provide a path to a vcf file with --vcf\n")
		return
	}

	if _, err := os.Stat(c.String("vcf")); os.IsNotExist(err) {
		glog.Errorf("error: the specified vcf file does not exist: %s\n", c.String("vcf"))
		return
	}

	if (len(c.String("sam")) == 0) {
		glog.Errorf("error: must provide a path to a sam file with --sam\n")
		return
	}

	if _, err := os.Stat(c.String("sam")); os.IsNotExist(err) {
		glog.Errorf("error: the specified sam file does not exist: %s\n", c.String("sam"))
		return
	}

	if (len(c.String("output")) == 0) {
		glog.Errorf("error: must provide a path to an output destination with --output\n")
		return
	}

	fmt.Println("Reading variants...")
	vcf, _ := util.ReadVariants(c.String("vcf"))

	fmt.Println("Parsing sam to variant links...")
	util.VariantLinksFromSam(c.String("sam"), c.String("output"), vcf)

}


