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
	LeafNodes []*MerkleNode
}

// GetRoot returns root
func (m MerkleTree) GetRoot() string {
	return m.Root
}

func (m *MerkleTree) setRoot(r string) {
	m.Root = r
}

// GetLeafNodes return leafnodes
func (m MerkleTree) GetLeafNodes() []*MerkleNode {
	return m.LeafNodes
}

// SetLeafNodes return leafnodes
func (m *MerkleTree) SetLeafNodes(l []*MerkleNode) {
	m.LeafNodes = l
}

//Build builds merkle tree from leafs to root, hashed the root and sets it as the root of the merkletree
func (m *MerkleTree) Build() error {
	if m.GetRoot() != "" {
		return ErrTreeRebuildAttempt
	}
	if len(m.GetLeafNodes()) > MaxTreeJobs {
		return ErrTooMuchLeafNodes
	}
	var shrink = m.GetLeafNodes()
	for len(shrink) != 1 {
		var levelUp []*MerkleNode
		if len(shrink)%2 == 0 {
			for i := 0; i < len(shrink); i += 2 {
				parent, err := merge(*shrink[i], *shrink[i+1])
				if err != nil {
					return err
				}
				levelUp = append(levelUp, parent)
			}
		} else {
			shrink = append(shrink, shrink[len(shrink)-1]) //duplicate last to balance tree
			for i := 0; i < len(shrink); i += 2 {
				parent, err := merge(*shrink[i], *shrink[i+1])
				if err != nil {
					return err
				}
				levelUp = append(levelUp, parent)
			}
		}
		m.setRoot(shrink[0].GetHash())
		shrink = levelUp
	}
	return nil
}

//VerifyTree returns true if tree is verified
func (m MerkleTree) VerifyTree() bool {
	t, _ := NewMerkleTree(m.GetLeafNodes())
	return t.GetRoot() == m.GetRoot()
}

//SearchNode returns true if node with hash exists
func (m MerkleTree) SearchNode(hash string) (*MerkleNode, error) {
	if len(m.GetLeafNodes()) == 0 {
		return nil, ErrLeafNodesEmpty
	}
	for _, n := range m.GetLeafNodes() {
		if n.GetHash() == hash {
			return n, nil
		}
	}
	return nil, ErrNodeDoesntExist
}

//SearchJob returns job from the tree
func (m MerkleTree) SearchJob(ID string) (*job.Job, error) {
	if len(m.GetLeafNodes()) == 0 {
		return nil, ErrLeafNodesEmpty
	}
	for _, n := range m.GetLeafNodes() {
		if n.GetJob().GetID() == ID {
			return &n.Job, nil
		}
	}
	return nil, ErrNodeDoesntExist
}

// NewMerkleTree returns empty merkletree
func NewMerkleTree(nodes []*MerkleNode) (*MerkleTree, error) {
	t := &MerkleTree{
		LeafNodes: nodes,
	}
	err := t.Build()
	if err != nil {
		return nil, err
	}
	return t, nil
}

//merges two nodes
func merge(left, right MerkleNode) (*MerkleNode, error) {
	job, err := MergeJobs(left, right)
	if err != nil {
		return nil, err
	}
	return NewNode(job, &left, &right)
}
