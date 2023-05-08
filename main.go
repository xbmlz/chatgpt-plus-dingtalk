package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xbmlz/chatgpt-dingtalk/internal/config"
	"github.com/xbmlz/chatgpt-dingtalk/internal/db"
	"github.com/xbmlz/chatgpt-dingtalk/internal/handlers"
	"github.com/xbmlz/chatgpt-dingtalk/pkg/logger"
)

// dingtalk bot

// chatgpt

func main() {
	config.Initialize()
	logger.Initialize(config.Instance.LogLevel)
	db.Initialize()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"msg": "pong"}) })
	r.POST("/", handlers.RootHandler)
	port := fmt.Sprintf(":%d", config.Instance.ServerPort)
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		logger.Info("🚀 Listening and serving HTTP on", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server ...")

	// 5秒后强制退出
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown: %s", err)
	}
	logger.Info("Server exiting!")
}
