
### 第10章 Secret(处理敏感信息) & Configmap

    敏感信息，直接保存在容器镜像中显然不妥，比如用户名、密码等。K8s提供的解决方案是Secret。
    
    Secret会以密文的方式存储数据，避免了在配置文件中保存敏感信息。Secret会以Volume的形式被mount到Pod，容器可通过文件的方式使用Secret中的敏感数据；
    
    此外，容器也可以"环境变量"的方式使用这些数据。Secret可通过命令行或YAML创建。
    
#### 1.创建secret--处理敏感信息

>四种方法创建secret

##### 1-1. --from-literal : --from-literal对应一个条目

    #创建secret
    [root@node1 ~]# kubectl create secret generic mysecret --from-literal=username=syf --from-literal=password=123
    secret/mysecret created
    
    #查询secret
    [root@node1 ~]# kubectl get secret
    NAME                  TYPE                                  DATA   AGE
    default-token-gtn9w   kubernetes.io/service-account-token   3      50d
    mysecret              Opaque                                2      21s

##### 1-2. --from-file:每个文件对应一个条目

     #向文件中写入密码
    [root@node1 ~]# echo -n 123456 > ./pas
     #向文件中写入用户
    [root@node1 ~]# echo -n kaixinyufeng > ./user
     #创建secret
    [root@node1 ~]# kubectl create secret generic mysecret2 --from-file=./user --from-file=./pas
    secret/mysecret2 created
     #查询secret
    [root@node1 ~]# kubectl get secret
    NAME                  TYPE                                  DATA   AGE
    default-token-gtn9w   kubernetes.io/service-account-token   3      50d
    mysecret              Opaque                                2      2m33s
    mysecret2             Opaque                                2      4s
     
##### 1-3. 通过--from-env-file:文件中的key=value对应一个条目

    # cat << EOF > env.txt 是覆盖模式，cat <<EOF >>env.txt是追加模式
    # 向文件env.txt中写
    [root@node1 ~]# cat << EOF > env.txt
    > username=yufeng
    > password=666
    > EOF
    
    #从env.txt文件中读
    [root@node1 ~]# cat env.txt 
    username=yufeng
    password=666   
    
    #创建secret
    [root@node1 ~]# kubectl create secret generic mysecret3 --from-env-file=env.txt
    secret/mysecret3 created
    
    #查询secret
    [root@node1 ~]# kubectl get secret
    NAME                  TYPE                                  DATA   AGE
    default-token-gtn9w   kubernetes.io/service-account-token   3      50d
    mysecret              Opaque                                2      10m
    mysecret2             Opaque                                2      7m47s
    mysecret3             Opaque                                2      2m27s
    
##### 1-4. 通过YAML配置文件

**secret里面存储的数据必须是通过base64编码后的结果**

    #先获取base64
    [root@node1 ~]# echo -n shiyufeng | base64
    c2hpeXVmZW5n
    [root@node1 ~]# echo -n 8989 | base64
    ODk4OQ==
    
    #编辑yaml文件
    [root@node1 k8s]# cat secret.yml 
    apiVersion: v1
    kind: Secret
    metadata:
      name: mysecret4
    data:
      username: c2hpeXVmZW5n   #base64后的值
      password: ODk4OQ==
    
    #创建secret
    [root@node1 k8s]# kubectl apply -f secret.yml 
    secret/mysecret4 created
    
    #查询secret
    [root@node1 k8s]# kubectl get secret
    NAME                  TYPE                                  DATA   AGE
    default-token-gtn9w   kubernetes.io/service-account-token   3      50d
    mysecret              Opaque                                2      14m
    mysecret2             Opaque                                2      11m
    mysecret3             Opaque                                2      6m31s
    mysecret4             Opaque                                2      10s

    #查看secret明细
    [root@node1 k8s]# kubectl describe secret mysecret2
    Name:         mysecret2
    Namespace:    default
    Labels:       <none>
    Annotations:  <none>
    
    Type:  Opaque
    
    Data
    ====
    pas:   6 bytes
    user:  12 bytes
    
    #查看secret具体内容
    [root@node1 k8s]# kubectl edit secret mysecret4
    
    #显示内容
    apiVersion: v1
    data:
      password: ODk4OQ==
      username: c2hpeXVmZW5n
    kind: Secret
    metadata:
      annotations:
        kubectl.kubernetes.io/last-applied-configuration: |
          {"apiVersion":"v1","data":{"password":"ODk4OQ==","username":"c2hpeXVmZW5n"},"kind":"Secret","metadata":{"annotations":{},"name":"mysecret4","namespace":"default"}}
      creationTimestamp: "2019-05-18T07:37:45Z"
      name: mysecret4
      namespace: default
      resourceVersion: "366645"
      selfLink: /api/v1/namespaces/default/secrets/mysecret4
      uid: d3a228c1-793f-11e9-a4b7-525400261060
    type: Opaque
    
    #应用base64解码
    [root@node1 k8s]# echo -n c2hpeXVmZW5n| base64 --decode
    shiyufeng

#### 2.在Pod中使用Secret

##### 2-1.Volume方式

    ##编辑yaml
    [root@node1 k8s]# cat mypod-secret.yml 
    apiVersion: v1
    kind: Pod
    metadata:
      name: mypod
    spec:
      containers:
      - name: mypod
        image: busybox
        args:
          - /bin/sh
          - -c
          - sleep 10; touch /tmp/healthy; sleep 30000
        volumeMounts:
        - name: foo
          mountPath: "/etc/foo"    # 在容器内部的该路径下
          readOnly: true
      volumes:
      - name: foo
        secret:
          secretName: mysecret     # 指定有前面创建的mysecret

    ##创建Pod
    [root@node1 k8s]# kubectl apply -f mypod-secret.yml 
    pod/mypod created
    
    ##查询pod
    [root@node1 k8s]# kubectl get pod
    NAME                     READY   STATUS    RESTARTS   AGE
    mypod                    1/1     Running   0          15s
    mysql-7686899cf9-rkg7q   1/1     Running   1          4d5h
    mysql-client             0/1     Error     0          4d5h
    
    ##进入容器
    [root@node1 k8s]# kubectl exec -it mypod sh
    / # cd /etc/foo   ##进入
    
    /etc/foo # ls
    password  username
    
    /etc/foo # cat username   ##可直接查看内容，是明文
    syf/etc/foo # cat password
    123/etc/foo # exit

>k8s会在指定路径下为每条敏感信息创建一个文件.文件名是数据条目的Key, 
>/etc/foo/username和 etc/foo/password, value是以明文的形式存放在文件中

     ###删除Pod
    [root@node1 k8s]# kubectl delete pod mypod
    pod "mypod" deleted

>可自定义存放数据的文件名，配置文件如下改动：这时，数据将存放在/etc/foo/my-group/my-username中

    [root@node1 k8s]# cat mypod-secret2.yml 
    apiVersion: v1
    kind: Pod
    metadata:
      name: mypod
    spec:
      containers:
      - name: mypod
        image: busybox
        args:
          - /bin/sh
          - -c
          - sleep 10; touch /tmp/healthy; sleep 30
        volumeMounts:
        - name: foo
          mountPath: "/etc/foo"
          readOnly: true
      volumes:
      - name: foo
        secret:
          secretName: mysecret  #指定secret名称
          items:
          - key: username
            path: my-group/my-username  #自定义存放数据文件名
          - key: password
            path: my-group/my-password

    [root@node1 k8s]# kubectl apply -f mypod-secret2.yml 
    pod/mypod created
    
    [root@node1 k8s]# kubectl get pod
    NAME                     READY   STATUS    RESTARTS   AGE
    mypod                    1/1     Running   2          2m16s
    mysql-7686899cf9-rkg7q   1/1     Running   1          4d5h
    mysql-client             0/1     Error     0          4d5h
    
    [root@node1 k8s]# kubectl exec -it mypod sh
    / # cd /etc/foo/
    /etc/foo # ls
    my-group
    
    [root@node1 k8s]# kubectl delete pod mypod
    pod "mypod" deleted
    
**以Volume方式使用secret支持动态更新：Secret更新后，容器中的数据也会更新**

##### 2-2.环境变量方式

>通过volume方式使用secret，容器必须从文件读取数据，稍显麻烦。K8s支持通过环境变量使用secret。

    [root@node1 k8s]# cat mypod_env.yml 
    apiVersion: v1
    kind: Pod
    metadata:
      name: mypod
    spec:
      containers:
      - name: mypod
        image: busybox
        args:
          - /bin/sh
          - -c
          - sleep 10; touch /tmp/healthy; sleep 30000
        env:
          - name: SECRET_USERNAME        # 环境变量名字
            valueFrom:                   
              secretKeyRef:
                name: mysecret           # 从哪个secret来
                key: username            # key
          - name: SECRET_PASSWORD
            valueFrom:
              secretKeyRef:
                name: mysecret
                key: password
                
    [root@node1 k8s]# kubectl apply -f mypod_env.yml 
    pod/mypod created
    
    [root@node1 k8s]# kubectl get pod
    NAME                     READY   STATUS    RESTARTS   AGE
    mypod                    1/1     Running   0          26s
    mysql-7686899cf9-rkg7q   1/1     Running   1          4d5h
    mysql-client             0/1     Error     0          4d6h
    
    通过环境变量SECRET_USERNAME 和 SECRET_PASSWORD就可以读取到secret的数据.
    
**注意：环境变量的方式不支持Secret动态更新**

#### 3.ConfigMap

    Secret可以为Pod提供密码、Token、私钥等敏感数据；对于一些非敏感数据，比如一些配置信息，则可以用ConfigMap
    
    configMap的使用方式与Secret非常类似，不同的是数据是以明文的形式存放。
    
    创建方式也是四种（和Secret一样，不同的是kubectl create Secret 变成 kubectl create Configmap; 另外一个不同是：在yml文件中使用的时候，Secret 换成Configmap）


--------------------------------------------

**体胖还需勤跑步，人丑就要多读书!!! --开心玉凤**