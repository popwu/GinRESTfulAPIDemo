package main

import (
	"fmt"
	"os"

	"apidemo/dbconfig"
	"apidemo/routes"
	"apidemo/store"
	"apidemo/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载环境变量
	err := utils.LoadEnv()
	if err != nil {
		fmt.Println("Error loading environment variables:", err)
		return
	}

	// 连接数据库
	err = dbconfig.DBConnect()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	// 初始化数据库
	err = store.AutoMigrate()
	if err != nil {
		fmt.Println("Error migrating database:", err)
		return
	}

	// 启动 Gin 服务器
	r := gin.Default()
	routes.SetupRouter(r)
	err = r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		fmt.Println("Error starting Gin server:", err)
		return
	}

}
