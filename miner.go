package main

import (
	"crypto/sha256"
//	"math/big"
	"fmt"
	"strconv"
//	"encoding/hex"
//	"encoding/binary"
	"time"
)

// minerID int // this may honestly be completely unnecessary

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

func computeHash(str string, difficulty int) ([32]byte, int) {
	var nonce int = 0
	hashResult := sha256.Sum256([]byte(strconv.Itoa(nonce) + str))

	for true {
		hashResult = sha256.Sum256([]byte(strconv.Itoa(nonce) + str))
		if verify(hashResult[:], difficulty) == true {
			break
		}
		nonce++
	}

	return hashResult, nonce
}

func miner(ch chan Block) // this is the "main" of any miner

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


