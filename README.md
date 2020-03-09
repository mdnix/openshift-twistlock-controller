# Kubernetes Twistlock Controller

The Kubernetes Twistlock Controller is a controller written in Go. The current version uses the client-go library in version 1.17.0.
A configuration file named config.yaml allows the selective control of Kubernetes resources and if needed, own handle functions can be defined for these resources. The handler defines how the controller should behave when an object is added, updated or deleted.

This controller was written to allow a team of developers to easily work with OpenShift and Twistlock. This means that developers can login to OpenShift and Twistlock with the same user and see the same namespace view on both platforms.
The controller uses the RoleBinding objects to create various Twistlock Collections. This allows you to restrict access to different namespaces in the Twistlock Console.

The controller can run outside of the cluster or on the cluster as a Pod.



## Usage

### Configuration

The following config will enable a watch for RoleBindings. The Twistlock handler will take care of the RoleBinding objects:
```yaml
resources:
  pod: false
  deployment: false
  replicationcontroller: false
  replicaset: false
  daemonset: false
  services: false
  secret: false
  configmap: false
  rolebinding: true
handler:
  name: Twistlock
```

### Running the Controller with an out-of-cluster-config:
```bash
export KUBECONFIG=/path/to/config
./twistlock-controller
```

### Running the Controller with an in-cluster-config:
The Controller can be deployed on Kubernetes as a Pod. The application then makes use of the ServiceAccount and its respective ServiceAccount Token.
Creating a new ServiceAccount and adding a the cluster-reader ClusterRole is the recommended way:
```bash
oc create sa twistlock-cluster-reader -n mgt-infra-controllers
oc adm policy add-cluster-role-to-user cluster-reader -z twistlock-cluster-reader -n mgt-infra-controllers
```

## Prerequisites

The controller needs an etcd cluster to store the rolebindings for processing at a later point in time.
A cluster consisting of 3 instances should be already available and ready to use for various controllers inside the namespace mgt-infra-controllers.
If that's not the case you can initialize a cluster as follows:

### Create an ImageStream and pull the latest etcd image from Red Hat:
```bash
oc create -f is-etcd.yaml -n mgt-infra-controllers
```

### Initialize the cluster:
```bash
oc create -f deploy.yaml -n mgt-infra-controllers
```
This will create a StatefulSet (3 replicas), Service and PersistentVolumeClaims/PersistentVolumes 

## Installation

### Create the ImageStream
```bash
oc create -f is-controller.yaml -n mgt-infra-controllers
```

### Add Git Source Secret
```bash
oc create source-secret.yaml -n mgt-infra-controllers
```

### Add Twistlock Secret
```bash
oc create twistlock-secret.yaml -n mgt-infra-controllers
```

### Add BuildConfig
```bash
oc create bc.yaml -n mgt-infra-controllers
```

### Add DeploymentConfig
```bash
oc create dc.yaml -n mgt-infra-controllers
```

### Add Service for Controller Health Checks
```bash
oc create svc-health -n mgt-infra-controllers
```

### Starting a build
A build can be started directly from the command prompt:
```bash
oc start-build kubernetes-twistlock-controller --follow -n mgt-infra-controllers
```
As soon as the build is done a new Deployment of the image will be rolled out due to ImageChange Triggers.

### Controller internals
![alt text](https://tfs-prod.service.raiffeisen.ch:8081/tfs/RCH/87c91262-4a39-419e-93de-93b48f0f3d84/_apis/git/repositories/417f93a8-48d9-4309-af3b-fdef7ce71fe1/Items?path=%2Ftools%2Fkubernetes-twistlock-controller%2Fdrawing%2Farchitecture.jpeg&versionDescriptor%5BversionOptions%5D=0&versionDescriptor%5BversionType%5D=0&versionDescriptor%5Bversion%5D=master&download=false&resolveLfs=true&%24format=octetStream&api-version=5.0-preview.1)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
