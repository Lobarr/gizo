package merkletree

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gizo-network/gizo/helpers"

	"github.com/gizo-network/gizo/job"
)

// MerkleNode nodes that make a merkletree
type MerkleNode struct {
	Hash  string //hash of a job struct
	Job   job.Job
	Left  *MerkleNode
	Right *MerkleNode
}

// GetHash returns hash
func (n MerkleNode) GetHash() string {
	return n.Hash
}

// generates hash value of merklenode
func (n *MerkleNode) setHash() error {
	l, err := helpers.Serialize(n.Left)
	if err != nil {
		return err
	}
	r, err := helpers.Serialize(n.Right)
	if err != nil {
		return err
	}

	jobBytes, err := helpers.Serialize(n.Job)
	if err != nil {
		return err
	}
	headers := bytes.Join([][]byte{l, r, jobBytes}, []byte{})
	if err != nil {
		return err
	}
	hash := sha256.Sum256(headers)
	n.Hash = hex.EncodeToString(hash[:])
	return nil
}

// GetJob returns job
func (n MerkleNode) GetJob() job.Job {
	return n.Job
}

// SetJob setter for job
func (n *MerkleNode) SetJob(j job.Job) {
	n.Job = j
}

// GetLeftNode return leftnode
func (n MerkleNode) GetLeftNode() MerkleNode {
	return *n.Left
}

// SetLeftNode setter for leftnode
func (n *MerkleNode) SetLeftNode(l MerkleNode) {
	n.Left = &l
}

// GetRightNode return rightnode
func (n MerkleNode) GetRightNode() MerkleNode {
	return *n.Right
}

//SetRightNode setter for rightnode
func (n *MerkleNode) SetRightNode(r MerkleNode) {
	n.Right = &r
}

//IsLeaf checks if the merklenode is a leaf node
func (n *MerkleNode) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}

//IsEmpty check if the merklenode is empty
func (n *MerkleNode) IsEmpty() (bool, error) {
	empty, err := n.GetJob().IsEmpty()
	if err != nil {
		return false, err
	}
	return n.Right == nil && n.Left == nil && empty && n.GetHash() == "", nil
}

//IsEqual check if the input merklenode equals the merklenode calling the function
func (n MerkleNode) IsEqual(x MerkleNode) (bool, error) {
	nBytes, err := helpers.Serialize(n)
	if err != nil {
		return false, err
	}
	xBytes, err := helpers.Serialize(x)
	if err != nil {
		return false, err
	}
	return bytes.Equal(nBytes, xBytes), nil
}

//NewNode returns a new merklenode
func NewNode(j job.Job, lNode, rNode *MerkleNode) (*MerkleNode, error) {
	n := &MerkleNode{
		Left:  lNode,
		Right: rNode,
		Job:   j,
	}
	err := n.setHash()
	return n, err
}

//MergeJobs merges two jobs into one
func MergeJobs(x, y MerkleNode) (job.Job, error) {
	xTask, err := x.GetJob().GetTask()
	if err != nil {
		return job.Job{}, err
	}
	yTask, err := y.GetJob().GetTask()
	if err != nil {
		return job.Job{}, err
	}
	return job.Job{
		ID:        x.GetJob().GetID() + y.GetJob().GetID(),
		Hash:      x.GetJob().GetHash() + y.GetJob().GetHash(),
		Execs:     append(x.GetJob().GetExecs(), y.GetJob().GetExecs()...),
		Task:      helpers.Encode64(bytes.Join([][]byte{[]byte(xTask), []byte(yTask)}, []byte{})),
		Signature: x.GetJob().GetSignature() + y.GetJob().GetSignature(),
	}, nil
}
