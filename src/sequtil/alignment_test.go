package util

import (
	"testing"
	"reflect"
)

func TestParseSAMFlag(t *testing.T) {

	if f, _ := parseSAMFlag("1"); !f.multiseg {
		t.Errorf("")
	}

	if f, _ := parseSAMFlag("2"); !f.allpropper {
		t.Errorf("")
	}

	if f, _ := parseSAMFlag("4"); !f.unmapped {
		t.Errorf("")
	}

	if f, _ := parseSAMFlag("8"); !f.nextunmapped {
		t.Errorf("")
	}

	if f, _ := parseSAMFlag("2048"); !f.supplementary {
		t.Errorf("")
	}

	if f, _ := parseSAMFlag("2049"); (!f.supplementary || !f.multiseg) {
		t.Errorf("")
	}

}

func TestParseCIGAR(t *testing.T) {

	test := map[string]CIGAR {
		"100M":CIGAR{CIGARCode{100, "M"}},
		"":	CIGAR{},// not allowed
		//"H100": CIGAR{}, // Not allowed, expand to handle later
		//"100M100H100M": CIGAR{}, // Not allowed, expand to handle later
		"100H100M": CIGAR{CIGARCode{100, "H"}, CIGARCode{100, "M"}},
		"100S100M100H": CIGAR{CIGARCode{100, "S"}, CIGARCode{100, "M"}, CIGARCode{100, "H"}},
	}

	for k, v := range test {

		c, e := parseCIGAR(k)
		if e != nil {
			if !reflect.DeepEqual(v, CIGAR{}) {
				t.Errorf("test error: unexpected CIGAR error case: %s\n", k)
			}
		} else {
			if !reflect.DeepEqual(c, v) {
				t.Errorf("test error: parsed and expected CIGAR structures do not match: %s, %s", c, v)
			}
		}

	}

}

func TestParseSAMLine(t *testing.T) {

	l1 := "SRR927086.7	2048	12	20766468	60	100H81M19H	*	0	0	CAAACGTGTGCACATCCNNGAGAGCCGTGAGCAACTTGCTCAGCANACNNCTCANCTTCCANGNCNTTCNCAAGCCCAGAG	<<<???@?@?@@@???<%%33=>???@??????????????????%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%	NM:i:10	MD:Z:17T0T26C2A0G4C6G1C1A3A11	AS:i:61	XS:i:0	SA:Z:10,52681560,+,100M100S,60,1;"
	sf, _ := parseSAMFlag("2048")
	a1 := Alignment{"SRR927086.7", sf, "12", 20766468, 60,
		CIGAR{CIGARCode{100, "H"}, CIGARCode{81, "M"}, CIGARCode{19, "H"}},
		"*", 0, 0, 
		"CAAACGTGTGCACATCCNNGAGAGCCGTGAGCAACTTGCTCAGCANACNNCTCANCTTCCANGNCNTTCNCAAGCCCAGAG",
		"<<<???@?@?@@@???<%%33=>???@??????????????????%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%",
	}

	l2 := "SRR927086.6	0	11	95204105	0	99S101M	*	0	0	CAGGACATANGCGNNNGCAAGGACTTCATGTCCAAAACACCAAAAGCAATGGCAACAAAAGCCAAAATTGACAAATGAGATCTAATTAAACTAAAGAGCTTCTATATCTCTGTTTTGGTACCAGTACCATGCTGTTTTGGTTACTGTAGCCTTGTAGTATAGTTTGAAGTCAGGTAGTGTGATGCCTCCAGCTTTGTTCN	%%%%%%%%%%%%%%%%EEEHCHHHHDC@CDIIGGHHFGDGEGFEGGIIIH@IIHGIIIIIGIIIIHBCIIHIIHF>FIHHBHGBGIIFFFBD?BDDD?@@B@>B@B?B>@B@BBCCCCABDDDDB;FECFFIIIIFCGF@FBBGIFDBIIIIGEIGBDDFFEFGEGEG9IFFF;EFF<CCA:F?FCBDDFFDDDDDB:1NM:i:1	MD:Z:100T0	AS:i:100	XS:i:100	SA:Z:8,98312769,-,103M97S,0,4;"
	a2 := Alignment{
		"SRR927086.6", SAMFlag{}, "11", 95204105, 0, 
		CIGAR{CIGARCode{99, "S"}, CIGARCode{101, "M"}, }, 
		"*", 0, 0,
		"CAGGACATANGCGNNNGCAAGGACTTCATGTCCAAAACACCAAAAGCAATGGCAACAAAAGCCAAAATTGACAAATGAGATCTAATTAAACTAAAGAGCTTCTATATCTCTGTTTTGGTACCAGTACCATGCTGTTTTGGTTACTGTAGCCTTGTAGTATAGTTTGAAGTCAGGTAGTGTGATGCCTCCAGCTTTGTTCN", 
		"%%%%%%%%%%%%%%%%EEEHCHHHHDC@CDIIGGHHFGDGEGFEGGIIIH@IIHGIIIIIGIIIIHBCIIHIIHF>FIHHBHGBGIIFFFBD?BDDD?@@B@>B@B?B>@B@BBCCCCABDDDDB;FECFFIIIIFCGF@FBBGIFDBIIIIGEIGBDDFFEFGEGEG9IFFF;EFF<CCA:F?FCBDDFFDDDDDB:1NM:i:1",
	}

	check := map[string]Alignment{
		l1: a1,
		l2: a2,
	}

	for k, v := range check {
		a, e := parseSAMLine(k)
		if e != nil {
			t.Errorf("test error: error when attempting to parse sam line: %s\n", e)
		}
		if !reflect.DeepEqual(a, v) {
			t.Errorf("test error: parsed sam line does not match key")
		}
	}

}

func TestGenomePositions(t *testing.T) {

	a1 := Alignment{"SRR927086.7", SAMFlag{}, "12", 20766468, 60,
		CIGAR{CIGARCode{100, "H"}, CIGARCode{4, "M"}, CIGARCode{19, "H"}},
		"*", 0, 0, 
		"CAAA",
		"<<<???@?@?@@@???<%%33=>???@??????????????????%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%",
	}

	a2 := Alignment{
		"SRR927086.6", SAMFlag{}, "11", 95204105, 0, 
		CIGAR{CIGARCode{1, "S"}, CIGARCode{3, "M"}, }, 
		"*", 0, 0,
		"CAGG", 
		"%%%%%%%%%%%%%%%%EEEHCHHHHDC@CDIIGGHHFGDGEGFEGGIIIH@IIHGIIIIIGIIIIHBCIIHIIHF>FIHHBHGBGIIFFFBD?BDDD?@@B@>B@B?B>@B@BBCCCCABDDDDB;FECFFIIIIFCGF@FBBGIFDBIIIIGEIGBDDFFEFGEGEG9IFFF;EFF<CCA:F?FCBDDFFDDDDDB:1NM:i:1",
	}
	a := []Alignment{a1, a2}

	gp1 := map[int]string {
		20766468: "C",
		20766469: "A",
		20766470: "A",
		20766471: "A",
	}

	gp2 := map[int]string {
		95204106: "A",
		95204107: "G",
		95204108: "G",
	}
	gp := []map[int]string{gp1, gp2}

	for i,v := range gp {
		check, _ := GenomePositions(a[i])
		if !reflect.DeepEqual(check, v) {
			t.Errorf("")
		}
	}

}

func TestGetVariants(t *testing.T) {

	vars := Variants{
		"1":{

			20766468: variant{ 	// position 0 in 0-based
				"1", "20766468", "rs1", "A", "T", "0", "PASS", "SNP:NN",
			},
			20766470: variant{
				"1", "20766470", "rs2", "G", "A", "0", "PASS", "SNP:NN",
			},
			95204106: variant{		// position 2 in 0-based
				"1", "95204106", "rs3", "C", "G", "0", "PASS", "SNP:NN",
			},
			95204107: variant{		// position 2 in 0-based
				"1", "95204107", "rs4", "T", "G", "0", "PASS", "SNP:NN",
			},
		},
	}

	gp1 := map[int]string {
		20766468: "C",
		20766469: "A",
		20766470: "A",
		20766471: "A",
	}

	gp2 := map[int]string {
		95204106: "C",
		95204107: "G",
		95204108: "G",
	}

	varpos1 := GetVariants("1", gp1, &vars)
	varpos2 := GetVariants("1", gp2, &vars)

	key1 := map[int]string{20766468:"N", 20766470:"A"}
	key2 := map[int]string{95204106:"R", 95204107:"A"}

	if !reflect.DeepEqual(key1, varpos1) || !reflect.DeepEqual(key2, varpos2) {
		t.Errorf("test error: error filtering positions with variants")
	}

}

func TestTabulateVariantLinks(t *testing.T) {

	vp1 := map[int]string{20766468:"R", 20766470:"A"}
	vp2 := map[int]string{95204106:"R", 95204107:"A"}

	l := NewLinks()
	l.TabulateVariantLinks("1", vp1, vp2)
	l.TabulateVariantLinks("1", vp1, vp1)

	check := map[string]map[string]float64 {
		"1_20766468": map[string]float64 {
			"1_20766468": 0.0,
			"1_20766470": -1.0,
			"1_95204106": 1.0,
			"1_95204107": -1.0,
		},
		"1_20766470": map[string]float64 {
			"1_20766468": -1.0,
			"1_20766470": 0.0,
			"1_95204106": -1.0,
			"1_95204107": 1.0,
		}, 
		"1_95204106": map[string]float64 {
			"1_20766468": 1.0,
			"1_20766470": -1.0,
			"1_95204106": 0.0,
			"1_95204107": 0.0,
		}, 
		"1_95204107": map[string]float64 {
			"1_20766468": -1.0,
			"1_20766470": 1.0,
			"1_95204106": 0.0,
			"1_95204107": 0.0,
		},
	}

	for k1, v1 := range check {
		for k2, v2 :=  range v1 {
			link := l.Get(l.ID(k1), l.ID(k2))
			if link != v2 {
				t.Errorf("test error: incorrect inference of variant phase")
			}
		}
	}

}

func TestVariantLinksFromSAM(t *testing.T) {

}

func TestScaffoldLinksFromSam(t *testing.T) {

}


