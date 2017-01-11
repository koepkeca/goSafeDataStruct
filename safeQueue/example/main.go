package main

import (
	"fmt"

	"github.com/koepkeca/goSafeDataStruct/safeQueue"
)

const (
	//The GO Race detector can manage a max of 8192 concurrent routines.
	//Read more about this here: https://golang.org/doc/articles/race_detector.html
	nbrOfRoutines = 7654 //The GO Race detector can manage a max of 8192 concurrent routines.
)

func main() {
	S := safeQueue.New()
	for i := 0; i < nbrOfRoutines; i++ {
		go func(j int) {
			S.Enqueue(j)
		}(i)
	}
	fmt.Printf("%d elements", S.Len())
	next := S.Dequeue()
	for next != nil {
		fmt.Printf("%d\n", next.(int))
		next = S.Dequeue()
	}
}
