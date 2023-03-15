package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"goChat/controllers"
	"goChat/repositories"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	//rateLimiter := &middleware.RateLimiter{
	//	MaxRequests: 10000,
	//	WindowSize:  time.Minute,
	//	IPMap: make(map[string]int64),
	//}
	//r.Use(rateLimiter.LimitHandler)
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static/")
	// 添加 POST /postChatMsg/:username 路由
	r.POST("/postChatStreamMsg/:token", controllers.PostChatStreamMsgHandler)
	r.POST("/postChatMsg/:token", controllers.SendMessage)
	r.GET("/", controllers.Index)
	// 写一个服务监听
	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    5 * 60 * time.Second,
		WriteTimeout:   5 * 60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// 启动服务
	go server.ListenAndServe()
	// 优雅退出
	gracefulExitServer(server)
}

func gracefulExitServer(server *http.Server) {
	// 使用缓存的channel；建议用1；详情看Uber Go style；其他情况酌情使用缓冲大小
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	// 接收信号
	sig := <-ch
	log.Println("获取一个系统信号", sig)
	//关闭所有服务
	repositories.CloseAllServer()
	// 设置当前时间
	nowTime := time.Now()
	// 设置为5秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 最后关闭
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Println("err", err)
	}
	log.Println("-----exited-----", time.Since(nowTime))
}
