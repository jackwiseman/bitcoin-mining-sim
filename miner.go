package main

import (
	"time"
	"math/rand"
)

// verify a hash by comparing the amount of leading 0s in its
// []byte representation where each i is two hex values
func verify(hash []byte, d int) (bool) {
	if d % 2 == 0 {
		for i :=0; i < (d / 2); i++ {
			if hash[i] != 0 {
				return false
			}
		}
	} else {
		for i :=0; i < ((d + 1) / 2); i++ {
			if ((i + 1) >= ((d + 1) / 2)) {
				if hash[i] > 15 {
					return false
				}
				return true
			} else {
				if hash[i] != 0 {
					return false
				}
			}
		}
	}
	return true
}

// pick a random Uint64 as the nonce, and hash this along with the block
// if it doesn't have the correct number of leading 0s, add 1 to nonce and try again
// at the same time, we need to check if we should quit this puzzle, which is done with
// the select statement, finally send the block, the time it took to solve, and the id of the
// miner through the chan and exit (unless its a faulty block sent)
func solvePuzzle(block Block, sendCH chan struct{Block; string; int}, quit chan int, id int) {
	rand.Seed(time.Now().UnixNano())
	block.Header.Nonce = rand.Uint64()

	hash := hashBlock(block)
	startTime := time.Now()


	for {
		select {
		case <-quit:
			return
		default:
			hash = hashBlock(block)

			if verify(hash, block.Header.Bits) {
				// if the miner sends a bad block, instead of a good one, it will just start over and try again
				if(rand.Float64() >= SUCCESS_PROBABILITY) { // 95% chance to send a "good" block
					block.Header.Nonce = rand.Uint64()
					endTime := time.Now()
					sendCH <- struct{Block; string; int}{block, endTime.Sub(startTime).String(), id}
				} else {
					endTime := time.Now()
					sendCH <- struct{Block; string; int}{block, endTime.Sub(startTime).String(), id}
					<-quit
					return
				}
			} else {
				block.Header.Nonce++
			}
		}
	}
}

// if this miner has solved the puzzle or quits, block until it receives a new block
func miner(receiveCH chan Block, sendCH chan struct{Block; string; int}, quit chan int, id int) {
	for {
		solvePuzzle(<-receiveCH, sendCH, quit, id)
	}
}
