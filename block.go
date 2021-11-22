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
	Hash []byte // this may need to be changed
	Pointer *Block
}

// Returns a new block with a random transaction as its data,
// and nothing set as its header
// TODO: because this sends the transaction as well, anyone can edit the transaction before sending back, perhaps the logger verifies based on this already set transaction?
func NewBlock() Block { // basically SetNull() in block.cpp
	rand.Seed(time.Now().UnixNano())
	sender := strconv.Itoa(rand.Intn(999999999))
	recipient := strconv.Itoa(rand.Intn(999999999))
	b := Block{0xD9B4BEF9, Header{HashPointer{}, time.Now(), DIFFICULTY, 0}, sender + " -> " + recipient}
	return b
}

func doubleSHA256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	h2 := sha256.New()
	h2.Write(h.Sum(nil))
	return h2.Sum(nil)
}

func BlockToBytes(block Block) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(block)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}


func UInt64ToBytes(i uint64) []byte {
	slice := make([]byte, 8)
	binary.LittleEndian.PutUint64(slice, i)
	return slice
}

func hashBlock(block Block) []byte {
	x := BlockToBytes(block)
	n := UInt64ToBytes(block.Header.Nonce)
	n = append(n, x...)
	return doubleSHA256(n)
}
