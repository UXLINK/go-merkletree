package db

import (
	"context"
	"errors"
)

// ErrNotFound is used by the implementations of the interface db.Storage for
// when a key is not found in the storage
var ErrNotFound = errors.New("key not found")

type TreeNode struct {
	MtAddress string
	Data      string
	Hash      string
	Level     int
	LevelNo   int
}

type NodePos struct {
	Level   int
	LevelNo int
}

type Storage interface {
	Insert(ctx context.Context, node *TreeNode) error
	Update(ctx context.Context, node *TreeNode) error
	FindRootNode(ctx context.Context, address string) (*TreeNode, error)
	FindMaxNoOfLeaf(ctx context.Context, address string) (int, error)
	FindOneByLeafData(ctx context.Context, address string, data string) (*TreeNode, error)
	FindMultiTreeNode(ctx context.Context, address string, nodePoses []*NodePos) ([]*TreeNode, error)
	FindNodesByLevel(ctx context.Context, address string, level int) ([]*TreeNode, error)
}
