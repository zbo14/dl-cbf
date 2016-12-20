package dl_cbf

import (
	"crypto/sha1"
	"math/rand"
	"testing"
)

// TODO: add more tests

func RandBytes(maxsize int) []byte {
	size := rand.Intn(maxsize)
	bytes := make([]byte, size)
	for i, _ := range bytes {
		bytes[i] = byte(rand.Intn(256))
	}
	return bytes
}

func GenerateDatas(size int) [][]byte {
	datas := make([][]byte, size)
	for i, _ := range datas {
		datas[i] = RandBytes(32)
	}
	return datas
}

var ht32, _ = NewHashTable(32, 4, 10, 10000, sha1.New())
var datas = GenerateDatas(10000000)

// 32-bit framework

// Testing

func TestHt32(t *testing.T) {
	var idxs []int
	// Add
	for i := 0; i < 10; i++ {
		idx, success := ht32.Add(datas[i])
		if !success {
			t.Error("Could not add data")
		}
		idxs = append(idxs, idx)
	}
	// Lookup
	for i := 0; i < 10; i++ {
		idx, success := ht32.Lookup(datas[i])
		if !success {
			t.Error("Could not find data")
		}
		if idx != idxs[i] {
			t.Errorf("Expected idx=%d; got idx=%d\n", idxs[i], idx)
		}
	}
	// Delete
	_, success := ht32.Delete(datas[0])
	if !success {
		t.Error("Could not delete data")
	}
	// Try to lookup
	_, success = ht32.Lookup(datas[0])
	if success {
		t.Error("Should not find deleted data")
	}

	// Add data multiple times
	for i := 0; i < 10; i++ {
		_, success = ht32.Add(datas[100])
		if !success {
			t.Error("Could not add data")
		}
	}

	// Check count
	_, count := ht32.GetCount(datas[100])
	if count != 10 {
		t.Errorf("Expected count=%d; got count=%d\n", 10, count)
	}
}

// Benching

func BenchmarkAdd32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ht32.Add(datas[i])
	}
}

func BenchmarkLookup32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ht32.Lookup(datas[i])
	}
}
