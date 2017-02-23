[![License BSD](https://img.shields.io/badge/License-BSD-blue.svg)](http://opensource.org/licenses/BSD-3-Clause)
[![Go Report Card](https://goreportcard.com/badge/github.com/bradclawsie/bloomfilter)](https://goreportcard.com/report/github.com/bradclawsie/bloomfilter)
[![GoDoc](https://godoc.org/github.com/bradclawsie/httpshutdown?status.svg)](http://godoc.org/github.com/bradclawsie/bloomfilter)
[![Build Status](https://travis-ci.org/bradclawsie/bloomfilter.png)](https://travis-ci.org/bradclawsie/bloomfilter)

## bloomfilter

This package implements a bloom filter in Go.

http://en.wikipedia.org/wiki/Bloom_filter

provides an explanation of what a bloom filter is. Essentially it is a probabilistic memebership
function with good size characteristics. For example, we may wish to read in the words from
the dictionary file into the filter. Over 99% of the words can be entered into the filter before
a collision occurs.

The approach in this package for hashing items into the filter is to obtain the 160 bit SHA1
hash of the original input item, which should give a good distribution. Then, this 160 bit
value is decomposed into five 32-bit integers which are then used as modulo (wrapping) offsets
into a BitSet (the storage mechanism used by this package). The bits at those offsets are set to
true (1).

Should an input token ever hash to five locations that are already set in the BitSet, it will
be considered a collision (false positive). Experiment with size settings that minimize collisions.

The size argument provided to the constructor is a desired size in bits for the BitSet used by the
bloom filter as a storage mechanism. This value is rounded up to a byte boundary.

