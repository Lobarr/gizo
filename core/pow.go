package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"math/big"
	"strconv"

	"github.com/gizo-network/gizo/helpers"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/kpango/glg"
)

var maxNonce = math.MaxInt64

//POW - consensus algorithm
type POW struct {
	difficulty uint8
	block      *Block
	target     *big.Int
	logger     *glg.Glg
}

//sets block
func (p *POW) setBlock(b *Block) {
	p.block = b
}

//GetBlock returns block
func (p POW) GetBlock() *Block {
	return p.block
}

//sets target difficult
func (p *POW) setTarget(t *big.Int) {
	p.target = t
}

//GetTarget returns target
func (p POW) GetTarget() *big.Int {
	return p.target
}

//GetDifficulty returns difficulty
func (p POW) GetDifficulty() uint8 {
	return p.difficulty
}

func (p *POW) setDifficulty(d uint8) {
	p.difficulty = d
}

//merges info and returns it as byttes
func (p POW) prepareData(nonce int) ([]byte, error) {
	tree := merkletree.MerkleTree{Root: p.GetBlock().GetHeader().GetMerkleRoot(), LeafNodes: p.GetBlock().GetNodes()}
	mBytes, err := json.Marshal(tree)
	if err != nil {
		p.logger.Fatal(err)
	}
	hashBytes, err := hex.DecodeString(p.block.GetHeader().GetPrevBlockHash())
	if err != nil {
		return nil, err
	}
	byBytes, err := hex.DecodeString(p.block.GetBy())
	if err != nil {
		return nil, err
	}
	data := bytes.Join(
		[][]byte{
			hashBytes,
			[]byte(strconv.FormatInt(p.GetBlock().GetHeader().GetTimestamp(), 10)),
			mBytes,
			[]byte(strconv.FormatInt(int64(nonce), 10)),
			[]byte(strconv.FormatInt(int64(p.GetBlock().GetHeight()), 10)),
			[]byte(strconv.FormatInt(int64(p.GetBlock().GetHeader().GetDifficulty().Int64()), 10)),
			byBytes,
		},
		[]byte{},
	)
	return data, nil
}

//Run looks for a hash that is less than the current target difficulty
func (p *POW) run() error {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	for nonce < maxNonce {
		data, err := p.prepareData(nonce)
		if err != nil {
			return err
		}
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(p.GetTarget()) == -1 {
			break
		} else {
			nonce++
		}
	}
	p.GetBlock().Header.setHash(hex.EncodeToString(hash[:]))
	p.GetBlock().Header.setNonce(uint64(nonce))
	return nil
}

//Validate - validates POW
func (p *POW) Validate() (bool, error) {
	var hashInt big.Int
	data, err := p.prepareData(int(p.GetBlock().GetHeader().GetNonce()))
	if err != nil {
		return false, err
	}
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(p.GetTarget()) == -1, nil
}

//NewPOW returns POW
func NewPOW(b *Block) *POW {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-b.GetHeader().GetDifficulty().Int64()))
	pow := &POW{
		target: target,
		block:  b,
		logger: helpers.Logger(),
	}
	return pow
}
