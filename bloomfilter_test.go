package bloomfilter

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestFilter(t *testing.T) {
	dictFile := ""
	dictFile = "./words.txt"
	if dictFile == "" {
		fmt.Printf("\n\n****\nset dictFile in TestFilter to be a full path to a dictionary file, and rerun\n****\n\n\n")
		return
	}
	var size uint32 = 800000
	bf := NewBloomFilter(size)
	dictBytes, dictErr := ioutil.ReadFile(dictFile)
	if dictErr != nil {
		e := fmt.Sprintf("%s\n", dictErr.Error())
		t.Errorf(e)
	}
	sep := []byte("\n")
	collisions := 0
	writes := 0
	wordBytes := bytes.Split(dictBytes, sep)
	for _, v := range wordBytes {
		sha1Ints, sha1Err := GetSHA1Ints((string(v)))
		if sha1Err != nil {
			e := fmt.Sprintf("%s\n", sha1Err.Error())
			t.Errorf(e)
		}
		collision, writeErr := bf.Write(sha1Ints)
		if writeErr != nil {
			e := fmt.Sprintf("%s\n", writeErr.Error())
			t.Errorf(e)
		}
		_, inFilter, readErr := bf.Read(sha1Ints)
		if readErr != nil {
			e := fmt.Sprintf("%s\n", readErr.Error())
			t.Errorf(e)
		}
		if !inFilter {
			e := fmt.Sprintf("%v sha1Ints do not all read as true after a write", sha1Ints)
			t.Errorf(e)
		}
		writes++
		if collision {
			collisions++
		}
	}

	fmt.Printf("dictionary insert:\n")
	fmt.Printf("writes: %d, collisions (false positives): %d\n", writes, collisions)
	rate := 100.0 - ((float64(collisions) / float64(writes)) * 100.0)
	fmt.Printf("filter of approxoimate size %d bits shows no collisions on %f pct of inserted dictionary words\n", size, rate)

	iterations := 1000000 // int(math.Pow(26,float64(strlen)))
	for strlen := 4; strlen < 9; strlen++ {
		fmt.Printf("\nrandom strings: %d iterations of rand strings of len %d\n", iterations, strlen)
		randCollisions := 0
		for i := 0; i < iterations; i++ {
			b := make([]byte, 20)
			rand.Read(b)
			en := base64.StdEncoding // or URLEncoding
			d := make([]byte, en.EncodedLen(len(b)))
			en.Encode(d, b)
			s := string(d)
			ss := s[0:strlen]
			sha1Ints, sha1Err := GetSHA1Ints(ss)
			if sha1Err != nil {
				e := fmt.Sprintf("%s\n", sha1Err.Error())
				t.Errorf(e)
			}
			_, inFilter, readErr := bf.Read(sha1Ints)
			if readErr != nil {
				e := fmt.Sprintf("%s\n", readErr.Error())
				t.Errorf(e)
			}
			if inFilter {
				randCollisions++
			}
		}
		fmt.Printf("collisions (false positives) on random strings: %d\n", randCollisions)
		rate = 100.0 - ((float64(randCollisions) / float64(iterations)) * 100.0)
		fmt.Printf("populated dictionary filter of approximate size %d bits shows no collisions on %f pct of random len %d words\n", size, rate, strlen)
	}
}
