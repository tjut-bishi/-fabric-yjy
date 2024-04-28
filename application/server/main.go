package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/oklog/run"

	"application/blockchain"
	"application/config"
	"application/pkg/cron"
	orm "application/pkg/mysql"
	"application/routers"

	"gorm.io/gorm"
)

func main() {
	// 接收系统信号, 通过context关系退出服务
	// ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	// defer stop()
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("时区设置失败 %s", err)
	}
	time.Local = timeLocal

	blockchain.Init()
	go cron.Init()

	endPoint := fmt.Sprintf("0.0.0.0:%d", 8888)
	server := &http.Server{
		Addr:              endPoint,
		Handler:           routers.InitRouter(),
		ReadHeaderTimeout: 5 * time.Second, // 5秒超时时间
	}
	// 初始化 config 配置
	conf, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	//
	orm, err := orm.NewOrm(&orm.MysqlConfig{
		Host:     conf.Mysql.Host,
		Username: conf.Mysql.Username,
		Passwd:   conf.Mysql.Passwd,
		DB:       conf.Mysql.DB,
		IsDebug:  conf.IsDebug,
	}, new(gorm.Config))
	if err != nil {
		panic(err)
	}
	// TODO: yjy 将manager和路由进行统一管理，阻塞管道达到所有携程持续执行
	// manager, _ := model.NewLoginAndRegisterManager(orm)
	// err = manager.CreateTables(context.Background())
	if err != nil {
		log.Printf("[info] create table user login err:%v,orm:%v", err, orm)
	}
	// 运行组，param1 为执行函数，param2 为回调函数，执行函数发生错误时，回调函数会执行
	g := &run.Group{}
	g.Add(func() error {
		log.Printf("[info] start http server listening %s", endPoint)
		return server.ListenAndServe()
	}, func(error error) {
		if err = server.Shutdown(context.Background()); err != nil {
			log.Printf("shutdown http server failed %s", err)
		}
	})
	// log.Printf("[info] start http server listening %s", endPoint)
	// if err := server.ListenAndServe(); err != nil {
	//	 log.Printf("start http server failed %s", err)
	// }
	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	if err = g.Run(); err != nil {
		log.Println("start http server err, will showdown")
		os.Exit(1)
	}
}
