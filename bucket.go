package dl_cbf

import (
	"bytes"
	"fmt"
)

type Bucket []*Cell

func NewBucket(size int) Bucket {
	b := make(Bucket, size)
	for i, _ := range b {
		b[i] = &Cell{}
	}
	return b
}

func (b Bucket) String(num int) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\n( B%d )\n", num))
	for _, c := range b {
		buf.WriteString(c.String())
	}
	return buf.String()
}

func (b Bucket) Members() int {
	m := 0
	for _, c := range b {
		m += int(c.count)
	}
	return m
}

func (b Bucket) Insert(fp Fingerprint) bool {

	// Inserts fp in matching or open cell
	// Returns false if neither are found or if
	// cell count exceeds threshold (shouldn't happen)

	j := -1
	for i, c := range b {
		switch CompareFingerprints(fp, c.fp) {
		case 0:
			if i == 0 || j < 0 {
				return false
			}
			return b[j].insert(fp)
		case 1:
			return c.insert(fp)
		case 2:
			if c.fp == nil {
				j = i
			}
			continue
		}

		// shouldn't get here
		panic("Unexpected result from comparison")
	}

	if j < 0 {
		return false
	}

	return b[j].insert(fp)
}

func (b Bucket) Has(fp Fingerprint) bool {
	for _, c := range b {
		switch CompareFingerprints(fp, c.fp) {
		case 0:
			return false
		case 1:
			return true
		case 2:
			continue
		}
		// shouldn't get here
		panic("Unexpected result from comparison")
	}
	return false
}

func (b Bucket) Remove(fp Fingerprint) bool {

	// Removes fp from matching cell
	// Returns false if cell is not found or
	// if cell count equals 0 (shouldn't happen)

	for _, c := range b {
		switch CompareFingerprints(fp, c.fp) {
		case 0:
			return false
		case 1:
			return c.remove()
		case 2:
			continue
		}
		// shouldn't get here
		panic("Unexpected result from comparison")
	}
	return false
}
