package main

import "github.com/sirupsen/logrus"

// Default handler implements Handler interface,
// print each event with JSON format
type Default struct {
}

// Init initializes handler configuration
// Do nothing for default handler
func (d *Default) Init(c Config) error {
	return nil
}

// ObjectCreated sends events on object creation
func (d *Default) ObjectCreated(obj interface{}) {
	logrus.Info("Default CREATE function invoked")
}

// ObjectDeleted sends events on object deletion
func (d *Default) ObjectDeleted(obj interface{}) {
	logrus.Info("Default DELETE function invoked")
}

// ObjectUpdated sends events on object updation
func (d *Default) ObjectUpdated(obj interface{}) {
	logrus.Info("Default UPDATE function invoked")

}
