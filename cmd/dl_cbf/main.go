package main

import (
	"fmt"
	"github.com/zbo14/dl_cbf"
)

func main() {
	ht, _ := dl_cbf.NewHashTable(32, 4, 5, 20, nil)
	data1 := []byte("hello world")
	data2 := []byte("hello universe")
	data3 := []byte("deadbeef")
	ht.Add(data1)
	ht.Add(data2)
	ht.Add(data3)
	ht.Add(data3)
	ht.Add(data3)
	ht.Delete(data3)
	_, count := ht.GetCount(data3)
	fmt.Println(count)
	fmt.Println(ht)
}
