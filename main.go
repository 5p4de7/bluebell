package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/router"
	"bluebell/settings"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click

// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	//*1.加载配置
	err := settings.Init()
	if err != nil {
		fmt.Printf("配置加载错误%v\n", err)
		return
	}
	//*2.加载日志
	err = logger.Init(settings.Conf.LogConf)
	if err != nil {
		fmt.Printf("日志加载错误%v\n", err)
	}
	defer zap.L().Sync() //*写入日志
	zap.L().Debug("日志加载成功")
	//*3.连接mysql
	if err := mysql.Init(settings.Conf.MysqlConf); err != nil {
		fmt.Printf("mysql数据库连接错误%v\n", err)
	}
	defer mysql.Close()
	//*4.连接redis
	if err := redis.Init(settings.Conf.RedisConf); err != nil {
		fmt.Printf("redis<UNK>%v\n", err)
	}
	defer redis.Close()
	//*5.注册路由
	r := router.SetUp()
	//*6.启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf("%d", &settings.Conf.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Info("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
