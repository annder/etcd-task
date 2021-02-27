package scheduler



func Init() {
	t := InitTaskChan()
	t.WatchTaskChange()
}
