package model

import (
	"gin-etcd-task/db"
	"gorm.io/gorm"
	"log"
)

// 创建任务表，修改某项值，然后重新更新
type Task struct {
	gorm.Model
	Name   string `json:"name"`
	Cron   string `json:"cron"`
	Bash   string `json:"bash"`
	Status int    `json:"status"`
}

func (t *Task) New(name, cron, bash string) (tt *Task, err error) {
	tt = &Task{
		Name:   name,
		Cron:   cron,
		Bash:   bash,
		Status: 1,
	}
	ok := db.Mysql.Create(tt)
	if ok.Error != nil {
		return
	}
	return
}

// 更新状态
func (t *Task) UpdateStatus(id, status int) bool {
	ok := db.Mysql.Model(&Task{}).
		Where("id = ?", id).
		Update("status", status)
	if ok.Error != nil {
		log.Fatal(ok.Error)
		return false
	}
	return true
}

// 获取所有的值
func (t *Task) GetAll() (r []Task, err error) {
	err = db.Mysql.Find(&r).Error
	return
}

func (t *Task) FindById(id int) (r Task, err error) {
	err = db.Mysql.Where("id = ?", id).First(&r).Error
	return
}

// 所有的更新
func (t *Task) Update(name, cron, bash string) bool {
	updateTask := &Task{
		Name: name,
		Cron: cron,
		Bash: bash,
	}
	ok := db.Mysql.Model(&Task{}).Updates(updateTask)
	if ok.Error != nil {
		return false
	}
	return true
}

func (t *Task) Del(id int) bool {
	ok := db.Mysql.Where("id = ?", id).Delete(&Task{})
	if ok.Error != nil {
		return false
	}
	return true
}
