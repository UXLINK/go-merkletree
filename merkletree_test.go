package merkletree

import (
	"context"
	"fmt"
	"github.com/UXUYLabs/go-merkletree/db/chache"
	"github.com/UXUYLabs/go-merkletree/keccak256"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"regexp"
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

func setupRedis() {
	ctx := context.Background()
	merkleTreeManager, err = NewMerkleTreeManager(ctx, chache.NewRedisStorage())
	if err != nil {
		fmt.Printf("NewMerkleTree err:%v\n", err)
		return
	}
}

func TestMemAppend(t *testing.T) {
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

	_ = tree.PrintTree()

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

	_ = tree.PrintTree()

	err = tree.AppendLeaf("0x7e533CF779A533eD8f9C1b8E5C3d7F936335ca54")
	if err != nil {
		fmt.Printf("AppendLeaf err :%v\n", err)
		return
	}

	_ = tree.PrintTree()

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

func TestMemAppend2(t *testing.T) {
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

	err = tree.AppendLeaf("0x7e533CF779A533eD8f9C1b8E5C3d7F936335ca54")
	if err != nil {
		fmt.Printf("AppendLeaf err :%v\n", err)
		return
	}

	_ = tree.PrintTree()

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

func TestRedisAppend1(t *testing.T) {
	setupRedis()

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

	err = tree.AppendLeaf("0x7e533CF779A533eD8f9C1b8E5C3d7F936335ca54")
	if err != nil {
		fmt.Printf("AppendLeaf err :%v\n", err)
		return
	}

	err = tree.AppendLeaf("0x63120cc1c7Bb0a42C2D77D27faB9EDd2560F9cA3")
	if err != nil {
		fmt.Printf("AppendLeaf err :%v\n", err)
		return
	}

	_ = tree.PrintTree()

	// proof
	proofes, err := tree.GenerateProof("0x63120cc1c7Bb0a42C2D77D27faB9EDd2560F9cA3")
	if err != nil {
		fmt.Printf("NewMerkleTree err:%v\n", err)
		return
	}

	for _, proof := range proofes {
		fmt.Printf("NewMerkleTree proofes:%+v\n", keccak256.Bytes2Hex(proof))
	}

	proof, err := tree.VerifyProof(proofes, "0x63120cc1c7Bb0a42C2D77D27faB9EDd2560F9cA3")
	//logx.WithContext(ctx).Infof("NewMerkleTree proof:%+v", proof)
	assert.Equal(t, proof, true)

	proof, err = tree.VerifyProof(proofes, "0xa6820eeA9B5BB08Ab1cD693128Bb85Ad460a8e6E")
	assert.Equal(t, proof, false)

	// proof
	proofes, err = tree.GenerateProof("0x7e533CF779A533eD8f9C1b8E5C3d7F936335ca54")
	if err != nil {
		fmt.Printf("NewMerkleTree err:%v\n", err)
		return
	}

	for _, proof := range proofes {
		fmt.Printf("NewMerkleTree proofes:%+v\n", keccak256.Bytes2Hex(proof))
	}

	proof, err = tree.VerifyProof(proofes, "0x7e533CF779A533eD8f9C1b8E5C3d7F936335ca54")
	//logx.WithContext(ctx).Infof("NewMerkleTree proof:%+v", proof)
	assert.Equal(t, proof, true)
}

func TestRedisAppend2(t *testing.T) {
	setupRedis()

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

	err = tree.AppendLeaf("0x7e533CF779A533eD8f9C1b8E5C3d7F936335ca54")
	if err != nil {
		fmt.Printf("AppendLeaf err :%v\n", err)
		return
	}

	_ = tree.PrintTree()

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

func TestRedis(t *testing.T) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

func TestStringCmp(t *testing.T) {

	str := "merkletree:tree:ddf:level:10:no:11"
	compileRegex := regexp.MustCompile("^merkletree:tree:ddf:level:(.*?):no:(.*?)$") // 正则表达式的分组，以括号()表示，每一对括号就是我们匹配到的一个文本，可以把他们提取出来。
	matchArr := compileRegex.FindStringSubmatch(str)                                 // FindStringSubmatch 方法是提取出匹配的字符串，然后通过[]string返回。我们可以看到，第1个匹配到的是这个字符串本身，从第2个开始，才是我们想要的字符串。
	fmt.Println("提取字符串内容：", matchArr[1], matchArr[2])                                // 输出：蜜桃乌龙茶

}
