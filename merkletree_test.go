package merkletree

import (
	"context"
	"fmt"
	"github.com/UXUYLabs/go-merkletree/keccak256"
	"github.com/stretchr/testify/assert"
	"testing"
)

var merkleTreeManager *MerkleTreeManager
var err error

func setup() {
	ctx := context.Background()
	merkleTreeManager, err = NewMemoryMerkleTreeManager(ctx)
	if err != nil {
		fmt.Printf("NewMerkleTree err:%v\n", err)
		return
	}
}

func TestAppend(t *testing.T) {
	setup()

	tree, err := merkleTreeManager.CreateMerkleTree("1637704523306766336")
	if err != nil {
		fmt.Printf("NewMerkleTree err:%v\n", err)
		return
	}

	err = tree.AppendLeaf("0x8b1b201E91966957f18bBcDDB520c53c521bF5cd")
	if err != nil {
		fmt.Printf("AppendLeaf err :%v\n", err)
		return
	}

	err = tree.AppendLeaf("0xeA726629EC5fe5cE300000d1a8c89B3054A22cE7")
	if err != nil {
		fmt.Printf("AppendLeaf err :%v\n", err)
		return
	}

	err = tree.AppendLeaf("0x00440DC3377A8a6b745aB5F92fD850b7c7291DdE")
	if err != nil {
		fmt.Printf("AppendLeaf err :%v\n", err)
		return
	}

	err = tree.AppendLeaf("0x63120cc1c7Bb0a42C2D77D27faB9EDd2560F9cA3")
	if err != nil {
		fmt.Printf("AppendLeaf err :%v\n", err)
		return
	}

	err = tree.AppendLeaf("0x7e533CF779A533eD8f9C1b8E5C3d7F936335ca54")
	if err != nil {
		fmt.Printf("AppendLeaf err :%v\n", err)
		return
	}

	_ = tree.printTree()

	// proof
	proofes, err := tree.GenerateProof("0x7e533CF779A533eD8f9C1b8E5C3d7F936335ca54")
	if err != nil {
		fmt.Printf("NewMerkleTree err:%v\n", err)
		return
	}

	for _, proof := range proofes {
		fmt.Printf("NewMerkleTree proofes:%+v\n", keccak256.Bytes2Hex(proof))
	}

	proof, err := tree.VerifyProof(proofes, "0x7e533CF779A533eD8f9C1b8E5C3d7F936335ca54")
	//logx.WithContext(ctx).Infof("NewMerkleTree proof:%+v", proof)
	assert.Equal(t, proof, true)

	proof, err = tree.VerifyProof(proofes, "0xa6820eeA9B5BB08Ab1cD693128Bb85Ad460a8e6E")
	assert.Equal(t, proof, false)
}
