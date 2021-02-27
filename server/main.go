package main

import (
	"gin-etcd-task/db"
	"gin-etcd-task/router"
	"gin-etcd-task/scheduler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	db.ETCDInit()
	db.MysqlInit()
	go scheduler.Init()

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/all", router.GetAllTask)
	r.POST("/add",router.AddTask)
	r.POST("/del",router.DelTask)
	r.POST("/update",router.UpdateStatus)
	log.Fatal(r.Run(":3000"))
}
