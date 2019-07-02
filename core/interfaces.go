package core

import (
	"github.com/asdine/storm"
	"math/big"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
)

//IBlockChain  interface for blockchain struct
type IBlockChain interface {
	GetBlockInfo(string) (*BlockInfo, error)
	GetPrevHash() (string, error)
	GetBlocksWithinMinute() ([]Block, error)
	GetBlockByHeight(int) (*Block, error)
	GetLatest15() ([]Block, error)
	GetLatestHeight() (int, error)
	GetLatestBlock() (*Block, error)
	GetNextHeight() (uint64, error)
	AddBlock(*Block) error
	FindJob(string) (*job.Job, error)
	FindExec(string, string) (*job.Exec, error)
	GetJobExecs(string) ([]job.Exec, error)
	FindMerkleNode(string) (*merkletree.MerkleNode, error)
	Verify() (bool, error)
	GetBlockHashes() ([]string, error)
	GetBlockHashesHex() ([]string, error)
	InitGenesisBlock(string) error
}

// IBlockHeader interface for blockheader
type IBlockHeader interface {
	GetTimestamp() int64
	GetPrevBlockHash() string
	GetMerkleRoot() string
	GetNonce() uint64
	GetDifficulty() *big.Int
	GetHash() string
}

//IBlock interface for block
type IBlock interface {
	GetHeader() IBlockHeader
	GetNodes() []merkletree.IMerkleNode
	GetHeight() uint64
	GetBy() string
	Export() error
	Import(string) error
	IsEmpty() bool
	DeleteFile() error
}

//IBlockIterator interface for blockiterator
type IBlockIterator interface {
	GetCurrent() []byte
	Next() (IBlock, error)
	NextBlockinfo() (IBlockInfo, error)
}

//IBlockInfo interface for blockinfo
type IBlockInfo interface {
}

type IStorm interface {
	...storm.Node
	...storm.KeyValueStore
	...storm.Query
	...storm.TypeStore
	...storm.Finder
}