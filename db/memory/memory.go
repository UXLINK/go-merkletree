package memory

import (
	"context"
	"fmt"
	"github.com/UXUYLabs/go-merkletree/db"
)

// Tree record a single trees
type Tree map[int]map[int]*db.TreeNode

// TreeMap record different trees
type TreeMap map[string]Tree

// DataMap is Position to a data node of the tree
type DataMap map[string]map[string]*db.TreeNode

type MemoryStorage struct {
	db.Storage
	dataMap DataMap
	treeMap TreeMap
}

func NewMemoryStorage() *MemoryStorage {
	dataMap := make(DataMap)
	treeMap := make(TreeMap)

	return &MemoryStorage{
		dataMap: dataMap,
		treeMap: treeMap,
	}
}

func (s *MemoryStorage) Insert(ctx context.Context, node *db.TreeNode) error {
	//fmt.Printf("Insert node %+v\n", node)
	// 记录数据
	if s.dataMap[node.MtAddress] == nil {
		s.dataMap[node.MtAddress] = make(map[string]*db.TreeNode)
	}
	s.dataMap[node.MtAddress][node.Data] = node

	tree := s.treeMap[node.MtAddress]
	if tree == nil {
		tree = make(Tree)
		s.treeMap[node.MtAddress] = tree
	}

	if len(tree) <= node.Level {
		tree[node.Level] = make(map[int]*db.TreeNode)
	}
	tree[int(node.Level)][int(node.LevelNo)] = node
	//fmt.Printf("Insert tree, Level:%d LevelNo:%d node:%+v\n", node.Level, node.LevelNo, tree[node.Level][node.LevelNo])
	return nil
}

func (s *MemoryStorage) Update(ctx context.Context, node *db.TreeNode) error {
	s.dataMap[node.MtAddress][node.Data] = node
	tree := s.treeMap[node.MtAddress]
	if tree == nil || len(tree) < node.Level+1 || len(tree[node.Level]) < (node.LevelNo+1) {
		return db.ErrNotFound
	}
	tree[node.Level][node.LevelNo] = node
	return nil
}

func (s *MemoryStorage) FindRootNode(ctx context.Context, address string) (*db.TreeNode, error) {
	tree := s.treeMap[address]
	if tree == nil {
		return nil, db.ErrNotFound
	}

	level := len(tree) - 1
	//fmt.Printf("FindRootNode level:%d, levelNo:%d\n", level, len(tree[level])-1)
	return tree[level][len(tree[level])-1], nil
}

func (s *MemoryStorage) FindMaxNoOfLeaf(ctx context.Context, address string) (int, error) {
	tree := s.treeMap[address]
	if tree == nil {
		return -1, db.ErrNotFound
	}

	//fmt.Printf("FindMaxNoOfLeaf MaxNo:%d\n", len(tree[0])-1)

	return len(tree[0]) - 1, nil
}

func (s *MemoryStorage) FindOneByLeafData(ctx context.Context, address string, data string) (*db.TreeNode, error) {
	treeMap := s.dataMap[address]
	if treeMap == nil {
		return nil, db.ErrNotFound
	}

	if treeMap[data] == nil {
		return nil, db.ErrNotFound
	}

	return treeMap[data], nil
}

func (s *MemoryStorage) FindMultiTreeNode(ctx context.Context, address string, nodePoses []*db.NodePos) ([]*db.TreeNode, error) {
	if nodePoses == nil || len(nodePoses) == 0 {
		fmt.Printf("FindMultiTreeNode invalid params\n")
		return nil, db.ErrNotFound
	}

	tree := s.treeMap[address]
	if tree == nil {
		fmt.Printf("FindMultiTreeNode tree is null, address:%s\n", address)
		return nil, db.ErrNotFound
	}

	var retTreeNodes []*db.TreeNode
	for _, pose := range nodePoses {
		if len(tree) > pose.Level && tree[pose.Level] != nil &&
			len(tree[pose.Level]) > pose.LevelNo && tree[pose.Level][pose.LevelNo] != nil {
			retTreeNodes = append(retTreeNodes, tree[pose.Level][pose.LevelNo])
		}
	}

	if retTreeNodes == nil || len(retTreeNodes) == 0 {
		fmt.Printf("FindMultiTreeNode retTreeNodes is null, address:%s\n", address)
		return nil, db.ErrNotFound
	}

	return retTreeNodes, nil
}

func (s *MemoryStorage) FindNodesByLevel(ctx context.Context, address string, level int) ([]*db.TreeNode, error) {
	tree := s.treeMap[address]
	if tree == nil || len(tree) < (level+1) {
		return nil, db.ErrNotFound
	}

	levelMap := tree[level]

	var retsz []*db.TreeNode
	for _, node := range levelMap {
		retsz = append(retsz, node)
	}

	return retsz, nil
}
