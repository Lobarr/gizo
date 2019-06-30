package merkletree

import "github.com/gizo-network/gizo/job"

// IMerkleNode interface of a merklenode
type IMerkleNode interface {
	GetHash() string
	GetJob() job.IJob
	SetJob(j job.IJob)
	GetLeftNode() IMerkleNode
	SetLeftNode(IMerkleNode)
	GetRightNode() IMerkleNode
	SetRightNode(MerkleNode)
	IsLeaf() bool
	IsEmpty() (bool, error)
}
