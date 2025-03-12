package main

import (
	"backend/controller"
	"backend/dao/mysql"
	"backend/dao/redis"
	"backend/logger"
	"backend/pkg/snowflake"
	"backend/routers"
	"backend/settings"
	"fmt"
)

// @host 127.0.0.1:8081
// @BasePath /api/v1/
func main() {
	//var confFile string
	//flag.StringVar(&confFile, "conf", "./conf/config.yaml", "Configuration file")
	//flag.Parse()
	// Load configuration
	if err := settings.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close() // Close database connection when the program exits

	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}

	defer redis.Close()
	// Snowflake algorithm to generate distributed ID
	if err := snowflake.Init(1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator Trans failed,err:%v\n", err)
		return
	}
	// Register routes
	r := routers.SetupRouter(settings.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
