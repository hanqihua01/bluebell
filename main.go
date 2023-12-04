package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/router"
	"bluebell/util/logger"
	"bluebell/util/settings"
	"bluebell/util/snowflake"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 加载配置
	var confFile string
	flag.StringVar(&confFile, "c", "./config.yaml", "config file path")
	flag.Parse() // 通过命令行参数执行配置文件
	if err := settings.Init(confFile); err != nil {
		fmt.Printf("init settings failed, err: %v\n", err)
		return
	}

	// 初始化日志
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}
	defer zap.L().Sync() // 确保所有日志写入日志文件中

	// 初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()

	// 初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	// 初始化雪花算法
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
	}

	// 注册路由
	r := router.Setup(settings.Conf.Mode)

	// 启动服务（优雅关机）
	// 配置服务器信息
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}
	// 启动goroutine运行服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 等待中断信号来优雅关闭服务器，为关闭服务器设置一个5秒的超时
	quit := make(chan os.Signal, 1)                      // 设置channel来存放中断信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 设置如果有信号，则放到quit通道里
	<-quit                                               // 没有信号前一直阻塞，收到信号后停止阻塞
	zap.L().Info("shutdown server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 创建一个5秒超时的context
	defer cancel()                                                          // 执行完毕或者超时后释放上下文资源
	if err := srv.Shutdown(ctx); err != nil {                               // 在上下文超时时限内进行Shutdown()操作
		zap.L().Fatal("server shutdown", zap.Error(err))
	}
	zap.L().Info("server exiting")
}
