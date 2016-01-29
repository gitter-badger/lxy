package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"

)

type Links struct {
	idKey map[string]int
	idKeyRev map[int]string
	maxid int
	data map[int]map[int]float64
}

func NewLinks() Links {
	l := Links{}
	l.idKey = make(map[string]int)
	l.idKeyRev = make(map[int]string)
	l.maxid = 0
	l.data = make(map[int]map[int]float64)
	return l
}

func (l *Links) addKey(key string) {
	maxid := (*l).maxid
	(*l).idKey[key] = maxid
	(*l).idKeyRev[maxid] = key
	(*l).maxid = maxid + 1
}

func (l *Links) IntIDs() []int {
	ret := make([]int, l.Size())
	i := 0
	for k, _ := range l.idKeyRev {
		ret[i] = k
		i+=1
	}
	return ret
}

func (l *Links) StringIDs() []string {
	ret := make([]string, l.Size())
	i := 0
	for k, _ := range l.idKey {
		ret[i] = k
		i += 1
	}
	return ret
}

func (l *Links) ID(key string) int {
	if _, ok := l.idKey[key]; !ok {
		l.addKey(key)
	}
	id := l.idKey[key]
	return id
}

func (l *Links) Set(id1, id2 int, val float64) {
	if id1 > id2 {
		id2, id1 = id1, id2
	}
	if _, ok := (*l).data[id1]; !ok {
		(*l).data[id1] = make(map[int]float64)
	}
	(*l).data[id1][id2] = val
}

func (l *Links) Add(id1, id2 int, val float64) {
	if id1 > id2 {
		id2, id1 = id1, id2
	}
	if _, ok := (*l).data[id1]; !ok {
		(*l).data[id1] = make(map[int]float64)
	}
	(*l).data[id1][id2] += val
}

func (l *Links) Get(id1, id2 int) float64 {
	if id1 > id2 {
		id2, id1 = id1, id2
	}
	return (*l).data[id1][id2]
}

func (l *Links) Print() {
    for k, _ := range (*l).data {
	    for k2, v := range (*l).data[k] {
	        fmt.Println(k, k2, v)
	    }
	}
}

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

func (l *Links) Size() int {
	return len(l.idKey)
}

func (l *Links) SubsetLinksByPrefix(prefix string) {
	
}

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

