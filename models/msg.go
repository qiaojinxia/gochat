package models

import (
	"log"
	"strconv"
)

// 定义聊天消息结构体
type ChatMessage struct {
	Message       string `json:"msg"`
	IsOpenContext string `json:"isOpenContext"`
	SystemRole    string `json:"systemRole"`
	APIKey        string `json:"apiKey"`
	Temperature   string `json:"temperature"`
}

func (i *ChatMessage) GetTemperature() float32 {
	if i.Temperature == "" {
		return 0.0
	}
	floatValue, err := strconv.ParseFloat(i.Temperature, 32) // 将string类型的值转换为float32类型
	if err != nil {
		log.Fatal(err)
	}
	return float32(floatValue)

}
