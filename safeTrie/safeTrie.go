//Package safeTrie implements a thread safe trie in go
package safeTrie

import (
	"fmt"
	"sort"
)

//lexicalKeys store the keys of a resulting depth search in lexical order.
//this is useful for obtaining a lexically sorted list from a given prefix
//search term. It implements the sort interface for a slice of runes.
type lexicalKeys []rune

func (lk lexicalKeys) Len() int           { return len(lk) }
func (lk lexicalKeys) Swap(i, j int)      { lk[i], lk[j] = lk[j], lk[i] }
func (lk lexicalKeys) Less(i, j int) bool { return lk[i] < lk[j] }

//SafeTrie is the sturcture that contains the channel used to
//communicate with the trie.
type SafeTrie struct {
	op chan (func(*trie))
}

//Insert will place the array data at v into the trie at location k
func (t *SafeTrie) Insert(k string, v []interface{}) (e error) {
	if k == "" {
		e = fmt.Errorf("Insert may not have empty key.")
		return
	}
	if len(v) == 0 {
		e = fmt.Errorf("Insert may not have empty data.")
		return
	}
	t.op <- func(st *trie) {
		curr := st.root
		for _, char := range k {
			if next, ok := curr.children[char]; !ok {
				nc := newNode()
				curr.children[char] = nc
				curr = nc
			} else {
				curr = next
			}
		}
		curr.data = append(curr.data, v...)
		return
	}
	return
}

//Get retreives the data (if any) at location k
func (t *SafeTrie) Get(k string) (v []interface{}, e error) {
	if k == "" {
		e = fmt.Errorf("Get requires a string to search for")
		return
	}
	ich := make(chan []interface{})
	t.op <- func(st *trie) {
		curr := st.root
		for _, char := range k {
			if next, ok := curr.children[char]; !ok {
				ich <- nil
				return
			} else {
				curr = next
			}
		}
		ich <- curr.data
		return
	}
	v = <-ich
	return
}

//Search makes a DFS Search for the term which begins with startAt
//if it is blank, the entire trie will be returned as a lexically sorted
//slice of interfaces.
func (t *SafeTrie) Search(startAt string) []interface{} {
	rch := make(chan []interface{})
	t.op <- func(st *trie) {
		curr := st.root
		for _, char := range startAt {
			next, ok := curr.children[char]
			if !ok {
				rch <- nil
				return
			}
			curr = next
		}
		rch <- curr.getDataBelow()
		return
	}
	return <-rch
}

//trie contains the locally available trie
type trie struct {
	root *trieNode
}

//trieNode is the structure that contains the node data
//for the trie
type trieNode struct {
	data     []interface{}
	children map[rune]*trieNode
}

//getDataBelow returns all of the data for all descendants of n
func (n *trieNode) getDataBelow() (d []interface{}) {
	if len(n.data) > 0 {
		d = append(d, n.data...)
	}
	tmpKeys := lexicalKeys{}
	for k, _ := range n.children {
		tmpKeys = append(tmpKeys, k)
	}
	sort.Sort(tmpKeys)
	for _, next := range tmpKeys {
		d = append(d, n.children[next].getDataBelow()...)
	}
	return
}

//newNode is the method for creating a new node for the trie
func newNode() (n *trieNode) {
	n = &trieNode{}
	n.children = make(map[rune]*trieNode)
	return
}

//loop is the method that runs the goroutine for the data structure
func (s *SafeTrie) loop() {
	t := &trie{}
	t.root = newNode()
	for op := range s.op {
		op(t)
	}
}

//Destroy stops the running go routine
func (s *SafeTrie) Destroy() {
	close(s.op)
	return
}

//New creates a new trie
func New() (t *SafeTrie) {
	t = &SafeTrie{make(chan func(*trie))}
	go t.loop()
	return
}
