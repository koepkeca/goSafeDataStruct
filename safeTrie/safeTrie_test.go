package safeTrie

import (
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
	_, e = trie.Get("testing")
	if e != nil {
		t.Fatal(e)
	}
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

func TestSearch(t *testing.T) {
	trie := New()
	e := trie.Insert("example", []interface{}{"example"})
	if e != nil {
		panic(e)
	}
	e = trie.Insert("authority", []interface{}{"authority"})
	if e != nil {
		panic(e)
	}
	e = trie.Insert("exam", []interface{}{"exam"})
	if e != nil {
		panic(e)
	}
	tmp := trie.Search("exa")
	if len(tmp) != 2 {
		t.Fatalf("Depth search failed, expected 2 results got %d", len(tmp))
	}
	trie.Destroy()
	return
}

func TestAscii(t *testing.T) {
	pullList := []string{"abc", "zxy", "def", "xxx", "abacab", "lmnop", "ponml"}
	ordList := []string{"abacab", "abc", "def", "lmnop", "ponml", "xxx", "zxy"}
	trie := New()
	for _, next := range pullList {
		e := trie.Insert(next, []interface{}{next})
		if e != nil {
			t.Fatal(e)
		}
	}
	tmp := trie.Search("")
	for idx, next := range ordList {
		if next != tmp[idx].(string) {
			t.Fatalf("Mismatch expected %s got %s", next, tmp[idx].(string))
		}
	}
	trie.Destroy()
	return
}

func TestUTF(t *testing.T) {
	pullList := []string{"こんにちは", "こんばんは", "今日"}
	trie := New()
	e := trie.Insert("こんにちは", []interface{}{"こんにちは"})
	if e != nil {
		t.Fatal(e)
	}
	e = trie.Insert("今日", []interface{}{"今日"})
	if e != nil {
		t.Fatal(e)
	}
	e = trie.Insert("こんばんは", []interface{}{"こんばんは"})
	if e != nil {
		t.Fatal(e)
	}
	tmp := trie.Search("")
	for i, rlt := range tmp {
		if pullList[i] != rlt.(string) {
			t.Fatalf("Mismatch expected %s got %s", pullList[i], rlt.(string))
		}
	}
	trie.Destroy()
	return
}

func TestEmptyInvalidResult(t *testing.T) {
	trie := New()
	tmp := trie.Search("")
	if len(tmp) != 0 {
		t.Fatal("Empty trie got search result??")
	}
	trie.Destroy()
	trie = New()
	trie.Insert("aaa",[]interface{}{"aaa"})
	tmp = trie.Search("bbbbb")
	if len(tmp) != 0 {
		t.Fatalf("Trie search with no result obtained a result??")
	}
	trie.Destroy()
	return
}
