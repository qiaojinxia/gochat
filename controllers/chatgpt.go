package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"goChat/conf"
	"goChat/models"
	"goChat/utils"
	"io"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

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
		chatMsg.APIKey = conf.DefaultOpenAIKey
	}
	if chatMsg.SystemRole == "" {
		chatMsg.SystemRole = conf.DefaultRole
	}

	log.Printf("requestParam: isOpenContext(%s) systemRole(%s)  apiKey(%s) temperature(%s)", chatMsg.IsOpenContext, chatMsg.SystemRole, chatMsg.APIKey, chatMsg.Temperature)
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
		err := utils.SendChatCompletionStream(allContexts, chatMsg, msgChan, stopChan)
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
			log.Printf("%s", v)
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
	log.Println()
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
	log.Printf("提问:%s\n", msg)
	token := c.Param("token")
	isOpenContext, _ := c.GetPostForm("isOpenContext")
	systemRole, _ := c.GetPostForm("systemRole")
	apiKey, _ := c.GetPostForm("apiKey")
	if apiKey == "" {
		apiKey = conf.DefaultOpenAIKey
	}
	if systemRole == "" {
		systemRole = conf.DefaultRole
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
	result, err := utils.SendChat(allContexts, apiKey, systemRole)
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
