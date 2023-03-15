package models

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"goChat/repositories"

	"sort"
	"strconv"
)

func GetChatGptUserKey(token string) string {
	return fmt.Sprintf("%s:%s", "ChatGptUserKey", token)
}
func GetChatGptUserCountKey(token string) string {
	return fmt.Sprintf("%s:%s", "ChatGptUserCountKey", token)
}

func SetUserChatContext(token string, context string) error {
	pool := repositories.GetRedisPool()
	conn := pool.Get()
	ct, err := redis.Int(conn.Do("INCR", GetChatGptUserCountKey(token)))
	if err != nil {
		return err
	}
	_, err = conn.Do("HMSET", GetChatGptUserKey(token), strconv.Itoa(ct), context)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", GetChatGptUserCountKey(token), 60*60*60)
	return err
}

func GetUserChatContext(token string) ([]string, error) {
	pool := repositories.GetRedisPool()
	conn := pool.Get()
	data, err := redis.StringMap(conn.Do("HGETALL", GetChatGptUserKey(token)))
	if err != nil {
		return nil, err
	}
	// 获取所有键
	keys := make([]int, len(data))
	i := 0
	for k := range data {
		keys[i], _ = strconv.Atoi(k)
		i++
	}
	// 对键进行排序
	sort.Ints(keys)
	// 遍历已排序的键并将它们添加到输出数组
	sortedKeys := make([]string, len(data))
	for i, k := range keys {
		sortedKeys[i] = data[strconv.Itoa(k)]
	}
	return sortedKeys, nil
}
