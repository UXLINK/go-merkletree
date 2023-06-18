package chache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UXUYLabs/go-merkletree/db"
	"github.com/redis/go-redis/v9"
	"regexp"
	"strconv"
)

const (
	RedisTreeKeys string = "merkletree:tree:%s:level:%d:*"
	RedisTree     string = "merkletree:tree:%s:level:%d:no:%d"
	RedisTreeNode string = "merkletree:tree:%s:node:%s"

	RedisInfoRegex string = "^merkletree:tree:(.*?):level:(.*?):no:(.*?)$"
)

type RedisStorage struct {
	db.Storage
	redisClient *redis.Client
}

func NewRedisStorage() *RedisStorage {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisStorage{
		redisClient: redisClient,
	}
}

func (s *RedisStorage) Insert(ctx context.Context, node *db.TreeNode) error {
	err := s.redisClient.Set(ctx, getRedisNodeKey(node.MtAddress, node.Data), node.ToString(), 0).Err()
	if err != nil {
		fmt.Printf("Insert Set RedisNode err. err:%+v\n", err)
		return err
	}

	err = s.redisClient.Set(ctx, getRedisTreeKey(node.MtAddress, node.Level, node.LevelNo), node.ToString(), 0).Err()
	if err != nil {
		fmt.Printf("Insert Set RedisNode err. err:%+v\n", err)
		return err
	}

	return nil
}

func (s *RedisStorage) Update(ctx context.Context, node *db.TreeNode) error {
	err := s.redisClient.Set(ctx, getRedisNodeKey(node.MtAddress, node.Data), node.ToString(), 0).Err()
	if err != nil {
		fmt.Printf("Insert Set RedisNode err. err:%+v\n", err)
		return err
	}

	err = s.redisClient.Set(ctx, getRedisTreeKey(node.MtAddress, node.Level, node.LevelNo), node.ToString(), 0).Err()
	if err != nil {
		fmt.Printf("Insert Set RedisNode err. err:%+v\n", err)
		return err
	}

	return nil
}

func (s *RedisStorage) FindRootNode(ctx context.Context, address string) (*db.TreeNode, error) {
	level := 0
	keys, err := s.redisClient.Keys(ctx, getRedisTreeKeysKey(address, level)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, db.ErrNotFound
		}
		fmt.Printf("FindRootNode get Keys err. err:%+v\n", err)
		return nil, err
	}

	if keys == nil || len(keys) == 0 {
		return nil, db.ErrNotFound
	}

	rootKeys := keys
	for keys != nil && len(keys) > 0 {
		level++
		keys, err = s.redisClient.Keys(ctx, getRedisTreeKeysKey(address, level)).Result()
		if err != nil && err != redis.Nil {
			fmt.Printf("FindRootNode get Keys err. err:%+v\n", err)
			return nil, err
		}

		if err == redis.Nil || len(keys) == 0 {
			continue
		}
		rootKeys = keys
	}

	_, level, levelNo, err := getInfoFromRedisKey(rootKeys[0])
	if err != nil {
		fmt.Printf("FindRootNode getInfoFromRedisKey err. err:%+v\n", err)
		return nil, err
	}

	val, err := s.redisClient.Get(ctx, getRedisTreeKey(address, level, levelNo)).Result()
	if err != nil {
		fmt.Printf("FindRootNode Get err. err:%+v\n", err)
		return nil, err
	}

	var node db.TreeNode
	if err = json.Unmarshal([]byte(val), &node); err != nil {
		fmt.Printf("FindRootNode Unmarshal err. err:%+v\n", err)
		return nil, err
	}

	return &node, nil
}

func (s *RedisStorage) FindMaxNoOfLeaf(ctx context.Context, address string) (int, error) {
	keys, err := s.redisClient.Keys(ctx, getRedisTreeKeysKey(address, 0)).Result()
	if err != nil {
		if err == redis.Nil {
			return -1, db.ErrNotFound
		}
		fmt.Printf("FindMaxNoOfLeaf get Keys err. err:%+v\n", err)
		return -1, err
	}

	if keys == nil || len(keys) == 0 {
		fmt.Printf("FindMaxNoOfLeaf  Keys is nil.\n")
		return -1, db.ErrNotFound
	}

	retLevelNo := 0
	for _, key := range keys {
		_, _, levelNo, err := getInfoFromRedisKey(key)
		if err != nil {
			fmt.Printf("FindMaxNoOfLeaf getInfoFromRedisKey err. err:%+v\n", err)
			return -1, err
		}

		if levelNo > retLevelNo {
			retLevelNo = levelNo
		}
	}

	return retLevelNo, nil
}

func (s *RedisStorage) FindOneByLeafData(ctx context.Context, address string, data string) (*db.TreeNode, error) {
	val, err := s.redisClient.Get(ctx, getRedisNodeKey(address, data)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, db.ErrNotFound
		}
		fmt.Printf("FindOneByLeafData Get keys err. err:%+v\n", err)
		return nil, err
	}

	var node db.TreeNode
	if err = json.Unmarshal([]byte(val), &node); err != nil {
		fmt.Printf("FindOneByLeafData Unmarshal err. err:%+v\n", err)
		return nil, err
	}

	return &node, nil
}

func (s *RedisStorage) FindMultiTreeNode(ctx context.Context, address string, nodePoses []*db.NodePos) ([]*db.TreeNode, error) {
	if nodePoses == nil || len(nodePoses) == 0 {
		fmt.Printf("FindMultiTreeNode invalid params\n")
		return nil, db.ErrNotFound
	}

	var retTreeNodes []*db.TreeNode
	for _, pose := range nodePoses {
		val, err := s.redisClient.Get(ctx, getRedisTreeKey(address, pose.Level, pose.LevelNo)).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			fmt.Printf("FindMultiTreeNode Get err. err:%+v\n", err)
			return nil, err
		}

		var node db.TreeNode
		if err = json.Unmarshal([]byte(val), &node); err != nil {
			fmt.Printf("FindMultiTreeNode Unmarshal err. err:%+v\n", err)
			return nil, err
		}
		retTreeNodes = append(retTreeNodes, &node)
	}

	if retTreeNodes == nil || len(retTreeNodes) == 0 {
		fmt.Printf("FindMultiTreeNode retTreeNodes is null, address:%s\n", address)
		return nil, db.ErrNotFound
	}

	return retTreeNodes, nil
}

func (s *RedisStorage) FindNodesByLevel(ctx context.Context, address string, level int) ([]*db.TreeNode, error) {
	keys, err := s.redisClient.Keys(ctx, getRedisTreeKeysKey(address, level)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, db.ErrNotFound
		}
		fmt.Printf("FindNodesByLevel Get keys err. err:%+v\n", err)
		return nil, err
	}

	if keys == nil || len(keys) == 0 {
		fmt.Printf("FindNodesByLevel keys is nil\n")
		return nil, db.ErrNotFound
	}

	var retsz []*db.TreeNode
	for _, key := range keys {
		val, err := s.redisClient.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			fmt.Printf("FindNodesByLevel Get err. err:%+v\n", err)
			return nil, err
		}

		var node db.TreeNode
		if err = json.Unmarshal([]byte(val), &node); err != nil {
			fmt.Printf("FindNodesByLevel Unmarshal err. err:%+v\n", err)
			return nil, err
		}

		retsz = append(retsz, &node)
	}

	return retsz, nil
}

func getRedisNodeKey(address string, data string) string {
	return fmt.Sprintf(RedisTreeNode, address, data)
}

func getRedisTreeKey(address string, level, levelNo int) string {
	return fmt.Sprintf(RedisTree, address, level, levelNo)
}

func getRedisTreeKeysKey(address string, level int) string {
	return fmt.Sprintf(RedisTreeKeys, address, level)
}

func getInfoFromRedisKey(key string) (string, int, int, error) {

	compileRegex := regexp.MustCompile(RedisInfoRegex) // 正则表达式的分组，以括号()表示，每一对括号就是我们匹配到的一个文本，可以把他们提取出来。
	matchArr := compileRegex.FindStringSubmatch(key)   // FindStringSubmatch 方法是提取出匹配的字符串，然后通过[]string返回。我们可以看到，第1个匹配到的是这个字符串本身，从第2个开始，才是我们想要的字符串。
	if matchArr == nil || len(matchArr) < 4 {
		return "", 0, 0, errors.New("key resolve err")
	}

	level, err := strconv.Atoi(matchArr[2])
	if err != nil {
		return "", 0, 0, err
	}

	levelNo, err := strconv.Atoi(matchArr[3])
	if err != nil {
		return "", 0, 0, err
	}

	return matchArr[1], level, levelNo, nil
}
