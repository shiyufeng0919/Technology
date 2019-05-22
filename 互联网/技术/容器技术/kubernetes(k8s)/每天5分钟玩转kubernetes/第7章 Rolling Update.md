
### 第7章 Rolling Update

    滚动更新是一次只更新一小部分副本，成功后再更新更多的副本。最终完成所有副本的更新。
    
    滚动更新最大好处是零停机，整个更新过程始终有副本在运行，从而保证了业务的持续性。
    
#### 1.实践

>目标：将初始镜像httpd:2.2.31更新到httpd:2.2.32
    
    ###Step1,编写httpd.yml
    [root@node1 k8s]# cat httpd.yml 
    apiVersion: apps/v1beta1
    kind: Deployment
    metadata:
      name: httpd
    spec:
      replicas: 3
      template:
        metadata:
          labels:
            run: httpd
        spec:
          containers:
          - name: httpd
            image: httpd:2.2.31   #httpd镜像
            ports:
            - containerPort: 80

    ###Step2,部署
    [root@node1 k8s]# kubectl get deployment httpd -o wide
    NAME    READY   UP-TO-DATE   AVAILABLE   AGE     CONTAINERS   IMAGES         SELECTOR
    httpd   3/3     3            3           2m49s   httpd        httpd:2.2.31   run=httpd
    
    ###Step3,查看Deployment,3个副本均已部署完成
    [root@node1 k8s]# kubectl get replicaset -o wide
    NAME               DESIRED   CURRENT   READY   AGE     CONTAINERS   IMAGES         SELECTOR
    httpd-76cfb94bf4   3         3         3       2m56s   httpd        httpd:2.2.31   pod-template-hash=76cfb94bf4,run=httpd
    
    ###Step3,查看pod，三个副本pod均已启动
    [root@node1 k8s]# kubectl get pod -o wide
    NAME                     READY   STATUS    RESTARTS   AGE    IP              NODE    NOMINATED NODE   READINESS GATES
    httpd-76cfb94bf4-qw4c6   1/1     Running   0          3m6s   10.244.135.57   node3   <none>           <none>
    httpd-76cfb94bf4-w48pp   1/1     Running   0          3m6s   10.244.104.24   node2   <none>           <none>
    httpd-76cfb94bf4-wgxp2   1/1     Running   0          3m6s   10.244.104.23   node2   <none>           <none>
    
    上述部署过程如下：
    (1).创建Deployment httpd
    (2).创建replicaset httpd-76cfb94bf4
    (3).创建三个pod
    (4).当前镜像为httpd:2.2.31
    
    ###修改配置文件，将镜像由httpd:2.2.31修改为httpd:2.2.32
    [root@node1 k8s]# kubectl apply -f httpd.yml 
    deployment.apps/httpd configured
    
    ###查看deployment,发现版本已经更新
    [root@node1 k8s]# kubectl get deployment httpd -o wide
    NAME    READY   UP-TO-DATE   AVAILABLE   AGE   CONTAINERS   IMAGES         SELECTOR
    httpd   3/3     3            3           13m   httpd        httpd:2.2.32   run=httpd
    
    ###此时查看replicaset，可见httpd-76cfb94bf4已被httpd-6cf6bf9f57替代
    [root@node1 k8s]# kubectl get replicaset -o wide
    NAME               DESIRED   CURRENT   READY   AGE    CONTAINERS   IMAGES         SELECTOR
    httpd-6cf6bf9f57   3         3         3       3m7s   httpd        httpd:2.2.32   pod-template-hash=6cf6bf9f57,run=httpd
    httpd-76cfb94bf4   0         0         0       10m    httpd        httpd:2.2.31   pod-template-hash=76cfb94bf4,run=httpd
    
    ###查看pod,发现httpd-6cf6bf9f57管理了三个新的pod,先前httpd-76cfb94bf4里已经没有任何pod
    [root@node1 k8s]# kubectl get pod -o wide
    NAME                     READY   STATUS    RESTARTS   AGE     IP              NODE    NOMINATED NODE   READINESS GATES
    httpd-6cf6bf9f57-9vgn5   1/1     Running   0          3m10s   10.244.135.58   node3   <none>           <none>
    httpd-6cf6bf9f57-ctmqw   1/1     Running   0          13s     10.244.135.59   node3   <none>           <none>
    httpd-6cf6bf9f57-glppv   1/1     Running   0          97s     10.244.104.25   node2   <none>           <none>
    
    
    ###查看具体部署过程
    [root@node1 k8s]# kubectl describe deployment httpd
    Name:                   httpd
    Namespace:              default
    CreationTimestamp:      Sat, 11 May 2019 04:52:47 +0000
    Labels:                 run=httpd
    Annotations:            deployment.kubernetes.io/revision: 2
                            kubectl.kubernetes.io/last-applied-configuration:
                              {"apiVersion":"apps/v1beta1","kind":"Deployment","metadata":{"annotations":{},"name":"httpd","namespace":"default"},"spec":{"replicas":3,"...
    Selector:               run=httpd
    Replicas:               3 desired | 3 updated | 3 total | 3 available | 0 unavailable
    StrategyType:           RollingUpdate
    MinReadySeconds:        0
    RollingUpdateStrategy:  25% max unavailable, 25% max surge
    Pod Template:
      Labels:  run=httpd
      Containers:
       httpd:
        Image:        httpd:2.2.32
        Port:         80/TCP
        Host Port:    0/TCP
        Environment:  <none>
        Mounts:       <none>
      Volumes:        <none>
    Conditions:
      Type           Status  Reason
      ----           ------  ------
      Available      True    MinimumReplicasAvailable
      Progressing    True    NewReplicaSetAvailable
    OldReplicaSets:  <none>
    NewReplicaSet:   httpd-6cf6bf9f57 (3/3 replicas created)
    Events:
      Type    Reason             Age    From                   Message
      ----    ------             ----   ----                   -------
      Normal  ScalingReplicaSet  13m    deployment-controller  Scaled up replica set httpd-76cfb94bf4 to 3 
      Normal  ScalingReplicaSet  6m38s  deployment-controller  Scaled up replica set httpd-6cf6bf9f57 to 1    ##新replica增加一个pod总数为1
      Normal  ScalingReplicaSet  5m5s   deployment-controller  Scaled down replica set httpd-76cfb94bf4 to 2  ##旧replica减小一个pod总数为2
      Normal  ScalingReplicaSet  5m5s   deployment-controller  Scaled up replica set httpd-6cf6bf9f57 to 2    ##新replica增加一个pod总数为2
      Normal  ScalingReplicaSet  3m41s  deployment-controller  Scaled down replica set httpd-76cfb94bf4 to 1  ##旧replica减小一个pod总数为1
      Normal  ScalingReplicaSet  3m41s  deployment-controller  Scaled up replica set httpd-6cf6bf9f57 to 3    ##新replica增加一个pod总数为3
      Normal  ScalingReplicaSet  3m40s  deployment-controller  Scaled down replica set httpd-76cfb94bf4 to 0  ##旧replica减小一个pod总数为0
      
    **上述：每次只更新替换一个Pod**
    
    另：每次替换的pod数量是可以定制的，k8s提供了两个参数maxSurge和maxUnavailable来精细控制pod的替换数量
    
#### 2.回滚

    kubectl apply每次更新应用时，k8s均会记录下当前的配置，保存为一个revision(版次).这样可以回滚到某个特定的revision.
    
    默认配置下，k8s只会保留最近的几个revision,可在Deployment配置文件中通过revisionHistoryLimit属性增加revision数量
    
>实践回滚功能。应用三个配置文件：

>httpd.v1.yml(httpd:2.4.16),httpd.v2.yml(httpd:2.4.17),httpd.v3.yml(httpd:2.4.18)

    [root@node1 k8s]# cat httpd.v1.yml
    apiVersion: apps/v1beta1
    kind: Deployment
    metadata:
      name: httpd
    spec:
      replicas: 3
      template:
        metadata:
          labels:
            run: httpd
        spec:
          containers:
          - name: httpd
            image: httpd:2.4.16  #v2和v3修改此处镜像版本号即可
            ports:
            - containerPort: 80
            
    ###部署: --record作用：将当前命令记录到revision记录中，即知道每个revison对应哪个配置文件(方便定位)
    [root@node1 k8s]# kubectl apply -f httpd.v1.yml --record
    deployment.apps/httpd created
    
    ###查看Deployment部署完成
    [root@node1 k8s]# kubectl get deployment httpd -o wide
    NAME    READY   UP-TO-DATE   AVAILABLE   AGE     CONTAINERS   IMAGES         SELECTOR
    httpd   3/3     3            3           3m35s   httpd        httpd:2.4.16   run=httpd
    
    ###更新v1->v2
    [root@node1 k8s]# kubectl apply -f httpd.v2.yml --record
    deployment.apps/httpd configured
    [root@node1 k8s]# kubectl get deployment httpd -o wide
    NAME    READY   UP-TO-DATE   AVAILABLE   AGE     CONTAINERS   IMAGES         SELECTOR
    httpd   3/3     1            3           5m11s   httpd        httpd:2.4.17   run=httpd
    
    ###更新v2->v3
    root@node1 k8s]# kubectl apply -f httpd.v3.yml --record
    deployment.apps/httpd configured
    [root@node1 k8s]# kubectl get deployment httpd -o wide
    NAME    READY   UP-TO-DATE   AVAILABLE   AGE     CONTAINERS   IMAGES         SELECTOR
    httpd   3/3     1            3           6m10s   httpd        httpd:2.4.18   run=httpd
    
    ###查看revision的历史记录
    [root@node1 k8s]# kubectl rollout history deployment httpd
    deployment.extensions/httpd 
    REVISION  CHANGE-CAUSE  ##CHANGE-CAUSE即为--record的结果
    1         kubectl apply --filename=httpd.v1.yml --record=true
    2         kubectl apply --filename=httpd.v2.yml --record=true
    3         kubectl apply --filename=httpd.v3.yml --record=true
    
    
    ###执行回滚,v3->v1
    [root@node1 k8s]# kubectl rollout undo deployment httpd --to-revision=1
    deployment.extensions/httpd rolled back
    
    ###查看Deployment已经恢复到v1版本(2.4.16)
    [root@node1 k8s]# kubectl get deployment httpd -o wide
    NAME    READY   UP-TO-DATE   AVAILABLE   AGE     CONTAINERS   IMAGES         SELECTOR
    httpd   3/3     3            3           8m20s   httpd        httpd:2.4.16   run=httpd
    
>滚动更新：采用渐进的方式逐步替换旧版本，若更新不如预期，可回滚操作到更新前状态

--------------------------------------------

**体胖还需勤跑步，人丑就要多读书!!! --开心玉凤**