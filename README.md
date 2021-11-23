# Golang Bitcoin Mining Simulation

### Usage
1. run `go build` in the current directory
2. run `./main [# blocks] [# miners] [difficulty] [probability] [show info? t/f]`
- *\# blocks* : blocks to be mined
- *\# miners* : number of goroutines to be solving a given puzzle at one time
- *difficulty* : bits field of block header (leading number of 0s for a hash)
- *probability* : floating point number that a miner will send the correct solution
- *show info?* : either "t" or "f" denoting whether to display messages about duplicates or faulty blocks
### Design / Flow
- logger.go
	- works as the main file of the program, which spawns goroutines of the miners, and keeps the tamper-resistant log by storing one value `head` which points to the most recent, verified block
	- contains methods to check if a block is in the chain, add a block to the chain, send blocks to the miners, and verify the solutions it receives
- miner.go
	- does the puzzle solving by arbitrarily picking a random nonce value which is a UInt64, adding this together with the block interpreted as a `[]byte`, and double hashes it using sha256. If the result has `bits` leading zeros, it will send a block back to the logger with probability `probability` of it being the correct solution. At the same time as mining, the miner also listens for a new block, terminating the current `solvePuzzle()` method and working on the new one
- block.go
	- holds the structures of `Block`, `Header`, and `HashPointer`, all implemented similarly to how bitcoin does it
	- `Block` contains fields containing the bitcoin magic number of `0xD9B4BEF9`, the `Header` object, and a transaction which is stored as a string of two random ids and an arrow, ie `1223345 -> 4929560` which is completely arbitrary, but used to model the bitcoin block
	- `Header` contains the `HashPointer`, along with a `Timestamp` of its creation, `Bits` which holds the difficulty, and `Nonce` which starts at 0, but is randomly assigned by miners to give each of them a different base nonce from which to start incrementing
	- `HashPointer` holds the values that it infers, both a `Hash` of the previous block as a `[]byte` and a `Pointer` to the previous block
	- also contains methods
		- `NewBlock()` which returns a new block for logger to send out
		- `doubleSHA256()` which computes the hash for a given `[]byte` twice in the same way bitcoin does, and some helper functions to convert the `Block` structure into a `[]byte` and the `Nonce` into a `[]byte`
