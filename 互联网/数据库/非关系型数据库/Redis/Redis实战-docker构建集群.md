# Docker下redis的主从配置

参考[Docker下redis的主从配置](https://blog.csdn.net/qq_36183935/article/details/80794167)

## 1.拉取redis镜像

    docker pull daocloud.io/library/redis:latest
    
## 2.启动3个redis容器服务，分别使用到6379、6380、6381端口

    # --name指定容器名称  -p指定宿主机与虚拟机端口映射  -d指定镜像
    $ docker run --name redis-6379 -p 6379:6379 -d daocloud.io/library/redis
    
    $ docker run --name redis-6380 -p 6380:6379 -d daocloud.io/library/redis
    
    $ docker run --name redis-6381 -p 6381:6379 -d daocloud.io/library/redis
   
## 3.查看容器

    shiyufeng:redis shiyufeng$ docker ps | grep redis
    5f0d05fcc01c        daocloud.io/library/redis   "docker-entrypoint.s…"   8 seconds ago       Up 8 seconds        0.0.0.0:6381->6379/tcp               redis-6381
    13ad84d8671c        daocloud.io/library/redis   "docker-entrypoint.s…"   15 seconds ago      Up 14 seconds       0.0.0.0:6380->6379/tcp               redis-6380
    2090161b82d2        daocloud.io/library/redis   "docker-entrypoint.s…"   23 seconds ago      Up 23 seconds       0.0.0.0:6379->6379/tcp               redis-6379

## 4.测试容器,成功

     shiyufeng:redis shiyufeng$ docker exec -it redis-6379 bash
     root@2090161b82d2:/data# redis-cli
     127.0.0.1:6379> set key-a value-a
     OK
     127.0.0.1:6379> get key-a
     "value-a"
     127.0.0.1:6379> quit
     root@2090161b82d2:/data# exit
     exit
     
## 5.redis集群配置

### 5.1 查看容器内网ip地址

    $ docker inspect redis-6379   //NetworkSettings->IPAddress->172.17.0.4
    $ docker inspect redis-6380   //172.17.0.5
    $ docker inspect redis-6381   //172.17.0.6
    
### 5.2 进入容器内部。查看当前redis角色(主/从)

    #当前三个容器redis角色均为master
    root@2090161b82d2:/data# redis-cli
    127.0.0.1:6379> info replication
    # Replication
    role:master
    connected_slaves:0
    master_replid:6226836fa0d3de254aa0e776ab3b22f11462a89e
    master_replid2:0000000000000000000000000000000000000000
    master_repl_offset:0
    second_repl_offset:-1
    repl_backlog_active:0
    repl_backlog_size:1048576
    repl_backlog_first_byte_offset:0
    repl_backlog_histlen:0
    
### 5.3 使用redis-cli命令修改redis-6380、redis-6381的主机为172.17.0.4:6379

    shiyufeng:data shiyufeng$ docker exec -it redis-6380 bash
    root@13ad84d8671c:/data# redis-cli 
    127.0.0.1:6379> SLAVEOF 172.17.0.4 6379
    OK
    
    shiyufeng:~ shiyufeng$ docker exec -it redis-6381 bash
    root@5f0d05fcc01c:/data# redis-cli
    127.0.0.1:6379> SLAVEOF 172.17.0.4 6379
    OK
    
### 5.4 查看redis-6379是否已经拥有2个从机，connected_slaves:2

    127.0.0.1:6379> info replication
    # Replication
    role:master
    connected_slaves:2
    slave0:ip=172.17.0.5,port=6379,state=online,offset=84,lag=0
    slave1:ip=172.17.0.6,port=6379,state=online,offset=84,lag=0
    master_replid:39b96e101a7c100cf4073ad916653a30730e9f42
    master_replid2:0000000000000000000000000000000000000000
    master_repl_offset:84
    second_repl_offset:-1
    repl_backlog_active:1
    repl_backlog_size:1048576
    repl_backlog_first_byte_offset:1
    repl_backlog_histlen:84
 
### 5.5 配置Sentinel哨兵
 
    进入3台redis容器内部进行配置，在容器根目录里面创建sentinel.conf文件
    
    文件内容为：sentinel monitor mymaster 172.17.0.4 6379 1 
   
    #分别配置redis-6379 / redis-6380 / redis-6381 
    $ docker exec -it redis-6379 bash
    
    root@2090161b82d2:/data# cd / && touch sentinel.conf 
    
    #若提示无vim,解决方案：$ apt-get update && apt-get install vim
    root@2090161b82d2:/# vim sentinel.conf 
    
    #添加内容：sentinel monitor mymaster 172.17.0.4 6379 1 
    
    #启动redis哨兵
    root@2090161b82d2:/# redis-sentinel ./sentinel.conf 
    343:X 01 May 2019 06:58:44.878 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
    343:X 01 May 2019 06:58:44.878 # Redis version=5.0.4, bits=64, commit=00000000, modified=0, pid=343, just started
    343:X 01 May 2019 06:58:44.878 # Configuration loaded
                    _._                                                  
               _.-``__ ''-._                                             
          _.-``    `.  `_.  ''-._           Redis 5.0.4 (00000000/0) 64 bit
      .-`` .-```.  ```\/    _.,_ ''-._                                   
     (    '      ,       .-`  | `,    )     Running in sentinel mode
     |`-._`-...-` __...-.``-._|'` _.-'|     Port: 26379
     |    `-._   `._    /     _.-'    |     PID: 343
      `-._    `-._  `-./  _.-'    _.-'                                   
     |`-._`-._    `-.__.-'    _.-'_.-'|                                  
     |    `-._`-._        _.-'_.-'    |           http://redis.io        
      `-._    `-._`-.__.-'_.-'    _.-'                                   
     |`-._`-._    `-.__.-'    _.-'_.-'|                                  
     |    `-._`-._        _.-'_.-'    |                                  
      `-._    `-._`-.__.-'_.-'    _.-'                                   
          `-._    `-.__.-'    _.-'                                       
              `-._        _.-'                                           
                  `-.__.-'                                               
    
    343:X 01 May 2019 06:58:44.880 # WARNING: The TCP backlog setting of 511 cannot be enforced because /proc/sys/net/core/somaxconn is set to the lower value of 128.
    343:X 01 May 2019 06:58:44.887 # Sentinel ID is c6d5f58f947160b43ee1996e66fd0ed5f0d87274
    343:X 01 May 2019 06:58:44.887 # +monitor master mymaster 172.17.0.4 6379 quorum 1
    343:X 01 May 2019 06:58:44.888 * +slave slave 172.17.0.5:6379 172.17.0.5 6379 @ mymaster 172.17.0.4 6379
    343:X 01 May 2019 06:58:44.889 * +slave slave 172.17.0.6:6379 172.17.0.6 6379 @ mymaster 172.17.0.4 6379
    343:X 01 May 2019 07:03:18.428 * +sentinel sentinel df353b7623d86d3c7d445a56841d506740151a12 172.17.0.6 26379 @ mymaster 172.17.0.4 6379
    343:X 01 May 2019 07:03:34.573 * +sentinel sentinel 062e6cbb381fb222d0e4f0ed5b70b28bbac47a10 172.17.0.5 26379 @ mymaster 172.17.0.4 6379

**Sentinel哨兵配置完毕**

### 5.6 测试-关闭master

    $ docker stop redis-6379
    
    #redis-6381容器,发现原redis-6379变成了从机，redis-6381变成了主机(随机)
    root@5f0d05fcc01c:/# redis-sentinel ./sentinel.conf 
    329:X 01 May 2019 07:03:16.380 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
    329:X 01 May 2019 07:03:16.380 # Redis version=5.0.4, bits=64, commit=00000000, modified=0, pid=329, just started
    329:X 01 May 2019 07:03:16.380 # Configuration loaded
                    _._                                                  
               _.-``__ ''-._                                             
          _.-``    `.  `_.  ''-._           Redis 5.0.4 (00000000/0) 64 bit
      .-`` .-```.  ```\/    _.,_ ''-._                                   
     (    '      ,       .-`  | `,    )     Running in sentinel mode
     |`-._`-...-` __...-.``-._|'` _.-'|     Port: 26379
     |    `-._   `._    /     _.-'    |     PID: 329
      `-._    `-._  `-./  _.-'    _.-'                                   
     |`-._`-._    `-.__.-'    _.-'_.-'|                                  
     |    `-._`-._        _.-'_.-'    |           http://redis.io        
      `-._    `-._`-.__.-'_.-'    _.-'                                   
     |`-._`-._    `-.__.-'    _.-'_.-'|                                  
     |    `-._`-._        _.-'_.-'    |                                  
      `-._    `-._`-.__.-'_.-'    _.-'                                   
          `-._    `-.__.-'    _.-'                                       
              `-._        _.-'                                           
                  `-.__.-'                                               
    
    329:X 01 May 2019 07:03:16.382 # WARNING: The TCP backlog setting of 511 cannot be enforced because /proc/sys/net/core/somaxconn is set to the lower value of 128.
    329:X 01 May 2019 07:03:16.389 # Sentinel ID is df353b7623d86d3c7d445a56841d506740151a12
    329:X 01 May 2019 07:03:16.389 # +monitor master mymaster 172.17.0.4 6379 quorum 1
    329:X 01 May 2019 07:03:16.390 * +slave slave 172.17.0.5:6379 172.17.0.5 6379 @ mymaster 172.17.0.4 6379
    329:X 01 May 2019 07:03:16.391 * +slave slave 172.17.0.6:6379 172.17.0.6 6379 @ mymaster 172.17.0.4 6379
    329:X 01 May 2019 07:03:18.283 * +sentinel sentinel c6d5f58f947160b43ee1996e66fd0ed5f0d87274 172.17.0.4 26379 @ mymaster 172.17.0.4 6379
    329:X 01 May 2019 07:03:34.573 * +sentinel sentinel 062e6cbb381fb222d0e4f0ed5b70b28bbac47a10 172.17.0.5 26379 @ mymaster 172.17.0.4 6379
    329:X 01 May 2019 07:10:03.597 # +new-epoch 1
    329:X 01 May 2019 07:10:03.597 # +vote-for-leader 062e6cbb381fb222d0e4f0ed5b70b28bbac47a10 1
    329:X 01 May 2019 07:10:03.607 # +sdown master mymaster 172.17.0.4 6379
    329:X 01 May 2019 07:10:03.607 # +odown master mymaster 172.17.0.4 6379 #quorum 1/1
    329:X 01 May 2019 07:10:03.607 # Next failover delay: I will not start a failover before Wed May  1 07:16:04 2019  #发现异常
    329:X 01 May 2019 07:10:03.607 # +sdown sentinel c6d5f58f947160b43ee1996e66fd0ed5f0d87274 172.17.0.4 26379 @ mymaster 172.17.0.4 6379
    329:X 01 May 2019 07:10:04.836 # +config-update-from sentinel 062e6cbb381fb222d0e4f0ed5b70b28bbac47a10 172.17.0.5 26379 @ mymaster 172.17.0.4 6379
    329:X 01 May 2019 07:10:04.836 # +switch-master mymaster 172.17.0.4 6379 172.17.0.6 6379 #切换主机由4更为6
    329:X 01 May 2019 07:10:04.837 * +slave slave 172.17.0.5:6379 172.17.0.5 6379 @ mymaster 172.17.0.6 6379
    329:X 01 May 2019 07:10:04.837 * +slave slave 172.17.0.4:6379 172.17.0.4 6379 @ mymaster 172.17.0.6 6379
    329:X 01 May 2019 07:10:34.849 # +sdown slave 172.17.0.4:6379 172.17.0.4 6379 @ mymaster 172.17.0.6 6379