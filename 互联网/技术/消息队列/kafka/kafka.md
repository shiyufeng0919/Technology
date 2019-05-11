# Apache Kafka -- 分布式流处理平台

[kafka中文文档](http://kafka.apachecn.org/documentation.html#producerapi)

# Apache Kafka -- 掘金好文

### 一。 [消息中间件如何实现每秒几十万的高并发写入](https://mp.weixin.qq.com/s?__biz=MzU0OTk3ODQ3Ng==&mid=2247484700&idx=1&sn=fbfdb57ea53882828e4e3bd0b3b61947&chksm=fba6ed1fccd16409c43baa7f941e522d97a72e63e4139f663b327c606c6bb5dfe516b6f61424&scene=21#wechat_redirect)

    Kafka是高吞吐低延迟的高并发、高性能的消息中间件，在大数据领域有极为广泛的运用。配置良好的Kafka集群甚至可以做到每秒几十万、上百万的超高并发写入。
    
+ 页面缓存技术+磁盘顺序写(写入) **kafka核心架构原理**

![image](resources/images/1.jpg)

+ 零拷贝技术(消费)

![image](resources/images/2.jpg)

性能优化后：

![image](resources/images/3.jpg)

----------------------------------

### 二。[写入消息中间件的数据，如何保证不丢失](https://mp.weixin.qq.com/s/wbqA9vZOCQ0M_N9Q0NXWVg)

+ kafka分布式存储架构

![image](resources/images/4.jpg)

+ kafka高可用架构

![image](resources/images/5.jpg)

+ kafka写入数据丢失问题

![image](resources/images/6.jpg)

+ kafka的核心机制--ISR机制

    ISR是kafka自动维护和监控哪些follower及时的跟上了leader的数据同步
    
+ kafka写入的数据保证不丢失

![image](resources/images/7.jpg)

**若写入失败，则让生产者不停的重试，直到kafka恢复正常**

----------------------------------

### 三。[什么是消息队列](https://mp.weixin.qq.com/s?__biz=MzI4Njg5MDA5NA==&mid=2247485080&idx=1&sn=f223feb9256727bde4387d918519766b&chksm=ebd74799dca0ce8fa46223a33042a79fc16ae6ac246cb8f07e63a4a2bdce33d8c6dc74e8bd20&token=1439272449&lang=zh_CN&scene=21#wechat_redirect)

参考：[github-3y](https://github.com/ZhongFuCheng3y/3y)

+ 消息队列(MQ(Message Queue)-中间件

![](resources/images/7.jpg)

+（一）为什么要用消息队列

1.解耦 

    A提供B,C，D服务，B和C可能需要A提供某数据，但D不需要。则可以用消息队列方式，谁要谁取
    
 2.异步

![](resources/images/9.jpg)

3.削峰/限流

![](resources/images/10.jpg)

+（二）消息队列问题

    JDK实现的队列都是简单的内存队列
    
1。高可用--集群/分布式(不能够单机)

2。数据丢失问题

![](resources/images/11.jpg)

    Redis可以将数据持久化磁盘上，万一Redis挂了，还能从磁盘从将数据恢复过来。同样地，消息队列中的数据也需要存在别的地方，这样才尽可能减少数据的丢失。

3。消费者怎么得到消息队列的数据

    生产者将数据放到消息队列中，消息队列有数据了，主动叫消费者去拿(俗称push)
    
    消费者不断去轮训消息队列，看看有没有新的数据，如果有就消费(俗称pull)
    
----------------------------------

[Kafka简明教程](https://zhuanlan.zhihu.com/p/37405836)

[消息队列使用的四种场景介绍，有图有解析，一看就懂](https://zhuanlan.zhihu.com/p/55712984))

[消息队列设计精要](https://zhuanlan.zhihu.com/p/21479556)

[消息队列的使用场景是怎样的](https://www.zhihu.com/question/34243607)






