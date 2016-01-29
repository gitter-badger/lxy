package util

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"strconv"
)

type Variants map[string]map[int]variant

type variant struct {
	Chrom string // chromosome
	Pos string // position
	ID string // identifier
	Ref string // reference base(s)
	Alt string // alternate base(s)
	Qual string // quality
	Filter string // filter status (PASS vs. semicolon-delmited list of failing filters)
	Info string // additional information (semicolon-delimited list)
	// Additional fields ignored for now
}

func (v *variant) String() string {
	return strings.Join([]string{v.Chrom, v.Pos, v.ID, v.Ref, v.Alt, v.Qual, v.Filter, v.Info}, "\t")
}

func readVariant(line string) (variant, error) {
	arr := strings.Split(line, "\t")
	if len(arr) < 8 {return variant{}, fmt.Errorf("Variant line does not include the required minimum number of fields.")}
	return variant{arr[0], arr[1], arr[2], arr[3], 
		arr[4], arr[5], arr[6], arr[7],}, nil
}

func (vars *Variants) add(v variant) {
	chr := v.Chrom
	pos, _ := strconv.Atoi(v.Pos)
	if _, ok := (*vars)[chr]; !ok {
		(*vars)[chr] = make(map[int]variant)
	}
	(*vars)[chr][pos] = v
}

// ReadVariants loads a set of variants specified in a VCF file into a Variants struct
func ReadVariants(path string) (Variants, error) {

	vars := Variants{}

	in, err := os.Open(path)
	if err != nil {
		fmt.Errorf("Couldn't open VCF file %s\n", path)
		return vars, err
	}
	defer in.Close()
	s := bufio.NewScanner(in)

	for s.Scan() {

		line := s.Text()

		if string(line[0]) == "#" {
			// todo
		} else {
			v, _ := readVariant(line)
			vars.add(v)
		}

	}

	return vars, err

}

// TODO: write variants in appropriate sorted order
// WriteVariants writes a variant object to a VCF file on disk
func WriteVariants(v Variants, path string) error {

	MkdirForFile(path)

	out, err := os.Create(path)
	if err != nil {
		fmt.Errorf("Couldn't open VCF file %s\n", path)
	}
	defer out.Close()

	// Print VCF body
	for _, value := range v {
		for _, value2 := range value {
			out.WriteString(value2.String() + "\n")
		}
	}

	return err

}
