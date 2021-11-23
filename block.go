package main

import (
	"math/rand"
	"strconv"
	"time"
	"crypto/sha256"
	"encoding/binary"
	"bytes"
	"encoding/gob"
)

type Block struct {
	Magic int
	Header Header
	Tx string
}

type Header struct {
	HashPrevBlock HashPointer
	Timestamp time.Time // used instead for time to mine
	Bits int // difficulty
	Nonce uint64
}

type HashPointer struct {
	Hash []byte
	Pointer *Block
}

// Returns a new block with a random transaction as its data,
// and nothing set as its header
func NewBlock() Block {
	rand.Seed(time.Now().UnixNano())
	// we don't really care about sender and recipient this is just to give it some data
	// formatted as a string "somenumber -> somenumber"
	sender := strconv.Itoa(rand.Intn(999999999))
	recipient := strconv.Itoa(rand.Intn(999999999))
	b := Block{0xD9B4BEF9, Header{HashPointer{}, time.Now(), DIFFICULTY, 0}, sender + " -> " + recipient} // 0xD9B4BEF9 is the magic number in bitcoin
	return b
}

// hash []byte using sha256 twice, as this is how bitcoin does it
func doubleSHA256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	h2 := sha256.New()
	h2.Write(h.Sum(nil))
	return h2.Sum(nil)
}

// write a block into a []byte for hashing
func BlockToBytes(block Block) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(block)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// write a UInt64 into []byte for hashing
// primarily used for nonce value
func UInt64ToBytes(i uint64) []byte {
	slice := make([]byte, 8)
	binary.LittleEndian.PutUint64(slice, i)
	return slice
}

// helper function to hash a block
// utilizes the above three functions together to do so
func hashBlock(block Block) []byte {
	x := BlockToBytes(block)
	n := UInt64ToBytes(block.Header.Nonce)
	n = append(n, x...)
	return doubleSHA256(n)
}
