package main

import (
	"context"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
)

func getEtcdClient() (*clientv3.Client, error) {
	ep0, exists := os.LookupEnv("ETCD_CONN_0")
	if !exists {
		logrus.Panic("ETCD_CONN_0 variable not set!")
	}
	ep1, exists := os.LookupEnv("ETCD_CONN_1")
	if !exists {
		logrus.Panic("ETCD_CONN_1 variable not set!")
	}
	ep2, exists := os.LookupEnv("ETCD_CONN_2")
	if !exists {
		logrus.Panic("ETCD_CONN_2 variable not set!")
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{ep0, ep1, ep2},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logrus.Fatal(err)
		return nil, err
	}
	return cli, nil
}

func kvPut(k, v string) (*clientv3.PutResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	putResp, err := etcdClient.Put(ctx, k, v)
	cancel()
	if err != nil {
		return nil, err
	}
	return putResp, nil

}

func kvGet(k string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	getResp, err := etcdClient.Get(ctx, k)
	cancel()
	if err != nil {
		return nil, err
	}
	return getResp, nil

}

func kvDel(k string) (*clientv3.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	delResp, err := etcdClient.Delete(ctx, k)
	cancel()
	if err != nil {
		return nil, err
	}
	return delResp, nil

}
