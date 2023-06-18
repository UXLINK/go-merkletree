package merkletree

import (
	"context"
	"errors"
	"fmt"
	"github.com/UXUYLabs/go-merkletree/db"
	"github.com/UXUYLabs/go-merkletree/keccak256"
	"github.com/shopspring/decimal"
	"log"
	"os"
	"regexp"
)

// MerkleTree is the structure for the Merkle tree.
type MerkleTree struct {
	Logger
	ctx       context.Context
	mtAddress string
	storage   db.Storage
}

func NewMerkleTree(ctx context.Context, storage db.Storage, mtAddress string) (*MerkleTree, error) {
	return &MerkleTree{
		Logger:    PrintfLogger(log.New(os.Stdout, "merkleTree: ", log.LstdFlags)),
		ctx:       ctx,
		mtAddress: mtAddress,
		storage:   storage,
	}, nil
}

func (t *MerkleTree) AppendLeaf(data string) error {
	if !IsAddress(data) {
		return errors.New("data address invalid.")
	}

	// 1. 查询是否已有，直接返回
	leaf, err := t.getLeafNodeByData(data)
	if err != nil {
		t.Error("AppendLeaf getLeafNodeByData err: ", err)
		return err
	}

	if leaf != nil {
		return nil
	}

	// 2. 生成叶子
	leaf = &db.TreeNode{
		MtAddress: t.mtAddress,
		Data:      data,
		Hash:      keccak256.Bytes2Hex(keccak256.HashLeaf(data[2:])),
		Level:     0,
		LevelNo:   0,
	}
	t.Info("AppendLeaf leaf: ", leaf)

	// 3. 找到对应的branch
	branches, leaf, err := t.doNewTreeBranches(leaf)

	for _, branch := range branches {
		t.Info("AppendLeaf branche: ", branch)
	}
	// 4. 关联到branch，并修改branch的hash值
	if len(branches) == 0 {
		return nil
	}

	root, err := t.storage.FindRootNode(t.ctx, t.mtAddress)
	if err != nil && err != db.ErrNotFound {
		t.Error("AppendLeaf FindRootNode err: ", err)
		return err
	}

	hash := ""
	levelNo := leaf.LevelNo
	i := 0
	needCreateRoot := true
	for ; i < len(branches); i++ {
		branch := branches[i]
		if branch == nil {
			branch = make(map[int]*db.TreeNode)
		}
		t.Info("AppendLeaf branch: ", i, branch)
		// 计算叶子节点hash
		if i == 0 {
			if !isEven(levelNo) && branch[levelNo-1] != nil {
				hash = keccak256.Bytes2Hex(keccak256.HashBranch(branch[levelNo-1].Hash, branch[levelNo].Hash))
			} else {
				hash = branch[levelNo].Hash
			}
		} else {
			// 父节点hash，如果没有，新创建
			if branch[levelNo] == nil {
				branchNode := &db.TreeNode{
					MtAddress: t.mtAddress,
					Hash:      hash,
					Level:     i,
					LevelNo:   levelNo,
				}
				err = t.storage.Insert(t.ctx, branchNode)
				if err != nil {
					t.Error("AppendLeaf Insert err: ", err)
					return err
				}
				branch[levelNo] = branchNode
			} else {
				branch[levelNo].Hash = hash
				err = t.storage.Update(t.ctx, branch[levelNo])
				if err != nil {
					t.Error("AppendLeaf Update err: ", err)
					return err
				}
			}

			if !isEven(levelNo) {
				if branch[levelNo-1] != nil {
					hash = keccak256.Bytes2Hex(keccak256.HashBranch(branch[levelNo-1].Hash, branch[levelNo].Hash))
				}
			} else {
				if branch[levelNo+1] != nil {
					hash = keccak256.Bytes2Hex(keccak256.HashBranch(branch[levelNo].Hash, branch[levelNo+1].Hash))
				}

				if root.Level == i && levelNo == 0 && branch[levelNo+1] == nil {
					needCreateRoot = false
				}
			}
		}

		levelNoDecimal := decimal.NewFromInt(int64(levelNo)).Div(decimal.NewFromInt32(2)).RoundFloor(0)
		levelNo = int(levelNoDecimal.IntPart())
	}

	// 新建根节点
	if needCreateRoot {
		rootNode := &db.TreeNode{
			MtAddress: t.mtAddress,
			Hash:      hash,
			Level:     i,
			LevelNo:   levelNo,
		}
		err = t.storage.Insert(t.ctx, rootNode)
		if err != nil {
			t.Error("AppendLeaf Insert err: ", err)
			return err
		}
	}

	return nil
}

func (t *MerkleTree) getLeafNodeByData(data string) (*db.TreeNode, error) {
	leaf, err := t.storage.FindOneByLeafData(t.ctx, t.mtAddress, data)
	if err != nil && err != db.ErrNotFound {
		t.Error("getLeafNodeByData FindOneByLeafData err: ", err)
		return nil, err
	}

	return leaf, nil
}

func (t *MerkleTree) GetRootNode() (*db.TreeNode, error) {
	rootNode, err := t.storage.FindRootNode(t.ctx, t.mtAddress)
	if err != nil && err != db.ErrNotFound {
		t.Error("GetRootNode FindRootNode err: ", err)
		return nil, err
	}

	return rootNode, nil
}

func (t *MerkleTree) GenerateProof(data string) ([][]byte, error) {
	leaf, err := t.getLeafNodeByData(data)
	if err != nil && err != db.ErrNotFound {
		t.Error("GenerateProof FindOneByLeafData err: ", err)
		return nil, err
	}

	if leaf == nil {
		return make([][]byte, 0), nil
	}

	referTree, err := t.getReferTreeByLeaf(leaf)
	if err != nil {
		t.Error("GenerateProof getReferTreeByLeaf err: ", err)
		return nil, err
	}

	var retSz [][]byte
	maxlevelNoDecimal := decimal.NewFromInt(int64(leaf.LevelNo))
	for i := 0; i < len(referTree); i++ {
		branch := referTree[i]
		levelNo := int(maxlevelNoDecimal.IntPart())
		if !isEven(levelNo) {
			if branch[levelNo-1] != nil {
				retSz = append(retSz, keccak256.Hex2Bytes(branch[levelNo-1].Hash))
			}
		} else {
			if branch[levelNo+1] != nil {
				retSz = append(retSz, keccak256.Hex2Bytes(branch[levelNo+1].Hash))
			}
		}

		maxlevelNoDecimal = maxlevelNoDecimal.Div(decimal.NewFromInt32(2)).RoundFloor(0)

	}

	return retSz, nil
}

func (t *MerkleTree) VerifyProof(proofs [][]byte, user string) (bool, error) {
	if !IsAddress(user) {
		return false, nil
	}
	hash := keccak256.HashLeaf(user[2:])
	if len(proofs) > 0 {
		for _, proof := range proofs {
			t.Info("VerifyProof Hash: ", keccak256.Bytes2Hex(hash), keccak256.Bytes2Hex(proof))
			hash = keccak256.HashByteBranch(hash, proof)
		}
	}

	// 对比根节点
	node, err := t.storage.FindRootNode(t.ctx, t.mtAddress)
	if err != nil && err != db.ErrNotFound {
		t.Error("VerifyProof FindOneByLeafData err: ", err)
		return false, err
	}

	t.Info("VerifyProof hash: ", node.Hash, keccak256.Bytes2Hex(hash))
	if node.Hash == keccak256.Bytes2Hex(hash) {
		return true, nil
	}

	return false, nil
}

// 返回数据中map为每层的对应相关数据，数组为层级
func (t *MerkleTree) doNewTreeBranches(leaf *db.TreeNode) ([]map[int]*db.TreeNode, *db.TreeNode, error) {
	maxLevelNo, err := t.storage.FindMaxNoOfLeaf(t.ctx, t.mtAddress)
	if err != nil && err != db.ErrNotFound {
		t.Error("doNewTreeBranches FindMaxNoOfLeaf err: ", err)
		return nil, nil, err
	}

	// 创建新叶子
	leaf.LevelNo = maxLevelNo + 1
	err = t.storage.Insert(t.ctx, leaf)
	if err != nil {
		t.Error("doNewTreeBranches Insert err: ", err)
		return nil, nil, err
	}

	// 树为空的时候
	retSz, err := t.getReferTreeByLeaf(leaf)
	if err != nil {
		t.Error("doNewTreeBranches getReferTreeByLeaf err: ", err)
		return nil, nil, err
	}

	return retSz, leaf, nil
}

func (t *MerkleTree) getReferTreeByLeaf(leaf *db.TreeNode) ([]map[int]*db.TreeNode, error) {

	root, err := t.storage.FindRootNode(t.ctx, t.mtAddress)
	if err != nil && err != db.ErrNotFound {
		t.Error("getReferTreeByLeaf FindOneByLeafData err: ", err)
		return nil, err
	}

	nodePoses := calcNodeBranches(leaf.LevelNo, root.Level)
	if len(nodePoses) == 0 {
		levelMap := make(map[int]*db.TreeNode, 1)
		levelMap[leaf.LevelNo] = leaf

		retSz := make([]map[int]*db.TreeNode, 1)
		retSz[0] = levelMap
		return retSz, nil
	}

	/*for _, pose := range nodePoses {
		t.Info("getReferTreeByLeaf info  nodePose: ", pose)
	}*/

	//t.Info("getReferTreeByLeaf info  leaf: ", leaf, nodePoses)
	treeNodes, err := t.storage.FindMultiTreeNode(t.ctx, t.mtAddress, nodePoses)
	if err != nil {
		t.Error("getReferTreeByLeaf FindMultiTreeNode err: ", err)
		return nil, err
	}

	/*for _, node := range treeNodes {
		t.Info("getReferTreeByLeaf treeNodes info  node: ", node)
	}*/

	//t.Info("getReferTreeByLeaf info  leaf: ", leaf, nodePoses, treeNodes)
	retSz := make([]map[int]*db.TreeNode, treeNodes[len(treeNodes)-1].Level+1)
	for _, node := range treeNodes {
		levelNodes := retSz[node.Level]
		if levelNodes == nil {
			levelNodes = make(map[int]*db.TreeNode)
			retSz[node.Level] = levelNodes
		}

		levelNodes[node.LevelNo] = node
	}
	return retSz, nil
}

func (t *MerkleTree) PrintTree() error {
	root, err := t.storage.FindRootNode(t.ctx, t.mtAddress)
	if err != nil && err != db.ErrNotFound {
		t.Error("PrintTree FindOneByLeafData err: ", err)
		return err
	}

	level := root.Level
	t.Error("PrintTree root: ", root)
	for i := level; i >= 0; i-- {
		nodes, err := t.storage.FindNodesByLevel(t.ctx, t.mtAddress, i)
		if err != nil {
			t.Error("PrintTree FindOneByLeafData err: ", err)
			return err
		}

		for _, node := range nodes {
			t.Error(fmt.Sprintf("PrintTree level=%d, levelNo=%d", node.Level, node.LevelNo), node)
		}
	}

	return nil
}

func calcNodeBranches(maxLevelNo, rootLevel int) []*db.NodePos {
	// 还要包涵兄弟节点，用于计算hash
	var retSz []*db.NodePos

	retSz = append(retSz, &db.NodePos{Level: 0, LevelNo: maxLevelNo})
	if !isEven(maxLevelNo) {
		retSz = append(retSz, &db.NodePos{Level: 0, LevelNo: maxLevelNo - 1})
	} else {
		retSz = append(retSz, &db.NodePos{Level: 0, LevelNo: maxLevelNo + 1})
	}

	maxlevelNoDecimal := decimal.NewFromInt(int64(maxLevelNo))
	level := 1
	for maxlevelNoDecimal.IsPositive() || maxlevelNoDecimal.IsZero() {
		// 如果为1，则结束
		if level > rootLevel {
			break
		}
		maxlevelNoDecimal = maxlevelNoDecimal.Div(decimal.NewFromInt32(2)).RoundFloor(0)
		retSz = append(retSz, &db.NodePos{Level: level, LevelNo: int(maxlevelNoDecimal.IntPart())})

		// 每个层级的兄弟节点
		if !isEven(int(maxlevelNoDecimal.IntPart())) {
			retSz = append(retSz, &db.NodePos{Level: level, LevelNo: int(maxlevelNoDecimal.IntPart()) - 1})
		} else {
			retSz = append(retSz, &db.NodePos{Level: level, LevelNo: int(maxlevelNoDecimal.IntPart()) + 1})
		}

		level++
	}

	return retSz
}

// 判断是偶数
func isEven(num int) bool {
	if num%2 == 0 {
		return true
	}
	return false
}

func IsAddress(userAddress string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(userAddress)
}
