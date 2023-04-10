package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"goChat/conf"
	"goChat/controllers"
	"goChat/middleware"
	"goChat/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conf.InitGlobalConfig("config.ini")
	server.InitServers()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	//rateLimiter := &middleware.RateLimiter{
	//	MaxRequests: 10000,
	//	WindowSize:  time.Minute,
	//	IPMap: make(map[string]int64),
	//}
	//r.Use(rateLimiter.LimitHandler)
	r.Use(middleware.CorsMiddleware())
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static/")
	r.Static("/data", "data/")
	// 添加 POST /postChatMsg/:username 路由
	r.POST("/postChatStreamMsg/:token", controllers.PostChatStreamMsgHandler)
	r.POST("/postChatMsg/:token", controllers.SendMessage)
	r.GET("/", controllers.Index)
	r.GET("/downloadExcel/:name", controllers.GetExcel)
	r.POST("/uploadExcel", controllers.SaveExcel)
	r.GET("/chatExcel", controllers.Excel)
	r.POST("/chatOptExcel", controllers.ChatExcelGenerate)
	r.GET("/chatExcel/about", controllers.ChatExcelAbout)
	r.GET("/chatExcel/tutorial", controllers.Tutorial)
	// 写一个服务监听
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", conf.GetConfig("gin", "port")),
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

func gracefulExitServer(h *http.Server) {
	// 使用缓存的channel；建议用1；详情看Uber Go style；其他情况酌情使用缓冲大小
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	// 接收信号
	sig := <-ch
	log.Println("get exit single", sig)
	//关闭所有服务
	server.CloseAllServer()
	// 设置当前时间
	nowTime := time.Now()
	// 设置为5秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 最后关闭
	defer cancel()
	err := h.Shutdown(ctx)
	if err != nil {
		log.Println("err", err)
	}
	log.Println("-----exited-----", time.Since(nowTime))
}
