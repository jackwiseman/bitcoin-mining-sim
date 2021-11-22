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
var NUM_MINERS  int = 59
var NUM_BLOCKS int = 20
var NUM_MINED int = 0
var DIFFICULTY int = 1

func appendBlock(block Block) {
	if(head.pointer != nil) {
		block.hashPrevBlock.pointer = head.pointer
	}
	head.pointer = &block
	head.hash = sha256.Sum256([]byte(block.transaction))
}

func broadcastBlock(unsolvedCH chan Block) {
	block := NewBlock()
	fmt.Println("Sending out block " + strconv.Itoa(NUM_MINED + 1) + " with transaction " + block.transaction)
	for i := 0; i < NUM_MINERS; i++ {
		unsolvedCH <- block
	}
}

func broadcastSolved(quit chan int) {
	for i := 0; i < NUM_MINERS; i++ {
		quit <- 1
	}
}

func printChain() {
	i := head
	for i.pointer != nil {
		fmt.Printf("miner#%d -- %d : %s (Time elapsed: %s)\n", i.pointer.minerID, i.pointer.nonce, i.pointer.transaction, i.pointer.duration)
		i = i.pointer.hashPrevBlock
	}
}

func main() {
	unsolvedCH := make(chan Block, NUM_MINERS)
	candidateCH := make(chan Block, NUM_MINERS)
	quitCH := make(chan int, NUM_MINERS)

	for i := 0; i < NUM_MINERS; i++ {
		go miner(unsolvedCH, candidateCH, quitCH, DIFFICULTY, i)
	}
	// send first block
	broadcastBlock(unsolvedCH)

	for NUM_MINED < NUM_BLOCKS {
		block := <-candidateCH // both should prob be named better
		// if verify(block) { }
		fmt.Println("Recieved: " + block.transaction + " from " + strconv.Itoa(block.minerID))

		appendBlock(block)
		NUM_MINED++
		fmt.Println(NUM_MINED < NUM_BLOCKS)
		broadcastSolved(quitCH)
		broadcastBlock(unsolvedCH)
		fmt.Println(len(quitCH))
	}
	printChain()
}
