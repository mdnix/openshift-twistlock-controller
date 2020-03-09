package main

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

var etcdClient *clientv3.Client

const maxRetries = 10

var serverStartTime time.Time
var twc *TwistlockConfig

var configPath *string

const twgrpAPI = "/api/v1/groups"
const twcollAPI = "/api/v1/collections"

// Config struct generated from config.yaml
type Config struct {
	Resources struct {
		Pod                   bool
		Deployment            bool
		Replicationcontroller bool
		Replicaset            bool
		Daemonset             bool
		Services              bool
		Secret                bool
		Configmap             bool
		Rolebinding           bool
	} `yaml:"resources"`
	Handler struct {
		Name string
	} `yaml:"handler"`
}

// Handler is implemented by any handler.
// The Handle method is used to process event
type Handler interface {
	Init(c Config) error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(obj interface{})
}

// Event indicate the informerEvent
type Event struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
	newObj       *rbacv1.RoleBinding
	oldObj       *rbacv1.RoleBinding
}

// Controller struct
type Controller struct {
	logger       *logrus.Entry
	clientset    kubernetes.Interface
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	eventHandler Handler
}

// Rolebinding struct, used to create a Twistlock Collection and Group
type Rolebinding struct {
	Name      string
	Namespace string
	Group     []string
	CN        []string
	Action    string
	Role      string
}

// TwistlockConfig struct, used to make API calls to console
type TwistlockConfig struct {
	User     string
	Password string
	Host     string
}

// TwistlockGroup struct, used to generate a Group json object
type TwistlockGroup struct {
	CN    string
	Group string
	Role  string
}

// TwistlockCollection struct, used to generate a Collection json object
type TwistlockCollection struct {
	CN        string
	Namespace string
}

// CollectionAPI needed to get JSON object from API
type CollectionAPI struct {
	Name        string   `json:"name"`
	Color       string   `json:"color"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Containers  []string `json:"containers"`
	Hosts       []string `json:"hosts"`
	Labels      []string `json:"labels"`
	Services    []string `json:"services"`
	Functions   []string `json:"funtions"`
	Namespaces  []string `json:"namespaces"`
	AppIDs      []string `json:"appIDs"`
}

// GroupAPI needed to get JSON object from API
type GroupAPI struct {
	GroupName    string   `json:"groupName"`
	User         []string `json:"user"`
	LastModified string   `json:"lastModified"`
	Owner        string   `json:"owner"`
	LdapGroup    bool     `json:"ldapGroup"`
	SamlGroup    bool     `json:"samlGroup"`
	Role         string   `json:"role"`
	ID           string   `json:"_id"`
	Projects     []string `json:"projects"`
	GroupID      string   `json:"groupId"`
	Collections  []string `json:"collections"`
}
