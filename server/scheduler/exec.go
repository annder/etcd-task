package scheduler

import (
	"log"
	"os/exec"
	"runtime"
)

// 执行脚本
func Exec(cmd string) {
	sys := runtime.GOOS
	var (
		err error
		c   []byte
	)
	if sys == "windows" {
		c, err = exec.Command("cmd", "/C", cmd).Output()
	}
	if sys == "linux" {
		c, err = exec.Command("/bin/bash", "-c", cmd).Output()
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(c))
}
