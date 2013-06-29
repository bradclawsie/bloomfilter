package bloomfilter

import (
	"testing"
	"fmt"
	"io/ioutil"
	"bytes"
	"crypto/rand"
	"encoding/base64"
)

func TestFilter(t *testing.T) {
	dict_file := "" 
	// SET THIS LINE TO *YOUR* DICT FILE
	//dict_file = "/usr/share/dict/american-english" // ubuntu
	if dict_file == "" {
		fmt.Printf("\n\n****\nset dict_file in TestFilter to be a full path to a dictionary file, and rerun\n****\n\n\n")
		return
	}
	var size uint32 = 800000
	bf := NewBloomFilter(size)
	dict_bytes,dict_err := ioutil.ReadFile(dict_file)
	if dict_err != nil {
		e := fmt.Sprintf("%s\n",dict_err.Error())
		t.Errorf(e)
	}
	sep := []byte("\n")
	collisions := 0
	writes := 0
	word_bytes := bytes.Split(dict_bytes,sep)
	for _,v := range word_bytes {
		sha1_ints,sha1_err := GetSHA1_ints((string(v)))
		if sha1_err != nil {
			e := fmt.Sprintf("%s\n",sha1_err.Error())
			t.Errorf(e)
		}
		collision,write_err := bf.Write(sha1_ints)
		if write_err != nil {
			e := fmt.Sprintf("%s\n",write_err.Error())
			t.Errorf(e)
		}
		_,in_filter,read_err := bf.Read(sha1_ints)
		if read_err != nil {
			e := fmt.Sprintf("%s\n",read_err.Error())
			t.Errorf(e)
		}		
		if (!in_filter) {
			e := fmt.Sprintf("%v sha1_ints do not all read as true after a write",sha1_ints)
			t.Errorf(e)
		}
		writes++
		if collision {
			collisions++
		}
	}

	fmt.Printf("dictionary insert:\n")
	fmt.Printf("writes: %d, collisions (false positives): %d\n",writes,collisions)
	rate := 100.0 - ((float64(collisions)/float64(writes)) * 100.0)
	fmt.Printf("filter of approxoimate size %d bits shows no collisions on %f pct of inserted dictionary words\n",size,rate)

	iterations := 1000000 // int(math.Pow(26,float64(strlen))) 
	for strlen := 4; strlen < 9; strlen++ {
		fmt.Printf("\nrandom strings: %d iterations of rand strings of len %d\n",iterations,strlen)
		rand_collisions := 0
		for i := 0; i < iterations; i++ {
			b := make([]byte,20) 
			rand.Read(b) 
			en := base64.StdEncoding // or URLEncoding 
			d := make([]byte, en.EncodedLen(len(b))) 
			en.Encode(d, b) 
			s := string(d)
			ss := s[0:strlen]
			sha1_ints,sha1_err := GetSHA1_ints(ss)
			if sha1_err != nil {
				e := fmt.Sprintf("%s\n",sha1_err.Error())
				t.Errorf(e)
			}
			_,in_filter,read_err := bf.Read(sha1_ints)
			if read_err != nil {
				e := fmt.Sprintf("%s\n",read_err.Error())
				t.Errorf(e)
			}		
			if in_filter {
				rand_collisions++
			}
		}
		fmt.Printf("collisions (false positives) on random strings: %d\n",rand_collisions)
		rate = 100.0 - ((float64(rand_collisions)/float64(iterations)) * 100.0)
		fmt.Printf("populated dictionary filter of approximate size %d bits shows no collisions on %f pct of random len %d words\n",size,rate,strlen)
	}
}