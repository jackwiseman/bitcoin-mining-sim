package main

import (
	"crypto/sha256"
)

/* Functions to implement

- verify: check incoming puzzle solves
- append: add block to log
- broadcast: tell all miners about new block

Other

this holds the channels?
*/

head HashPointer
var NUM_MINERS  int = 3

func appendBlock(block Block) {
	if(head.pointer != nil) {
		block.hashPrevBlock.pointer = head
	}
	head.pointer = &block
	head.hash = sha256.Sum256([]byte(block.transaction))
}

func broadcastBlock(unsolvedCH chan Block, unsolved Block) {
	for i = 0; i < NUM_MINERS; i++ {
		candidateCH <- unsolved
	}
}

func main() {
	unsolvedCH := make(chan Block)
	candidateCH := make(chan Block)

	for i = 0; i < NUM_MINERS; i++ {
		go miner(unsolvedCH, candidateCH)
	}

	broadcastBlock(unsolvedCH, newBlock())

	for {
		block := <-candidateCH // both should prob be named better
		// if verify(block) { }
		appendBlock(block)
		broadcastBlock(unsolvedCH, newBlock()) // this is sending a pointer to a block
	}
}
