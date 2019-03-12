# [框架gin vs beego](https://www.imooc.com/video/18638)

go框架|支持MVC|支持正则路由|支持session|性能
:---:|:---:|:---:|:---:|:---:
beego|支持            |支持   | 支持|低
gin  |不支持，手动编写MVC|不支持 |不支持(依赖第三方支持session包)|高

## 一。beego框架 (大型项目应用)

**beego是典型的MVC框架**

## 二。gin框架

gin框架获取post请求的所有参数

-----------------------------------------

# golang实用的第三方库

## viper方便好用的Golang配置库(用于读取配置文件参数)

**[Github-viper](https://github.com/spf13/viper) | [Viper示例](resources/gosrc/opFiles/viper.txt)**

    安装:go get github.com/spf13/viper

## golang访问mysql数据库

### 一.标准库database/sql和mysql的驱动"github.com/go-sql-driver/mysql"

[示例](https://blog.csdn.net/lengyuezuixue/article/details/79148762)

### 二.sqlx框架

### 三.gorm框架

+ 支持的数据库有：mysql、postgre、sqlite、sqlserver

+ 文档 [github](https://github.com/jinzhu/gorm) | [gorm](http://gorm.io/)

### 四.xorm框架

## golang日志管理

+ log(原生)

+ [beego logs](resources/gosrc/logs/logs-beego.txt)

    go get github.com/astaxie/beego

+ logrus

## golang定时任务

+ [robfig/cron计划任务](https://www.cnblogs.com/zuxingyu/p/6023919.html)

    [go get -u github.com/robfig/cron]((https://studygolang.com/articles/10967))


## govendor包管理工具



-----------------------------------------

# golang项目应用

## [golang实现文件的上传和下载](resources/gosrc/opFiles/文件上传和下载.txt)

## [JWT实现权限验证(Json web Token)](https://www.cnblogs.com/kaixinyufeng/p/9651304.html)

## [cron定时任务](resources/gosrc/jobs/cron.txt)

**业务需求：定时调用prometheus接口获取监控指标.**

**实现方案：1。独立启一个goroutine,与main并行执行 2。通过配置文件配置执行的方法与时间 3。go get github.com/robfig/cron**

-----------------------------------------

# golang实用工具

## [json-> go struct](https://mholt.github.io/json-to-go/)

    可将json格式直接转成结构体(json->go从内向外拆分)