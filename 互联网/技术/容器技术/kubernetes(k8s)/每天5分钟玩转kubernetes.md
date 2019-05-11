2019.5.11 开心玉凤==>目标：

    1。kubernetes重要概念和架构
    
    2。学习kubernetes如何编排容器，包括：优化资源利用，高可用，滚动更新，网络插件，服务发现，监控，数据管理，日志管理等
    
教程[完成创建kubernetes集群，部署应用，访问应用，扩展应用，更新应用等](https://kubernetes.io/docs/tutorials/#basics)

博客[每天5分钟玩转k8s](https://mp.sohu.com/profile?xpt=Y2xvdWRtYW42QHNvaHUuY29t&_f=index_pagemp_2&spm=smpc.content.author.3.1557473409950ErIE9g3)

### 第1章 先把kubernetes跑起来

#### 1。启动minikube

Step1.[打开k8s教程菜单](https://kubernetes.io/docs/tutorials/kubernetes-basics/create-cluster/cluster-interactive/)

![img](resources/images/7.jpg)

Step2.启动minikube

    $ minikube version #查看minikube版本
    
    $ kubectl get nodes #获取节点信息。此时已创建好了一个单节点的k8s集群
    
    $ kubectl cluster-info  #查看集群信息
    
![img](resources/images/8.jpg)

#### 2。部署应用

    # 通过kubectl run部署一个应用，名为kubernetes-bootcamp | Docker镜像通过--image指定 | --port设置应用对外服务的端口
    $ kubectl run kubernetes-bootcamp \
    > --image=docker.io/jocatalin/kubernetes-bootcamp:v1 \
    > --port=8080
    
    k8s术语：
    (1).Deployment可理解为应用
    (2).Pod是容器的集合(一般会将紧密相关的一组容器放到一个pod中，同一个pod所有容器共享IP地址和Port间。即它们在一个network namespace中)
        Pod是k8s调度的最小单位，同一Pod中容器始终被一起调度
        
    $ kubectl get pods #查看当前Pod
    
![img](resources/images/9.jpg)

#### 3。访问应用

默认情况下，所有Pod只能在集群内部访问。为能够从外部访问应用，需将容器8080端口映射到节点的端口

    #暴露应用的8080端口
    $ kubectl expose deployment/kubernetes-bootcamp \
    > --type="NodePort" \
    > --port 8080
    
    #查看应用被映射到节点的哪个端口
    $ kubectl get services
    
    #查看节点名称：minikube
    $ kubectl get nodes 
    
    #访问应用
    $ curl minikube:30753 #节点名称:暴露端口
    
![img](resources/images/10.jpg)

#### 4。Scale应用

##### (1)增加副本

默认情况下只会运行一个副本

    #查看副本数
    $ kubectl get deployments
    NAME                  READY   UP-TO-DATE   AVAILABLE   AGE
    kubernetes-bootcamp   1/1     1            1           3m30s
    
    #副本数加到3个
    $ kubectl scale deployments/kubernetes-bootcamp --replicas=3
    
    #查看副本数
    $ kubectl get deployments
    
    #查看当前Pod增加到3个
    $ kubectl get pods
    
    #访问应用,每次执行均会将请求发送到不同Pod，3个副本轮询处理，实现了负载均衡
    $ curl minikube:30753 
    
![img](resources/images/11.jpg)   

##### (2)减少副本

     #将副本数减少到2个(其中一个副本会被删除)
     kubectl scale deployments/kubernetes-bootcamp --replicas=2
     
![img](resources/images/12.jpg)     

     #多次执行,直到仅有2个副本
     $ kubectl get pods
     NAME                                   READY   STATUS    RESTARTS   AGE
     kubernetes-bootcamp-5bf4d5689b-5hz8p   1/1     Running   0          58s
     kubernetes-bootcamp-5bf4d5689b-6sfsc   1/1     Running   0          59s
     
     #访问应用（多次执行，会将请求负载发到2个副本）
     $ curl minikube:30753  

#### 5。滚动更新

##### (1).升级镜像版本

    #升级镜像版本v1->v2
    $kubectl set image deployments/kubernetes-bootcamp kubernetes-bootcamp=jocatalin/kubernetes-bootcamp:v2 deployment.extensions/kubernetes-bootcamp image updated
    
![img](resources/images/13.jpg)  

##### (2).回退版本

    #回退v2->v1
    $ kubectl rollout undo deployments/kubernetes-bootcamp
    
    #观察滚动更新过程
    $ kubectl get pods
    
    #访问应用，验证版本是否已经回退到v1
    $ curl minikube:30753
    Hello Kubernetes bootcamp! | Running on: kubernetes-bootcamp-6c5cfd894b-hzfrj | v=1

-------------------------------------------------

### 第2章 重要概念

#### 1.Cluster

    Cluster是计算，存储和网络资源的集合，k8s利用这些资源运行各种基于容器的应用

#### 2.Master

    Master是Cluster的大脑，职责是调度(决定将应用放在哪里运行)。
    
    Master运行在linux操作系统，可以是物理机/虚拟机

#### 3.Node

    Node职责是运行容器应用
    
    Node由Master管理
    
    Node负责监控并汇报容器的状态，并根据Master要求管理容器的生命周期.
    
    Node运行在linux操作系统，可以是物理机/虚拟机

#### 4.Pod

    Pod是k8s最小工作单元
    
    每个Pod包含一个/多个容器
    
    Pod中的容器会作为一个整体被Master调度到一个Node上运行
    
**Pod两种使用方式：**

(1).运行单一容器(将单个容器封装成Pod)

(2).运行多个容器(相关联容器放到一个Pod中)

#### 5.Controller

**k8s不会直接创建Pod，而是通过Controller来管理Pod。**

    Controller中定义了Pod的部署特性，如：有几个副本/在什么样的Node上运行..
    
**k8s提供了多种Controller**

    (1).Deployment(常用)
    
    如：创建Deployment部署应用，Deployment可管理Pod多个副本，并确保Pod按期望状态运行
    
    (2).ReplicaSet(实现Pod多副本管理)
    
    Deployment是通过ReplicaSet来管理Pod的多个副本，通常不需要直接使用ReplicaSet
    
    (3).DaemonSet
    
    用于每个Node最多只运行一个Pod副本场景
    
    (4).StatefuleSet
    
    保证Pod的每个副本在整个生命周期中名称不变(上述滚动更新名称会改变)。保证副本按固定顺序启动，更新，删除
    
    (5).Job
    
    用于运行结束就删除的应用。（其他Controller中的Pod一般长期持续运行）
    
#### 6.Service

问题：Deployment可部署多个副本，每个Pod都有自己IP(动态，销毁/重启IP会改变),外界如何访问副本？

    k8s Service定义了外界访问一组特定Pod的方式
    
    Service有自己的IP和端口
    
    Service为Pod提供了负载均衡
    
**k8s运行容器(Pod)与访问容器(Pod)由Controller和Service执行**

#### 7.Namespace

问题：当有多个用户/项目组使用同一个k8s Cluster,则如何将他们创建的Controller和Pod等资源分开？

    Namespace
    
    Namespace可将一个物理的Cluster逻辑划分成多个虚拟Cluster,每个Cluster即为一个Namespace.不同Namespace里的资源是完合隔离的。
    
    k8s默认创建两个Namespace
    $ kubectl get namespace / $ kubectl get ns
    NAME          STATUS   AGE
    default       Active   60s
    kube-system   Active   60s
    
-------------------------------------------------

### 第3章 部署k8s集群-kubeadm

部署三台虚拟机器,构建k8s集群

    node1 192.168.56.101   ###master节点
    node2 192.168.56.102   ###node节点
    node3 192.168.56.103   ###node节点

-------------------------------------------------

### 第4章 k8s架构

![img](resources/images/14.jpg)  

-------------------------------------------------

### 第5章 运行应用

    k8s通过各种Controller来管理Pod的生命周期。为满足不同业务场景，k8s开发了Deployment,ReplicaSet,DaemonSet,StatefuleSet,Job等多种Controller。(kind值指定)
    
#### 1。Deployment

##### 1-1.运行Deployment

     ###创建2个副本nginx-deployment  --image指定镜像
     [root@node1 ~]# kubectl run nginx-deployment --image=nginx:1.7.9 --replicas=2
     kubectl run --generator=deployment/apps.v1 is DEPRECATED and will be removed in a future version. Use kubectl run --generator=run-pod/v1 or kubectl create instead.
     deployment.apps/nginx-deployment created
     
     ------------------------
     
     ###查看nginx-deployment的状态，输出显示两个副本正常运行
     [root@node1 ~]# kubectl get deployment nginx-deployment
     NAME               READY   UP-TO-DATE   AVAILABLE   AGE
     nginx-deployment   2/2     2            2           2m33s
     
     
     ###查看deployment详细信息
     [root@node1 ~]# kubectl describe deployment
     Name:                   nginx-deployment
     Namespace:              default
     CreationTimestamp:      Fri, 10 May 2019 01:44:16 +0000
     Labels:                 run=nginx-deployment
     Annotations:            deployment.kubernetes.io/revision: 1
     Selector:               run=nginx-deployment
     Replicas:               2 desired | 2 updated | 2 total | 2 available | 0 unavailable
     StrategyType:           RollingUpdate
     MinReadySeconds:        0
     RollingUpdateStrategy:  25% max unavailable, 25% max surge
     Pod Template:
       Labels:  run=nginx-deployment
       Containers:
        nginx-deployment:
         Image:        nginx:1.7.9
         Port:         <none>
         Host Port:    <none>
         Environment:  <none>
         Mounts:       <none>
       Volumes:        <none>
     Conditions:
       Type           Status  Reason
       ----           ------  ------
       Available      True    MinimumReplicasAvailable
       Progressing    True    NewReplicaSetAvailable
     OldReplicaSets:  <none>
     NewReplicaSet:   nginx-deployment-578fb949d8 (2/2 replicas created)
     Events: ###Events是Deployment的日志,记录了ReplicaSet启动过程,同时验证了Deployment通过ReplicaSet来管理Pod
       Type    Reason             Age    From                   Message
       ----    ------             ----   ----                   -------
       Normal  ScalingReplicaSet  3m10s  deployment-controller  Scaled up replica set nginx-deployment-578fb949d8 to 2
       
     ------------------------
     
      ###查看两个副本已就绪,Name即为上述Events中的Message(Scaled up replica set "nginx-deployment-578fb949d8" to 2)
      [root@node1 ~]# kubectl get replicaset
       NAME                          DESIRED   CURRENT   READY   AGE
       nginx-deployment-578fb949d8   2         2         2       15m
       
       
      ###查看详细信息
      [root@node1 ~]# kubectl describe replicaset
      Name:           nginx-deployment-578fb949d8
      Namespace:      default
      Selector:       pod-template-hash=578fb949d8,run=nginx-deployment
      Labels:         pod-template-hash=578fb949d8
                      run=nginx-deployment
      Annotations:    deployment.kubernetes.io/desired-replicas: 2
                      deployment.kubernetes.io/max-replicas: 3
                      deployment.kubernetes.io/revision: 1
      Controlled By:  Deployment/nginx-deployment  ###指明此ReplicaSet是由Deployment nginx-deployment创建的
      Replicas:       2 current / 2 desired   ###Deployment创建Replicas
      Pods Status:    2 Running / 0 Waiting / 0 Succeeded / 0 Failed
      Pod Template:
        Labels:  pod-template-hash=578fb949d8
                 run=nginx-deployment
        Containers:
         nginx-deployment:
          Image:        nginx:1.7.9
          Port:         <none>
          Host Port:    <none>
          Environment:  <none>
          Mounts:       <none>
        Volumes:        <none>
      Events: ###两个副本日志
        Type    Reason            Age   From                   Message
        ----    ------            ----  ----                   -------
        Normal  SuccessfulCreate  14m   replicaset-controller  Created pod: nginx-deployment-578fb949d8-4h4sg
        Normal  SuccessfulCreate  14m   replicaset-controller  Created pod: nginx-deployment-578fb949d8-6bfjf
        
      ###上述Events->Message创建两个副本,对象命名: 名称=父对象名称+随机字符串或数字
        
     ------------------------
     
     ###查看pod
     [root@node1 ~]# kubectl get pod
     NAME                                READY   STATUS    RESTARTS   AGE
     nginx-deployment-578fb949d8-4h4sg   1/1     Running   0          23m
     nginx-deployment-578fb949d8-6bfjf   1/1     Running   0          23m
     
     
     ###查看更详细信息
     [root@node1 ~]# kubectl describe pod
     Name:               nginx-deployment-578fb949d8-4h4sg
     Namespace:          default
     Priority:           0
     PriorityClassName:  <none>
     Node:               node2/192.168.1.32
     Start Time:         Fri, 10 May 2019 01:44:16 +0000
     Labels:             pod-template-hash=578fb949d8
                         run=nginx-deployment
     Annotations:        cni.projectcalico.org/podIP: 10.244.104.6/32
     Status:             Running
     IP:                 10.244.104.6
     Controlled By:      ReplicaSet/nginx-deployment-578fb949d8 ###指明此Pod是由ReplicaSet/nginx-deployment-578fb949d8创建的
     Containers:
       nginx-deployment:
         Container ID:   docker://78b084d9b37b0c7b41851d2564c3265d0ea9e4d20abad03c1587d8e71c855a66
         Image:          nginx:1.7.9
         Image ID:       docker-pullable://nginx@sha256:e3456c851a152494c3e4ff5fcc26f240206abac0c9d794affb40e0714846c451
         Port:           <none>
         Host Port:      <none>
         State:          Running
           Started:      Fri, 10 May 2019 01:46:48 +0000
         Ready:          True
         Restart Count:  0
         Environment:    <none>
         Mounts:
           /var/run/secrets/kubernetes.io/serviceaccount from default-token-gtn9w (ro)
     Conditions:
       Type              Status
       Initialized       True 
       Ready             True 
       ContainersReady   True 
       PodScheduled      True 
     Volumes:
       default-token-gtn9w:
         Type:        Secret (a volume populated by a Secret)
         SecretName:  default-token-gtn9w
         Optional:    false
     QoS Class:       BestEffort
     Node-Selectors:  <none>
     Tolerations:     node.kubernetes.io/not-ready:NoExecute for 300s
                      node.kubernetes.io/unreachable:NoExecute for 300s
     Events:
       Type     Reason          Age                From               Message
       ----     ------          ----               ----               -------
       Normal   Scheduled       24m                default-scheduler  Successfully assigned default/nginx-deployment-578fb949d8-4h4sg to node2
       Warning  Failed          23m                kubelet, node2     Failed to pull image "nginx:1.7.9": rpc error: code = Unknown desc = context canceled
       Warning  Failed          23m                kubelet, node2     Error: ErrImagePull
       Normal   SandboxChanged  23m                kubelet, node2     Pod sandbox changed, it will be killed and re-created.
       Normal   BackOff         23m (x3 over 23m)  kubelet, node2     Back-off pulling image "nginx:1.7.9"
       Warning  Failed          23m (x3 over 23m)  kubelet, node2     Error: ImagePullBackOff
       Normal   Pulling         22m (x2 over 24m)  kubelet, node2     pulling image "nginx:1.7.9"
       Normal   Pulled          22m                kubelet, node2     Successfully pulled image "nginx:1.7.9"
       Normal   Created         22m                kubelet, node2     Created container
       Normal   Started         22m                kubelet, node2     Started container
     
     
     Name:               nginx-deployment-578fb949d8-6bfjf
     Namespace:          default
     Priority:           0
     PriorityClassName:  <none>
     Node:               node3/192.168.1.33
     Start Time:         Fri, 10 May 2019 01:44:16 +0000
     Labels:             pod-template-hash=578fb949d8
                         run=nginx-deployment
     Annotations:        cni.projectcalico.org/podIP: 10.244.135.1/32
     Status:             Running
     IP:                 10.244.135.1
     Controlled By:      ReplicaSet/nginx-deployment-578fb949d8
     Containers:
       nginx-deployment:
         Container ID:   docker://cf0aa2d9467a653879a2e4f795463e454f7e9a676d018e44a539ac3e985fb46a
         Image:          nginx:1.7.9
         Image ID:       docker-pullable://nginx@sha256:e3456c851a152494c3e4ff5fcc26f240206abac0c9d794affb40e0714846c451
         Port:           <none>
         Host Port:      <none>
         State:          Running
           Started:      Fri, 10 May 2019 01:45:17 +0000
         Ready:          True
         Restart Count:  0
         Environment:    <none>
         Mounts:
           /var/run/secrets/kubernetes.io/serviceaccount from default-token-gtn9w (ro)
     Conditions:
       Type              Status
       Initialized       True 
       Ready             True 
       ContainersReady   True 
       PodScheduled      True 
     Volumes:
       default-token-gtn9w:
         Type:        Secret (a volume populated by a Secret)
         SecretName:  default-token-gtn9w
         Optional:    false
     QoS Class:       BestEffort
     Node-Selectors:  <none>
     Tolerations:     node.kubernetes.io/not-ready:NoExecute for 300s
                      node.kubernetes.io/unreachable:NoExecute for 300s
     Events: ###日志,记录Pod启动过程
       Type    Reason     Age   From               Message
       ----    ------     ----  ----               -------
       Normal  Scheduled  24m   default-scheduler  Successfully assigned default/nginx-deployment-578fb949d8-6bfjf to node3
       Normal  Pulling    24m   kubelet, node3     pulling image "nginx:1.7.9"
       Normal  Pulled     23m   kubelet, node3     Successfully pulled image "nginx:1.7.9"
       Normal  Created    23m   kubelet, node3     Created container
       Normal  Started    23m   kubelet, node3     Started container
       
      ------------------------
      
      ###上述创建nginx-deployment默认在default命名空间下
      [root@node1 nginx]# kubectl get deployments -n default
      NAME               READY   UP-TO-DATE   AVAILABLE   AGE
      nginx-deployment   2/2     2            2           44m
       
      ------------------------
      
      ###总结：
      
      A.用户通过kubectl创建Deployment
      B.Deployment创建ReplicaSet
      C.ReplicaSet创建Pod
     
      kubectl->Deployment->ReplicaSet->Pod
      

##### 1-2.命令 VS 配置文件

k8s支持两种创建资源方式：

1.kubectl命令直接创建

简单，适合临时测试和实验
  
    ###命令行中通过参数指定资源的属性
    $ kubectl run nginx-deployment --image=nginx:1.7.9 --replicas=2
    
    ###删除命令创建的deployment
    [root@node1 nginx]# kubectl delete deployment/nginx-deployment
    deployment.extensions "nginx-deployment" deleted
    
    ###查看应用
    [root@node1 nginx]# kubectl get deployment
    No resources found.

2.通过配置文件/kubectl apply创建

配置文件提供了创建资源模版，能够重复部署，方便管理，适合正式，跨环境，规模化部署

    Step1:创建nginx.yaml,资源的属性写在配置文件中，文件格式为YAML
    [root@node1 nginx]# cat nginx.yaml 
    apiVersion: extensions/v1beta1   ###当前配置格式的版本
    kind: Deployment    ###要创建的资源类型
    metadata: ###该资源的元数据
      name: nginx-deployment     ###必需的元数据项
      namespace: kube-public     ###指定命名空间，不指定会应用默认default命名空间
    spec:     ###该deployment的规格说明
      replicas: 2    ###指明副本数量，默认为1
      template:      ###定义pod的模版，此处为配置文件重要部分
        metadata:    ###定义pod的元数据，至少定义一个label
          labels:
            app: web_server ###label的key和value任意指定
        spec:        ###描述pod规格,此部分定义pod中每一个容器的属性
          containers:
          - name: nginx  ###必需
            image: nginx:1.7.9   ###必需
            
    ----------------------------------------------------------
    
    Step2:$ kubectl apply -f nginx.yaml  ###配置文件方式启动应用
    
    [root@node1 nginx]# kubectl apply -f nginx.yaml 
    deployment.extensions/nginx-deployment created
    
    [root@node1 nginx]# kubectl get deployment -n kube-public
    NAME               READY   UP-TO-DATE   AVAILABLE   AGE
    nginx-deployment   2/2     2            2           76m
    
    [root@node1 nginx]# kubectl get replicaset -n kube-public
    NAME                          DESIRED   CURRENT   READY   AGE
    nginx-deployment-65998d8886   2         2         2       76m
    
    [root@node1 nginx]# kubectl get pod -o wide -n kube-public
    NAME                                READY   STATUS    RESTARTS   AGE   IP             NODE    NOMINATED NODE   READINESS GATES
    nginx-deployment-65998d8886-hb7rk   1/1     Running   0          77m   10.244.135.2   node3   <none>           <none>
    nginx-deployment-65998d8886-zfvc8   1/1     Running   0          77m   10.244.104.7   node2   <none>           <none>

    ###以上，Deployment,ReplicaSet,Pod均已就绪
    
**kubectl apply:创建，更新资源**

    kubectl apply不但能创建k8s资源，也可对k8s资源更新
    
    k8s也提供了几个类似命令：kubectl create / kubectl replace / kubectl edit / kubectl patch  
    
    以上，尽量应用kubectl apply
    
**kubectl delete:删除资源**

    [root@node1 nginx]# kubectl delete -f nginx.yaml 
    deployment.extensions "nginx-deployment" deleted

    [root@node1 nginx]# kubectl get deployment -n kube-public
    No resources found.

##### 1-3.伸缩:在线增加/减少Pod副本数

    ### Step0.修改nginx.yaml，应用默认namespace且副本数量修改为5(2个副本时，在node2和node3上各跑一个副本)
    [root@node1 nginx]# cat nginx.yaml 
    apiVersion: extensions/v1beta1
    kind: Deployment
    metadata:
      name: nginx-deployment
    spec:
      replicas: 5
      template:
        metadata:
          labels:
            app: web_server
        spec:
          containers:
          - name: nginx
            image: nginx:1.7.9
            
    ### Step1.创建/更新资源
    [root@node1 nginx]# kubectl apply -f nginx.yaml 
    deployment.extensions/nginx-deployment created
    
    ### Step2.查看pod,增加3个副本被调度到node2和node3
    [root@node1 nginx]# kubectl get pod -o wide 
    NAME                                READY   STATUS    RESTARTS   AGE   IP              NODE    NOMINATED NODE   READINESS GATES
    nginx-deployment-65998d8886-c7r4j   1/1     Running   0          53s   10.244.135.3    node3   <none>           <none>
    nginx-deployment-65998d8886-lp45t   1/1     Running   0          53s   10.244.104.10   node2   <none>           <none>
    nginx-deployment-65998d8886-nq7k2   1/1     Running   0          53s   10.244.135.4    node3   <none>           <none>
    nginx-deployment-65998d8886-wcs9r   1/1     Running   0          53s   10.244.104.8    node2   <none>           <none>
    nginx-deployment-65998d8886-xk6wk   1/1     Running   0          53s   10.244.104.9    node2   <none>           <none>

    ### Step3.减少副本数量为3，再更新资源，查看pod
    [root@node1 nginx]# kubectl apply -f nginx.yaml 
    deployment.extensions/nginx-deployment configured

    [root@node1 nginx]# kubectl get pod -o wide
    NAME                                READY   STATUS    RESTARTS   AGE     IP              NODE    NOMINATED NODE   READINESS GATES
    nginx-deployment-65998d8886-lp45t   1/1     Running   0          5m16s   10.244.104.10   node2   <none>           <none>
    nginx-deployment-65998d8886-wcs9r   1/1     Running   0          5m16s   10.244.104.8    node2   <none>           <none>
    nginx-deployment-65998d8886-xk6wk   1/1     Running   0          5m16s   10.244.104.9    node2   <none>           <none>
    
**注：安全考虑，默认配置下k8s不会将Pod调度到Master节点，若将k8s-master也当作Node使用，则**

    #Master作为Node节点应用
    $ kubectl taint node node1 node-role.kubernetes.io/master -
    
    #恢复Master Only状态
    $ kubectl taint node node1 node-role.kubernetes.io/master="":NoSchedule
    
##### 1-4.Failover 

    ###模拟机器故障，如node2节点故障(关闭node2)
    [root@node1 nginx]# kubectl get node
    NAME    STATUS     ROLES    AGE   VERSION
    node1   Ready      master   42d   v1.13.3
    node2   NotReady   <none>   42d   v1.13.3   ###node2节点被关闭了
    node3   Ready      <none>   42d   v1.13.3
    
    ###查看pod信息,此时pod被调度到node3运行
    [root@node1 nginx]# kubectl get pod -o wide
    NAME                                READY   STATUS        RESTARTS   AGE   IP             NODE    NOMINATED NODE   READINESS GATES
    nginx-deployment-65998d8886-br7cv   1/1     Running       0          65s   10.244.135.7   node3   <none>           <none>
    nginx-deployment-65998d8886-lp45t   0/1     Terminating   0          83m   <none>         node2   <none>           <none>
    nginx-deployment-65998d8886-tdx8l   1/1     Running       0          65s   10.244.135.6   node3   <none>           <none>
    nginx-deployment-65998d8886-w6jfv   1/1     Running       0          65s   10.244.135.5   node3   <none>           <none>
    nginx-deployment-65998d8886-wcs9r   0/1     Terminating   0          83m   <none>         node2   <none>           <none>
    nginx-deployment-65998d8886-xk6wk   0/1     Terminating   0          83m   <none>         node2   <none>           <none>
    
    ###重启node2后，会删除Terminating状态的pod,不过已经运行的pod不会重新调度回node2
    [root@node1 nginx]# kubectl get pod -o wide
    NAME                                READY   STATUS    RESTARTS   AGE    IP             NODE    NOMINATED NODE   READINESS GATES
    nginx-deployment-65998d8886-br7cv   1/1     Running   0          2m7s   10.244.135.7   node3   <none>           <none>
    nginx-deployment-65998d8886-tdx8l   1/1     Running   0          2m7s   10.244.135.6   node3   <none>           <none>
    nginx-deployment-65998d8886-w6jfv   1/1     Running   0          2m7s   10.244.135.5   node3   <none>           <none>
    
    ###删除nginx-deployment
    [root@node1 nginx]# kubectl delete deployment nginx-deployment
    deployment.extensions "nginx-deployment" deleted
    
##### 1-5.用label控制pod的位置

    Scheduler会将Pod调度到所有可用的Node.
    
    但有时希望将Pod部署到指定Node,如：将大量磁盘I/O的Pod部署到配置了SSD的Node 或 Pod需要GPU,需要运行在配置了GPU的节点上
    
    以上：k8s通过label实现该功能
    
    -----------
    
    label是key-value对，各种资源均可设置label,灵活添加各种自定义属性
    
    ###标注node3是配置了SSD的节点
    [root@node1 nginx]# kubectl label node node3 disktype=ssd
    node/node3 labeled
    
    ###查看节点的label
    [root@node1 nginx]# kubectl get node --show-labels
    NAME    STATUS   ROLES    AGE   VERSION   LABELS
    node1   Ready    master   42d   v1.13.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/hostname=node1,node-role.kubernetes.io/master=
    node2   Ready    <none>   42d   v1.13.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/hostname=node2
    node3   Ready    <none>   42d   v1.13.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,disktype=ssd,kubernetes.io/hostname=node3  ###自定义disktype=ssd的label
    
    ###设置disktype这个自定义label，即可指定将pod部署到node3
    ###编辑nginx.yaml
    [root@node1 nginx]# cat nginx.yaml
    apiVersion: extensions/v1beta1
    kind: Deployment
    metadata:
      name: nginx-deployment
    spec:
      replicas: 3
      template:
        metadata:
          labels:
            app: web_server
        spec:
          containers:
          - name: nginx
            image: nginx:1.7.9
          nodeSelector:     ###增加，指定将此pod部署到具有label disktype=ssd的node上
            disktype: ssd
            
    ###部署deployment并查看pod运行节点(全部部署到了node3)
    [root@node1 nginx]# kubectl apply -f nginx.yaml 
    deployment.extensions/nginx-deployment created
    [root@node1 nginx]# kubectl get deployment
    NAME               READY   UP-TO-DATE   AVAILABLE   AGE
    nginx-deployment   3/3     3            3           9s

    [root@node1 nginx]# kubectl get pod -o wide
    NAME                                READY   STATUS    RESTARTS   AGE   IP              NODE    NOMINATED NODE   READINESS GATES
    nginx-deployment-7c75d8cdf6-mftr2   1/1     Running   0          19s   10.244.135.8    node3   <none>           <none>
    nginx-deployment-7c75d8cdf6-td57j   1/1     Running   0          19s   10.244.135.10   node3   <none>           <none>
    nginx-deployment-7c75d8cdf6-vbxnq   1/1     Running   0          19s   10.244.135.9    node3   <none>           <none>
    
    
    ###删除label disktype
    [root@node1 nginx]# kubectl label node node3 disktype-
    node/node3 labeled
    [root@node1 nginx]# kubectl get node --show-labels
    NAME    STATUS   ROLES    AGE   VERSION   LABELS
    node1   Ready    master   42d   v1.13.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/hostname=node1,node-role.kubernetes.io/master=
    node2   Ready    <none>   42d   v1.13.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/hostname=node2
    node3   Ready    <none>   42d   v1.13.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/hostname=node3  ###自定义disktype=ssd已经被删除
    [root@node1 nginx]# kubectl get pod -o wide
    NAME                                READY   STATUS    RESTARTS   AGE   IP              NODE    NOMINATED NODE   READINESS GATES
    nginx-deployment-7c75d8cdf6-mftr2   1/1     Running   0          99s   10.244.135.8    node3   <none>           <none>
    nginx-deployment-7c75d8cdf6-td57j   1/1     Running   0          99s   10.244.135.10   node3   <none>           <none>
    nginx-deployment-7c75d8cdf6-vbxnq   1/1     Running   0          99s   10.244.135.9    node3   <none>           <none>
    
    注：
    
    上述删除label disktype,pod并不会重新部署，依然在node3节点上运行
    
    除非nginx.yaml中删除nodeSelector设置，然后通过kubectl apply重新部署。k8s会删除之前的pod并调度和运行新pod
    
    ###删除nodeSelector,重新部署
    [root@node1 nginx]# kubectl apply -f nginx.yaml 
    deployment.extensions/nginx-deployment configured
    [root@node1 nginx]# kubectl get pod -o wide
    NAME                                READY   STATUS    RESTARTS   AGE   IP              NODE    NOMINATED NODE   READINESS GATES
    nginx-deployment-65998d8886-h5482   1/1     Running   0          17s   10.244.104.11   node2   <none>           <none>
    nginx-deployment-65998d8886-jz5rq   1/1     Running   0          17s   10.244.135.11   node3   <none>           <none>
    nginx-deployment-65998d8886-l9blg   1/1     Running   0          16s   10.244.104.12   node2   <none>           <none>

#### 2。DaemonSet    

    Deployment部署的副本Pod会分布在各个Node上，每个Node都可能运行好几个副本

    DaemonSet不同处在于：每个Node上最多只能运行一个副本
    
>DaemonSet典型应用场景

    (1).在集群的每个节点上运行存储Daemon,如: glusterd / ceph
    (2).在每个节点上运行日志收集Daemon,如: flunentd / logstash
    (3).在每个节点上运行监控Daemon,如: prometheus Node Exporter / collectd
    
>k8s应用DaemonSet运行系统组件

    ###若不指定namespace,则默认namespace为default
    [root@node1 nginx]# kubectl get daemonset --namespace=kube-system
    NAME          DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR                 AGE
    calico-node   3         3         3       3            3           beta.kubernetes.io/os=linux   42d
    kube-proxy    3         3         3       3            3           <none>                        42d
    
    ###上述calico-node和kube-proxy分别负责在每个节点上运行calico和kube-proxy组件
    
    ### calico网络，也可以是flannel网络(常用)
    
    ### kube-proxy组件，用于将service代理映射到pod(外界访问->service->kube-proxy->pod)
    
##### 2-1。kube-flannel-ds

>部署flannel网络

    ###部署flannel网络
    [root@node1 k8s]# kubectl apply -f  https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
    podsecuritypolicy.extensions/psp.flannel.unprivileged created
    clusterrole.rbac.authorization.k8s.io/flannel created
    clusterrolebinding.rbac.authorization.k8s.io/flannel created
    serviceaccount/flannel created
    configmap/kube-flannel-cfg created
    daemonset.extensions/kube-flannel-ds-amd64 created
    daemonset.extensions/kube-flannel-ds-arm64 created
    daemonset.extensions/kube-flannel-ds-arm created
    daemonset.extensions/kube-flannel-ds-ppc64le created
    daemonset.extensions/kube-flannel-ds-s390x created
    
    ###查看k8s应用DaemonSet运行的系统组件
    [root@node1 k8s]# kubectl get daemonset --namespace=kube-system
    NAME                      DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR                     AGE
    calico-node               3         3         3       3            3           beta.kubernetes.io/os=linux       42d
    kube-flannel-ds-amd64     3         3         3       3            3           beta.kubernetes.io/arch=amd64     4m26s
    kube-flannel-ds-arm       0         0         0       0            0           beta.kubernetes.io/arch=arm       4m26s
    kube-flannel-ds-arm64     0         0         0       0            0           beta.kubernetes.io/arch=arm64     4m26s
    kube-flannel-ds-ppc64le   0         0         0       0            0           beta.kubernetes.io/arch=ppc64le   4m26s
    kube-flannel-ds-s390x     0         0         0       0            0           beta.kubernetes.io/arch=s390x     4m26s
    kube-proxy                3         3         3       3            3           <none>                            42d
    
    ###以上：flannel的DaemonSet即定义在kube-flannel.yaml中
    $ wget https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
    
[kube-flannel.yaml](resources/kube-flannel.yml)

    ###截取一小部分说明。DaemonSet的语法结构同Deployment。如：
    ---
    apiVersion: extensions/v1beta1
    kind: DaemonSet  ###指定kind为DaemonSet
    metadata:
      name: kube-flannel-ds-amd64
      namespace: kube-system
      labels:
        tier: node
        app: flannel
    spec:
      template:
        metadata:
          labels:
            tier: node
            app: flannel
        spec:
          hostNetwork: true  ###指定hostNetwork,表指定Pod直接使用的是Node网络，相当于docker run --network=host
          nodeSelector:
            beta.kubernetes.io/arch: amd64
          tolerations:
          - operator: Exists
            effect: NoSchedule
          serviceAccountName: flannel
          initContainers:
          - name: install-cni
            image: quay.io/coreos/flannel:v0.11.0-amd64
            command:
            - cp
            args:
            - -f
            - /etc/kube-flannel/cni-conf.json
            - /etc/cni/net.d/10-flannel.conflist
            volumeMounts:
            - name: cni
              mountPath: /etc/cni/net.d
            - name: flannel-cfg
              mountPath: /etc/kube-flannel/
          containers:    ###定义了运行flannel服务的容器
          - name: kube-flannel
            image: quay.io/coreos/flannel:v0.11.0-amd64
    
##### 2-2。kube-proxy

    ###由于无法拿到kube-proxy的YAML文件，只能运行下述命令查看配置：
    $ kubectl edit daemonset kube-proxy --namespace=kube-system

[kube-proxy.yaml](resources/kube-proxy.yaml)

##### 2-3。运行自己的DaemonSet

 >以Prometheus Node Exporter为例演示用户如何运行自己的DaemonSet
 
    Prometheus是流行的系统监控方案
    
    Node Exporter是Prometheus的agent,以Daemon的形式运行在每个被监控节点上
    
    ###方式一：直接在Docker中运行Node Exporter容器
    [root@node1 k8s]# docker run -d \
    > -v "/proc:/host/proc" \
    > -v "/sys:/host/sys" \
    > -v "/:/rootfs" \
    > --net=host \
    > prom/node-exporter \
    > --path.procfs /host/proc \
    > --path.sysfs /host/sys \
    > --collector.filesystem.ignored-mount-points "^/(sys|proc|dev|host|etc)($|/)"
    
    ###方式二：转换为DaemonSet的YAML配置文件node_exporter.yml
    
[node_exporter.yaml](resources/node_exporter.yml)   

    ###部署应用
    [root@node1 k8s]# kubectl apply -f node_exporter.yml 
    daemonset.extensions/node-exporter-daemonset created
    
    ###查看Pod信息,node2和node3上分别运行了1个pod
    [root@node1 k8s]# kubectl get pod -o wide | grep node-exporter
    node-exporter-daemonset-2xhzq       1/1     Running   0          3m45s   192.168.1.33    node3   <none>           <none>
    node-exporter-daemonset-h299s       1/1     Running   0          3m45s   192.168.1.32    node2   <none>           <none>

#### 3。Job

容器按运行时间分为两类：

>服务类容器:持续提供服务，需一直运行(如:Http Server/Daemon等)

    k8s的Deployment,ReplicaSet,DaemonSet都用于管理服务类容器

>工作类容器:一次性任务(如：批处理程序，完成后容器就退出)

    Job管理工作类容器
    
创建简是单的Job配置文件[myjob.yml](resources/myjob.yml)

    [root@node1 k8s]# cat myjob.yml 
    apiVersion: batch/v1  #当前Job的apiVersion
    kind: Job #指明当前资源的类型为Job
    metadata:
      name: myjob
    spec:
      template:
        metadata:
          name: myjob
        spec:
          containers:
          - name: hello
            image: busybox
            command: ["echo","hello k8s Job!"]
          restartPolicy: Never #指定什么情况需重启容器，对于Job只能设置Never/OnFailure.对于其他Controller(如Deployment),可设置为Always

    ###启动Job
    [root@node1 k8s]# kubectl apply -f myjob.yml 
    job.batch/myjob created
    
    ###查看Job状态(按预期启动了一个Pod，并已成功执行)
    [root@node1 k8s]# kubectl get job
    NAME    COMPLETIONS   DURATION   AGE
    myjob   1/1           14s        6m20s
    
    ###查看Pod状态
    [root@node1 k8s]# kubectl get pod
    NAME                                READY   STATUS      RESTARTS   AGE
    myjob-k7sfw                         0/1     Completed   0          7m27s
    
    ###查看Pod的标准输出
    [root@node1 k8s]# kubectl logs myjob-k7sfw
    hello k8s Job!
    
##### 3-1.Pod失败情况

>验证1：myjob.yml/restartPolicy为Never

    ###删除job
    [root@node1 k8s]# kubectl delete -f myjob.yml 
    job.batch "myjob" deleted
   
    ###模拟创建Pod失败，修改myjob.yml中
    command: ["invalide command","hello k8s Job!"]  ##使其异常
    
    ###运行新的Job并查看状态
    [root@node1 k8s]# kubectl get job
    NAME    COMPLETIONS   DURATION   AGE
    myjob   0/1           13s        13s
    
    ###此时，COMPLETIONS的Pod数量为0，查看Pod状态,状态均不正常
    [root@node1 k8s]# kubectl get pod -o wide
    NAME                                READY   STATUS               RESTARTS   AGE    IP              NODE    NOMINATED NODE   READINESS GATES
    myjob-g877f                         0/1     ContainerCreating    0          13s    <none>          node3   <none>           <none>
    myjob-gxk4v                         0/1     ContainerCannotRun   0          29s    10.244.135.13   node3   <none>           <none>

    ###查看某个pod的启动日志
    $ kubectl describe pod myjob-gxk4v
    
kubectl get pod -o wide会查询多个失败Pod?

    当第一个Pod启动时，容器失败退出。
    
    根据restartPolicy: Never，此时失败容器不会被重启。但Job DESIRED的Pod是1.当前COMPLETIONS为0不满足。所以Job Controller会启动新的Pod，直到COMPLETIONS为1
    
若设置restartPolicy: OnFailure?

    [root@node1 k8s]# kubectl apply -f myjob.yml 
    job.batch/myjob created
    
    [root@node1 k8s]# kubectl get job
    NAME    COMPLETIONS   DURATION   AGE
    myjob   0/1           4s         4s      ###COMPLETIONS数量为0
    
    [root@node1 k8s]# kubectl get pod 
    NAME                                READY   STATUS              RESTARTS   AGE
    myjob-ltnrm                         0/1     RunContainerError   2          75s  ##此时RESTARTS数量为2
    
    ###上述说明：OnFailure配置生效，容器失败后会重启，重启次数为2
    
    ###删除job
    [root@node1 k8s]# kubectl delete -f myjob.yml 
    job.batch "myjob" deleted

##### 3-2.Job的并行性

>希望同时运行多个Pod，提高Job的执行效率 

    [root@node1 k8s]# cat myjob.yml 
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: myjob
    spec:
      completions: 2  #每次并行执行2个pod,直到总数有2个(可任意设置)pod成功完成.默认1
      parallelism: 2  #每次并行执行2个Pod,默认1
      template:
        metadata:
          name: myjob
        spec:
          containers:
          - name: hello
            image: busybox
            command: ["echo","hello k8s Job!"]
          restartPolicy: OnFailure

    ###查看job
    [root@node1 k8s]# kubectl get job
    NAME    COMPLETIONS   DURATION   AGE
    myjob   2/2           16s        17s
    
    ###查看pod,AGE一样，说明是并行执行
    [root@node1 k8s]# kubectl get pod
    NAME                                READY   STATUS      RESTARTS   AGE
    myjob-t8sfj                         0/1     Completed   0          21s
    myjob-zlfq2                         0/1     Completed   0          21s
    
    ###需要并行处理场景：批处理程序，每个副本(Pod)都会从任务池中读取任务并执行，副本越多，执行时间越短，效率越高(均可用job执行)

##### 3-3.定时Job
    
>linux中有cron程序定时执行任务，k8s的CronJob提供了类似功能，可定时执行Job

    #定时Job
    [root@node1 k8s]# cat mycronjob.yml 
    apiVersion: batch/v2alpha1  #batch/v2alpha1当前CronJob的apiVersion
    kind: CronJob #指明当前资源类型为CronJob
    metadata:
      name: hello
    spec:
      schedule: "*/1 * * * *" #指定什么时候运行Job,格式与linux cron一致,每分钟启动一次
      jobTemplate: #定义job的模版
        spec:
          template:
            spec:
              containers:
              - name: hello
                image: busybox
                command: ["echo","hello k8s CronJob!"]
              restartPolicy: OnFailure
    
    
    ###创建CronJob
    [root@node1 k8s]# kubectl apply -f mycronjob.yml 
    error: unable to recognize "mycronjob.yml": no matches for kind "CronJob" in version "batch/v2alpha1"
    
    ###此时报错：因为k8s默认没有enable CronJob功能，需在kube-apiserver中加入这个功能
    [root@node1 k8s]# vim /etc/kubernetes/manifests/kube-apiserver.yaml 
    - kube-apiserver
    - --runtime-config=batch/v2alpha1=true #add此行
    
    ###重启kubelet服务,此时kubelet会重启kube-apiservice Pod
    [root@node1 k8s]# systemctl restart kubelet.service
    
    ###通过kubectl api-versions确认kube-apiserver现已支持batch/v2alpha1
    [root@node1 k8s]# kubectl api-versions
    admissionregistration.k8s.io/v1beta1
    apiextensions.k8s.io/v1beta1
    apiregistration.k8s.io/v1
    apiregistration.k8s.io/v1beta1
    apps/v1
    apps/v1beta1
    apps/v1beta2
    authentication.k8s.io/v1
    authentication.k8s.io/v1beta1
    authorization.k8s.io/v1
    authorization.k8s.io/v1beta1
    autoscaling/v1
    autoscaling/v2beta1
    autoscaling/v2beta2
    batch/v1
    batch/v1beta1
    batch/v2alpha1   ###支持
    certificates.k8s.io/v1beta1
    coordination.k8s.io/v1beta1
    crd.projectcalico.org/v1
    events.k8s.io/v1beta1
    extensions/v1beta1
    networking.k8s.io/v1
    policy/v1beta1
    rbac.authorization.k8s.io/v1
    rbac.authorization.k8s.io/v1beta1
    scheduling.k8s.io/v1beta1
    storage.k8s.io/v1
    storage.k8s.io/v1beta1
    v1
    
    ###再次创建CronJob
    [root@node1 k8s]# kubectl apply -f mycronjob.yml 
    cronjob.batch/hello created
    
    ###查看cronjob状态
    [root@node1 k8s]# kubectl get cronjob
    NAME    SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
    hello   */1 * * * *   False     0        <none>          12s
    
    ###查看job执行情况,每分钟会启动一个Job
    [root@node1 k8s]# kubectl get jobs
    NAME               COMPLETIONS   DURATION   AGE
    hello-1557481620   1/1           10s        2m28s
    hello-1557481680   1/1           10s        88s
    hello-1557481740   1/1           10s        28s
    
    ###查看pod
    [root@node1 k8s]# kubectl get pod
    NAME                                READY   STATUS              RESTARTS   AGE
    hello-1557481740-268mk              0/1     Completed           0          3m4s
    hello-1557481800-g9drk              0/1     Completed           0          2m4s
    
    ###查看pod运行日志
    [root@node1 k8s]# kubectl logs hello-1557481800-g9drk
    hello k8s CronJob!
    
    ###删除cronjob
    [root@node1 k8s]# kubectl delete -f mycronjob.yml 
    cronjob.batch "hello" deleted
    
-------------------------------------------------

### 第6章 通过service访问pod

>每个pod均有自己的IP地址，当动态创建和销毁Pod时，IP会改变，如何保证客户端访问并找到对应服务--k8s Service

#### 1.创建Service

    k8s从逻辑上代表了一组Pod,具体是哪些Pod，则是由label来挑选
    
    service有自己IP,且该IP不变
    
    客户端只需访问Service的IP,k8s负责建立和维护Service与Pod的映射关系

>Step1:创建Deployment

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
            run: httpd ###指定label,下述创建service会指定挑选该label的pod作为service的后端
        spec:
          containers:
          - name: httpd
            image: httpd
            ports:
            - containerPort: 80
            
    [root@node1 k8s]# kubectl apply -f httpd.yml 
    deployment.apps/httpd created
    
    httpd-8c6c4bd9b-v8hq7   0/1     ContainerCreating   0          2m13s   <none>          node3   <none>           <none>
    [root@node1 nginx]# kubectl get pod -o wide
    NAME                    READY   STATUS    RESTARTS   AGE     IP              NODE    NOMINATED NODE   READINESS GATES
    httpd-8c6c4bd9b-49x55   1/1     Running   0          2m23s   10.244.135.51   node3   <none>           <none>
    httpd-8c6c4bd9b-qztdj   1/1     Running   0          2m23s   10.244.104.15   node2   <none>           <none>
    httpd-8c6c4bd9b-v8hq7   1/1     Running   0          2m23s   10.244.135.52   node3   <none>           <none>
    
    ###Pod分配了各自的I，这些IP只能被k8s cluster中的容器和节点访问
    [root@node1 nginx]# curl 10.244.135.51
    <html><body><h1>It works!</h1></body></html>

>Step2:创建Service

    ###创建service的配置文件
    [root@node1 k8s]# cat httpd-service.yml 
    apiVersion: v1     ###v1是Service的apiVersion
    kind: Service      ###指明当前资源类型为Service
    metadata: 
      name: httpd-svc  ###service名字
    spec:
      selector:        ###指明挑选label为run:httpd的Pod作为Service的后端
        run: httpd
      ports:           ###将Service的8080端口映射到Pod的80端口,使用TCP协议
      - protocol: TCP
        port: 8080       
        targetPort: 80
        
    ###创建service
    [root@node1 k8s]# kubectl apply -f httpd-service.yml 
    service/httpd-svc created
    
    ###查询service
    [root@node1 k8s]# kubectl get service
    NAME         TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)    AGE
    httpd-svc    ClusterIP   10.1.82.198   <none>        8080/TCP   6s
    kubernetes   ClusterIP   10.1.0.1      <none>        443/TCP    42d
    
    ###httpd-svc分配到一个cluster-ip 10.1.82.198,通过该ip访问后端httpd pod
    [root@node1 k8s]# curl 10.1.82.198:8080
    <html><body><h1>It works!</h1></body></html>
    
    ###除自行创建的httpd-svc,还有一个service(kubernetes).cluster内部通过这个Service访问k8s API Server
    
    ###查看httpd-svc与pod对应关系
    [root@node1 k8s]# kubectl describe service httpd-svc
    Name:              httpd-svc
    Namespace:         default
    Labels:            <none>
    Annotations:       kubectl.kubernetes.io/last-applied-configuration:
                         {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"name":"httpd-svc","namespace":"default"},"spec":{"ports":[{"port":8080,"...
    Selector:          run=httpd
    Type:              ClusterIP
    IP:                10.1.82.198
    Port:              <unset>  8080/TCP
    TargetPort:        80/TCP
    Endpoints:         10.244.104.15:80,10.244.135.51:80,10.244.135.52:80  ###列出了三个Pod的IP和端口
    Session Affinity:  None
    Events:            <none>
    
  **Pod的IP是在容器中配置的**
    
  **Service的cluster ip配置在哪里？cluster-ip是如何映射到pod ip? --iptables**
  
#### 2.Cluster IP底层实现

>Cluster IP是一个虚拟IP，是由K8s节点上的iptables规则管理的。

    ###先查询service和pod，方便与下述说明对应理解
    ###新创建的service(httpd-svc),关注：Cluster-ip值
    [root@node1 k8s]# kubectl get service
    NAME         TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)    AGE
    httpd-svc    ClusterIP   10.1.82.198   <none>        8080/TCP   72m  
    kubernetes   ClusterIP   10.1.0.1      <none>        443/TCP    42d
    
    ###新创建的httpd-pod，关注IP值
    [root@node1 k8s]# kubectl get pod -o wide
    NAME                    READY   STATUS    RESTARTS   AGE   IP              NODE    NOMINATED NODE   READINESS GATES
    httpd-8c6c4bd9b-49x55   1/1     Running   0          82m   10.244.135.51   node3   <none>           <none>
    httpd-8c6c4bd9b-qztdj   1/1     Running   0          82m   10.244.104.15   node2   <none>           <none>
    httpd-8c6c4bd9b-v8hq7   1/1     Running   0          82m   10.244.135.52   node3   <none>           <none>

    --------------------------------------------
    
    ###打印当前节点的iptables规则
    [root@node1 k8s]# iptables-save > iptables-rules.txt
    
    ###截取与httpd-svc Cluter IP相关信息
    -A KUBE-SERVICES ! -s 10.244.0.0/16 -d 10.1.82.198/32 -p tcp -m comment --comment "default/httpd-svc: cluster IP" -m tcp --dport 8080 -j KUBE-MARK-MASQ
    -A KUBE-SERVICES -d 10.1.82.198/32 -p tcp -m comment --comment "default/httpd-svc: cluster IP" -m tcp --dport 8080 -j KUBE-SVC-RL3JAE4GN7VOGDGP

    上述两条规则含义:
    (1).如果cluster内的Pod(源地址来自10.244.0.0/16)要访问httpd-svc,则允许
    (2).其他源地址访问httpd-svc,跳转到规则KUBE-SVC-RL3JAE4GN7VOGDGP
    
    
    ###KUBE-SVC-RL3JAE4GN7VOGDGP规则如下所示：
    -A KUBE-SVC-RL3JAE4GN7VOGDGP -m statistic --mode random --probability 0.33332999982 -j KUBE-SEP-Y5YLOO6DS3UF7LBB
    -A KUBE-SVC-RL3JAE4GN7VOGDGP -m statistic --mode random --probability 0.50000000000 -j KUBE-SEP-DTMMTXPGHKXJGUTZ
    -A KUBE-SVC-RL3JAE4GN7VOGDGP -j KUBE-SEP-FQRCOM4I5Y4H2JWZ
    
    ###上述三个跳转规则说明：
    (1).1/3的概率跳转到规则 KUBE-SEP-Y5YLOO6DS3UF7LBB
    (2).1/3的概率(剩下2/3的一半)跳转到规则KUBE-SEP-DTMMTXPGHKXJGUTZ
    (3).1/3的概率跳转到规则KUBE-SEP-FQRCOM4I5Y4H2JWZ
        
    ###上述三个跳转规则如下所示
    -A KUBE-SEP-Y5YLOO6DS3UF7LBB -s 10.244.104.15/32 -j KUBE-MARK-MASQ
    -A KUBE-SEP-Y5YLOO6DS3UF7LBB -p tcp -m tcp -j DNAT --to-destination 10.244.104.15:80
    -A KUBE-SEP-DTMMTXPGHKXJGUTZ -s 10.244.135.51/32 -j KUBE-MARK-MASQ
    -A KUBE-SEP-DTMMTXPGHKXJGUTZ -p tcp -m tcp -j DNAT --to-destination 10.244.135.51:80
    -A KUBE-SEP-FQRCOM4I5Y4H2JWZ -s 10.244.135.52/32 -j KUBE-MARK-MASQ
    -A KUBE-SEP-FQRCOM4I5Y4H2JWZ -p tcp -m tcp -j DNAT --to-destination 10.244.135.52:80
    
    ###上述即将请求分别转发到后端的三个pod(由$ kubectl get pod -o wide 可以查看到每个pod对应IP )
    
    总结：iptables将访问service的流量转发到后端pod,而且使用类似轮询的负载均衡策略
    
    另：Cluster的每一个节点都配置了相同的iptables规则，这样就确保了整个cluster都能够通过Service的clusterIP访问Service
    
**主要说明：service的ClusterIp如何与Pod的IP建立联系(iptables)**
    
#### 3.DNS访问Service

>在cluster中，除了可以通过ClusterIp访问Service,k8s还提供了更为方便的DNS访问

> kubeadm部署时会默认安装kube-dns组件

    [root@node1 k8s]# kubectl get deployment --namespace=kube-system
    NAME                      READY   UP-TO-DATE   AVAILABLE   AGE
    calico-kube-controllers   1/1     1            1           42d
    coredns                   2/2     2            2           42d
    
    coredns是一个DNS服务器。每当有新的Service被创建，kube-dns会添加该Service的DNS记录。
    
    Cluster中的Pod可以通过<SERVICE_NAME><NAMESPACE_NAME>访问service
    
    ###在一个临时busybox pod中验证DNS有效性
    [root@node1 ~]# kubectl run busybox --rm -ti --image=busybox /bin/sh
    kubectl run --generator=deployment/apps.v1 is DEPRECATED and will be removed in a future version. Use kubectl run --generator=run-pod/v1 or kubectl create instead.
    If you don't see a command prompt, try pressing enter.
    / # wget httpd-svc.default:8080
    Connecting to httpd-svc.default:8080 (10.1.82.198:8080)
    index.html           100% |*******************************************************************************************************|    45  0:00:00 ETA
    / # wget httpd-svc:8080   ###未连接成功(因defautl为默认ns，可省略)
    / # nslookup httd-svc     ###查看httpd-svc的DNS信息,后续未连接成功,正常能够看到DNS服务器,如httpd-svc.default.svc.cluster.local，它是httpd-svc的完整域名
    Server:		10.1.0.10
    Address:	10.1.0.10:53
    
 >示例：在kube-public的namespace中部署Service httpd2-svc
 
     ###YAML文件注意书写，否则yaml转json会有问题
     [root@node1 k8s]# cat httpd2-svc.yml 
     apiVersion: apps/v1beta1
     kind: Deployment  ##指定资源为Deployment
     metadata:
       name: httpd2
       namespace: kube-public  ##指定命名空间
     spec:
       replicas: 3
       template:
         metadata:
           labels:   ##指定labels，service根据labels决定哪一组pod绑定为一个service
             run: httpd2
         spec:
           containers:
           - name: httpd2
             image: httpd
             ports:
             - containerPort: 80
     
     ---   ##多个资源定义在一个yaml中，可用该分割符分割
     apiVersion: v1
     kind: Service  ##指定资源类型为Service
     metadata:
       name: httpd2-svc
       namespace: kube-public
     spec:
       selector:    ##指定label，将标识有该label标签绑定到同一个service
         run: httpd2
       ports:
       - protocol: TCP
         port: 8080
         targetPort: 80
         
      
      ###创建资源
      [root@node1 k8s]# kubectl apply -f httpd2-svc.yml 
      deployment.apps/httpd2 created
      service/httpd2-svc created
      
      ###查看kube-public命名空间下的service
      [root@node1 k8s]# kubectl get service --namespace=kube-public
      NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
      httpd2-svc   ClusterIP   10.1.178.235   <none>        8080/TCP   28s
      
      ###在busybox pod中访问httpd2-svc
      [root@node1 k8s]# kubectl run busybox --rm -ti --image=busybox /bin/sh
      kubectl run --generator=deployment/apps.v1 is DEPRECATED and will be removed in a future version. Use kubectl run --generator=run-pod/v1 or kubectl create instead.
      If you don't see a command prompt, try pressing enter.
      / # wget httpd2-svc.kube-public:8080   ##访问[服务名][命名空间][端口]
      Connecting to httpd2-svc.kube-public:8080 (10.1.178.235:8080)
      index.html           100% |*******************************************************************************************************|    45  0:00:00 ETA
      / # exit

#### 4.外网如何访问Service

    除了Cluster内部可访问Service，希望应用的Service能够暴露给Cluster外部。
    
>k8s提供了多种类型Service,默认ClusterIP

(1).ClusterIP(默认)

    Service通过Cluster内部的IP对外提供服务，只有Cluster内的节点和Pod可访问。上述Service都是ClusterIP.

(2).NodePort

    Service通过Cluster节点的静态端口对外提供服务。
    
    Cluster外部可通过<NodeIP>:<NodePort>访问Service  ###NodeIp本机节点IP
    
    示例：
    [root@node1 k8s]# cat httpd-service.yml 
    apiVersion: v1
    kind: Service
    metadata: 
      name: httpd-svc
    spec:
      type: NodePort  ##增加
      selector:
        run: httpd
      ports:
      - protocol: TCP
        port: 8080
        targetPort: 80
     
    ###查询未添加type: NodePort前的service
    [root@node1 k8s]# kubectl get service
    NAME         TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)    AGE
    httpd-svc    ClusterIP   10.1.82.198   <none>        8080/TCP   17h
    kubernetes   ClusterIP   10.1.0.1      <none>        443/TCP    43d
    
    ###重新创建Service资源
    [root@node1 k8s]# kubectl apply -f httpd-service.yml 
    service/httpd-svc configured
    ###查询service
    [root@node1 k8s]# kubectl get service
    NAME         TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)          AGE
    httpd-svc    NodePort    10.1.82.198   <none>        8080:31775/TCP   17h
    kubernetes   ClusterIP   10.1.0.1      <none>        443/TCP          43d
    
    ###上述说明
    (1)EXTERNAL-IP若为<nodes>表示可通过Cluster每个节点自身的IP访问Service
    (2)PORT(S)为8080:31775。其中8080是监听的端口，31775是节点上监听的端口(k8s会从30000~32767中分配一个可用的端口).每个节点均会监听此端口，并将请求转发给Service
    
    $ yum install net-tools
    ###查看监听的31775端口
    [root@node1 k8s]# netstat -lnp|grep 31775
    tcp6       0      0 :::31775                :::*                    LISTEN      8361/kube-proxy  
    
    ###测试NodePort是否正常工作[主机IP][端口]
    [root@node1 k8s]# curl 192.168.1.32:31775   
    <html><body><h1>It works!</h1></body></html>
    
    ###任意节点均可以成功访问
    [root@node2 ~]# curl 192.168.1.32:31775
    <html><body><h1>It works!</h1></body></html>
    
 >k8s是如何将<NodeIP>:<NodePort>映射到Pod?
 
    同ClusterIp，借助iptables
    
    ###查看iptables规则
    [root@node2 ~]# iptables-save > iptables-rules1.txt
    [root@node2 ~]# vim iptables-rules1.txt  ### 输入/httpd-svc查询
    
    与ClusterIP相比，每个节点的iptables中均增加了下述两条规则
    -A KUBE-NODEPORTS -p tcp -m comment --comment "default/httpd-svc:" -m tcp --dport 31775 -j KUBE-MARK-MASQ
    -A KUBE-NODEPORTS -p tcp -m comment --comment "default/httpd-svc:" -m tcp --dport 31775 -j KUBE-SVC-RL3JAE4GN7VOGDGP
    
    规则含义为：访问当前节点21775端口的请求会应用规则KUBE-SVC-RL3JAE4GN7VOGDGP
    
    KUBE-SVC-RL3JAE4GN7VOGDGP规则
    -A KUBE-SVC-RL3JAE4GN7VOGDGP -m statistic --mode random --probability 0.33332999982 -j KUBE-SEP-67JJISW5PAMK6276
    -A KUBE-SVC-RL3JAE4GN7VOGDGP -m statistic --mode random --probability 0.50000000000 -j KUBE-SEP-UASYEDCNCKK6XFPM
    -A KUBE-SVC-RL3JAE4GN7VOGDGP -j KUBE-SEP-DUUX3NU476QIDE2Z
    
    其作用就是负载均衡到每一个pod
    
>NodePort是默认随机选择，可用nodePort指定某个特定端口

    [root@node1 k8s]# cat httpd-service.yml 
    apiVersion: v1
    kind: Service
    metadata: 
      name: httpd-svc
    spec:
      type: NodePort
      selector:
        run: httpd
      ports:
      - protocol: TCP
        nodePort: 30000 ##nodePort指定特定端口(节点上监听的端口)
        port: 8080      ##port是ClusterIP上监听的端口
        targetPort: 80  ##targetPort是Pod监听的端口
        
    ###重新部署资源
    [root@node1 k8s]# kubectl apply -f httpd-service.yml 
    service/httpd-svc configured
    
    ###查看service,发现httpd-svc的节点端口已是指定的30000
    [root@node1 k8s]# kubectl get service
    NAME         TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)          AGE
    httpd-svc    NodePort    10.1.82.198   <none>        8080:30000/TCP   18h
    kubernetes   ClusterIP   10.1.0.1      <none>        443/TCP          43d
    
    ###查看某个特定service
    [root@node1 k8s]# kubectl get service httpd-svc
    NAME        TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)          AGE
    httpd-svc   NodePort   10.1.82.198   <none>        8080:30000/TCP   18h

(3).LoadBalancer

    Service利用cloud provider特有的load balancer对外提供服务
    
    cloud provider负责将load balancer的流量导向Service.
    
    目前支持的cloud provider有GCP,AWS,Azur等
    
-------------------------------------------------

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
    
-------------------------------------------------

### 第8章 Health Check

    强大的自愈能力是k8s这类容器编排引擎的一个重要特性。
    
    自愈默认实现方式是自动重启发生故障的容器。
    
    除此之外，用户可得用Liveness和Readiness探测机制设置更精细的健康检查，进而实现如下需求：
    (1).零停机部署
    (2).避免部署无效的镜像
    (3).更加安全的回滚升级
    
#### 1.默认的健康检查

>k8s默认健康检查机制

    每个容器启动时都会执行一个进程，此进程由Dockerfile的CMD或ENTRYPOINT指定。
    
    若进程退出时返回码非零，则认为容器发生故障，k8s会根据restartPolicy重启容器.
    
>模拟容器发生故障

    [root@node1 k8s]# cat healthcheck.yml 
    apiVersion: v1
    kind: Pod   ##资源类型为Pod
    metadata:
      labels:
        test: healthcheck
      name: healthcheck
    spec:
      restartPolicy: OnFailure  ##失败后重启容器
      containers:
      - name: healthcheck
        image: busybox
        args:
        - /bin/sh
        - -c
        - sleep 10; exit 1
    [root@node1 k8s]# kubectl apply -f healthcheck.yml 
    pod/healthcheck created
    
    [root@node1 k8s]# kubectl get pod healthcheck
    NAME          READY   STATUS    RESTARTS   AGE
    healthcheck   1/1     Running   1          34s
    [root@node1 k8s]# kubectl get pod healthcheck
    NAME          READY   STATUS             RESTARTS   AGE
    healthcheck   0/1     CrashLoopBackOff   1          56s
    
    ###当前容器已经重启2次了
    [root@node1 k8s]# kubectl get pod healthcheck
    NAME          READY   STATUS    RESTARTS   AGE
    healthcheck   1/1     Running   2          62s
    
    上述：容器进程返回值非0，k8s则认为容器发生故障，需要重启。
    
    有很多情况发生了故障，但进程并不会退出。如访问web服务器时显示500，可能是系统超载，也可能是资源死锁。此时httpd进程并没有异常退出。
    在此种情况下，重启容器可能是最直接，最有效的解决方案。
    
    如何利用health check机制来处理这类场景--Liveness探测
    
    ###删除pod
    [root@node1 k8s]# kubectl delete pod healthcheck
    pod "healthcheck" deleted

#### 2.Liveness探测

>Liveness探测可让用户自定义判断容器是否健康的条件，若探测失败，k8s会重启容器

    ###Step1,编写liveness.yml文件
    [root@node1 k8s]# cat liveness.yml 
    apiVersion: v1
    kind: Pod
    metadata:
      labels:
        test: liveness
      name: liveness
    spec:
      restartPolicy: OnFailure  ##重启容器
      containers:
      - name: liveness
        image: busybox
        args:
        - /bin/sh
        - -c
        - touch /tmp/healthy; sleep 30; rm -rf /tmp/healthy; sleep 600
        livenessProbe:
          exec:
            command:
            - cat    ##探测方法通过cat命令检查/tmp/healthy文件是否存在，若命令执行成功返加0,k8s认为本次liveness探测成功，非0即liveness探测失败，
            - /tmp/healthy
          initialDelaySeconds: 10  ##10指定容器启动10s之后开始执行liveness探测(需大于应用启动时间)
          periodSeconds: 5 ##指定每5s执行一次liveness探测,k8s若连续执行3次liveness均失败，则会杀掉并重启容器
          
    说明：上述模拟：启动进程首先先创建/tmp.healthy,30s后删除，若文件存在，则容器处于正常状态，否则发生故障
    
    [root@node1 k8s]# kubectl apply -f liveness.yml 
    pod/liveness created
    
    
    [root@node1 k8s]# kubectl get pod liveness
    NAME       READY   STATUS    RESTARTS   AGE
    liveness   1/1     Running   1          102s
    
    [root@node1 k8s]# kubectl describe pod liveness
    Name:               liveness
    Namespace:          default
    Priority:           0
    PriorityClassName:  <none>
    Node:               node3/192.168.1.33
    Start Time:         Sat, 11 May 2019 06:00:11 +0000
    Labels:             test=liveness
    Annotations:        cni.projectcalico.org/podIP: 10.244.135.2/32
                        kubectl.kubernetes.io/last-applied-configuration:
                          {"apiVersion":"v1","kind":"Pod","metadata":{"annotations":{},"labels":{"test":"liveness"},"name":"liveness","namespace":"default"},"spec":...
    Status:             Running
    IP:                 10.244.135.2
    Containers:
      liveness:
        Container ID:  docker://84645e0ef024ca0bc873397ae73833441e1f0672c150b016ed91f324a24b267b
        Image:         busybox
        Image ID:      docker-pullable://busybox@sha256:260d47cd183d41ccd7e80f8e4b3ada9a249a395f96d6f9895f9e274c6f856d3b
        Port:          <none>
        Host Port:     <none>
        Args:
          /bin/sh
          -c
          touch /tmp/healthy; sleep 30; rm -rf /tmp/healthy; sleep 600
        State:          Running
          Started:      Sat, 11 May 2019 06:01:41 +0000
        Last State:     Terminated
          Reason:       Error
          Exit Code:    137
          Started:      Sat, 11 May 2019 06:00:20 +0000
          Finished:     Sat, 11 May 2019 06:01:32 +0000
        Ready:          True
        Restart Count:  1
        Liveness:       exec [cat /tmp/healthy] delay=10s timeout=1s period=5s #success=1 #failure=3
        Environment:    <none>
        Mounts:
          /var/run/secrets/kubernetes.io/serviceaccount from default-token-gtn9w (ro)
    Conditions:
      Type              Status
      Initialized       True 
      Ready             True 
      ContainersReady   True 
      PodScheduled      True 
    Volumes:
      default-token-gtn9w:
        Type:        Secret (a volume populated by a Secret)
        SecretName:  default-token-gtn9w
        Optional:    false
    QoS Class:       BestEffort
    Node-Selectors:  <none>
    Tolerations:     node.kubernetes.io/not-ready:NoExecute for 300s
                     node.kubernetes.io/unreachable:NoExecute for 300s
    Events:
      Type     Reason     Age                 From               Message
      ----     ------     ----                ----               -------
      Normal   Scheduled  106s                default-scheduler  Successfully assigned default/liveness to node3
      Warning  Unhealthy  56s (x3 over 66s)   kubelet, node3     Liveness probe failed: cat: can't open '/tmp/healthy': No such file or directory
      Normal   Pulling    26s (x2 over 106s)  kubelet, node3     pulling image "busybox"
      Normal   Killing    26s                 kubelet, node3     Killing container with id docker://liveness:Container failed liveness probe.. Container will be killed and recreated.
      Normal   Pulled     17s (x2 over 98s)   kubelet, node3     Successfully pulled image "busybox"
      Normal   Created    17s (x2 over 98s)   kubelet, node3     Created container
      Normal   Started    17s (x2 over 98s)   kubelet, node3     Started container

    ###删除pod
    [root@node1 k8s]# kubectl delete pod liveness
    pod "liveness" deleted

#### 3.Readliness探测

    除了Liveness探测，k8s健康检查机制还包括：Readliness探测
    
    liveness探测机制可告诉k8s什么时候通过重启容器治愈
    
    readliness探测告诉k8s什么时候可以将容器加入到Service的负载均衡池中,对外提供服务
    
>实战

    [root@node1 k8s]# cat readiness.yml 
    apiVersion: v1
    kind: Pod
    metadata:
      labels:
        test: readiness
      name: readiness
    spec:
      restartPolicy: OnFailure
      containers:
      - name: readiness
        image: busybox
        args:
        - /bin/sh
        - -c
        - touch /tmp/healthy; sleep 30; rm -rf /tmp/healthy; sleep 600
        readinessProbe:
          exec:
            command:
            - cat
            - /tmp/healthy
          initialDelaySeconds: 10  ##10指定容器启动10s之后开始执行readiness探测(需大于应用启动时间)
          periodSeconds: 5   ##指定每5s执行一次readiness探测
          
    [root@node1 k8s]# kubectl apply -f readiness.yml 
    pod/readiness created
    
    [root@node1 k8s]# kubectl get pod readiness
    NAME        READY   STATUS    RESTARTS   AGE
    readiness   0/1     Running   0          13s
    [root@node1 k8s]# kubectl get pod readiness
    NAME        READY   STATUS    RESTARTS   AGE
    readiness   1/1     Running   0          24s
    [root@node1 k8s]# kubectl get pod readiness
    NAME        READY   STATUS    RESTARTS   AGE
    readiness   0/1     Running   0          57s
    
    上述Pod readiness的READY状态经历如下变化：
    (1)。刚被创建时，READY状态为不可用
    (2)。15s后(initialDelaySeconds+periodSeconds)，第一次进行Readiness探测并成功返回，设置READY为可用
    (3)。30s后，/tmp/healthy被删除，连续3次readiness探测失败后READY不可用
    
    
    #$##查看Readiness探测失败日志
    [root@node1 k8s]# kubectl describe pod readiness
    
>liveness和Readiness探测比较

    (1)liveness和readiness探测是两种健康检查机制，若不特意配置，k8s将对两种探测采取相同的默认行为。即通过判断容器启动进程的返回值是否为0判断探测是否成功
    
    (2)两种探测的配置方法完全一样，支持的配置参数也一样。
       不同之处在于探测失败后的行为：
       A.liveness探测是重启容器
       B.Readiness探测则是将容器设置为不可用，不接收service转发的请求
       
    (3)liveness和readiness探测是独立执行的，两者没有依赖。即可单独使用，也可同时使用。
       A.用liveness探测判断容器是否需要重启以实现自愈
       B.用readiness探测判断容器是否已准备好对外提供服务
    
#### 4.Health Check在Scale up中的应用

    当执行Scale up操作时，新副本会作为backend被添加到service的负载均衡中，与已有副本一起处理客户请求。
    
    应用readiness探测可判断容器是否就绪，避免将请求发送到还未准备好的backend
    
>实战

    [root@node1 k8s]# cat scaleup-readiness.yml 
    apiVersion: apps/v1beta1
    kind: Deployment
    metadata:
      name: web
    spec:
      replicas: 3
      template:
        metadata:
          labels: 
            run: web
        spec:
          containers:
          - name: web
            image: myhttpd
            ports:
            - containerPort: 8080
            readinessProbe:
              httpGet:  ##不同于exec,应用httpGet探测方法，k8s对该方法探测成功判断条件是http请求的返回代码在200-400间
                scheme: HTTP  ##scheme指定协议，支持http/https
                path: /health  ##path指定访问路径
                port: 8080     ##指定端口
              initialDelaySeconds: 10  ##容器启动10s后开始探测,若http://[container_ip]:8080/health返回代码不是200-400，则容器未就绪，不接收service web-svc请求
              periodSeconds: 5  ##每隔5s探测一次,代码返回200-400表容器就绪，可加入到web-svc负载均衡。开始处理客户请求。(若连续发生3次失败，容器会从负载均衡中移除，直到探测成功重新加入)
    
    ---
    apiVersion: v1
    kind: Service
    metadata:
      name: web-svc
    spec:
      selector:
        run: web
      ports:
      - protocol: TCP
        port: 8080
        targetPort: 80
        
    [root@node1 k8s]# kubectl apply -f scaleup-readiness.yml 
    deployment.apps/web created
    service/web-svc created
    
    ###对于http://[container_ip]:8080/health，应用则可实现自己的判断逻辑
    
    

#### 5.Health Check在滚动更新中的应用

    在Rolling Update时，若没有Health Check,则产生问题
    
    新副本有问题，但未异常退出，默认Health Check机制认为容器已准备就绪，进而逐步替换现有副本。
    
    当所有旧副本被替换后，整个应用将无法处理请求，无法对外提供服务。
    
    若正确配置了Heath Check，则新副本只有通过readiness探测才会被添加到Service.如果没有探测，现有副本不会被全部替换。业务依然正常运行
    
>实例：Health check在Rolling update中的应用

    [root@node1 k8s]# cat app.v1.yml 
    apiVersion: apps/v1beta1
    kind: Deployment
    metadata:
      name: app
    spec:
      replicas: 10
      template:
        metadata:
          labels:
            run: app
        spec:
          containers:
          - name: app
            image: busybox
            args:
            - /bin/sh
            - -c
            - sleep 10; touch /tmp/healthy; sleep 30000
            readinessProbe:
              exec:
                command:
                - cat
                - /tmp/healthy
              initialDelaySeconds: 10
              periodSeconds: 5
    [root@node1 k8s]# kubectl apply -f app.v1.yml --record
    deployment.apps/app created
    
    [root@node1 k8s]# kubectl get deployment app
    NAME   READY   UP-TO-DATE   AVAILABLE   AGE
    app    10/10   10           10          90s
    [root@node1 k8s]# kubectl get pod
    NAME                   READY   STATUS    RESTARTS   AGE
    app-56878b4676-5qlf9   1/1     Running   0          94s
    app-56878b4676-5vwg5   1/1     Running   0          94s
    app-56878b4676-7xcnf   1/1     Running   0          94s
    app-56878b4676-h57mr   1/1     Running   0          94s
    app-56878b4676-njgrt   1/1     Running   0          94s
    app-56878b4676-prfvp   1/1     Running   0          94s
    app-56878b4676-rsch5   1/1     Running   0          94s
    app-56878b4676-s2bkk   1/1     Running   0          94s
    app-56878b4676-sxbk2   1/1     Running   0          94s
    app-56878b4676-tlqkz   1/1     Running   0          94s
    
    
    ###滚动更新，创建app.v2.yml
    $ cp app.v1.yml app.v2.yml
    [root@node1 k8s]# cat app.v2.yml 
    apiVersion: apps/v1beta1
    kind: Deployment
    metadata:
      name: app
    spec:
      replicas: 10
      template:
        metadata:
          labels:
            run: app
        spec:
          containers:
          - name: app
            image: busybox
            args:
            - /bin/sh
            - -c
            - sleep 3000   ###修改此处
            readinessProbe:
              exec:
                command:
                - cat
                - /tmp/healthy
              initialDelaySeconds: 10
              periodSeconds: 5
              
      [root@node1 k8s]# kubectl apply -f app.v2.yml --record
      deployment.apps/app configured
      
      [root@node1 k8s]# kubectl get deployment app
      NAME   READY   UP-TO-DATE   AVAILABLE   AGE
      app    8/10    5            8           3m15s
      [root@node1 k8s]# kubectl get pod
      NAME                   READY   STATUS    RESTARTS   AGE
      app-56878b4676-5qlf9   1/1     Running   0          5m51s   ###旧副本从最初10个减少到8个
      app-56878b4676-5vwg5   1/1     Running   0          5m51s
      app-56878b4676-7xcnf   1/1     Running   0          5m51s
      app-56878b4676-h57mr   1/1     Running   0          5m51s
      app-56878b4676-njgrt   1/1     Running   0          5m51s
      app-56878b4676-prfvp   1/1     Running   0          5m51s
      app-56878b4676-sxbk2   1/1     Running   0          5m51s
      app-56878b4676-tlqkz   1/1     Running   0          5m51s
      app-74cd448d85-7zfp2   0/1     Running   0          3m5s     ###最后5个是新副本，处于NOT READY状态
      app-74cd448d85-cscpt   0/1     Running   0          3m5s
      app-74cd448d85-rxwrb   0/1     Running   0          3m5s
      app-74cd448d85-xx4cl   0/1     Running   0          3m5s
      app-74cd448d85-zrc28   0/1     Running   0          3m5s
      
      ###查看deployment部署日志
      [root@node1 k8s]# kubectl describe deployment app
      
      上述：新副本始终无法通过readiness探测，Health check屏蔽了有缺陷的副本，同时保留了大部分旧副本。业务没有因更新失败受影响
      
      
      ###回滚，查看历史版本
      [root@node1 k8s]# kubectl rollout history deployment app
      deployment.extensions/app 
      REVISION  CHANGE-CAUSE
      1         kubectl apply --filename=app.v1.yml --record=true
      2         kubectl apply --filename=app.v2.yml --record=true
      
      ###恢复到1版本
      [root@node1 k8s]# kubectl rollout undo deployment app --to-revision=1
      deployment.extensions/app rolled back

      [root@node1 k8s]# kubectl get deployment app
      NAME   READY   UP-TO-DATE   AVAILABLE   AGE
      app    10/10   10           10          24m
      
  >为什么新创建的副本数是5个，同时只销毁了2个旧副本?
  
    因为滚动更新参数maxSurge和maxUnavailable来控制副本替换的数量
    
  **maxSurge**
  
    此参数控制滚动更新过程中副本总数超过DESIRED的上限。
    
    maxSurege可以是具体的整数(如：3)，也可以是百分比，向上取整。maxSurge默认值为25%
    
    上述：DESIRED值为10,副本最大值roundUp(10+10*25%)=13,因此上述是8+5=13
    
    maxSurge值越大初始创建的新副本数量就越多
  
  **maxUnavailable**

    此参数控制滚动更新过程中，不可用的副本占DESIRED的最大比例
    
    maxUnavailable可以是具体的整数(如：3)，也可以是百分比，向下取整。maxUnavailable默认值为25%
    
    上述：DESIRED值为10，则可用副本数至少(10-10*25%)=8,因此AVAILABLE值为8
    
    maxUnavailable值越大，初始销毁旧副本数量就越多
    
  >定制maxSurge和maxUnavailable
  
    [root@node1 k8s]# cat app.v1.yml 
    apiVersion: apps/v1beta1
    kind: Deployment
    metadata:
      name: app
    spec:
      strategy:  ##定制maxSurge和maxUnavailable
        rollingUpdate:
          maxSurge: 35%       ##最大创建副本数，可以为数字，也可以为百分比
          maxUnavailable: 35%  ##不可用副本数占DESIRED的最大比例
      replicas: 10
      template:
        metadata:
          labels:
            run: app
        spec:
          containers:
          - name: app
            image: busybox
            args:
            - /bin/sh
            - -c
            - sleep 10; touch /tmp/healthy; sleep 30000
            readinessProbe:
              exec:
                command:
                - cat
                - /tmp/healthy
              initialDelaySeconds: 10
              periodSeconds: 5
  
-------------------------------------------------

### 第9章 数据管理


-------------------------------------------------
    



      
     
     





    
    



  
