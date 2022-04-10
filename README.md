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

## Bootstrap project

1. Init Go lang projet:
 `go mod init github.com/mstrielnikov/k8-operator-sample`

2. Bootsrap Kubebuilder project:
 `kubebuilder init --domain demo.domain --repo github.com/mstrielnikov/k8-operator-sample`

3. Generate controller template:
 `kubebuilder create api --group scale --version v1 --kind DemoDeployment`

4. Create K8s test env with Kind:
  `kind create cluster --name demo-cluster`

## Install CRD and controller to cluster

1. Generate manifests for CRDs (see [this](https://book.kubebuilder.io/cronjob-tutorial/running-webhook.html) if you need to autogenerate webhooks):
 `make manifests`

2. Install CRDs to cluster:
 `make install`

3. Build and run controller:
 `make run`

## Deploy controller for CRD

1. Build controller:
 `make docker-build IMG=mstrielnikov/demodeployment-controller:v1`

2. Load image to cluster:
 `kind load docker-image mstrielnikov/demodeployment-controller:v1 --name demo-cluster`

3. Deploy controller to cluster:
 `make deploy IMG=mstrielnikov/demodeployment-controller:v1`

## Demo

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

* Check if CRD API deployed: `kubectl api-resources | grep DemoDeployment` 
```
NAME                       SHORTNAMES   APIVERSION                    NAMESPACED   KIND
demodeployments                         scale.mstrielnikov/v1         true         DemoDeployment
```

* Check is resource running: `kubectl get demodeployment.scale.mstrielnikov/demodeployment-sample` 
```
NAME                    AGE
demodeployment-sample   19s
```

* Find and check appropriate pods are running: `kubectl get pods -A | grep demo`
```
NAMESPACE                    NAME                                                      READY   STATUS    RESTARTS        AGE
default                      demodeployment-sample-88fd7557-tdvr6                      1/1     Running   0               29s
default                      demodeployment-sample-88fd7557-vt9t8                      1/1     Running   0               29s
```

* Describe pods: `kubectl describe pod demodeployment-sample-88fd7557-tdvr6`
```yaml
Name:         demodeployment-sample-88fd7557-tdvr6
Namespace:    default
Priority:     0
Node:         kind-control-plane/172.18.0.2
Start Time:   Sun, 10 Apr 2022 21:21:21 +0300
Labels:       app=demodeployment-sample
              controller=DemoDeploymentController
              pod-template-hash=88fd7557
Annotations:  <none>
Status:       Running
IP:           10.244.0.17
IPs:
  IP:           10.244.0.17
Controlled By:  ReplicaSet/demodeployment-sample-88fd7557
Containers:
  nginx:
    Container ID:   containerd://5b4048dab643e36a08d5772ffa3f789e8227b53e239abb8607c8c6064294f438
    Image:          nginx:1.14.2
    Image ID:       docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
    Port:           80/TCP
    Host Port:      0/TCP
    State:          Running
      Started:      Sun, 10 Apr 2022 21:21:22 +0300
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-qdn5p (ro)
Conditions:
  Type              Status
  Initialized       True 
  Ready             True 
  ContainersReady   True 
  PodScheduled      True 
Volumes:
  kube-api-access-qdn5p:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  15m   default-scheduler  Successfully assigned default/demodeployment-sample-88fd7557-tdvr6 to kind-control-plane
  Normal  Pulled     15m   kubelet            Container image "nginx:1.14.2" already present on machine
  Normal  Created    15m   kubelet            Created container nginx
  Normal  Started    15m   kubelet            Started container nginx
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
