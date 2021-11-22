package main

import (
	"crypto/sha256"
	"strconv"
	"time"
	"math/rand"
	"encoding/hex"
	"strings"
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

func verifySLOW(hash []byte, d int) (bool) {
	str := hex.EncodeToString(hash)
	res := strings.Split(str, "")
	for i := 0; i < d; i++ {
		if res[i] != "0" { return false }
	}
	return true
}

func doubleSHA256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	h2 := sha256.New()
	h2.Write(h.Sum(nil))
	return h2.Sum(nil)
}


func computeHash(block Block, sendCH chan Block, quit chan int, difficulty int, id int) {
	rand.Seed(time.Now().UnixNano())
	block.nonce = rand.Intn(99999999999999999) // this will obv have to be adjusted
	hash := doubleSHA256([]byte(strconv.Itoa(block.nonce) + block.transaction))
	startTime := time.Now()

	for {
		select {
		case <-quit:
			return
		default:
			hash = doubleSHA256([]byte(strconv.Itoa(block.nonce) + block.transaction))

			if verify(hash, difficulty) {
				endTime := time.Now()
				block.duration = endTime.Sub(startTime)
				block.minerID = id

				sendCH <- block
				<-quit
				return
			} else {
				block.nonce++
			}
		}
	}
	// add probability to send a faulty hash
}

func miner(receiveCH chan Block, sendCH chan Block, quit chan int, difficulty int, id int) {

	for {
		computeHash(<-receiveCH, sendCH, quit, difficulty, id)
	}
}
