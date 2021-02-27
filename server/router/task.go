package router

import (
	"context"
	"encoding/json"
	"gin-etcd-task/app"
	"gin-etcd-task/common"
	"gin-etcd-task/db"
	"gin-etcd-task/model"
	"gin-etcd-task/scheduler"
	"gin-etcd-task/utils"
	"github.com/gin-gonic/gin"
	"log"
)

func AddTask(ctx *gin.Context) {
	name := ctx.DefaultPostForm("name", "")
	cron := ctx.DefaultPostForm("cron", "")
	bash := ctx.DefaultPostForm("bash", "")
	mt := model.Task{}
	createOk, err := mt.New(name, cron, bash)
	if err != nil {
		app.Err(ctx, 500, "服务器内部错误")
		log.Fatal(err.Error())
	}
	app.Ok(ctx, createOk, "添加成功")
}

// 这里需要重构一下
func DelTask(ctx *gin.Context) {
	id := ctx.DefaultPostForm("id", "0")
	toIdInt := utils.Int(id)
	t := model.Task{}
	if toIdInt <= 0 {
		app.Err(ctx, 302, "请输入正确的ID")
		return
	}
	r, err := t.FindById(toIdInt)
	if err != nil {
		app.Err(ctx, 500, "服务器内部错误")
		log.Fatal(err)
	}
	tm := scheduler.TaskChanMange{}
	ok := tm.EndCron(r.Name)
	if !ok {
		app.Err(ctx, 500, "服务器内部错误")
		log.Fatal(err)
		return
	}
	if t.Del(toIdInt) {
		app.Ok(ctx, gin.H{}, "删除成功")
	} else {
		app.Ok(ctx, gin.H{}, "删除失败")
	}
}

func GetAllTask(ctx *gin.Context) {
	t := model.Task{}
	result, err := t.GetAll()
	if err != nil {
		app.Err(ctx, 500, "服务器内部错误")
		log.Fatal(err)
	}
	app.Ok(ctx, gin.H{"list": result}, "获取成功")
}

func UpdateStatus(ctx *gin.Context) {
	id := utils.Int(ctx.DefaultPostForm("id", "0"))
	status := utils.Int(ctx.DefaultPostForm("status", "0"))
	if status == 0 || id == 0 {
		app.Err(ctx, 300, "请传入正确的id和status")
	}
	t := model.Task{}

	ok := t.UpdateStatus(id, status)

	if ok {
		res, err := t.FindById(id)
		if err != nil {
			app.Err(ctx, 500, "服务器内部错误")
			log.Fatal(err)
		}
		switch status {
		case 2:
			resjson, err := json.Marshal(res)
			if err != nil {
				log.Fatal(err)
			}
			ok, err := db.Kv.Put(context.TODO(), common.ETCD_TASK_NAME+res.Name, string(resjson))
			if ok != nil && ok.Header.Revision > 0 && err == nil {
				app.Ok(ctx, gin.H{}, "修改成功")
			} else {
				app.Err(ctx, 500, "服务器内部错误")
			}
		case 1:
			tcm := scheduler.TaskChanMange{}
			ok :=  tcm.EndCron(res.Name)
			if ok {
				app.Ok(ctx, gin.H{}, "修改成功")
			} else {
				app.Err(ctx, 500, "服务器内部错误")
			}
		}
	} else {
		app.Err(ctx, 500, "修改失败，服务器内部失败")
	}
}
