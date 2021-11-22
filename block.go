package main

import (
//	"crypto/sha256"
	"math/rand"
	"strconv"
	"time"
)

/* Functions to implement

*/

type Block struct {
	hashPrevBlock HashPointer// this may need to be a pointer?
	transaction string
	nonce int// probably has to account for larger nums
	duration time.Duration
	minerID int
}

type HashPointer struct {
	hash [32]byte // this may need to be changed
	pointer *Block
}

// Returns a new block with a random transaction as its data,
// and nothing set as its header
// TODO: because this sends the transaction as well, anyone can edit the transaction before sending back, perhaps the logger verifies based on this already set transaction?
func NewBlock() Block { // basically SetNull() in block.cpp
	rand.Seed(time.Now().UnixNano())
	sender := strconv.Itoa(rand.Intn(999999999))
	recipient := strconv.Itoa(rand.Intn(999999999))
	b := Block{HashPointer{}, sender + " -> " + recipient, 0, 0, 0}
	return b
}
