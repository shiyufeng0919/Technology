### [GO语言标准库文档中文版](https://studygolang.com/pkgdoc)

-----------------------------------------

# [框架gin vs beego](https://www.imooc.com/video/18638)

go框架|支持MVC|支持正则路由|支持session|性能
:---:|:---:|:---:|:---:|:---:
beego|支持            |支持   | 支持|低
gin  |不支持，手动编写MVC|不支持 |不支持(依赖第三方支持session包)|高

### 一。beego框架 (大型项目应用)

**beego是典型的MVC框架**

### 二。gin框架

gin框架获取post请求的所有参数

-----------------------------------------

# golang实用的第三方库

### viper方便好用的Golang配置库(用于读取配置文件参数)

**[Github-viper](https://github.com/spf13/viper) | [Viper示例](resources/gosrc/opFiles/viper.txt)**

    安装:go get github.com/spf13/viper

### golang访问mysql数据库

#### 一.标准库database/sql和mysql的驱动"github.com/go-sql-driver/mysql"

[示例](https://blog.csdn.net/lengyuezuixue/article/details/79148762)

#### 二.sqlx框架

#### 三.[gorm框架](http://gorm.book.jasperxu.com/advanced.html#sb)

[示例](resources/gosrc/orm/gorm.txt)

+ 支持的数据库有：mysql、postgre、sqlite、sqlserver

+ 文档 [github](https://github.com/jinzhu/gorm) | [gorm](http://gorm.io/)

#### 四.xorm框架

### golang日志管理

+ log(原生)

+ [beego logs](resources/gosrc/logs/logs-beego.txt)

    go get github.com/astaxie/beego

+ logrus

### golang定时任务

+ [robfig/cron计划任务](https://www.cnblogs.com/zuxingyu/p/6023919.html)

  [go get -u github.com/robfig/cron](https://studygolang.com/articles/10967)


### govendor包管理工具

-----------------------------------------

# golang项目应用

### [golang实现文件的上传和下载](resources/gosrc/opFiles/文件上传和下载.txt)

### [JWT实现权限验证(Json web Token)](https://www.cnblogs.com/kaixinyufeng/p/9651304.html)

### [cron定时任务](resources/gosrc/jobs/cron.txt)

### [http-get请求](resources/gosrc/http/http-get.txt)

### [http-post请求](resources/gosrc/http/http-post.txt)

### [https请求-验证证书](resources/gosrc/http/http-put.txt)

### [golang读写文件](https://www.jianshu.com/p/7790ca1bc8f6)

### [golang中struct、json、map互相转化](https://blog.csdn.net/xiaoquantouer/article/details/80233177)

[开心玉凤](https://github.com/shiyufeng0919/goLang/blob/master/src/basicGrammer/%E5%B7%A5%E4%BD%9C%E5%B8%B8%E7%94%A8/struct%26Json%26map%26yaml.go)

### [golang yaml与json配置文件互转](https://studygolang.com/articles/15148)

[参考](https://studygolang.com/articles/2946)

### [websocket-golang.org/x/net/websocke](./resources/gosrc/websocket/websocket示例1.txt)

### [Centrifugo(实时消息服务器)]()

## golang+redis

1.[博文:Go实战--golang中使用redis(redigo和go-redis/redis)](https://blog.csdn.net/wangshubo1989/article/details/75050024)
2.[redigo官方文档](https://godoc.org/github.com/garyburd/redigo/redis)

-----------------------------------------

# golang实用工具

### [json-> go struct](https://mholt.github.io/json-to-go/)

    可将json格式直接转成结构体(json->go从内向外拆分)