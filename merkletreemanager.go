package merkletree

import (
	"context"
	"github.com/UXUYLabs/go-merkletree/db"
	"github.com/UXUYLabs/go-merkletree/db/memory"
	"log"
	"os"
)

// MerkleTree is the structure for the Merkle tree.
type MerkleTreeManager struct {
	Logger
	ctx     context.Context
	storage db.Storage
}

func NewMemoryMerkleTreeManager(ctx context.Context) (*MerkleTreeManager, error) {
	return &MerkleTreeManager{
		Logger:  PrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags)),
		ctx:     ctx,
		storage: memory.NewMemoryStorage(),
	}, nil
}
func NewMerkleTreeManager(ctx context.Context, storage db.Storage) (*MerkleTreeManager, error) {
	return &MerkleTreeManager{
		Logger:  PrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags)),
		ctx:     ctx,
		storage: storage,
	}, nil
}

func (mm *MerkleTreeManager) CreateMerkleTree(mtAddress string) (*MerkleTree, error) {
	tree, err := NewMerkleTree(mm.ctx, mm.storage, mtAddress)
	if err != nil {
		mm.Error("CreateMerkleTree err:%v", err)
		return nil, err
	}

	return tree, nil
}
