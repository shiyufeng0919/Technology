>1.Redis (error) NOAUTH Authentication required.解决方法

    问题出现原因：redis.conf配置文件中，配置了「requirepass syf」导致:

    [root@jdchain1 bin]# ./redis-cli -h 127.0.0.1 -p 6379
    127.0.0.1:6379> select 15
    (error) NOAUTH Authentication required.
    
    解决方案：
    
    127.0.0.1:6379> auth 'syf'
    OK
    127.0.0.1:6379> select 15
    OK
    127.0.0.1:6379[15]> exit
