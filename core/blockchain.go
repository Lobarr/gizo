package core

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/asdine/storm"
	"github.com/thoas/go-funk"

	"github.com/gizo-network/gizo/helpers"

	"github.com/boltdb/bolt"
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
	db     storm.Node
	mu     *sync.RWMutex
	logger helpers.ILogger
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
func (bc *BlockChain) getDB() storm.TypeStore {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	return bc.db
}

//sets the db
func (bc *BlockChain) setDB(db storm.TypeStore) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.db = db
}

// return a BlockChainIterator to loop throught the blockchain
func (bc *BlockChain) iterator() *BlockChainIterator {
	return &BlockChainIterator{
		current: bc.getTip(),
		db:      bc.getDB(),
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

//InitGenesisBlock creates genesis block
func (bc *BlockChain) InitGenesisBlock(nodeID string) error {
	return bc.AddBlock(GenesisBlock(nodeID))
}

//CreateBlockChain initializes a db, set's the tip to GenesisBlock and returns the blockchain
func CreateBlockChain(nodeID string, loggerArg helpers.ILogger, dbArg storm.Node) IBlockChain {
	var logger helpers.ILogger
	var dbFile string
	var tip []byte
	var err error
	var db storm.Node

	if loggerArg == nil {
		logger = helpers.Logger()
	} else {
		logger = loggerArg
	}

	if os.Getenv("ENV") == "dev" {
		dbFile = path.Join(IndexPathDev, fmt.Sprintf(IndexDB, nodeID))
	} else {
		dbFile = path.Join(IndexPathProd, fmt.Sprintf(IndexDB, nodeID))
	}

	if helpers.FileExists(dbFile) {
		db, err = storm.Open(dbFile, storm.BoltOptions(0600, &bolt.Options{Timeout: time.Second * 2})) //TODO: handle closing db
		if err != nil {
			logger.Fatal(err)
		}
		err = db.Node.Get("config", "latestHash", &tip)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		db, err = storm.Open(dbFile, storm.BoltOptions(0600, &bolt.Options{Timeout: time.Second * 2}))
		if err != nil {
			logger.Fatal(err)
		}
		db.Bolt.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucket([]byte(BlockBucket))
			if err != nil {
				logger.Fatal(err)
			}
			return
		})
	}

	logger.Debug("Core: Creating blockchain database")
	return &BlockChain{
		tip:    tip,
		db:     db,
		mu:     &sync.RWMutex{},
		logger: logger,
	}
}
