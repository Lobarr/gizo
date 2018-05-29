package core

import (
	"encoding/hex"

	"github.com/boltdb/bolt"
	"github.com/kpango/glg"
)

//BlockChainIterator - a way to loop through the blockchain (from newest block to oldest block)
type BlockChainIterator struct {
	current []byte
	db      *bolt.DB
}

//sets current block
func (i *BlockChainIterator) setCurrent(c []byte) {
	i.current = c
}

//GetCurrent returns current block
func (i BlockChainIterator) GetCurrent() []byte {
	return i.current
}

// Next returns the next block in the blockchain
func (i *BlockChainIterator) Next() (*Block, error) {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		blockinfoBytes := b.Get(i.GetCurrent())
		blockinfo := DeserializeBlockInfo(blockinfoBytes)
		block = blockinfo.GetBlock()
		return nil
	})
	if err != nil {
		return nil, err
	}
	current, err := hex.DecodeString(block.GetHeader().GetPrevBlockHash())
	i.setCurrent(current)
	return block, nil
}

// Next returns the next blockinfo in the blockchain - more lightweight
func (i *BlockChainIterator) NextBlockinfo() (*BlockInfo, error) {
	var block *BlockInfo
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		blockinfoBytes := b.Get(i.GetCurrent())
		blockinfo := DeserializeBlockInfo(blockinfoBytes)
		block = blockinfo
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
	current, err := hex.DecodeString(block.GetHeader().GetPrevBlockHash())
	if err != nil {
		return nil, err
	}
	i.setCurrent(current)
	return block, nil
}
