bloomfilter
===========

## About

This package attempts to implement a bloom filter in Go.

http://en.wikipedia.org/wiki/Bloom_filter

provides an explanation of what a bloom filter is. Essentially it is a probabilistic memebership
function with good size characteristics. For example, we may wish to read in the words from
the dictionary file and then test words that users enter to see if they are valid. The bloom filter
can test this with over 99% accuracy using only 100k in a data structure.

The approach in this package for hashing items into the filter is to obtain the 160 bit SHA1
hash of the original input item, which should give a good distribution. Then, this 160 bit
value is decomposed into five 32-bit integers which are then used as modulo (wrapping) offsets
into a BitSet (the storage mechanism used by this package). The bits at those offsets are set to
true (1).

Should an input token ever hash to five locations that are already set in the BitSet, it will
be considered a collision (false positive). Experiment with size settings that minimize collisions.

The size argument provided to the constructor is a desired size in bits for the BitSet used by the
bloom filter as a storage mechanism. This value is rounded up to a byte boundary.

## Installing

   $ go get github.com/bradclawsie/bloomfilter

## Docs

   $ go doc github.com/bradclawsie/bloomfilter

## Examples

The included unit test file contains an example use case of reading in a dict file from
a local path. You will need to edit the test and set that to run it.


