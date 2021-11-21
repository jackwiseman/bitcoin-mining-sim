package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

/* Functions to implement

- verify: check incoming puzzle solves
- append: add block to log
- broadcast: tell all miners about new block

Other

this holds the channels?
*/

var head HashPointer
var NUM_MINERS  int = 3
var NUM_BLOCKS int = 10
var NUM_MINED int = 0

func appendBlock(block Block) {
	if(head.pointer != nil) {
		block.hashPrevBlock.pointer = head.pointer
	}
	head.pointer = &block
	head.hash = sha256.Sum256([]byte(block.transaction))
}

func broadcastBlock(unsolvedCH chan Block, unsolved Block) {
	fmt.Println("Sending out block " + strconv.Itoa(NUM_MINED + 1) + "...")
	for i := 0; i < NUM_MINERS; i++ {
		unsolvedCH <- unsolved
	}
}

func main() {
	unsolvedCH := make(chan Block)
	candidateCH := make(chan Block)

	for i := 0; i < NUM_MINERS; i++ {
		go computeHash(unsolvedCH, candidateCH)
	}
	
	b := NewBlock()
	broadcastBlock(unsolvedCH, b)

	for NUM_MINED < NUM_BLOCKS {
		block := <-candidateCH // both should prob be named better
		// if verify(block) { }
		fmt.Println("Got a potential block!")
		fmt.Println(strconv.Itoa(len(candidateCH)))

		appendBlock(block)
		NUM_MINED++
		b = NewBlock()
		fmt.Println(NUM_MINED < NUM_BLOCKS)
		broadcastBlock(unsolvedCH, b) // this is sending a pointer to a block
	}
}
