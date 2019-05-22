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
              

--------------------------------------------

**体胖还需勤跑步，人丑就要多读书!!! --开心玉凤**