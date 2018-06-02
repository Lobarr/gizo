package core

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	funk "github.com/thoas/go-funk"

	"github.com/gizo-network/gizo/helpers"
	"github.com/gizo-network/gizo/job"

	"github.com/gizo-network/gizo/core/merkletree"

	"github.com/kpango/glg"

	"github.com/boltdb/bolt"
	"github.com/jinzhu/now"
)

var (
	//ErrUnverifiedBlock when unable to verify block
	ErrUnverifiedBlock = errors.New("Unverified block cannot be added to the blockchain")
	//ErrJobNotFound when unable to find job in bc
	ErrJobNotFound = errors.New("Job not found")
	//ErrBlockNotFound when unable to finc block in the bc
	ErrBlockNotFound = errors.New("Blockinfo not found")
)

//BlockChain - singly linked list of blocks
type BlockChain struct {
	tip    []byte //! hash of latest block in the blockchain
	db     *bolt.DB
	mu     *sync.RWMutex
	logger *glg.Glg
}

//returns the blockinfo of the latest block in the blockchain
func (bc *BlockChain) getTip() []byte {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.tip
}

//sets the tip
func (bc *BlockChain) setTip(t []byte) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.tip = t
}

//returns the db
func (bc *BlockChain) getDB() *bolt.DB {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	return bc.db
}

//sets the db
func (bc *BlockChain) setDB(db *bolt.DB) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.db = db
}

//GetBlockInfo returns the blockinfo of a particular block from the db
func (bc *BlockChain) GetBlockInfo(hash string) (*BlockInfo, error) {
	var blockinfo *BlockInfo
	err := bc.getDB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		hashBytes, err := hex.DecodeString(hash)
		if err != nil {
			return err
		}
		blockinfoBytes := b.Get(hashBytes)
		if blockinfoBytes != nil {
			err = helpers.Deserialize(blockinfoBytes, &blockinfo)
			if err != nil {
				return err
			}
		} else {
			blockinfo = nil
		}
		return nil
	})
	if err != nil {
		return nil, err //! handle db failure error
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
			} else {
				blocks = append(blocks, *block)
			}
		} else {
			break
		}
	}
	return blocks, nil
}

//GetLatestHeight returns the height of the latest block to the blockchain
func (bc *BlockChain) GetLatestHeight() (int, error) {
	var lastBlock *BlockInfo
	err := bc.getDB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		lastBlockBytes := b.Get(bc.getTip())
		err := helpers.Deserialize(lastBlockBytes, &lastBlock)
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
		err := helpers.Deserialize(lastBlockBytes, &lastBlock)
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
		biBytes, err := helpers.Serialize(blockinfo)
		if err != nil {
			return err
		}
		err = b.Put(hashBytes, biBytes)
		if err != nil {
			return err
		}

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
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// return a BlockChainIterator to loop throught the blockchain
func (bc *BlockChain) iterator() *BlockChainIterator {
	return &BlockChainIterator{
		current: bc.getTip(),
		db:      bc.getDB(),
	}
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

//Verify verifies the blockchain
func (bc *BlockChain) Verify() (bool, error) {
	bci := bc.iterator()
	for {
		block, err := bci.Next()
		if err != nil {
			return false, err
		}
		if block.GetHeight() == 0 {
			return true, nil
		}
		verify, err := block.VerifyBlock()
		if err != nil {
			return false, err
		}
		if verify == false {
			return false, nil
		}
	}
}

//GetBlockHashes returns all the hashes of all the blocks in the current bc
func (bc *BlockChain) GetBlockHashes() ([]string, error) {
	var hashes []string
	bci := bc.iterator()
	for {
		block, err := bci.NextBlockinfo()
		if err != nil {
			return []string{}, err
		}
		hashes = append(hashes, block.GetHeader().GetHash())
		if block.GetHeight() == 0 {
			break
		}
	}
	return hashes, nil
}

//GetBlockHashesHex returns hashes (hex) of all the blocks in the bc
func (bc *BlockChain) GetBlockHashesHex() ([]string, error) {
	var hashes []string
	bci := bc.iterator()
	for {
		block, err := bci.NextBlockinfo()
		if err != nil {
			return []string{}, nil
		}
		hashes = append(hashes, block.GetHeader().GetHash())
		if block.GetHeight() == 0 {
			break
		}
	}
	return funk.ReverseStrings(hashes), nil
}

//CreateBlockChain initializes a db, set's the tip to GenesisBlock and returns the blockchain
func CreateBlockChain(nodeID string) *BlockChain {
	logger := helpers.Logger()
	logger.Log("Core: Creating blockchain database")
	InitializeDataPath()
	var dbFile string
	if os.Getenv("ENV") == "dev" {
		dbFile = path.Join(IndexPathDev, fmt.Sprintf(IndexDB, nodeID[len(nodeID)/2:])) //half the length of the node id
	} else {
		dbFile = path.Join(IndexPathProd, fmt.Sprintf(IndexDB, nodeID[len(nodeID)/2:])) //half the length of the node id
	}
	if helpers.FileExists(dbFile) {
		var tip []byte
		logger.Info("Core: Using existing blockchain")
		db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
		if err != nil {
			logger.Fatal(err)
		}
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BlockBucket))
			tip = b.Get([]byte("l"))
			return nil
		})
		if err != nil {
			logger.Fatal(err)
		}
		return &BlockChain{
			tip:    tip,
			db:     db,
			mu:     &sync.RWMutex{},
			logger: helpers.Logger(),
		}
	}
	genesis := GenesisBlock(nodeID)
	genesisHash, err := hex.DecodeString(genesis.GetHeader().GetHash())
	if err != nil {
		logger.Fatal(err)
	}
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
	if err != nil {
		logger.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(BlockBucket))
		if err != nil {
			logger.Fatal(err)
		}
		blockinfo := BlockInfo{
			Header:    genesis.GetHeader(),
			Height:    genesis.GetHeight(),
			TotalJobs: uint(len(genesis.GetNodes())),
			FileName:  genesis.fileStats().Name(),
			FileSize:  genesis.fileStats().Size(),
		}
		blockinfoBytes, err := helpers.Serialize(blockinfo)
		if err != nil {
			logger.Fatal(err)
		}

		if err = b.Put(genesisHash, blockinfoBytes); err != nil {
			logger.Fatal(err)
		}

		//latest block on the chain
		if err = b.Put([]byte("l"), genesisHash); err != nil {
			logger.Fatal(err)
		}
		return nil
	})
	if err != nil {
		logger.Fatal(err)
	}
	bc := &BlockChain{
		tip:    genesisHash,
		db:     db,
		mu:     &sync.RWMutex{},
		logger: logger,
	}
	return bc
}
