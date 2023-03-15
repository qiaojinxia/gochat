package utils

import (
	"context"
	"errors"
	"fmt"
	goGpt "github.com/sashabaranov/go-openai"
	"goChat/models"
	"io"
	"log"
	"sync"
)

var chatGptClients = make(map[string]*goGpt.Client)
var lk sync.RWMutex

func GetClient(key string) *goGpt.Client {
	lk.RLock()
	if client, exists := chatGptClients[key]; exists {
		lk.RUnlock()
		return client
	}
	lk.RUnlock()
	// 你的api密钥
	chatGptClient := goGpt.NewClient(key)
	lk.Lock()
	defer lk.Unlock()
	chatGptClients[key] = chatGptClient
	return chatGptClient
}

func SendChat(allContexts []string, key, systemRole string) (string, error) {
	req := goGpt.ChatCompletionRequest{
		Model:       goGpt.GPT3Dot5Turbo0301, // 选择的模型
		MaxTokens:   2048,
		N:           1,
		Stop:        nil,
		Temperature: 0.5,
	}
	messages := []goGpt.ChatCompletionMessage{
		{Role: "system", Content: systemRole},
	}
	for _, v := range allContexts {
		messages = append(messages, goGpt.ChatCompletionMessage{Role: "user", Content: v})
	}
	req.Messages = messages
	ctx := context.Background()
	c := GetClient(key)
	resp, err := c.CreateChatCompletion(ctx, req) // 发起接口调用
	if err != nil {
		return "", err
	}
	log.Println("Created: 	 ", resp.Created)
	log.Println("Model:	 ", resp.Model)
	log.Println("Message:	 ", resp.Choices[0].Message.Content)
	log.Println("Usage: 	 ", resp.Usage)
	return resp.Choices[0].Message.Content, nil
}

func SendChatCompletionStream(allContexts []string, msgObj *models.ChatMessage,
	buffer chan string, stopChan chan bool) error {
	defer func() {
		stopChan <- true
	}()

	req := goGpt.ChatCompletionRequest{
		Model:       goGpt.GPT3Dot5Turbo0301, // 选择的模型
		MaxTokens:   2048,
		Stream:      true,
		Temperature: msgObj.GetTemperature(),
	}
	messages := []goGpt.ChatCompletionMessage{
		{Role: "system", Content: msgObj.SystemRole},
	}
	for _, v := range allContexts {
		messages = append(messages, goGpt.ChatCompletionMessage{Role: "user", Content: v})
	}
	req.Messages = messages
	ctx := context.Background()
	c := GetClient(msgObj.APIKey)
	stream, err := c.CreateChatCompletionStream(ctx, req) // 发起接口调用

	if err != nil {
		fmt.Printf("CompletionStream error: %v\n", err)
		return err
	}
	defer stream.Close()
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Println("Stream finished")
			break
		}
		if err != nil {
			log.Printf("Stream error: %v\n", err)
			return err
		}
		buffer <- response.Choices[0].Delta.Content
	}
	return nil
}
