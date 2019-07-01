package core

import (
	"encoding/hex"
	"encoding/json"

	"github.com/boltdb/bolt"
)

//BlockChainIterator - a way to loop through the blockchain (from newest block to oldest block)
type BlockChainIterator struct {
	current []byte
	db      IBolt
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
func (i *BlockChainIterator) Next() (IBlock, error) {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		blockinfoBytes := b.Get(i.GetCurrent())
		var blockinfo *BlockInfo
		err := json.Unmarshal(blockinfoBytes, &blockinfo)
		if err != nil {
			return err
		}
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

// NextBlockinfo returns the next blockinfo in the blockchain - more lightweight
func (i *BlockChainIterator) NextBlockinfo() (IBlockInfo, error) {
	var block *BlockInfo
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		blockinfoBytes := b.Get(i.GetCurrent())
		var blockinfo *BlockInfo
		err := json.Unmarshal(blockinfoBytes, &blockinfo)
		if err != nil {
			return err
		}
		block = blockinfo
		return nil
	})
	if err != nil {
		return nil, err
	}
	current, err := hex.DecodeString(block.GetHeader().GetPrevBlockHash())
	if err != nil {
		return nil, err
	}
	i.setCurrent(current)
	return block, nil
}

//NewBlockchainIterator create new blockchain iterator
func NewBlockchainIterator(blockTip []byte, boltDB IBolt) IBlockIterator {
	return &BlockChainIterator{
		current: blockTip,
		db: boltDB
	}
}
