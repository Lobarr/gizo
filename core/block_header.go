package core

import "math/big"

//BlockHeader holds the header of the block
type BlockHeader struct {
	Timestamp     int64  `storm:"index"`
	PrevBlockHash string ``
	MerkleRoot    string `storm:"index"` // hash of merkle root
	Nonce         uint64
	Difficulty    *big.Int
	Hash          string `storm:"index"`
}

//GetTimestamp returns timestamp
func (blockHeader BlockHeader) GetTimestamp() int64 {
	return blockHeader.Timestamp
}

//sets the timestamp
func (blockHeader *BlockHeader) setTimestamp(t int64) {
	blockHeader.Timestamp = t
}

//GetPrevBlockHash returns previous block hash
func (blockHeader BlockHeader) GetPrevBlockHash() string {
	return blockHeader.PrevBlockHash
}

//sets prevblockhash
func (blockHeader *BlockHeader) setPrevBlockHash(h string) {
	blockHeader.PrevBlockHash = h
}

//GetMerkleRoot returns merkleroot
func (blockHeader BlockHeader) GetMerkleRoot() string {
	return blockHeader.MerkleRoot
}

//sets merkleroot
func (blockHeader *BlockHeader) setMerkleRoot(merkleRoot string) {
	blockHeader.MerkleRoot = merkleRoot
}

//GetNonce returns the nonce
func (blockHeader BlockHeader) GetNonce() uint64 {
	return blockHeader.Nonce
}

//sets the nonce
func (blockHeader *BlockHeader) setNonce(nonce uint64) {
	blockHeader.Nonce = nonce
}

//GetDifficulty returns difficulty
func (blockHeader BlockHeader) GetDifficulty() *big.Int {
	return blockHeader.Difficulty
}

//sets the difficulty
func (blockHeader *BlockHeader) setDifficulty(d big.Int) {
	blockHeader.Difficulty = &d
}

//GetHash returns hash
func (blockHeader BlockHeader) GetHash() string {
	return blockHeader.Hash
}

//sets hash
func (blockHeader *BlockHeader) setHash(h string) {
	blockHeader.Hash = h
}
