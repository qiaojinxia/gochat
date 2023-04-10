package utils

import (
	"github.com/segmentio/ksuid"
)

func GenerateID() (string, error) {
	// 生成一个新的KSUID
	id, err := ksuid.NewRandom()
	if err != nil {
		// 处理错误
		return "", err
	}
	return id.String(), nil
}
