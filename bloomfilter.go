// A Bloom Filter implementation using SHA1 for hashing.
package bloomfilter

import (
	"io"
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"github.com/bradclawsie/bitset"
)

// BloomFilter is implemented using the bitset package.
type BloomFilter struct {
	bitset *bitset.BitSet
}

// NewBloomFilter will construct a new BloomFilter intended to model n bits.
// The BitSet constructor will round that number up to 
// the next byte boundary. The BitSet should be adequately compact.
// Values written into the bloom filter will use modulo to determine
// the index to set...meaning, overflow indexes will wrap.
// The BitSet is already concurrent safe through the use of RWMutex.
// Note: each entry into the filter sets five values, so having
// n below be less than five is nonsensical
func NewBloomFilter(n uint32) (*BloomFilter) {
	b := new(BloomFilter)
	b.bitset = bitset.NewBitSet(n)
	return b
}

// New is an alias for NewBloomFilter.
func New(n uint32) (*BloomFilter) {
	return NewBloomFilter(n)
}

// A SHA1 is 160 bits which we can decompose into 5 32-bit ints
type SHA1_ints [5]uint32

// The filter values corresponding to offsets derived from the SHA1-ints
type FilterVals [5]bool

// GetSHA1_ints will calculate the sha1 hash of a string. From this 160 bit hash, the five 32 bit ints are returned.
func GetSHA1_ints(s string) (SHA1_ints,error) {
	h := sha1.New()
	io.WriteString(h,s)
	sha1_bytes := h.Sum(nil)
	j := 4
	k := 5
	var sha1_ints SHA1_ints
	for i := 0; i < k; j += 4 {
	 	tb := sha1_bytes[i*4:j]
		// convert it into a 32 bit int
		tbuf := bytes.NewBuffer(tb)
		var u32 uint32
		err := binary.Read(tbuf,binary.LittleEndian,&u32)
		if err != nil {
			var empty_ints SHA1_ints
			return empty_ints,err
		}
		sha1_ints[i] = u32
	 	i++
	}
	return sha1_ints,nil
}

// Size will return the size of the underlying BitSet. May be greater than
// the arg provided to the constructor...the BitSet package rounds
// up to a byte boundary.
func (b *BloomFilter) Size() int {
	return b.bitset.Size()
}

// Write shall enter a true (1) value into the underlying BitSet at the
// modulo offsets described by the sha1_ints (five 32-bit ints).
// Returns a boolean indicating if there was a collision in the filter
// (meaning all indexes to be set were already set to true)
func (b *BloomFilter) Write(sha1_ints SHA1_ints) (bool,error) {
	l := uint32(b.bitset.Size())
	// warn if the filter positions have already been written
	collision := true
	for _,v := range sha1_ints {
		j := v % l
		existing_at_j,get_err := b.bitset.GetBitN(int(j))
		if get_err != nil {
			return false,get_err
		}
		collision = collision && existing_at_j
		set_err := b.bitset.SetBitN(int(j))
		if set_err != nil {
			return false,set_err
		}
	}
	return collision,nil
}

// Read the filter values for the modulo offsets for the SHA1_ints, and also
// send back a convenience bool to indicate if they were all true or not
func (b *BloomFilter) Read(sha1_ints SHA1_ints) (FilterVals,bool,error) {
	l := uint32(b.bitset.Size())
	var fv FilterVals
	all := true
	var get_err error
	for i,v := range sha1_ints {
		fv[i],get_err = b.bitset.GetBitN(int(v % l))
		if get_err != nil {
			return fv,false,get_err
		}
		all = all && fv[i]
	}
	return fv,all,nil
} 
