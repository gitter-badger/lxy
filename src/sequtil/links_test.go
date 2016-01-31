package util

import (
	"testing"
	"reflect"
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

	l := NewLinks()
	l.addKey("hi")
	l.addKey("hello")
	arr := l.IntIDs()
	key := []int{0,1}
	if !reflect.DeepEqual(arr, key) {
		t.Errorf("test links: the entity ids assigned upon creation of two new entities did not match the expectation, %d", arr)
	}

}

func TestStringIDs(t *testing.T) {

}

func TestID(t *testing.T) {


}

func TestSet(t *testing.T) {

	l := NewLinks()
	l.addKey("hi")
	l.addKey("hello")
	e := l.Set(0, 1, 0.01)
	if e != nil {
		t.Errorf("%s", e)
	}

}

func TestAdd(t *testing.T) {


}

func TestGet(t *testing.T) {


}

func TestPrint(t *testing.T) {


}

func TestWriteLinks(ot *testing.T) {


}

func TestSizeLinks(t *testing.T) {

	l := NewLinks()
	l.addKey("hi")
	l.addKey("hello")
	if l.Size() != 2 {
		t.Errorf("sequtil/links: unexpected link set size")
	}

}

func TestSubsetLinksByPrefix(t *testing.T) {


}

func TestDecode(t *testing.T) {


}

func TestSubset(t *testing.T) {


}



