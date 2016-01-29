package util

import (
	"os"
	"fmt"
	"strings"
	"bufio"
	"strconv"

	//"github.com/codegangsta/cli"
)

type Alignment struct {
	qname string // Query template NAME
	flag SAMFlag // bitwise FLAG
	rname string // Reference sequence NAME
	pos int // 1-based leftmost mapping POSition
	mapq int // MAPping Quality
	cigar CIGAR // CIGAR string
	rnext string // Ref. name of the mate/next read
	pnext int // Position of the mate/next read
	tlen int // observed Template LENgth
	seq string // segment SEQuence
	qual string // ASCII of Phred-scaled base QUALity+33ASCII of Phred-scaled base QUALity+33
}

type SAMFlag struct {
	multiseg bool // template having multiple segments in sequencing
	allpropper bool // each segment properly aligned according to the aligner
	unmapped bool // segment unmapped
	nextunmapped bool // next segment in the template unmapped
	reversecomp bool // SEQ being reverse complemented
	nextreversecomp bool // SEQ of the next segment in the template being reverse complemented
	first bool // the first segment in the template
	last bool // the last segment in the template
	secondary bool // secondary alignment
	nopass bool // not passing filters, such as platform/vendor quality controls
	duplicate bool // PCR or optical duplicate
	supplementary bool // supplementary alignment
}


func parseSAMFlag(encodedString string) (SAMFlag, error) {

	encoded, _ := strconv.Atoi(encodedString)

	return SAMFlag{
		(encoded & 0x1 != 0),
		(encoded & 0x2 != 0),
		(encoded & 0x4 != 0),
		(encoded & 0x8 != 0),
		(encoded & 0x10 != 0),
		(encoded & 0x20 != 0),
		(encoded & 0x40 != 0),
		(encoded & 0x80 != 0),
		(encoded & 0x100 != 0),
		(encoded & 0x200 != 0),
		(encoded & 0x400 != 0),
		(encoded & 0x800 != 0),
	}, nil

}

type CIGAR []CIGARCode

type CIGARCode struct {
	value int
	code string
}

func (c *CIGAR) Add(val int, code string) {
	cc := CIGARCode{val, code}
	(*c) = append(*c, cc)
}



// if it's star, handle that appropriately <- **
func parseCIGAR(cigarString string) (CIGAR, error) {
	// for now, must contain only S, H, or M
	///disallow := map[string]int{"I":1, "D":1, "N":1, "P":1, "=":1, "X":1,}
	cigar := CIGAR{}
	lastIndex := 0
	for i, c := range cigarString {
		if c > 57 {
			//if _, ok := disallow[string(c)]; ok {
				// it's a disallowed code, return error
			//	return CIGAR{}, fmt.Errorf("Found unsupported code in CIGAR string.")
			//}
		// Process the allowed CIGAR code
		val, err := strconv.Atoi(cigarString[lastIndex:i])
		if err != nil {
			return CIGAR{}, fmt.Errorf("Error parsing CIGAR string.")
		}
		cigar.Add(val, string(c))
		lastIndex = i + 1
		}
	}
	/*for _, cig := range cigar {
		fmt.Println(cig)
	}*/
	return cigar, nil
}


func parseSAMLine(line string) (Alignment, error) {
	arr := strings.Split(line, "\t")
	if len(arr) < 11 {
		return Alignment{}, fmt.Errorf("Incorrect length alignment")
	}

	qname := arr[0]

	samflag, esf := parseSAMFlag(arr[1])
	if esf != nil {
		return Alignment{}, esf
	}

	rname := arr[2]
	if rname == "*" {
		// didn't align, return null alignment for now
		return Alignment{}, fmt.Errorf("Read did not align")
	}

	pos, ep := strconv.Atoi(arr[3])
	if ep != nil {
		return Alignment{}, ep
	}

	mapq, em := strconv.Atoi(arr[4])
	if em != nil {
		return Alignment{}, em
	}

	c, ec := parseCIGAR(arr[5])
	if ec != nil {
		return Alignment{}, ec
	}

	rnext := arr[6]

	pnext, epn := strconv.Atoi(arr[7])
	if epn != nil {
		if arr[7] == "*" {
			pnext = -1
		} else {
			return Alignment{}, epn				
		}
	}

	tlen, etl := strconv.Atoi(arr[8])
	if etl != nil {
		return Alignment{}, etl
	}

	seq := arr[9]

	qual := arr[10]

	a := Alignment{qname, samflag, 
		rname, pos, mapq, c, rnext,
		pnext, tlen, seq, qual,}

	return a, nil

}

func GenomePositions(a Alignment) (map[int]string, error) {
	
	hits := map[int]string{}
	offset := 0
	firstH := 1
	for _, c := range a.cigar {
		if c.code == "M"{
			for i := 0; i < c.value; i++ {
				ind := (a.pos + offset)
				hits[ind] = string(a.seq[offset])
				offset += 1
			}
		} else if c.code == "H" {
			if firstH == 0 {
				break
			}
			firstH = 0
		} else if c.code == "S" {
			for i := 0; i < c.value; i++ {
				offset += 1
			}
		} else {
			return map[int]string{}, fmt.Errorf("Unsupported code in CIGAR.")
		}
	}

	return hits, nil

}

func GetVariants(chrom string, gp map[int]string, vars *Variants) map[int]string {
	
	if _, ok := (*vars)[chrom]; !ok {
		//e := fmt.Errorf("Variant set didn't contain any variants for contig %s\n", chrom)
		//fmt.Println(e)
		//fmt.Printf("couldn't find chromosome %s in variant set\n", chrom)
		return map[int]string{} // no variants for this chromosome...
	} else {
		//fmt.Printf("found chromosome %s in variant set\n", chrom)
	}

	for pos, v := range gp {
		if _, ok := (*vars)[chrom][pos]; !ok {
			//fmt.Printf("didn't find position %d in variant set\n", pos)
			delete(gp, pos)
		} else {
			// mark whether the variant is reference, alternate, or something else
			//fmt.Println((*vars)[chrom][pos])
			if v == (*vars)[chrom][pos].Ref {
				gp[pos] = "R"
				//fmt.Printf("Reference: %s %s %d %s\n", v, chrom, pos, (*vars)[chrom][pos][3])
			} else if v == (*vars)[chrom][pos].Alt {
				gp[pos] = "A"
				//fmt.Printf("Alternate: %s %s %d %s\n", v, chrom, pos, (*vars)[chrom][pos][4])
			} else {
				gp[pos] = "N"
				//fmt.Printf("Other: %s %s %d\n", v, chrom, pos)
			}
		}
	}

	return gp

}

// VariantLinksFromSam parses a sam file, constructing a Links object
// representing simple counts of association between variants.
func VariantLinksFromSam(samPath, outPath string, vars Variants) {

	out, err1 := os.Create(outPath)
	if err1 != nil {
		fmt.Printf("Couldn't open output file (%s) for reading: %s\n", outPath, err1)
	}
	defer out.Close()

	in, err2 := os.Open(samPath)
	if err2 != nil {
		fmt.Printf("Couldn't open input file (%s) for reading: %s\n", samPath, err2)
	}
	defer in.Close()

	links := NewLinks()
	currentID := ""
	hits := []Alignment{}
	s := bufio.NewScanner(in)
	linkedVariants := 0
	lineCount := 0
	balance := 0
	//maxlines := 1000000

	for s.Scan() {
		line := s.Text()
		lineCount += 1
		/*if lineCount >  maxlines {
			break
		}*/
		if string(line[0]) != "@" {

			a, e := parseSAMLine(s.Text())
			if e != nil {
				//fmt.Printf("error parsing sam line, continuing: %s, %s\n", e, line)
				hits = []Alignment{}
				currentID = ""
				continue 
			}

			if (a.qname != currentID) {
				if (len(hits) == 2) && (len(currentID) > 0) {
					if hits[0].rname == hits[1].rname { // same chromosome

						gp1, egp1 := GenomePositions(hits[0])
						gp2, egp2 := GenomePositions(hits[1])
						if (egp1 != nil) || (egp2 != nil) {
							//fmt.Println("continuing")
							//fmt.Println(gp1)
							//fmt.Println(gp2)
							hits = []Alignment{}
							currentID = ""
							continue
						}
						//fmt.Println(hits[0].rname)
						varPos1 := GetVariants(hits[0].rname, gp1, &vars)
						varPos2 := GetVariants(hits[0].rname, gp2, &vars)
						// requires both on same chromosome, otherwise bug in following
						ct1, bal1 := links.TabulateVariantLinks(hits[0].rname, varPos1, varPos2)
						ct2, bal2 := links.TabulateVariantLinks(hits[0].rname, varPos1, varPos1)
						ct3, bal3 := links.TabulateVariantLinks(hits[0].rname, varPos2, varPos2)
						ct := ct1 + ct2 + ct3
						bal := bal1 + bal2 + bal3
						if ct > 0{
							//fmt.Println(varPos1)
							//fmt.Println(varPos2)
							//fmt.Println(hits)
							linkedVariants += ct
							balance += bal
							fmt.Printf("linked a total of %d variants in %d lines with balance %d\n", linkedVariants, lineCount, balance)
						}						
						//fmt.Println(varPos1)
						//fmt.Println(varPos2)
					}
				}
				currentID = a.qname
				hits = []Alignment{}
			}

			hits = append(hits, a)

		}
	}

	links.Write(out)

}



// ScaffoldLinksFromSam parses a sam file, constructing a Links object
// representing simple counts of association between contigs.
func ScaffoldLinksFromSam(samPath, outPath string) {

	// assumes sam is sorted by id
	// for now, building links map in memory but could
	// also write all to disk then sort | uniq -c

	out, err1 := os.Create(outPath)
	if err1 != nil {
		fmt.Printf("Couldn't open output file (%s) for reading: %s\n", outPath, err1)
	}
	defer out.Close()

	in, err2 := os.Open(samPath)
	if err2 != nil {
		fmt.Printf("Couldn't open input file (%s) for reading: %s\n", samPath, err2)
	}
	defer in.Close()

	links := NewLinks()
	currentID := ""
	chrHits := []string{}
	//posHits := []int{}
	s := bufio.NewScanner(in)

	for s.Scan() {
		line := s.Text()
		if string(line[0]) != "@" {

			arr := strings.Split(s.Text(), "\t")

			if arr[0] != currentID {

				if len(chrHits) == 2 {
					id1 := links.ID(chrHits[0])
					id2 := links.ID(chrHits[1])
					links.Add(id1, id2, 1)
				}

				currentID = arr[0]
				chrHits = []string{}
				//posHits = []string{}

			}

			chrHits = append(chrHits, arr[2])
			//posHits = append(posHits, arr[3])

		}
	}

	links.Write(out)

}



func FilterAlignmentIntersectingVariants() {
	
}



