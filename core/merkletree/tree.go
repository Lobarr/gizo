package merkletree

import (
	"errors"

	"github.com/gizo-network/gizo/job"
)

var (
	//ErrNodeDoesntExist occurs when merklenode doesn't exist
	ErrNodeDoesntExist = errors.New("core/merkletree: node doesn't exist")
	//ErrLeafNodesEmpty occurs when merkletree with empty leafnodes is attempted to be built
	ErrLeafNodesEmpty = errors.New("core/merkletree: leafnodes is empty")
	//ErrTreeNotBuilt occurs when tree hasn't been built
	ErrTreeNotBuilt = errors.New("core/merkletree: tree hasn't been built")
	//ErrTreeRebuildAttempt occurs when tree rebuild is attempted
	ErrTreeRebuildAttempt = errors.New("core/merkletree: attempt to rebuild tree")
	//ErrTooMuchLeafNodes occurs when merkletree is attempted to be built with more block than allowed
	ErrTooMuchLeafNodes = errors.New("core/merkletree: length of leaf nodes is greater than 24")
	//ErrJobDoesntExist occurs when job doesn't exist
	ErrJobDoesntExist = errors.New("core/merkletree: job doesn't exist")
)

// MerkleTree tree of jobs
type MerkleTree struct {
	Root      string // hash of built tree
	LeafNodes []IMerkleNode
}

// GetRoot returns root
func (merkleTree MerkleTree) GetRoot() string {
	return merkleTree.Root
}

func (merkleTree *MerkleTree) setRoot(root string) {
	merkleTree.Root = root
}

// GetLeafNodes return leafnodes
func (merkleTree MerkleTree) GetLeafNodes() []IMerkleNode {
	return merkleTree.LeafNodes
}

// SetLeafNodes return leafnodes
func (merkleTree *MerkleTree) SetLeafNodes(leafNodes []IMerkleNode) {
	merkleTree.LeafNodes = leafNodes
}

//Build builds merkle tree from leafs to root, hashed the root and sets it as the root of the merkletree
func (merkleTree *MerkleTree) Build() error {
	if merkleTree.GetRoot() != "" {
		return ErrTreeRebuildAttempt
	}
	if len(merkleTree.GetLeafNodes()) > MaxTreeJobs {
		return ErrTooMuchLeafNodes
	}
	var leafNodes = merkleTree.GetLeafNodes()
	for len(leafNodes) != 1 {
		var levelUp []IMerkleNode
		if len(leafNodes)%2 == 0 {
			for i := 0; i < len(leafNodes); i += 2 {
				parent, err := merge(leafNodes[i], leafNodes[i+1])
				if err != nil {
					return err
				}
				levelUp = append(levelUp, parent)
			}
		} else {
			leafNodes = append(leafNodes, leafNodes[len(leafNodes)-1]) //duplicate last to balance tree
			for i := 0; i < len(leafNodes); i += 2 {
				parent, err := merge(leafNodes[i], leafNodes[i+1])
				if err != nil {
					return err
				}
				levelUp = append(levelUp, parent)
			}
		}
		merkleTree.setRoot(leafNodes[0].GetHash())
		leafNodes = levelUp
	}
	return nil
}

//VerifyTree returns true if tree is verified
func (merkleTree MerkleTree) VerifyTree() bool {
	newTree, _ := NewMerkleTree(merkleTree.GetLeafNodes())
	return newTree.GetRoot() == merkleTree.GetRoot()
}

//SearchNode returns true if node with hash exists
func (merkleTree MerkleTree) SearchNode(hash string) (IMerkleNode, error) {
	if len(merkleTree.GetLeafNodes()) == 0 {
		return nil, ErrLeafNodesEmpty
	}
	for _, leafNode := range merkleTree.GetLeafNodes() {
		if leafNode.GetHash() == hash {
			return leafNode, nil
		}
	}
	return nil, ErrNodeDoesntExist
}

//SearchJob returns job from the tree
func (merkleTree MerkleTree) SearchJob(ID string) (job.IJob, error) {
	if len(merkleTree.GetLeafNodes()) == 0 {
		return nil, ErrLeafNodesEmpty
	}
	for _, leafNode := range merkleTree.GetLeafNodes() {
		if leafNode.GetJob().GetID() == ID {
			return leafNode.GetJob(), nil
		}
	}
	return nil, ErrNodeDoesntExist
}

// NewMerkleTree returns empty merkletree
func NewMerkleTree(nodes []IMerkleNode) (IMerkleTree, error) {
	merkleTree := &MerkleTree{
		LeafNodes: nodes,
	}
	err := merkleTree.Build()
	if err != nil {
		return nil, err
	}
	return merkleTree, nil
}

//merges two nodes
func merge(leftNode, rightNode IMerkleNode) (IMerkleNode, error) {
	job, err := MergeJobs(leftNode, rightNode)
	if err != nil {
		return nil, err
	}
	return NewNode(job, leftNode, rightNode)
}
