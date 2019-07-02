package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"time"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/helpers"
)

var (
	//ErrUnableToExport occurs when unable to write blockfile to disk
	ErrUnableToExport = errors.New("Unable to export block")
	//ErrHashModification occurs when hash of block is attempted to be modified
	ErrHashModification = errors.New("Attempt to modify hash value of block")
)

//Block - the foundation of blockchain
type Block struct {
	Header     IBlockHeader `storm:"index`
	Jobs       []merkletree.IMerkleNode
	Height     uint64 `storm:"id"`
	ReceivedAt int64  `storm:"index"` //time it was received
	By         string `storm:"index"` // id of node that generated block
}

//GetHeader returns the block header
func (block Block) GetHeader() IBlockHeader {
	return block.Header
}

//sets the block header
func (block *Block) setHeader(blockHeader IBlockHeader) {
	block.Header = blockHeader
}

//GetNodes retuns the block's merklenodes
func (block Block) GetNodes() []merkletree.IMerkleNode {
	return block.Jobs
}

//sets merklenodes
func (block *Block) setNodes(merkleNodes []merkletree.IMerkleNode) {
	block.Jobs = merkleNodes
}

//GetHeight returns the block height
func (block Block) GetHeight() uint64 {
	return block.Height
}

//sets the block height
func (block *Block) setHeight(h uint64) {
	block.Height = h
}

//GetBy returns id of node that mined block
func (block Block) GetBy() string {
	return block.By
}

func (block *Block) setBy(id string) {
	block.By = id
}

//Export writes block on disk
func (block Block) Export() error {
	InitializeDataPath()
	if block.IsEmpty() {
		return ErrUnableToExport
	}
	blockBytes, err := json.Marshal(block)
	if err != nil {
		return err
	}
	if os.Getenv("ENV") == "dev" {
		err = ioutil.WriteFile(path.Join(BlockPathDev, fmt.Sprintf(BlockFile, block.Header.GetHash())), []byte(helpers.Encode64(blockBytes)), os.FileMode(0755))
	} else {
		err = ioutil.WriteFile(path.Join(BlockPathProd, fmt.Sprintf(BlockFile, block.Header.GetHash())), []byte(helpers.Encode64(blockBytes)), os.FileMode(0755))
	}
	return err
}

//Import reads block file into memory
func (block *Block) Import(hash string) error {
	var readBlock []byte
	var err error
	if os.Getenv("ENV") == "dev" {
		readBlock, err = ioutil.ReadFile(path.Join(BlockPathDev, fmt.Sprintf(BlockFile, hash)))
	} else {
		readBlock, err = ioutil.ReadFile(path.Join(BlockPathProd, fmt.Sprintf(BlockFile, hash)))
	}
	if err != nil {
		return err //TODO: handle block doesn't exist by asking peer
	}
	blockBytes, err := helpers.Decode64(string(readBlock))
	if err != nil {
		return err
	}
	temp := &Block{}
	err = json.Unmarshal(blockBytes, &temp)
	if err != nil {
		return err
	}
	block.setHeader(temp.GetHeader())
	block.setHeight(temp.GetHeight())
	block.setNodes(temp.GetNodes())
	block.setBy(temp.GetBy())
	return nil
}

//returns the file stats of a blockfile
func (block Block) fileStats() os.FileInfo {
	var info os.FileInfo
	var err error
	if os.Getenv("ENV") == "dev" {
		info, err = os.Stat(path.Join(BlockPathDev, fmt.Sprintf(BlockFile, block.Header.GetHash())))
	} else {
		info, err = os.Stat(path.Join(BlockPathProd, fmt.Sprintf(BlockFile, block.Header.GetHash())))
	}
	if os.IsNotExist(err) {
		helpers.Logger().Fatal("Core: block file doesn't exist, blocks folder corrupted")
	}
	return info
}

//IsEmpty returns true is block is empty
func (block *Block) IsEmpty() bool {
	return block.Header.GetHash() == ""
}

//VerifyBlock verifies a block
func (block *Block) VerifyBlock() (bool, error) {
	pow := NewPOW(block)
	return pow.Validate()
}

//DeleteFile deletes block file on disk
func (block Block) DeleteFile() error {
	var err error
	if os.Getenv("ENV") == "dev" {
		err = os.Remove(path.Join(BlockPathDev, block.fileStats().Name()))
	} else {
		err = os.Remove(path.Join(BlockPathProd, block.fileStats().Name()))
	}
	if err != nil {
		return err
	}
	return nil
}

//NewBlock returns a new block
func NewBlock(merkleTree merkletree.IMerkleTree, prevBlockHash string, height uint64, difficulty uint8, by string) (IBlock, error) {
	block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: prevBlockHash,
			MerkleRoot:    merkleTree.GetRoot(),
			Difficulty:    big.NewInt(int64(difficulty)),
		},
		Jobs:   merkleTree.GetLeafNodes(),
		Height: height,
		By:     by,
	}
	pow := NewPOW(block)
	err := pow.run() //! mines block
	if err != nil {
		return nil, err
	}
	err = block.Export()
	if err != nil {
		return nil, err
	}
	return block, nil
}
