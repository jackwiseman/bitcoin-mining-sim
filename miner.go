package main

import (
	"crypto/sha256"
	"strconv"
	"time"
	"math/rand"
)

var MINER_ID int

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

func miner(receiveCH chan Block, sendCH chan Block, difficulty int, id int) {
	block := receiveCH
	rand.Seed(time.Now().UnixNano())
	block.nonce = rand.Intn(99999999999999) // this will obv have to be adjusted
	hashResult := sha256.Sum256([]byte(strconv.Itoa(block.nonce) + block.transaction))
	startTime := time.Now()

	for {
		select {
		case block = <-receiveCH: // if new block received
			block.nonce = rand.Intn(99999999999999) // this will obv have to be adjusted
			startTime = time.Now()
		default:
			hashResult = sha256.Sum256([]byte(strconv.Itoa(block.nonce) + block.transaction))

			if verify(hashResult[:], difficulty) {
				endTime := time.Now()
				block.duration = endTime.Sub(startTime).String()
				block.minerID = id

				sendCH <- block
			}
			block.nonce++
		}
	}
	// add probability to send a faulty hash
}

/*func computeHash(block Block, sendCH chan Block, quit chan int, difficulty int) {
	rand.Seed(time.Now().UnixNano())
	block.nonce = rand.Intn(99999999999999) // this will obv have to be adjusted
	hashResult := sha256.Sum256([]byte(strconv.Itoa(block.nonce) + block.transaction))
	startTime := time.Now()

	for {
		select {
		case <-quit:
			//fmt.Println(quit_value)
			return
		default:
			hashResult = sha256.Sum256([]byte(strconv.Itoa(block.nonce) + block.transaction))

			if verify(hashResult[:], difficulty) {
				endTime := time.Now()
				block.duration = endTime.Sub(startTime).String()
				block.minerID = MINER_ID

				sendCH <- block
				return
			} else {
				block.nonce++
			}
		}
	}
	// add probability to send a faulty hash
}

func miner(receiveCH chan Block, sendCH chan Block, quit chan int, difficulty int, id int) {
	MINER_ID = id

	for {
		select {
			case block := <-receiveCH: 
				go computeHash(block, sendCH, quit, difficulty)
			default:
				continue
			}
	}
}
*/
