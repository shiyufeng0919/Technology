### 概念简介

> Iaas

    IaaS（Infrastructure as a Service），即基础设施即服务
    
    是指把IT基础设施作为一种服务通过网络对外提供。在这种服务模型中，用户不用自己构建一个数据中心，而是通过租用的方式来使用基础设施服务，包括服务器、存储和网络等。
    
    Iaas提供给消费者的服务是对所有计算基础设施的利用，包括处理CPU、内存、存储、网络和其它基本的计算资源，用户能够部署和运行任意软件，包括操作系统和应用程序。
    
> Paas

    Platform-as-a-Service（平台即服务）提供给消费者的服务是把客户采用提供的开发语言和工具（例如Java，python, .Net等）开发的或收购的应用程序部署到供应商的云计算基础设施上去。
    
    客户不需要管理或控制底层的云基础设施，包括网络、服务器、操作系统、存储等，但客户能控制部署的应用程序，也可能控制运行应用程序的托管环境配置；

> Saas

    Software-as-a-Service（软件即服务）提供给客户的服务是运营商运行在云计算基础设施上的应用程序，用户可以在各种设备上通过客户端界面访问，如浏览器。消费者不需要管理或控制任何云计算基础设施，包括网络、服务器、操作系统、存储等等；
    
    
> Iaas & Paas & Saas区别

    SaaS 是软件的开发、管理、部署都交给第三方，不需要关心技术问题，可以拿来即用。普通用户接触到的互联网服务，几乎都是 SaaS，下面是一些例子。
    
    > 客户管理服务 Salesforce
    
    > 团队协同服务 Google Apps
    
    > 储存服务 Box
    
    > 储存服务 Dropbox
    
    > 社交服务 Facebook / Twitter / Instagram
    
    PaaS 提供软件部署平台（runtime），抽象掉了硬件和操作系统细节，可以无缝地扩展（scaling）。开发者只需要关注自己的业务逻辑，不需要关注底层。下面这些都属于 PaaS。
    
    > Heroku
    
    > Google App Engine
    
    > OpenShift
    
    IaaS 是云服务的最底层，主要提供一些基础资源。它与 PaaS 的区别是，用户需要自己控制底层，实现基础设施的使用逻辑。下面这些都属于 IaaS。
    
    > Amazon EC2
    
    > Digital Ocean
    
    > RackSpace Cloud
    
----------------------------

### 项目简介

Iaas层管理：主要通过前端监控虚拟机信息，包括：虚拟机Ip和资源(CPU,内存，磁盘大小等)；并可能过前端控制虚拟机启动/停止

> 项目开发
    
    开发工具：Goland  开发语言: golang
    
    框架：Gin框架(路由) + beego日志框架(记录日志) + viper读取配置文件 + govendor包管理工具 + ticker定时任务 + 
    
    应用程序包：
    1.go get github.com/gin-gonic/gin
    2.go get github.com/astaxie/beego/logs
    3.go get github.com/spf13/viper
    4.go get -u -v github.com/kardianos/govendor
    5.golang自身time包
    6.go get -u github.com/uruddarraju/virtualbox-go
    


-----------------------------------

### golang技术总结

[VBoxManage命令](https://www.virtualbox.org/manual/ch08.html#vboxmanage-list)

[golang-virtualbox创建虚拟机/控制虚拟机启停](https://github.com/uruddarraju/virtualbox-go)

[GoLang中 json、map、struct 之间的相互转化](https://www.cnblogs.com/liang1101/p/6741262.html)

[golang程序执行命令：VBoxManage list vms 处理返回的结果](resources/gosrc/DataDeal/VBoxManageListVms.go)

[VBoxManage list vms返回前端数据作为参数](resources/gosrc/DataDeal/OperateVM.go)

[golang time包ticker使用]






