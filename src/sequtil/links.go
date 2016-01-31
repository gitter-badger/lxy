package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
)

// TODO: It may be very well worth storing the entity link data
// in matrix from instead of in this way. The initial advantage
// of this approach may be a reduction in memory but may cost
// lookup speed, worth investigating as the phase inference 
// pipeline chokes when given large numbers of variants.

// Links stores a set of entity links, such as contigs or sequence
// variants.
type Links struct {

	// A map storing the mapping between entity names and entity ids
	idKey map[string]int

	// A map storing the mapping of entity ids to entity names
	idKeyRev map[int]string
	
	// The maximum id value
	maxid int
	
	// The entity ID / entity ID mapping, storing a real number association
	// value between the two.
	data map[int]map[int]float64

}

// NewLinks instantiates a new entity link object
func NewLinks() Links {
	l := Links{}
	l.idKey = make(map[string]int)
	l.idKeyRev = make(map[int]string)
	l.maxid = 0
	l.data = make(map[int]map[int]float64)
	return l
}

// addKey takes an entity name and assigns it an ID
func (l *Links) addKey(key string) {
	maxid := (*l).maxid
	(*l).idKey[key] = maxid
	(*l).idKeyRev[maxid] = key
	(*l).maxid = maxid + 1
}

// IntIDs returns an array of the integer IDs for the entities in a
// Links object.
func (l *Links) IntIDs() []int {

	// Allocate a slice to store the returned integer entity ids
	ret := make([]int, l.Size())

	// For each entity id in the entity id map, add it to the
	// entity id slice.
	i := 0
	for k, _ := range l.idKeyRev {
		ret[i] = k
		i+=1
	}
	return ret
}

// StringIDs returns an array of string IDs for the entities in a Links
// object.
func (l *Links) StringIDs() []string {

	// Allocate a slice to store the returned string entity names
	ret := make([]string, l.Size())
	
	// For each entity name in the entity name map, add it to the
	// entity name slice.
	i := 0
	for k, _ := range l.idKey {
		ret[i] = k
		i += 1
	}
	return ret
}

// ID returns the id for a specified entity name, creating a new one if
// the entity is not yet tracked in the Links object.
func (l *Links) ID(key string) int {

	// If the entity name is not in the entity name map, add it
	if _, ok := l.idKey[key]; !ok {
		l.addKey(key)
	}

	// Return the entity id for the specified entity name
	return l.idKey[key]

}

// Set sets the association value for a pair of entity integer ID's.
func (l *Links) Set(id1, id2 int, val float64) error {

	// TODO: Consider allowing this to operate from string instead of 
	// integer ids to make code using it more readable.
	// TODO: Consider whether an error should be returned when setting for
	// two IDs not present in the Link set.

	if _, ok := l.idKeyRev[id1]; !ok {
		return fmt.Errorf("sequtil/links: cannot set value for integer id not known to Links object, %d", id1)
	}

	if _, ok := l.idKeyRev[id2]; !ok {
		return fmt.Errorf("sequtil/links: cannot set value for integer id not known to Links object, %d", id2)
	}

	// To prevent duplication in the map, always store the largest
	// of two ids as the top level key.
	if id1 > id2 {
		id2, id1 = id1, id2
	}

	// If the top level key does not exist in the map, allocate a
	// submap for it.
	if _, ok := (*l).data[id1]; !ok {
		(*l).data[id1] = make(map[int]float64)
	}

	// Set the value
	(*l).data[id1][id2] = val

	return nil

}

// Add adds a specified float value to the association value for a pair
// of entity ids.
func (l *Links) Add(id1, id2 int, val float64) error {

	if _, ok := l.idKeyRev[id1]; !ok {
		return fmt.Errorf("sequtil/links: cannot add to value for integer id not known to Links object, %d", id1)
	}

	if _, ok := l.idKeyRev[id2]; !ok {
		return fmt.Errorf("sequtil/links: cannot add to value for integer id not known to Links object, %d", id2)
	}

	// To avoid duplication in the map, use the largest of the two id
	// values as the top-level key.
	if id1 > id2 {
		id2, id1 = id1, id2
	}

	// If the top-level key does not exist in the map, allocate a
	// submap for it.
	if _, ok := (*l).data[id1]; !ok {
		(*l).data[id1] = make(map[int]float64)
	}

	// Add the additional association value to the current association value.
	(*l).data[id1][id2] += val

	return nil

}

// Get returns the association value for a pair of entity ids from a
// Links object.
//
// Get will return an error if either of the specified entity ids are not
// present in the Links object. If they are both present but there is no
// association value entry, a zero value is returned.
func (l *Links) Get(id1, id2 int) (float64, error) {

	if _, ok := l.idKeyRev[id1]; !ok {
		return -1, fmt.Errorf("sequtil/links: cannot get value for integer id not known to Links object, %d", id1)
	}

	if _, ok := l.idKeyRev[id2]; !ok {
		return -1, fmt.Errorf("sequtil/links: cannot get value for integer id not known to Links object, %d", id2)
	}

	// Given that the map only contains one copy of each pair association
	// value and that the larger of the two is used as the key for the top
	// level of the map, set the tw
	if id1 > id2 {
		id2, id1 = id1, id2
	}

	if _, ok := (*l).data[id1]; !ok {
		return 0, nil
	} else if _, ok := (*l).data[id1][id2]; !ok {
		return 0, nil
	}

	return (*l).data[id1][id2], nil

}

// Print prints the full set of entity links in a Links object as 
// intID1, intID2, value triplets.
func (l *Links) Print() {

	// For each key in the top level map of the links data map
    for k, _ := range (*l).data {

    	// For each key in the second level of the links data map
	    for k2, v := range (*l).data[k] {

	    	// Print the id1, id2, value triplet
	        fmt.Println(k, k2, v)

	    }
	}
}

// Write takes an writable os.File pointer and writes the stringID1, stringID2, value
// triplets for each association value inthe referenced Links object.
func (l *Links) Write(out *os.File) {

	header := "#"
	for k, v := range (*l).idKey {
		header = header + " " + k + ":" + strconv.Itoa(v)
	}
	out.WriteString(header)

    for k, _ := range (*l).data {
	    for k2, v := range (*l).data[k] {
	        out.WriteString(fmt.Sprintf("%s %s %f\n", l.idKeyRev[k], l.idKeyRev[k2], v))
	    }
	}
}

// Size returns the size of the Links object.
func (l *Links) Size() int {
	return len(l.idKey)
}

func (l *Links) SubsetLinksByPrefix(prefix string) {
	
}

// Decode takes an array of entity integer ids and returns an array of entity
// string names.
//
// If the provided entity id array contains an id not tracked in the Links object,
// a non-nil error will be returned.
func (l *Links) Decode(in []int) ([]string, error) {
	out := make([]string, len(in))
	for i, j := range in {
		val, ok := (*l).idKeyRev[j]
		if !ok {
			return []string{}, fmt.Errorf("Error decoding links, an input id was not recognized.")
		}
		out[i] = val
	}
	return out, nil
}

func (l *Links) Subset(tag string) {

	for s, i := range (*l).idKey {

		// if first part of s doesnt match the tag
		tagplus := tag + "_"
		if s[:len(tagplus)] != tagplus {
			// for now not deleting all data, yes wasteful of memory #alpha
			delete((*l).idKeyRev, i)
			delete((*l).idKey, s)
		}

	}

}

// the map[int]string here indicates whether positions are reference or alternate
func (l *Links) TabulateVariantLinks(chrom string, fgp1, fgp2 map[int]string) (int, int) {

	ct := 0
	balance := 0

	for k1, _ := range fgp1 {
		for k2, _ :=  range fgp2 {
			if k1 != k2 {
				if k1 < k2 {
					_, ok1 := fgp1[k2]
					_, ok2 := fgp2[k1]
					if ok1 && ok2 {
						continue // avoiding double counting
					}
				}
				id1 := fmt.Sprintf("%s_%d", chrom, k1)
				id2 := fmt.Sprintf("%s_%d", chrom, k2)
				if (fgp1[k1] != "N") && (fgp2[k2] != "N") {
					if fgp1[k1] == fgp2[k2] {
						(*l).Add((*l).ID(id1), (*l).ID(id2), 1) // in phase
						balance += 1
						ct += 1
					} else {
						(*l).Add((*l).ID(id1), (*l).ID(id2), -1) // out of phase
						balance -= 1
						ct += 1
					}
				}
			}
		}
	}

	return ct, balance

}

func LoadLinks(linksPath string) (Links, error) {

	// TODO: better parsing instead of simple line by line...

	links := NewLinks()

	lf, err := os.Open(linksPath)
	defer lf.Close()
	if err != nil {
		return Links{}, fmt.Errorf("Couldn't open input file with path %s\n", linksPath)
	}

	s := bufio.NewScanner(lf)

	ct := 0

	for s.Scan() {

		line := s.Text()
		ct += 1
		if string(line[0]) != "#" {
			
			arr := strings.Split(line, " ")
			// if arr[0], arr[1] not in links.idKey, add them

			if _, ok := links.idKey[arr[0]]; !ok {
				links.addKey(arr[0])
			}
			if _, ok := links.idKey[arr[1]]; !ok {
				links.addKey(arr[1])
			}
			id1 := links.idKey[arr[0]]
			id2 := links.idKey[arr[1]]

			val, _ := strconv.ParseFloat(arr[2], 64)
			links.Set(id1, id2, val)
		}
		
	}

	if links.Size() <= 0 {
		return links, fmt.Errorf("read %d non-header lines from a links file and resulted in an empty links object.", ct)
	}

	return links, nil

}

