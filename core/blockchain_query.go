package core

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
	"github.com/jinzhu/now"
)

//GetBlockInfo returns the blockinfo of a particular block from the db
func (bc *BlockChain) GetBlockInfo(hash string) (*BlockInfo, error) {
	bc.logger.Debugf("Blockchain: getting block info of hash %s", hash)
	var blockinfo *BlockInfo
	err := bc.getDB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		hashBytes, err := hex.DecodeString(hash)
		if err != nil {
			return err
		}
		blockinfoBytes := b.Get(hashBytes)
		if blockinfoBytes != nil {
			err = json.Unmarshal(blockinfoBytes, &blockinfo)
			if err != nil {
				return err
			}
		} else {
			blockinfo = nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if blockinfo != nil {
		return blockinfo, nil
	}
	return nil, ErrBlockNotFound
}

//GetPrevHash returns the hash of the last block in the bc
func (bc BlockChain) GetPrevHash() (string, error) {
	b, err := bc.GetLatestBlock()
	if err != nil {
		return "", err
	}

	return b.GetHeader().GetHash(), nil
}

//GetBlocksWithinMinute returns all blocks in the db within the last minute
func (bc *BlockChain) GetBlocksWithinMinute() ([]Block, error) {
	var blocks []Block
	var err error
	now := now.New(time.Now())
	bci := bc.iterator()
	for {
		block, err := bci.Next()
		if err != nil {
			return []Block{}, err
		}
		if block.GetHeight() == 0 && block.GetHeader().GetTimestamp() > now.BeginningOfMinute().Unix() {
			blocks = append(blocks, *block)
			break
		} else if block.GetHeader().GetTimestamp() > now.BeginningOfMinute().Unix() {
			blocks = append(blocks, *block)
		} else {
			break
		}
	}
	return blocks, err
}

//GetBlockByHeight return block by height
func (bc *BlockChain) GetBlockByHeight(height int) (*Block, error) {
	bci := bc.iterator()
	for {
		block, err := bci.Next()
		if err != nil {
			return nil, err
		}
		if height != 0 && block.GetHeight() == 0 {
			return nil, ErrBlockNotFound
		} else if int(block.GetHeight()) == height {
			return block, nil
		}
	}
}

//GetLatest15 retuns the latest 15 blocks
func (bc *BlockChain) GetLatest15() ([]Block, error) {
	var blocks []Block
	bci := bc.iterator()
	for {
		if len(blocks) <= 15 {
			block, err := bci.Next()
			if err != nil {
				return []Block{}, err
			}
			if block.GetHeight() == 0 {
				blocks = append(blocks, *block)
				break
			}
			blocks = append(blocks, *block)
		}
		break
	}
	return blocks, nil
}

//GetLatestHeight returns the height of the latest block to the blockchain
func (bc *BlockChain) GetLatestHeight() (int, error) {
	var lastBlock *BlockInfo
	err := bc.getDB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		lastBlockBytes := b.Get(bc.getTip())
		err := json.Unmarshal(lastBlockBytes, &lastBlock)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return int(lastBlock.GetHeight()), nil
}

//GetLatestBlock returns the tip as a block
func (bc *BlockChain) GetLatestBlock() (*Block, error) {
	var lastBlock *BlockInfo
	err := bc.getDB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		lastBlockBytes := b.Get(bc.getTip())
		err := json.Unmarshal(lastBlockBytes, &lastBlock)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return lastBlock.GetBlock(), nil
}

//GetNextHeight returns the next height in the blockchain
func (bc BlockChain) GetNextHeight() (uint64, error) {
	b, err := bc.GetLatestBlock()
	if err != nil {
		return 0, err
	}
	return b.GetHeight() + 1, err
}

//AddBlock adds block to the blockchain
func (bc *BlockChain) AddBlock(block *Block) error {
	verify, err := block.VerifyBlock()
	if err != nil {
		return err
	}
	if verify == false {
		return ErrUnverifiedBlock
	}
	err = bc.getDB().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		hashBytes, err := hex.DecodeString(block.GetHeader().GetHash())
		if err != nil {
			return err
		}
		inDb := b.Get(hashBytes)
		if inDb != nil {
			bc.logger.Warn("Core: block exists in blockchain")
			return nil
		}

		blockinfo := BlockInfo{
			Header:    block.GetHeader(),
			Height:    block.GetHeight(),
			TotalJobs: uint(len(block.GetNodes())),
			FileName:  block.fileStats().Name(),
			FileSize:  block.fileStats().Size(),
		}
		biBytes, err := json.Marshal(blockinfo)
		if err != nil {
			return err
		}

		if err = b.Put(hashBytes, biBytes); err != nil {
			return err
		}

		if bc.getTip() != nil {
			latest, err := bc.GetBlockInfo(hex.EncodeToString(bc.getTip()))
			if err != nil {
				return err
			}
			if block.GetHeight() > latest.GetHeight() {
				if err := b.Put([]byte("l"), hashBytes); err != nil {
					return err
				}
				bc.setTip(hashBytes)
			}
		} else {
			if err = b.Put([]byte("l"), hashBytes); err != nil {
				return err
			}
			bc.setTip(hashBytes)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//FindJob returns the job from the blockchain
func (bc *BlockChain) FindJob(id string) (*job.Job, error) {
	var tree merkletree.MerkleTree
	bci := bc.iterator()
	for {
		block, err := bci.Next()
		if err != nil {
			return nil, err
		}
		if block.GetHeight() == 0 {
			return nil, ErrJobNotFound
		}
		tree.SetLeafNodes(block.GetNodes())
		found, err := tree.SearchJob(id)
		if found == nil && err != nil {
			continue
		}
		for i, exec := range found.GetExecs() {
			if exec.GetTimestamp() > now.BeginningOfDay().Unix() {
				found.Execs = append(found.Execs[:i], found.Execs[i+1:]...) //! removes execs older than a day
			}
		}
		return found, nil
	}
}

//FindExec finds exec in the bc
func (bc *BlockChain) FindExec(id string, hash string) (*job.Exec, error) {
	var tree merkletree.MerkleTree
	bci := bc.iterator()
	for {
		block, err := bci.Next()
		if err != nil {
			return nil, err
		}
		if block.GetHeight() == 0 {
			return nil, job.ErrExecNotFound
		}
		tree.SetLeafNodes(block.GetNodes())
		found, err := tree.SearchJob(id)
		if found == nil && err != nil {
			continue
		}
		exec, err := found.GetExec(hash)
		if err != nil {
			continue
		}
		return exec, nil
	}
}

//GetJobExecs returns all execs of a job
func (bc *BlockChain) GetJobExecs(id string) ([]job.Exec, error) {
	execs := []job.Exec{}
	var tree merkletree.MerkleTree
	bci := bc.iterator()
	for {
		block, err := bci.Next()
		if err != nil {
			return []job.Exec{}, err
		}
		if block.GetHeight() == 0 {
			return job.UniqExec(execs), nil
		}
		tree.SetLeafNodes(block.GetNodes())
		found, err := tree.SearchJob(id)
		if found == nil && err != nil {
			continue
		}
		execs = append(execs, found.GetExecs()...)
	}
}

//FindMerkleNode returns the merklenode from the blockchain
func (bc *BlockChain) FindMerkleNode(h string) (*merkletree.MerkleNode, error) {
	var tree merkletree.MerkleTree
	bci := bc.iterator()
	for {
		block, err := bci.Next()
		if err != nil {
			return nil, err
		}
		if block.GetHeight() == 0 {
			return nil, ErrJobNotFound
		}
		tree.SetLeafNodes(block.GetNodes())
		found, err := tree.SearchNode(h)
		if err != nil {
			return nil, err
		}
		return found, nil
	}
}
