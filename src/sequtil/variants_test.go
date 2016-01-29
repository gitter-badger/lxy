package util

import (
	"testing"
	//"fmt"
	"reflect"
)

func TestVariantIO(t *testing.T) {
	
	//fmt.Println("testing: util/WriteVariants and util/ReadVariants")

	path := "/tmp/lxy/test/testvariantio.vcf"

	v := Variants{
		"chr1": {
			1: variant{"chr1", "1", "rs1", "C", "A", "0", "PASS", "NN"},
			5: variant{"chr1", "5", "rs2", "C", "A", "0", "PASS", "NN"},
			10: variant{"chr1", "10", "rs3", "C", "A", "0", "PASS", "NN"},
		},
		"chr2": {
			1: variant{"chr2", "1", "rs4", "C", "A", "0", "PASS", "NN"},
			12: variant{"chr2", "12", "rs5", "C", "A", "0", "PASS", "NN"},
			4: variant{"chr2", "4", "rs6", "C", "A", "0", "PASS", "NN"},
		},
	}

	err := WriteVariants(v, path)
	if err != nil {
		t.Errorf("test error: TestVariantIO, WriteVariants failed\n")
	}

	v2, err2 := ReadVariants(path)
	if err2 != nil {
		t.Errorf("test error: TestVariantIO, ReadVariants failed\n")
	}

	if !reflect.DeepEqual(v, v2) {
		t.Errorf("test error: variant objects before and after reading do not match")
	}

}

