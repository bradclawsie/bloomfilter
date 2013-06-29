package bloomfilter

import (
	"testing"
	"fmt"
	"io/ioutil"
	"bytes"
)

func TestFilter(t *testing.T) {
	dict_file := "" 
	// SET THIS LINE TO *YOUR* DICT FILE
	// dict_file = "/usr/share/dict/american-english" // ubuntu
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
	// make sure gibberish is not found in the filter
	sha1_ints,sha1_err := GetSHA1_ints("azzxxxdddhhhu")
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
		t.Errorf("non dict word was found in filter?")
	}
	
	fmt.Printf("writes: %d, collisions (false positives): %d\n",writes,collisions)
	rate := 100.0 - ((float64(collisions)/float64(writes)) * 100.0)
	fmt.Printf("filter of size %d will be correct %f of the time\n",size,rate)
}