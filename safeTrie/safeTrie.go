//Package safeTrie implements a thread safe trie in go
package safeTrie

import (
	"fmt"
)

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
