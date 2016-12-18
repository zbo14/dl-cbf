package dl_cbf

import "fmt"

type Cell struct {
	fp    Fingerprint
	count uint8
}

func (c *Cell) String() string {
	return fmt.Sprintf("%v, count=%d | ", c.fp, c.count)
}

func (c *Cell) insert(fp Fingerprint) bool {
	if c.count == 255 {
		// cell count shouldn't reach
		// or exceed this threshold
		return false
	}
	c.count++
	if c.fp == nil {
		c.fp = fp
	}
	return true
}

func (c *Cell) remove() bool {
	if c.count == 0 {
		// shouldn't happen
		return false
	}
	c.count--
	if c.count == 0 {
		c.fp = nil
	}
	return true
}
