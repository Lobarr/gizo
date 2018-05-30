package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"time"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/helpers"

	"github.com/kpango/glg"
)

var (
	//ErrUnableToExport occurs when unable to write blockfile to disk
	ErrUnableToExport = errors.New("Unable to export block")
	//ErrHashModification occurs when hash of block is attempted to be modified
	ErrHashModification = errors.New("Attempt to modify hash value of block")
)

//Block - the foundation of blockchain
type Block struct {
	Header     BlockHeader
	Jobs       []*merkletree.MerkleNode
	Height     uint64
	ReceivedAt int64  //time it was received
	By         string // id of node that generated block
}

//GetHeader returns the block header
func (b Block) GetHeader() BlockHeader {
	return b.Header
}

//sets the block header
func (b *Block) setHeader(h BlockHeader) {
	b.Header = h
}

//GetNodes retuns the block's merklenodes
func (b Block) GetNodes() []*merkletree.MerkleNode {
	return b.Jobs
}

//sets merklenodes
func (b *Block) setNodes(j []*merkletree.MerkleNode) {
	b.Jobs = j
}

//GetHeight returns the block height
func (b Block) GetHeight() uint64 {
	return b.Height
}

//sets the block height
func (b *Block) setHeight(h uint64) {
	b.Height = h
}

//GetBy returns id of node that mined block
func (b Block) GetBy() string {
	return b.By
}

func (b *Block) setBy(id string) {
	b.By = id
}

//NewBlock returns a new block
func NewBlock(tree merkletree.MerkleTree, pHash string, height uint64, difficulty uint8, by string) *Block {
	block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: pHash,
			MerkleRoot:    tree.GetRoot(),
			Difficulty:    big.NewInt(int64(difficulty)),
		},
		Jobs:   tree.GetLeafNodes(),
		Height: height,
		By:     by,
	}
	pow := NewPOW(block)
	err := pow.run() //! mines block
	if err != nil {
		glg.Fatal(err)
	}
	err = block.Export()
	if err != nil {
		glg.Fatal(err)
	}
	return block
}

//Export writes block on disk
func (b Block) Export() error {
	InitializeDataPath()
	if b.IsEmpty() {
		return ErrUnableToExport
	}
	bBytes, err := helpers.Serialize(b)
	if err != nil {
		return err
	}
	if os.Getenv("ENV") == "dev" {
		err = ioutil.WriteFile(path.Join(BlockPathDev, fmt.Sprintf(BlockFile, b.Header.GetHash())), []byte(helpers.Encode64(bBytes)), os.FileMode(0555))
	} else {
		err = ioutil.WriteFile(path.Join(BlockPathProd, fmt.Sprintf(BlockFile, b.Header.GetHash())), []byte(helpers.Encode64(bBytes)), os.FileMode(0555))
	}
	return err
}

//Import reads block file into memory
func (b *Block) Import(hash string) {
	if b.IsEmpty() == false {
		glg.Warn("Overwriting umempty block")
	}
	var read []byte
	var err error
	if os.Getenv("ENV") == "dev" {
		read, err = ioutil.ReadFile(path.Join(BlockPathDev, fmt.Sprintf(BlockFile, hash)))
	} else {
		read, err = ioutil.ReadFile(path.Join(BlockPathProd, fmt.Sprintf(BlockFile, hash)))
	}
	if err != nil {
		glg.Fatal(err) //FIXME: handle block doesn't exist by asking peer
	}
	bBytes := helpers.Decode64(string(read))
	temp := &Block{}
	err = helpers.Deserialize(bBytes, &temp)
	if err != nil {
		glg.Fatal(err)
	}
	b.setHeader(temp.GetHeader())
	b.setHeight(temp.GetHeight())
	b.setNodes(temp.GetNodes())
	b.setBy(temp.GetBy())
}

//returns the file stats of a blockfile
func (b Block) fileStats() os.FileInfo {
	var info os.FileInfo
	var err error
	if os.Getenv("ENV") == "dev" {
		info, err = os.Stat(path.Join(BlockPathDev, fmt.Sprintf(BlockFile, b.Header.GetHash())))
	} else {
		info, err = os.Stat(path.Join(BlockPathProd, fmt.Sprintf(BlockFile, b.Header.GetHash())))
	}
	if os.IsNotExist(err) {
		glg.Fatal("Block file doesn't exist")
	}
	return info
}

//IsEmpty returns true is block is empty
func (b *Block) IsEmpty() bool {
	return b.Header.GetHash() == ""
}

//VerifyBlock verifies a block
func (b *Block) VerifyBlock() (bool, error) {
	glg.Info("Core: Verifying block - " + b.GetHeader().GetHash())
	pow := NewPOW(b)
	return pow.Validate()
}

//DeleteFile deletes block file on disk
func (b Block) DeleteFile() {
	glg.Info("Core: Deleting blockfile - " + b.GetHeader().GetHash())
	var err error
	if os.Getenv("ENV") == "dev" {
		err = os.Remove(path.Join(BlockPathDev, b.fileStats().Name()))
	} else {
		err = os.Remove(path.Join(BlockPathProd, b.fileStats().Name()))
	}
	if err != nil {
		glg.Fatal(err)
	}
}
