package main

import (
//	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"time"
	"runtime"
//	"math/rand"
)

var head HashPointer

var NUM_MINERS int
var NUM_BLOCKS int
var NUM_MINED int
var DIFFICULTY int
var SUCCESS_PROBABILITY float64
var PRINT_INFO bool

const(
	OK = 0
	DUP = 1
	FAULTY = 2
)

// return true if the given block is already in the
// chain, ie it's a duplicate left over from the channel
func existsInChain(block Block) (bool) {
	t := block.Tx
	i := head
	for i.Pointer != nil {
		if i.Pointer.Tx == t {
			return true
		}
		i = i.Pointer.Header.HashPrevBlock
	}
	return false
}

// adds a block to the chain
func appendBlock(block Block) {
	if(head.Pointer != nil) {
		block.Header.HashPrevBlock.Pointer = head.Pointer
	}
	head.Pointer = &block
	head.Hash = hashBlock(block)
}

// create and send a block to all miners
func broadcastBlock(unsolvedCH chan Block) {
	block := NewBlock()
	for i := 0; i < NUM_MINERS; i++ {
		unsolvedCH <- block
	}
}

// tell all miners to stop mining by passing
// an arbitrary into over the quit chan
func broadcastSolved(quit chan int) {
	for i := 0; i < NUM_MINERS; i++ {
		quit <- 1
	}
}

// make sure the block is both a correct solution and isn't
// already in the chain
func checkSolution(block Block) int {
	if !verify(hashBlock(block), block.Header.Bits) {
		return FAULTY
	}
	if existsInChain(block) {
		return DUP
	}
	return OK
}

// ./main [# blocks] [# miners] [difficulty] [probability] [show info? t/f]
func main() {

	runtime.GOMAXPROCS(10)

	// input handling
	if len(os.Args) != 6 {
		fmt.Println("Usage: ./main [# blocks] [# miners] [difficulty] [succss probability] [print info t/f")
		return
	} else {
		NUM_BLOCKS, _ = strconv.Atoi(os.Args[1])
		if NUM_BLOCKS <= 0 {
			fmt.Println("# blocks must be >= 1")
			return
		}
		NUM_MINERS, _ = strconv.Atoi(os.Args[2])
		if NUM_MINERS <= 0 {
			fmt.Println("# miners must be >= 1")
			return
		}
		DIFFICULTY, _ = strconv.Atoi(os.Args[3])
		if DIFFICULTY <= 0 {
			fmt.Println("Difficulty must be >= 1")
			return
		}
		SUCCESS_PROBABILITY, _ = strconv.ParseFloat(os.Args[4], 32)
		if SUCCESS_PROBABILITY <= 0 || SUCCESS_PROBABILITY > 1.0 {
			fmt.Println("Probability must be 0 < p <= 1")
			return
		}
		if os.Args[5] == "t" {
			PRINT_INFO = true
		} else {
			PRINT_INFO = false
		}
	}

	unsolvedCH := make(chan Block, NUM_MINERS) // used to send new blocks to miners
	candidateCH := make(chan struct {Block; string; int}, NUM_MINERS) // pair of potentially solved block and string representation of duration to mine
	quitCH := make(chan int) // notify miners to stop mining a solved block
	var totalTime float64 = 0

	for i := 0; i < NUM_MINERS; i++ {
		go miner(unsolvedCH, candidateCH, quitCH, i)
	}

	// send first block
	broadcastBlock(unsolvedCH)

	for NUM_MINED < NUM_BLOCKS {
		// block here until we receive a potentially solved block
		triple := <-candidateCH

		// parse the channel and convert the duration into seconds for averaging + total time
		block := triple.Block
		durStr := triple.string
		dur, _ := time.ParseDuration(durStr)
		sec := dur.Seconds()
		totalTime += sec
		id := triple.int

		switch checkSolution(block) {
			case DUP:
				if PRINT_INFO {
					fmt.Println("[INFO] Duplicate block found, skipping")
				}
				continue
			case FAULTY:
				if PRINT_INFO {
					fmt.Println("[INFO] Faulty block found, skipping")
				}
				continue
			case OK:
				appendBlock(block)
				NUM_MINED++
				broadcastSolved(quitCH)
				broadcastBlock(unsolvedCH)
		}
		fmt.Printf("[Block %d] %x - mined by miner#%d (%fs)\n", NUM_MINED, hashBlock(block),id, sec)
	}
	avgTime := totalTime / (float64(NUM_BLOCKS) + 1.0)
	fmt.Printf("-- Mined %d blocks with difficulty %d with %d miners in %fs (avg of %fs) --\n", NUM_BLOCKS, DIFFICULTY, NUM_MINERS, totalTime, avgTime)
}
