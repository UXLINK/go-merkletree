# Quickstart

## memory storage 

```go
import (
    "context"
    "fmt"
    "github.com/UXUYLabs/go-merkletree/keccak256"
)

func main() {
    ctx := context.Background()
    merkleTreeManager, err = NewMemoryMerkleTreeManager(ctx)
    if err != nil {
        fmt.Printf("NewMerkleTree err:%v\n", err)   
		return
    }

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

	_ = tree.printTree()

    // proof
    proofes, err := tree.GenerateProof("0xeA726629EC5fe5cE300000d1a8c89B3054A22cE7")
    if err != nil {
        fmt.Printf("NewMerkleTree err:%v\n", err)
        return
    }

	for _, proof := range proofes {
		fmt.Printf("NewMerkleTree proofes:%+v\n", keccak256.Bytes2Hex(proof))
	}

	proof, err := tree.VerifyProof(proofes, "0xeA726629EC5fe5cE300000d1a8c89B3054A22cE7")
    fmt.Printf("NewMerkleTree proof result:%+v\n", proof)
}
```

## redis storage
```go
import (
    "context"
    "fmt"
    "github.com/UXUYLabs/go-merkletree/db/chache"
    "github.com/UXUYLabs/go-merkletree/keccak256"
)

func main() {
    ctx := context.Background()
	merkleTreeManager, err = NewMerkleTreeManager(ctx, chache.NewRedisStorage())
	if err != nil {
        fmt.Printf("NewMerkleTree err:%v\n", err)
        return
	}

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

	_ = tree.printTree()

    // proof
    proofes, err := tree.GenerateProof("0xeA726629EC5fe5cE300000d1a8c89B3054A22cE7")
    if err != nil {
        fmt.Printf("NewMerkleTree err:%v\n", err)
        return
    }

	for _, proof := range proofes {
		fmt.Printf("NewMerkleTree proofes:%+v\n", keccak256.Bytes2Hex(proof))
	}

	proof, err := tree.VerifyProof(proofes, "0xeA726629EC5fe5cE300000d1a8c89B3054A22cE7")
    fmt.Printf("NewMerkleTree proof result:%+v\n", proof)
}
```
