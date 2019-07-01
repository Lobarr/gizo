package merkletree

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/gizo-network/gizo/helpers"

	"github.com/gizo-network/gizo/job"
)

// MerkleNode nodes that make a merkletree
type MerkleNode struct {
	Hash  string //hash of a job struct
	Job   job.IJob
	Left  IMerkleNode
	Right IMerkleNode
}

// GetHash returns hash
func (merkleNode MerkleNode) GetHash() string {
	return merkleNode.Hash
}

// generates hash value of merklenode
func (merkleNode *MerkleNode) setHash() error {
	leftNode, err := json.Marshal(merkleNode.Left)
	if err != nil {
		return err
	}
	rightNode, err := json.Marshal(merkleNode.Right)
	if err != nil {
		return err
	}
	jobBytes, err := json.Marshal(merkleNode.Job)
	if err != nil {
		return err
	}
	headers := bytes.Join([][]byte{leftNode, rightNode, jobBytes}, []byte{})
	hash := sha256.Sum256(headers)
	merkleNode.Hash = hex.EncodeToString(hash[:])
	return nil
}

// GetJob returns job
func (merkleNode MerkleNode) GetJob() job.IJob {
	return merkleNode.Job
}

// SetJob setter for job
func (merkleNode *MerkleNode) SetJob(j job.IJob) {
	merkleNode.Job = j
}

// GetLeftNode return leftnode
func (merkleNode MerkleNode) GetLeftNode() IMerkleNode {
	return merkleNode.Left
}

// SetLeftNode setter for leftnode
func (merkleNode *MerkleNode) SetLeftNode(leftNode IMerkleNode) {
	merkleNode.Left = leftNode
}

// GetRightNode return rightnode
func (merkleNode MerkleNode) GetRightNode() IMerkleNode {
	return merkleNode.Right
}

//SetRightNode setter for rightnode
func (merkleNode *MerkleNode) SetRightNode(rightNode IMerkleNode) {
	merkleNode.Right = rightNode
}

//IsLeaf checks if the merklenode is a leaf node
func (merkleNode *MerkleNode) IsLeaf() bool {
	return merkleNode.Left == nil && merkleNode.Right == nil
}

//IsEmpty check if the merklenode is empty
func (merkleNode *MerkleNode) IsEmpty() (bool, error) {
	empty, err := merkleNode.GetJob().IsEmpty()
	if err != nil {
		return false, err
	}
	return merkleNode.Right == nil && merkleNode.Left == nil && empty && merkleNode.GetHash() == "", nil
}

//IsEqual check if the input merklenode equals the merklenode calling the function
func (merkleNode MerkleNode) IsEqual(compareMerkleNode IMerkleNode) (bool, error) {
	merkleNodeBytes, err := json.Marshal(merkleNode)
	if err != nil {
		return false, err
	}
	compareMerkleNodeBytes, err := json.Marshal(compareMerkleNode)
	if err != nil {
		return false, err
	}
	return bytes.Equal(merkleNodeBytes, compareMerkleNodeBytes), nil
}

//NewNode returns a new merklenode
func NewNode(job job.IJob, leftNode, rightNode IMerkleNode) (IMerkleNode, error) {
	merkleNode := &MerkleNode{
		Left:  leftNode,
		Right: rightNode,
		Job:   job,
	}
	err := merkleNode.setHash()
	return merkleNode, err
}

//MergeJobs merges two jobs into one
func MergeJobs(x, y IMerkleNode) (job.IJob, error) {
	xTask, err := x.GetJob().GetTask()
	if err != nil {
		return nil, err
	}
	yTask, err := y.GetJob().GetTask()
	if err != nil {
		return nil, err
	}
	return &job.Job{
		ID:        x.GetJob().GetID() + y.GetJob().GetID(),
		Hash:      x.GetJob().GetHash() + y.GetJob().GetHash(),
		Execs:     append(x.GetJob().GetExecs(), y.GetJob().GetExecs()...),
		Task:      helpers.Encode64(bytes.Join([][]byte{[]byte(xTask), []byte(yTask)}, []byte{})),
		Signature: x.GetJob().GetSignature() + y.GetJob().GetSignature(),
	}, nil
}
