package core

import "math/big"

//BlockHeader holds the header of the block
type BlockHeader struct {
	Timestamp     int64
	PrevBlockHash string
	MerkleRoot    string // hash of merkle root
	Nonce         uint64
	Difficulty    *big.Int
	Hash          string
}

//GetTimestamp returns timestamp
func (bh BlockHeader) GetTimestamp() int64 {
	return bh.Timestamp
}

//sets the timestamp
func (bh *BlockHeader) setTimestamp(t int64) {
	bh.Timestamp = t
}

//GetPrevBlockHash returns previous block hash
func (bh BlockHeader) GetPrevBlockHash() string {
	return bh.PrevBlockHash
}

//sets prevblockhash
func (bh *BlockHeader) setPrevBlockHash(h string) {
	bh.PrevBlockHash = h
}

//GetMerkleRoot returns merkleroot
func (bh BlockHeader) GetMerkleRoot() string {
	return bh.MerkleRoot
}

//sets merkleroot
func (bh *BlockHeader) setMerkleRoot(mr string) {
	bh.MerkleRoot = mr
}

//GetNonce returns the nonce
func (bh BlockHeader) GetNonce() uint64 {
	return bh.Nonce
}

//sets the nonce
func (bh *BlockHeader) setNonce(n uint64) {
	bh.Nonce = n
}

//GetDifficulty returns difficulty
func (bh BlockHeader) GetDifficulty() *big.Int {
	return bh.Difficulty
}

//sets the difficulty
func (bh *BlockHeader) setDifficulty(d big.Int) {
	bh.Difficulty = &d
}

//GetHash returns hash
func (bh BlockHeader) GetHash() string {
	return bh.Hash
}

//sets hash
func (bh *BlockHeader) setHash(h string) {
	bh.Hash = h
}
