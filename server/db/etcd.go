package db

import (
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)
var Kv clientv3.KV
var Lease clientv3.Lease
var Watcher clientv3.Watcher


func ETCDInit(){
	client,err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: 1 * time.Second,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	Kv = clientv3.NewKV(client)
	Lease = clientv3.NewLease(client)

	Watcher = clientv3.NewWatcher(client)
}

