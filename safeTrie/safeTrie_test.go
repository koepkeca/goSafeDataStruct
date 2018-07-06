package safeTrie

import (
	"fmt"
	"testing"
)

func TestBasicTrie(t *testing.T) {
	trie := New()
	tt := []string{"example", "1", "2", "3"}
	td := make([]interface{}, len(tt))
	for k, v := range tt {
		td[k] = v
	}
	trie.Insert("test", td)
	trie.Destroy()
	return
}

func TestBasicRead(t *testing.T) {
	trie := New()
	tt := []int64{1, 2, 3, 4, 5}
	td := make([]interface{}, len(tt))
	for i, v := range tt {
		td[i] = v
	}
	e := trie.Insert("testing", td)
	rv, e := trie.Get("testing")
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(rv)
	trie.Destroy()
	return
}

func TestReadEmptyTrie(t *testing.T) {
	trie := New()
	v, e := trie.Get("example")
	if e != nil {
		t.Fatal(e)
	}
	if v != nil {
		t.Fatal("Got invalid read of empty trie")
	}
	trie.Destroy()
	return
}
