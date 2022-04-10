# k8s-operator-sample

## Usefull reading
* [K8s types and common machinery](https://iximiuz.com/en/posts/kubernetes-api-go-types-and-common-machinery/)
* [K8s official docs: develop simple controller](https://kubernetes.io/blog/2021/06/21/writing-a-controller-for-pod-labels/)
* [K8s official docs: extend API with CRDs](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/)

## Requirements

* [Go lang](https://go.dev/doc/install)
* [Kubebuilder](https://book.kubebuilder.io/quick-start.html) (note that latest Kubebuilder version does not works with Go 1.18)
* [Docker](https://www.docker.com/get-started/)
* [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) as a local K8s test env

## Bootstrap

1. Init Go lang projet:
 `go mod init github.com/mstrielnikov/k8-operator-sample`

2. Bootsrap Kubebuilder project:
 `kubebuilder init --domain demo.domain --repo github.com/mstrielnikov/k8-operator-sample`

3. Generate controller template:
 `kubebuilder create api --group scale --version v1 --kind DemoDeploument`

## Install CRD and controller to cluster

0. Build and run controller as a regular programm:
 `make run`

1. Generate manifests:
 `make manifests`

2. Install manifests to cluster:
 `make install`

## Deploy controller for CRD

0. Create K8s test env with Kind:
  `kind create cluster --name demo-cluster`

1. Build controller:
 `make docker-build IMG=mstrielnikov/demodeployment-controller:v1`

2. Load image to cluster:
 `kind load docker-image mstrielnikov/demodeployment-controller:v1 --name demo-cluster`

3. Run controller:
 `make deploy IMG=mstrielnikov/demodeployment-controller:v1`

## Deploy CRD in cluster

Create CRD: 
  `kubectl create -f config/samples/scale_v1_demodeployment.yaml` 

Test deployment `scale_v1_demodeployment.yaml` with content following:

```yaml
---
apiVersion: scale.mstrielnikov/v1
kind: DemoDeployment
metadata:
  name: demodeployment-sample
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```
___

# Demo

* Check if CRD API deployed `kubectl api-resources | grep DemoDeployment` 
```
NAME                       SHORTNAMES   APIVERSION                    NAMESPACED   KIND
demodeployments                         scale.mstrielnikov/v1         true         DemoDeployment
```

* Check is resource running `kubectl get demodeployment.scale.mstrielnikov/demodeployment-sample` 
```
NAME                    AGE
demodeployment-sample   19s
```

* Describe and validate CRD resource `kubectl get demodeployments.scale.mstrielnikov -o yaml`
```yaml
apiVersion: v1
items:
- apiVersion: scale.mstrielnikov/v1
  kind: DemoDeployment
  metadata:
    annotations:
      kubectl.kubernetes.io/last-applied-configuration: |
        {"apiVersion":"scale.mstrielnikov/v1","kind":"DemoDeployment","metadata":{"annotations":{},"name":"demodeployment-sample","namespace":"default"},"spec":{"replicas":2,"selector":{"matchLabels":{"app":"nginx"}},"template":{"metadata":{"labels":{"app":"nginx"}},"spec":{"containers":[{"image":"nginx:1.14.2","name":"nginx","ports":[{"containerPort":80}]}]}}}}
    creationTimestamp: "2022-04-09T21:11:34Z"
    generation: 4
    name: demodeployment-sample
    namespace: default
    resourceVersion: "42525"
    uid: 7f591933-409c-4603-8c9e-3c95612406ac
  spec:
    replicas: 2
    selector:
      matchLabels:
        app: nginx
    template:
      metadata: {}
      spec:
        containers:
        - image: nginx:1.14.2
          name: nginx
          ports:
          - containerPort: 80
            protocol: TCP
          resources: {}
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""

```

## Cleanup

1. Delete deployment `kubectl delete demodeployment demodeployment-sample`
2. Delete CRD `make uninstall`
3. Undeploy operator `make undeploy`