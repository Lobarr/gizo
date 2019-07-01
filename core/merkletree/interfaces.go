package merkletree

import "github.com/gizo-network/gizo/job"

// IMerkleNode interface of a merklenode
type IMerkleNode interface {
	GetHash() string
	GetJob() job.IJob
	SetJob(job.IJob)
	GetLeftNode() IMerkleNode
	SetLeftNode(IMerkleNode)
	GetRightNode() IMerkleNode
	SetRightNode(IMerkleNode)
	IsLeaf() bool
	IsEmpty() (bool, error)
}

// IMerkleTree interface for merkletree
type IMerkleTree interface {
	GetRoot() string
	GetLeafNodes() []IMerkleNode
	SetLeafNodes([]IMerkleNode)
	Build() error
	VerifyTree() bool
	SearchNode(string) (IMerkleNode, error)
	SearchJob(string) (job.IJob, error)
}
