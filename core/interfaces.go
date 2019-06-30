package core

import (
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
