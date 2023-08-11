package main

import (
	"server_demo/dao"
	"server_demo/routes"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	//初始化数据库
	dao.InitPostIndexMap()
	dao.InitTopicIndexMap()
	//配置路由
	router := gin.New()
	routes.SetupRoute(router)
	//启动
	err := router.Run(":8080")
	if err != nil {
		fmt.Println(err.Error())
	}
}