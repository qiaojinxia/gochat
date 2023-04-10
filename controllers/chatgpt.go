package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	goGpt "github.com/sashabaranov/go-openai"
	"goChat/conf"
	"goChat/models"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"unicode/utf8"
)

var chatGptClients = make(map[string]*goGpt.Client)
var lk sync.RWMutex

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"result": c.Param("content"),
	})
}

// PostChatStreamMsgHandler 流处理处理 POST /postChatStreamMsg/:username 请求
func PostChatStreamMsgHandler(c *gin.Context) {
	// 从请求主体中读取JSON数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// 将JSON数据解码为ChatMessage结构体
	chatMsg := &models.ChatMessage{}
	err = json.Unmarshal(body, chatMsg)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	log.Printf("提问:%s\n", chatMsg.Message)

	token := c.Param("token")

	if chatMsg.APIKey == "" {
		chatMsg.APIKey = conf.GetConfig("open_ai", "key")
	}
	if chatMsg.SystemRole == "" {
		chatMsg.SystemRole = conf.GetConfig("open_ai", "default_prompt")
	}

	log.Printf("requestParam: isOpenContext(%s) systemRole(%s)  apiKey(%s) temperature(%s) \n", chatMsg.IsOpenContext, chatMsg.SystemRole, chatMsg.APIKey, chatMsg.Temperature)
	allContexts := make([]string, 0)

	if chatMsg.IsOpenContext == "1" {
		allContexts, err = models.GetUserChatContext(token)
		if err != nil {
			c.Error(err)
			return
		}
		log.Println("历史上下文:", allContexts)
	}
	allContexts = append(allContexts, chatMsg.Message)

	// 启动一个新的Goroutine来逐步发送响应内容
	c.Header("Transfer-Encoding", "chunked")
	c.Header("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)

	msgChan := make(chan string, 1024)
	stopChan := make(chan bool)
	go func() {
		err := SendChatCompletionStream(allContexts, chatMsg, msgChan, stopChan)
		if err != nil {
			c.Error(err)
			return
		}
	}()

	builder := strings.Builder{}
	var shouldExit bool
	for !shouldExit {
		select {
		case v := <-msgChan:
			builder.WriteString(v)
			c.Writer.Write([]byte(fmt.Sprintf("%x\r\n%s\r\n", utf8.RuneCountInString(v), v)))
			c.Writer.Flush() // 注意需要刷新缓冲区
		case <-stopChan:
			if len(msgChan) != 0 {
				continue
			}
			shouldExit = true // 设置退出标志
		}
	}
	log.Printf("chatGptResponse:%s", builder.String())
	log.Println("talk fished ")
	c.Writer.Write([]byte("0\r\n\r\n")) // 最后一个空分块
	if err != nil {
		log.Fatal(err)
		return
	}
	if chatMsg.IsOpenContext == "1" {
		err = models.SetUserChatContext(token, chatMsg.Message)
		if err != nil {
			log.Fatal(err)
			return
		}
		err = models.SetUserChatContext(token, builder.String())
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func SendMessage(c *gin.Context) {
	var err error
	msg, _ := c.GetPostForm("msg")

	log.Printf("question:%s\n", msg)

	token := c.Param("token")
	isOpenContext, _ := c.GetPostForm("isOpenContext")
	systemRole, _ := c.GetPostForm("systemRole")
	apiKey, _ := c.GetPostForm("apiKey")
	if apiKey == "" {
		apiKey = conf.GetConfig("open_ai", "key")
	}
	if systemRole == "" {
		systemRole = conf.GetConfig("open_ai", "default_prompt")
	}
	log.Printf("requestParam: isOpenContext(%s) systemRole(%s)  apiKey(%s)", isOpenContext, systemRole, apiKey)
	allContexts := make([]string, 0)
	if isOpenContext == "1" {
		allContexts, err = models.GetUserChatContext(token)
		if err != nil {
			c.Error(err)
			return
		}
	}
	allContexts = append(allContexts, msg)
	result, err := SendChat(allContexts, apiKey, systemRole)
	if err != nil {
		c.Error(err)
		return
	}
	err = models.SetUserChatContext(token, msg)
	if err != nil {
		c.Error(err)
		return
	}
	err = models.SetUserChatContext(token, result)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

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
