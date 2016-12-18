package dl_cbf

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"hash"
	"hash/fnv"
)

type HashTable interface {
	IsHashTable()
	Add([]byte) (int, bool)
	ConcurrentLookup([]byte) (int, bool)
	Delete([]byte) (int, bool)
	Lookup([]byte) (int, bool)
	Size() int
}

type hashTable struct {
	// produces 16,32,64-bit fp
	fp func(int) Fingerprint

	// hash function should produce output
	// with more bits than fp
	hash hash.Hash

	rem func([]byte) int

	sub int //size of subtable

	table Table
}

func (_ *hashTable) IsHashTable() {}

func (ht *hashTable) String() string {
	return ht.table.String()
}

func rem16(h []byte) int {
	r := binary.BigEndian.Uint16(h)
	return int(r)
}

func rem32(h []byte) int {
	r := binary.BigEndian.Uint32(h)
	return int(r)
}

func rem64(h []byte) int {
	r := binary.BigEndian.Uint64(h)
	return int(r)
}

func fp16(i int) Fingerprint {
	if i < 0 {
		panic("Uint16 underflow")
	} else if i > 65535 {
		panic("Uint16 overflow")
	}
	return fingerprint16(i)
}

func fp32(i int) Fingerprint {
	if i < 0 {
		panic("Uint32 underflow")
	} else if i > 4294967295 {
		panic("Uint32 overflow")
	}
	return fingerprint32(i)
}

func fp64(i int) Fingerprint {
	if i < 0 {
		panic("Uint64 underflow")
	}
	return fingerprint64(i)
}

func NewHashTable(bits uint8, b, d, t int, hash hash.Hash) (HashTable, error) {

	// b - bucket size (# of cells)
	// d - # of subtables
	// t - table size

	var fp func(int) Fingerprint
	var rem func([]byte) int

	switch bits {
	case 16:
		fp = fp16
		rem = rem16
	case 32:
		fp = fp32
		rem = rem32
	case 64:
		fp = fp64
		rem = rem64
	default:
		return nil, errors.New("Unsupported #bits")
	}

	if t%d != 0 {
		return nil, errors.New("Number of subtables should evenly divide table size")
	}

	table := make(Table, t)
	for i, _ := range table {
		table[i] = NewBucket(b)
	}
	if hash == nil {
		hash = fnv.New64()
	}
	return &hashTable{
		fp:    fp,
		hash:  hash,
		rem:   rem,
		sub:   t / d,
		table: table,
	}, nil
}

func Hash(data []byte, hash hash.Hash) []byte {
	hash.Reset()
	hash.Write(data)
	h := hash.Sum(nil)
	return h
}

// Shuffles the hash so a different
// remainder is extracted each permutation
func Shuffle(h []byte) {
	copy(h, append(h[1:], h[0]))
	i := int(h[0]) % len(h)
	h[0], h[i] = h[i], h[0]
}

func (ht *hashTable) Size() int {
	return len(ht.table)
}

func (ht *hashTable) Add(data []byte) (int, bool) {
	h := Hash(data, ht.hash)
	var bucket, members int
	var fprint Fingerprint
	var found bool
	for i := 0; i < len(ht.table); i += ht.sub {
		r := ht.rem(h)
		b, fp := r%ht.sub+i, ht.fp(r)
		if ht.table[b].Has(fp) {
			if found {
				panic("Found multiple matching fingerprints")
			}
			bucket = b
			fprint = fp
			found = true
			// can we just break here?
			// Necessary to check for other matches??
		} else if !found {
			m := ht.table[b].Members()
			if m < members || i == 0 {
				// tie goes to the leftest bucket
				bucket = b
				fprint = fp
				members = m
			}
		}
		Shuffle(h)
	}
	ht.table[bucket].Insert(fprint) //check result
	return bucket, true
}

func (ht *hashTable) Lookup(data []byte) (int, bool) {
	h := Hash(data, ht.hash)
	bucket, found := 0, false
	for i := 0; i < len(ht.table); i += ht.sub {
		r := ht.rem(h)
		b, fp := r%ht.sub+i, ht.fp(r)
		if ht.table[b].Has(fp) {
			if found {
				// we found multiple matching fingerprints
				// for the data, this shouldn't happen
				panic("Found multiple matching fingerprints")
			}
			bucket = b
			found = true
		}
		Shuffle(h)
	}
	if !found {
		return -1, false
	}
	return bucket, true
}

func (ht *hashTable) ConcurrentLookup(data []byte) (int, bool) {
	h := Hash(data, ht.hash)
	ch := make(chan int)
	// This will spawn d goroutines (one for each subtable)
	for i := 0; i < len(ht.table); i += ht.sub {
		go func(h []byte, i int) {
			r := ht.rem(h)
			b, fp := r%ht.sub+i, ht.fp(r)
			if ht.table[b].Has(fp) {
				ch <- b
			} else if ht.sub+i >= len(ht.table) {
				ch <- -1
				close(ch)
			}
		}(h, i)
		Shuffle(h)
	}
	bucket := <-ch
	if bucket < 0 {
		return -1, false
	}
	if len(ch) > 0 {
		panic("Found multiple matching remainders")
	}
	return bucket, true
}

func (ht *hashTable) Delete(data []byte) (int, bool) {
	h := Hash(data, ht.hash)
	bucket, found := 0, false
	var fprint Fingerprint
	for i := 0; i < len(ht.table); i += ht.sub {
		r := ht.rem(h)
		b, fp := r%ht.sub+i, ht.fp(r)
		if ht.table[b].Has(fp) {
			if found {
				panic("Found multiple matching remainders")
			}
			bucket = b
			found = true
			fprint = fp
		}
		Shuffle(h)
	}
	if found {
		ht.table[bucket].Remove(fprint) //check result
		return bucket, true
	}
	return -1, false
}

type Table []Bucket

func (t Table) String() string {
	var buf bytes.Buffer
	buf.WriteString("------TABLE------")
	for i, b := range t {
		buf.WriteString(b.String(i))
	}
	return buf.String()
}
