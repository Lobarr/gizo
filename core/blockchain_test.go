package core

import (
	"os"
	"testing"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestNewBlockChain(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	bc := CreateBlockChain("74657374")
	assert.NotNil(t, bc)
}
func TestAddBlock(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("74657374")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, 1, 10, "74657374")
	err = bc.AddBlock(block)
	assert.NoError(t, err)
	latestHeight, err := bc.GetLatestHeight()
	assert.NoError(t, err)
	assert.Equal(t, 1, latestHeight)
}

func TestVerify(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("74657374")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "74657374")
	bc.AddBlock(block)
	verify, err := bc.Verify()
	assert.NoError(t, err)
	assert.True(t, verify)
}

func TestGetBlockInfo(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("74657374")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "74657374")
	bc.AddBlock(block)
	blockinfo, err := bc.GetBlockInfo(block.GetHeader().GetHash())
	assert.NoError(t, err)
	assert.NotNil(t, blockinfo)
}

func TestGetBlocksWithinMinute(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("74657374")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "74657374")
	bc.AddBlock(block)
	lastMinute, err := bc.GetBlocksWithinMinute()
	assert.NoError(t, err)
	assert.NotNil(t, lastMinute)
}

func TestGetLatestHeight(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("74657374")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "74657374")
	bc.AddBlock(block)
	latestHeight, err := bc.GetLatestHeight()
	assert.NoError(t, err)
	assert.NotNil(t, latestHeight)
}

func TestFindMerkleNode(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, pub := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	exec1, err := job.NewExec([]interface{}{2}, 5, job.NORMAL, 0, 0, 0, 0, pub, envs, "test")
	j.AddExec(*exec1)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("74657374")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "74657374")
	bc.AddBlock(block)
	f, err := bc.FindMerkleNode(nodes[10].GetHash())
	assert.NoError(t, err)
	assert.Equal(t, nodes[10].GetHash(), f.GetHash())
}

func TestFindJob(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("74657374")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "74657374")
	bc.AddBlock(block)
	f, err := bc.FindJob(nodes[4].GetJob().GetID())
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestFindExec(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, pub := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	exec1, err := job.NewExec([]interface{}{2}, 5, job.NORMAL, 0, 0, 0, 0, pub, envs, "test")
	j.AddExec(*exec1)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("74657374")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "74657374")
	bc.AddBlock(block)
	f, err := bc.FindExec(j.GetID(), nodes[4].GetJob().GetExecs()[0].GetHash())
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestGetJobExecs(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, pub := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	exec1, err := job.NewExec([]interface{}{2}, 5, job.NORMAL, 0, 0, 0, 0, pub, envs, "test")
	j.AddExec(*exec1)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("74657374")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "74657374")
	bc.AddBlock(block)
	e, err := bc.GetJobExecs(j.GetID())
	assert.NoError(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, exec1.GetBy(), e[0].GetBy())
}

func TestGetBlockHashes(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("74657374")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "74657374")
	err = bc.AddBlock(block)
	assert.NoError(t, err)
	hashes, err := bc.GetBlockHashes()
	assert.NoError(t, err)
	assert.NotNil(t, hashes)
}

func TestGetBlockHashesHex(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("74657374")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "74657374")
	err = bc.AddBlock(block)
	assert.NoError(t, err)
	hashes, err := bc.GetBlockHashesHex()
	assert.NoError(t, err)
	assert.NotNil(t, hashes)
}

func TestCreateBlockChain(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	bc := CreateBlockChain("74657374")
	assert.NotNil(t, bc)
}

func TestGetBlockByHeight(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	bc := CreateBlockChain("74657374")
	assert.NotNil(t, bc)
	b, err := bc.GetBlockByHeight(0)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, b.GetHeight(), uint64(0))

	b2, err := bc.GetBlockByHeight(20)
	assert.Error(t, err)
	assert.Nil(t, b2)
}

func TestLatest15(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("74657374")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "74657374")
	bc.AddBlock(block)
	latest, err := bc.GetLatest15()
	assert.NoError(t, err)
	assert.NotNil(t, latest)
}
