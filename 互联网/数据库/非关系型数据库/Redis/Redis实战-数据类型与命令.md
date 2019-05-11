## 一.Redis数据类型

### (一).string类型

##### 1.赋值(set key value)与取值(get key)

    127.0.0.1:6379[1]> set key1 123
    OK
    127.0.0.1:6379[1]> get key1
    "123"

##### 2.取值并赋值(getset key value)

    127.0.0.1:6379[1]> getset key1 321
    "123"
    127.0.0.1:6379[1]> get key1
    "321"
    
##### 3.设置获取多个键值(mset key value [key value...] mget key [key...])

    127.0.0.1:6379[1]> mset key1 123 key2 456 key3 789
    OK
    127.0.0.1:6379[1]> mget key1 key2 key3
    1) "123"
    2) "456"
    3) "789"
    127.0.0.1:6379[1]> mset key1 a key2 b
    OK
    127.0.0.1:6379[1]> mget key1 key2 key3
    1) "a"
    2) "b"
    3) "789"
    
##### 4.删除(del key)

    127.0.0.1:6379[1]> del key3
    (integer) 1
    127.0.0.1:6379[1]> get key3
    (nil)
    
##### 5.数值增减

+ 递增数字(incr key)

**应用：自增主键 商品编号、订单号采用 string 的递增数字特性生成**

    127.0.0.1:6379[1]> set num 1
    OK
    127.0.0.1:6379[1]> incr num
    (integer) 2
    127.0.0.1:6379[1]> incr num
    (integer) 3

+ 增加指定的整数(incrby key increment)


    127.0.0.1:6379[1]> incrby num 2
    (integer) 5
    127.0.0.1:6379[1]> incrby num 5
    (integer) 10

+ 递减数字(decr key)

    
    127.0.0.1:6379[1]> decr num
    (integer) 9

+ 减少指定的数值(decryby key decrement)


    127.0.0.1:6379[1]> decrby num 4
    (integer) 5
    
##### 6。向尾部追加值Append

    127.0.0.1:6379[1]> set key1-str hello
    OK
    127.0.0.1:6379[1]> append key1-str "redis"
    (integer) 10
    127.0.0.1:6379[1]> get key1-str
    "helloredis"
    127.0.0.1:6379[1]> strlen key1-str    #计算字符串总长度
    (integer) 10
    127.0.0.1:6379[1]> append key2-str "syf" #追加值不存在，则设置key-value值
    (integer) 3
    127.0.0.1:6379[1]> get key2-str
    "syf"
    127.0.0.1:6379[1]> append key2-str " learn redis"
    (integer) 15
    127.0.0.1:6379[1]> get key2-str
    "syf learn redis"

### (二).Hash 散列类型

##### 1。使用string问题

    假设有User对象以JSON序列化的形式存储到Redis中(User对象有id，username、password、age、name等属性)，
    
    存储的过程如下：保存、更新：User对象  json(string)  redis
    
    如果在业务上只是更新age属性，其他的属性并不做更新。若仍采用上边的方法在传输、处理时会造成资源浪费，因此hash可以很好的解决这个问题。

##### 2。hash散列类型

    hash是散列类型，它提供了字段和字段值的映射。字段值只能是字符串类型，不支持散列类型、集合类型等其它类型
    
**HSET命令**

    HSET命令不区分插入和更新操作，当执行插入操作时HSET命令返回1，执行更新操作时返回0
    
###### 1。一次只设置/获取一个字段值 语法：hset key field value / hget key field

    127.0.0.1:6379[1]> hset user username syf
    (integer) 1
    
    127.0.0.1:6379[1]> hget user username
    "syf"

###### 2。一次设置/获取多个字段值 语法：hmset key field value [field value...] / hmget key field [field...]

    127.0.0.1:6379[1]> hmset user age 31 sex girl
    OK
    
    127.0.0.1:6379[1]> hmget user username age sex
    1) "syf"
    2) "31"
    3) "girl"
    
    127.0.0.1:6379[1]> hmset user username xiaofeng class redis   #有该字段则覆盖值，否则写入值
    OK
    
    127.0.0.1:6379[1]> hgetall user
    1) "username"
    2) "xiaofeng"
    3) "age"
    4) "31"
    5) "sex"
    6) "girl"
    7) "class"
    8) "redis"

###### 3。当字段不存在时赋值，类似hset,区别在于如果字段存在，该命令不执行任何操作。 语法：hsetnx key field value

    127.0.0.1:6379[1]> hsetnx user username yufeng
    (integer) 0

###### 4。获取所有字段值 语法：hgetall key

    127.0.0.1:6379[1]> hgetall user
    1) "username"
    2) "xiaofeng"
    3) "age"
    4) "31"
    5) "sex"
    6) "girl"
    7) "class"
    8) "redis"
    
###### 5。删除字段 可以删除一个或多个字段，返回值是被删除的字段的个数。 语法：hdel key field [field...]

    127.0.0.1:6379[1]> hdel user age
    (integer) 1
    127.0.0.1:6379[1]> hdel user age class
    (integer) 1
    127.0.0.1:6379[1]> hgetall user
    1) "username"
    2) "xiaofeng"
    3) "sex"
    4) "girl"

###### 6。增加数字 语法：hincrby key field increment

    127.0.0.1:6379[1]> hincrby user age 2    #age已经被删除，则写入age赋值为2
    (integer) 2
    127.0.0.1:6379[1]> hgetall user
    1) "username"
    2) "xiaofeng"
    3) "sex"
    4) "girl"
    5) "age"
    6) "2"

    127.0.0.1:6379[1]> hincrby user age 30  #age已存在且值为2，在其上加30
    (integer) 32
    127.0.0.1:6379[1]> hgetall user
    1) "username"
    2) "xiaofeng"
    3) "sex"
    4) "girl"
    5) "age"
    6) "32"
    
###### 7。判断字段是否存在 语法：hexists key field
         
    127.0.0.1:6379[1]> hexists user age   
    (integer) 1
    127.0.0.1:6379[1]> hexists user class
    (integer) 0
    
###### 8。只获取字段名或字段值 语法： hkeys key / hvals key

    127.0.0.1:6379[1]> hkeys user
    1) "username"
    2) "sex"
    3) "age"
    
    127.0.0.1:6379[1]> hvals user
    1) "xiaofeng"
    2) "girl"
    3) "32"
    
###### 9。获取字段数量 语法：hlen key

    127.0.0.1:6379[1]> hlen user
    (integer) 3
    
### (三).List 类型

##### 1.ArrayList 和 LinkedList 的区别

    Arraylist是使用数组来存储数据，特点：查询快、增删慢
    
    Linkedlist是使用双向链表存储数据，特点：增删快、查询慢，但是查询链表两端的数据也很快
    
    Redis的list是采用链表来存储的，所以对于redis的list数据类型的操作，是操作list的两端数据来操作的
    
##### 2.命令

###### (1).向列表两端增加元素(左边lpush/右边rpush key value [value...])

    127.0.0.1:6379[1]> lpush list:1 1 2 3 #从左边向里推1,2,3三个元素,3在最左边
    (integer) 3
    127.0.0.1:6379[1]> rpush list:1 4 5 6 #从右边向里推4,5,6三个元素,6在最右边
    (integer) 6
    
###### (2).查看列表 LRANGE命令(lrange key start stop)

    查看列表 LRANGE命令是列表类型最常用的命令之一，获取列表中的某一片段，将返回start、stop之间的所有元素（包含两端的元素），索引从0开始。索引可以是负数，如：“-1”代表最后边的一个元素。
    
    127.0.0.1:6379[1]> lrange list:1 0 2 #查询索引从0开始到索引为2的元素
    1) "3"
    2) "2"
    3) "1"
    127.0.0.1:6379[1]> lrange list:1 0 -1 #查询从索引为0开始到最后一个的元素
    1) "3"
    2) "2"
    3) "1"
    4) "4"
    5) "5"
    6) "6"
    
###### (3).从列表两端弹出元素 (左边lpop/右边rpop key)

    #弹出分两步：第一步是将列表左边的元素从列表中移除；第二步是返回被移除的元素值

    127.0.0.1:6379[1]> lpop list:1  #从左侧弹出一个元素
    "3"
    127.0.0.1:6379[1]> rpop list:1  #从右侧弹出一个元素
    "6"
    127.0.0.1:6379[1]> lrange list:1 0 -1  #查看列表
    1) "2"
    2) "1"
    3) "4"
    4) "5"
    
###### (4).获取列表中元素的个数 语法：llen key

    127.0.0.1:6379[1]> llen list:1
    (integer) 4
    
###### (5).获得/设置指定索引的元素值(lindex key index / lset key index value)
   
    127.0.0.1:6379[1]> lindex list:1 2   #获取指定元素索引2的值
    "4"
    127.0.0.1:6379[1]> lset list:1 2 2   #设置指定元素索引2的值
    OK
    127.0.0.1:6379[1]> lindex list:1 2   #查询索引2的值
    "2"
    127.0.0.1:6379[1]> lrange list:1 0 -1  #查询从索引为0到最后一个元素的值
    1) "2"
    2) "1"
    3) "2"    #索引为2的值被改
    4) "5"
    
###### (6).只保留列表指定片段 指定范围和 lrange 一致 语法：ltrim key start stop

    127.0.0.1:6379[1]> ltrim list:1 0 2    #保留索引从0到2的值
    OK
    127.0.0.1:6379[1]> lrange list:1 0 -1  #查询所有列表的值
    1) "2"
    2) "1"
    3) "2"
    
###### (7).向列表中插入元素(linsert key before | after pivot value)

**向列表中插入元素 该命令首先会在列表中从左到右查找值为pivot的元素，然后根据第二个参数是BEFORE还是AFTER来决定将value插入到该元素的前面还是后面**

    127.0.0.1:6379[1]> linsert list:1 after 1 9 #从左向右查找1，并将9插入到1后
    (integer) 4
    127.0.0.1:6379[1]> lrange list:1 0 -1  
    1) "2"
    2) "1"
    3) "9"
    4) "2"

###### (8).将元素从一个列表转移到另一个列表 语法：rpoplpush source destination

    127.0.0.1:6379[1]> rpoplpush list:1 newlist
    "2"
    127.0.0.1:6379[1]> lrange newlist 0 -1
    1) "2"
    127.0.0.1:6379[1]> lrange list:1 0 -1
    1) "2"
    2) "1"
    3) "9"
    
###### (9).删除列表中指定的值 LREM命令(lrem key count value)

**LREM命令会删除列表中前count个值为value的元素，返回实际删除的元素个数。根据count值的不同，该命令的执行方式会有所不同:**

    当count>0时， LREM会从列表左边开始删除。
    当count<0时， LREM会从列表后边开始删除。
    当count=0时， LREM删除所有值为value的元素。
    
    127.0.0.1:6379[1]> lrange list:1 0 -1 #查询列表所有数据
    1) "2"
    2) "1"
    3) "9"
    127.0.0.1:6379[1]> lrem list:1 1 1    #从左边删除值为1的所有元素
    (integer) 1
    127.0.0.1:6379[1]> lrange list:1 0 -1 #查询列表数据
    1) "2"
    2) "9"
    
## [Redis命令](http://redisdoc.com/pubsub/subscribe.html)

### 发布订阅

    #启动二个终端，便于测试和查看
    $ docker exec -it redis-6381 bash
    
    #终端1：订阅通道 msg和chat_room
    127.0.0.1:6379> subscribe msg chat_room
    Reading messages... (press Ctrl-C to quit)
    1) "subscribe"  #返回值的类型，显示订阅成功
    2) "msg"        #订阅频道名字
    3) (integer) 1  #目前已订阅频道数量
    1) "subscribe"
    2) "chat_room"
    3) (integer) 2
    
    #终端2:向通道发布消息
    127.0.0.1:6379> publish msg "hello syf"
    (integer) 1
    127.0.0.1:6379> publish chat_room "hello everyone"
    (integer) 1
    
    #终端1：查看
    1) "message"    #返回值类型：信息
    2) "msg"        #来源(从哪个通道发过来)
    3) "hello syf"  #信息内容
    1) "message"
    2) "chat_room"
    3) "hello everyone"
    
    
参考:

[go语言之行--golang操作redis、mysql大全](https://www.cnblogs.com/wdliu/p/9330278.html)
    
     
    
    