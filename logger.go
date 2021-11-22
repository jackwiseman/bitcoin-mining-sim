package main

import (
//	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
//	"time"
)

/* Functions to implement

- verify: check incoming puzzle solves

*/

var head HashPointer

var NUM_MINERS int = 100
var NUM_BLOCKS int = 10
var NUM_MINED int = 0
var DIFFICULTY int = 7
//var DUPLICATES int = 0

const(
	OK = 0
	DUP = 1
	FAULTY = 2
)

func existsInChain(block Block) (bool) {
	t := block.Tx
	i := head
	for i.Pointer != nil {
		if i.Pointer.Tx == t {
			//DUPLICATES++
			return true
		}
		i = i.Pointer.Header.HashPrevBlock
	}
	return false
}

func appendBlock(block Block) {
	if(head.Pointer != nil) {
		block.Header.HashPrevBlock.Pointer = head.Pointer
	}
	head.Pointer = &block
	head.Hash = hashBlock(block)
}

func broadcastBlock(unsolvedCH chan Block) {
	block := NewBlock()
//	fmt.Println("Sending out block " + strconv.Itoa(NUM_MINED + 1) + " with Tx " + block.Tx)
	for i := 0; i < NUM_MINERS; i++ {
		unsolvedCH <- block
	}
}

func broadcastSolved(quit chan int) {
	for i := 0; i < NUM_MINERS; i++ {
		quit <- 1
	}
}

func checkSolution(block Block) int {
	if !verify(hashBlock(block), block.Header.Bits) {
		return FAULTY
	}
	if existsInChain(block) {
		return DUP
	}
	return OK
}

func printChain() {
	i := head
	for i.Pointer != nil {
//	fmt.Printf("miner#%d -- (%d%s) (Time elapsed: %s)\n", i.Pointer.minerID, i.Pointer.Nonce, i.Pointer.Tx, i.Pointer.Header.Timestamp)
	fmt.Printf("%d%s\n", i.Pointer.Header.Nonce, i.Pointer.Tx)
	//fmt.Printf("%x (%s)\n", doubleSHA256([]byte(strconv.Itoa(i.Pointer.Header.Nonce) + i.Pointer.Tx)), i.Pointer.Header.Timestamp)
		i = i.Pointer.Header.HashPrevBlock
//		totalTime = totalTime.Add(i.Pointer.Header.Timestamp)
	}
//	avg := time.totalTime.Seconds() / NUM_BLOCKS

//	fmt.Printf("%d Duplicates\n", DUPLICATES)
//	fmt.Println("Average time" + avg.String())
}

// ./main [# blocks] [# miners] [difficulty]
func main() {
	// input handling
	if len(os.Args) != 4 {
		fmt.Println("Usage: ./main [# blocks] [# miners] [difficulty]")
		return
	} else {
		NUM_BLOCKS, _ = strconv.Atoi(os.Args[1])
		NUM_MINERS, _ = strconv.Atoi(os.Args[2])
		DIFFICULTY, _ = strconv.Atoi(os.Args[3])
	}

	unsolvedCH := make(chan Block, NUM_MINERS)
	candidateCH := make(chan struct {Block; string}, NUM_MINERS)
	quitCH := make(chan int)

	for i := 0; i < NUM_MINERS; i++ {
		go miner(unsolvedCH, candidateCH, quitCH, i)
	}

	// send first block
	broadcastBlock(unsolvedCH)

	for NUM_MINED < NUM_BLOCKS {
		pair := <-candidateCH // both should prob be named better
		block := pair.Block
		dur := pair.string

		switch checkSolution(block) {
			case DUP:
				fmt.Println("[INFO] Duplicate block found, skipping")
				continue
			case FAULTY:
				fmt.Println("[INFO] Faulty block found, skipping")
				continue
			case OK:
				appendBlock(block)
				NUM_MINED++
				broadcastSolved(quitCH)
				broadcastBlock(unsolvedCH) // this is the same thing as saying its been solved
		}
		fmt.Printf("[Block %d] %x (%s)\n", NUM_MINED, hashBlock(block), dur)
	}

//	printChain()
}
