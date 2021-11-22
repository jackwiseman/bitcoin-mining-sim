package main

import (
	"time"
	"math/rand"
)

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


func solvePuzzle(block Block, sendCH chan struct{Block; string}, quit chan int, id int) {
	rand.Seed(time.Now().UnixNano())
	block.Header.Nonce = rand.Uint64()

	hashInput := UInt64ToBytes(block.Header.Nonce)
	blockAsBytes := BlockToBytes(block)
	hashInput = append(hashInput, blockAsBytes...)

	hash := hashBlock(block)
	startTime := time.Now()


	for {
		select {
		case <-quit:
			return
		default:
			hash = hashBlock(block)

			if verify(hash, block.Header.Bits) {
				endTime := time.Now()
				sendCH <- struct{Block; string}{block, endTime.Sub(startTime).String()}
				<-quit
				return
			} else {
				block.Header.Nonce++
			}
		}
	}
	// add probability to send a faulty hash
}

func miner(receiveCH chan Block, sendCH chan struct{Block; string}, quit chan int, id int) {

	for {
		solvePuzzle(<-receiveCH, sendCH, quit, id)
	}
}
