package main

import (
//	"crypto/sha256"
	"math/rand"
	"strconv"
	"time"
	"fmt"
)

/* Functions to implement

*/

type Block struct {
	hashPrevBlock HashPointer// this may need to be a pointer?
	transaction string
	nonce int// probably has to account for larger nums
}

type HashPointer struct {
	hash [32]byte // this may need to be changed
	pointer *Block
}

// Returns a new block with a random transaction as its data,
// and nothing set as its header
func NewBlock() *Block { // basically SetNull() in block.cpp
	rand.Seed(time.Now().UnixNano())
	sender := strconv.Itoa(rand.Intn(999999999))
	recipient := strconv.Itoa(rand.Intn(999999999))

	block := Block{HashPointer{}, sender + " -> " + recipient, 0}
	fmt.Println("Generated new block: " + block.transaction)

	return &block
}
