// Package bloomfilter implemented with SHA1 for hashing.
package bloomfilter

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"github.com/bradclawsie/bitset"
	"io"
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
func NewBloomFilter(n uint32) *BloomFilter {
	b := new(BloomFilter)
	b.bitset = bitset.NewBitSet(n)
	return b
}

// New is an alias for NewBloomFilter.
func New(n uint32) *BloomFilter {
	return NewBloomFilter(n)
}

// SHA1Ints is 160 bits which we can decompose into 5 32-bit ints.
type SHA1Ints [5]uint32

// FilterVals are filter values corresponding to offsets derived from the SHA1-ints.
type FilterVals [5]bool

// GetSHA1Ints will calculate the sha1 hash of a string.
// From this 160 bit hash, the five 32 bit ints are returned.
func GetSHA1Ints(s string) (SHA1Ints, error) {
	h := sha1.New()
	io.WriteString(h, s)
	sha1Bytes := h.Sum(nil)
	j := 4
	k := 5
	var sha1Ints SHA1Ints
	for i := 0; i < k; j += 4 {
		tb := sha1Bytes[i*4 : j]
		// convert it into a 32 bit int
		tbuf := bytes.NewBuffer(tb)
		var u32 uint32
		err := binary.Read(tbuf, binary.LittleEndian, &u32)
		if err != nil {
			var emptyInts SHA1Ints
			return emptyInts, err
		}
		sha1Ints[i] = u32
		i++
	}
	return sha1Ints, nil
}

// Size will return the size of the underlying BitSet. May be greater than
// the arg provided to the constructor...the BitSet package rounds
// up to a byte boundary.
func (b *BloomFilter) Size() int {
	return b.bitset.Size()
}

// Write shall enter a true (1) value into the underlying BitSet at the
// modulo offsets described by the sha1Ints (five 32-bit ints).
// Returns a boolean indicating if there was a collision in the filter
// (meaning all indexes to be set were already set to true)
func (b *BloomFilter) Write(sha1Ints SHA1Ints) (bool, error) {
	l := uint32(b.bitset.Size())
	// warn if the filter positions have already been written
	collision := true
	for _, v := range sha1Ints {
		j := v % l
		existingAtJ, getErr := b.bitset.GetBitN(int(j))
		if getErr != nil {
			return false, getErr
		}
		collision = collision && existingAtJ
		setErr := b.bitset.SetBitN(int(j))
		if setErr != nil {
			return false, setErr
		}
	}
	return collision, nil
}

// Read the filter values for the modulo offsets for the SHA1Ints, and also
// send back a convenience bool to indicate if they were all true or not
func (b *BloomFilter) Read(sha1Ints SHA1Ints) (FilterVals, bool, error) {
	l := uint32(b.bitset.Size())
	var fv FilterVals
	all := true
	var getErr error
	for i, v := range sha1Ints {
		fv[i], getErr = b.bitset.GetBitN(int(v % l))
		if getErr != nil {
			return fv, false, getErr
		}
		all = all && fv[i]
	}
	return fv, all, nil
}
