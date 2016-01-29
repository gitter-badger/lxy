package util

import (

	"testing"
)


func TestAddKey(t *testing.T) {

	l := NewLinks()

	l.addKey("hi")
	val, ok := l.idKey["hi"]
	if !ok {
		t.Errorf("test error: error adding key to links set")
	}

	if _, ok := l.idKeyRev[val]; !ok {
		t.Errorf("test error: the id for key added via links.addKey() was not found in links.idKeyRev")
	}

}


func TestIntIDs(t *testing.T) {

	
	
}


func TestIDs(t *testing.T) {}
func TestID(t *testing.T) {}
func TestSet(t *testing.T) {}
func TestAdd(t *testing.T) {}
func TestGet(t *testing.T) {}
func TestPrint(t *testing.T) {}
func TestWriteLinks(ot *testing.T) {}
func TestSizeLinks(t *testing.T) {}
func TestSubsetLinksByPrefix(t *testing.T) {}
func TestDecode(t *testing.T) {}
func TestSubset(t *testing.T) {}
