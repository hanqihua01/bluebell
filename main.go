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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// gin: Web框架
// viper: 配置文件读取
// zap: 日志生成
// lumberjack: 日志切割
// sqlx: MySQL连接
// go-redis: Redis连接
// snowflake: 分布式ID生成
// validator: 请求参数校验
// jwt-go: JWT用户认证
// makefile: 快速构建项目
// air: 文件热重载
func main() {
	// 加载配置
	var confFile string
	flag.StringVar(&confFile, "c", "./config.yaml", "config file path")
	flag.Parse() // 通过命令行参数执行配置文件
	if err := settings.Init(confFile); err != nil {
		fmt.Printf("config file initialization failed, err: %v\n", err)
		return
	}

	// 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("logger initialization failed, err: %v\n", err)
		return
	}
	defer zap.L().Sync() // 确保所有日志写入日志文件中

	// 初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("mysql initialization failed, err:%v\n", err)
		return
	}
	defer mysql.Close()

	// 初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("redis initialization failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	// 初始化雪花算法
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("snowflake initialization failed, err:%v\n", err)
	}

	// 注册路由
	r := router.Setup(settings.Conf.Mode)

	// 启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}
	go func() { // 启动goroutine运行服务器
		zap.L().Info("server is launching")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("server is wrong", zap.Error(err))
		}
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)                      // 设置channel来存放中断信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 设置如果有信号，则放到quit通道里
	<-quit                                               // 没有信号前一直阻塞，收到信号后停止阻塞
	zap.L().Info("server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 创建一个5秒超时的context
	defer cancel()                                                          // 执行完毕或者超时后释放上下文资源
	if err := srv.Shutdown(ctx); err != nil {                               // 在上下文超时时限内进行Shutdown()操作
		zap.L().Fatal("server shutting down is wrong", zap.Error(err))
	}
}
