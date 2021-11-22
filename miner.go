package main

import (
	"crypto/sha256"
	"strconv"
	"time"
	"math/rand"
)


// if  d = 4
// (d / 2 - 1)
// [12 34 56 78]
// i[1]

func verify(hash []byte, d int) (bool) {
	for i :=0; i < ((d + 1) / 2); i++ {
		if (i % 2 == 0) && hash[i] != 0 {
				return false
		} else {
			if(i + 1 >= ((d + 1) / 2) && hash[i] > 15) {
				return false
			} else {
				if(hash[i] != 0) {
					return false
				}
			}
		}
	}
	return true
}

// func verify(hash []byte, d int) (bool) {
// 	i := 0
// 	for i < d {
// 		if (d % 2 == 1 && i + 1 >= d) {
// 			if int(hash[i]) > 15 {
// 				return false
// 			} else {
// 				return true
// 			}
// 		}
// 		if int(hash[i]) != 0 {
// 			return false
// 		}
// 		i++
// 	}
// 	return true
// }


func computeHash(block Block, sendCH chan Block, quit chan int, difficulty int, id int) {
	rand.Seed(time.Now().UnixNano())
	block.nonce = rand.Intn(99999999999999999) // this will obv have to be adjusted
	hashResult := sha256.Sum256([]byte(strconv.Itoa(block.nonce) + block.transaction))
	startTime := time.Now()

	for {
		select {
		case <-quit:
			return
		default:
			hashResult = sha256.Sum256([]byte(strconv.Itoa(block.nonce) + block.transaction))

			if verify(hashResult[:], difficulty) {
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
/*		select {
			case block := <-receiveCH: 
				computeHash(block, sendCH, quit, difficulty)
			default:
				continue
			}*/
	}
}
