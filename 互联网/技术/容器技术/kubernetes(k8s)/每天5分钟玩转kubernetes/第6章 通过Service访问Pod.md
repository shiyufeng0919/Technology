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
    
--------------------------------------------

**体胖还需勤跑步，人丑就要多读书!!! --开心玉凤**