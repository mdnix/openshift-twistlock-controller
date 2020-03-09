package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	go getHealth()

	var err error
	etcdClient, err = getEtcdClient()
	if err != nil {
		logrus.Panicf("Unable to establish connection to etcd. Error: %s", err)
	}
	defer etcdClient.Close()
	twc, err := getTwistlockConfig()
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Println("TWCONFIG: ", twc)
	config := initConfig()
	logrus.Printf("%+v\n ", config)
	startController(config)
}
