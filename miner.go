package main

import (
	"crypto/sha256"
//	"math/big"
	"strconv"
//	"encoding/hex"
//	"encoding/binary"
)

// minerID int // this may honestly be completely unnecessary
var DIFFICULTY int = 3
var waitingForNewBlock = false

func verify(hash []byte, d int) (bool) {
	i := 0
	for i < d {
		if int(hash[i]) != 0 {
			return false
		}
		i++
	}
	return true
}

func computeHash(receiveCH chan Block, sendCH chan Block) {
	block := <-receiveCH
	hashResult := sha256.Sum256([]byte(strconv.Itoa(block.nonce) + block.transaction))

	for {
		if waitingForNewBlock || len(receiveCH) != 0 {
			block = <-receiveCH
			waitingForNewBlock = false
		}
		hashResult = sha256.Sum256([]byte(strconv.Itoa(block.nonce) + block.transaction))
		if verify(hashResult[:], DIFFICULTY) && len(receiveCH) == 0 {
			sendCH <- block
			waitingForNewBlock = true
		} else {
			block.nonce++
		}
	}
	// add probability to send a faulty hash
	// adds a whole bunch of issues with this logic like double sending a verified block
}

//func miner(receiveCH chan Block, sendCH chan Block) { // this is the "main" of any miner
//	go computeHash(receiveCH, sendCH)
//}

/*func main() {
	i := 1
	for i < 5 {
		startTime := time.Now()
		hash, nonce := computeHash("Jack Wiseman", i)
		endTime := time.Now()

		fmt.Println(nonce)
		fmt.Println(hash[:])
		fmt.Println(endTime.Sub(startTime))
		fmt.Printf("Difficulty: %d\n\n", i)

		i++
	}
}*/


