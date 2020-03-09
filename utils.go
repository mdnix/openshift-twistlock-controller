package main

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	_, configexists := os.LookupEnv("KUBECONFIG")

	if configexists {
		logrus.Info("Loading kubeconfig", os.Getenv("KUBECONFIG"))
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		if err != nil {
			panic(err.Error())
		}
		logrus.Infof("Kubeconfig %s initialized", os.Getenv("KUBECONFIG"))
	} else {
		logrus.Info("Loading ServiceAccount token")
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		logrus.Infof("ServiceAccount token initialized")
	}
	return kubernetes.NewForConfig(config)
}

func getTwistlockConfig() (*TwistlockConfig, error) {
	twUser, ok := os.LookupEnv("TWISTLOCK_USER")
	if !ok {
		logrus.Println("Env TWISTLOCK_USER not defined!")
	}
	twPass, ok := os.LookupEnv("TWISTLOCK_PASSWORD")
	if !ok {
		logrus.Println("Env TWISTLOCK_PASSWORD not defined!")
	}
	twHost, ok := os.LookupEnv("TWISTLOCK_HOST")
	if !ok {
		logrus.Println("Env TWISTLOCK_HOST not defined!")
	}

	if len(twUser) > 0 && len(twPass) > 0 && len(twHost) > 0 {
		twc = &TwistlockConfig{
			User:     twUser,
			Password: twPass,
			Host:     twHost,
		}
		return twc, nil
	}
	return nil, errors.New("Could not get Twistlock config")
}

// GetObjectMetaData returns metadata of a given k8s object
func GetObjectMetaData(obj interface{}) metav1.ObjectMeta {

	var objectMeta metav1.ObjectMeta

	switch object := obj.(type) {
	case *rbacv1.RoleBinding:
		objectMeta = object.ObjectMeta
	case *appsv1.Deployment:
		objectMeta = object.ObjectMeta
	case *apiv1.ReplicationController:
		objectMeta = object.ObjectMeta
	case *appsv1.ReplicaSet:
		objectMeta = object.ObjectMeta
	case *appsv1.DaemonSet:
		objectMeta = object.ObjectMeta
	case *apiv1.Service:
		objectMeta = object.ObjectMeta
	case *apiv1.Pod:
		objectMeta = object.ObjectMeta
	case *batchv1.Job:
		objectMeta = object.ObjectMeta
	case *apiv1.PersistentVolume:
		objectMeta = object.ObjectMeta
	case *apiv1.Namespace:
		objectMeta = object.ObjectMeta
	case *apiv1.Secret:
		objectMeta = object.ObjectMeta
	case *extv1beta1.Ingress:
		objectMeta = object.ObjectMeta
	}
	return objectMeta
}

func sliceContains(s []string, i string) bool {
	for _, a := range s {
		if a == i {
			return true
		}
	}
	return false
}

func sliceRemove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
