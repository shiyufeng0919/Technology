### [Maven技术视频](https://www.bjsxt.com/down/8654.html)

--------

### Maven学习笔记

#### 一。Ant和分布式介绍

##### 1.Ant:项目构建工具

>Ant: 项目的编译，运行，打包过程均通过Ant管理 

    项目.project文件中，是Ant配置信息(高版本Eclipse)
    
    需求描述：现有项目A和项目B，需要将项目A和项目B建立联系如何实现？
    
    Ant实现：若A项目调用B项目，则需先将B项目打jar包，A项目引入B项目jar包
    
    问题描述：后期版本迭代，B项目需要重新打jar包，A项目需将原jar包删除掉，再重新引入。
    
    解决方案：Maven很好的解决了该问题

##### 2.分布式

>传统项目部署

    创建一个App项目，通过不同包区分不同模块；将项目打war包，布署到服务器tomcat/webapp
    
    问题：高负载(高并发/高访问量)下，一个tomcat会出现效率降低，甚至宕机
    
    解决方案：分布式部署
    
>分布式项目部署

    将一个大项目，拆分成若干独立的小项目，拆分后的项目分别部署到远端服务器(Tomcat)上。服务器间可相互通信.
    
**分布式项目即把传统的项目包(jar包)，换成一个单独的项目**

--------

#### 二。Maven简介和配置

>Maven：是基于Ant的构建工具(Ant有的功能，Maven均有，又额外添加了其他功能)

    集成插件：Eclipse&IDEA集成了Maven,不自定义，会应用IDEA默认Maven配置。也可自定义，引入。
    
>Maven运行原理图

[Maven远程仓库](https://mvnrepository.com/)

    (1)配置本地仓库：Maven先到本地仓库(计算机的一个文件夹)找相关依赖包，找不到，则到远程中央仓库(网上的一个地址)下载(远程仓库存储一些jar包)，并存储到本地仓库.
    
        <localRepository>  #本地仓库
        
        <mirros></mirros>  #镜像
        
        <profile></profile>  #修改jdk版本(默认1.4/1.5)
    
    (2)Maven下载镜像速度慢，一般会配置国内镜像
    
       Maven项目永远都从本地仓库查找jar包
    
    (3)保证JDK版本和开发环境一致。
    
 --------

#### 三。Maven项目创建

>坐标：通过坐标可精确确认哪个jar

    坐标组成：
    
    (1)GroupID:公司名，公司网址倒写
    
    (2)AritifactID:项目名称
    
    (3)Version:版本
    
>pom: project object model项目对象模型 (把一个项目project当作对象看待，通过maven构建工具，可让对象(项目)和对象(项目)间产生关系)

>新建->project->Maven Project

    package：项目的类型，最终会被打包成什么类型
    
    (1).jar : java项目
    
    (2).war:  web项目
    
    (3).pom: 逻辑父项目（只要一个项目有子项目，就必须是pom类型）
    
> maven项目目录结构(jar类型)

    src/main/java 真实目录的快捷目录，写java代码
    
    src/main/resources 快捷目录，存放配置文件。里面所有文件最终会被编译放入到classes类路径
    
    src/test/java 写测试java代码
    
    src/test/resources 测试配置文件
    
    pom.xml: maven的配置文件
    
       配置1:当前项目所依赖的其他项目或jar包或插件
       
            项目A引用项目B，pom.xml配置依赖<dependencies><dependencie>引入项目B的坐标</dependencie></dependencies>
            
            此时，项目B改代码，则项目A会快速依赖，也会更改。而且能够看到源码
            
            依赖谁即下载谁的jar包
            
    
>命令：maven->install 将项目打成jar包并存储到本地仓库
    

--------

#### 四。项目和项目之间的关系


##### 1.依赖关系 --<dependencies>

    下载需要的其他jar包
    

##### 2.继承关系（可以不是一个项目，可继承别的项目）--<parent>

    父项目是pom类型
    
    子项目jar或war,若子项目还是其他项目的父项目，则子项目也是pom类型
    
    有继承关系后，子项目中会出现<parent>标签（若子项目中<groupId>和<version>与父项目相同，则子项目中可不配置）
    
    父项目pom.xml中是看不到有哪些子项目，在逻辑上具有父子项目关系
    
##### 3.聚合关系(必定为一个项目) --<modules>

    前提：是继承关系，父项目会把子项目包含到父项目中，同时子项目的类型必须是Maven Module，而不是Maven project
    
    新建聚合项目的子项目时，需点击父项目右建新建maven module
    
    具有聚合关系的父项目，在pom.xml中用<modules>
    
    具有聚合关系的子项目，在pom.xml中<parent>
    
 **聚合项目和继承项目区别：**

    在语义上，聚合项目，父项目和子项目关系性较强；单纯继承项目，父项目和子项目关系性较弱
    
--------

#### 五。创建war类型项目 (必须为标准web项目结构)

##### 1.创建maven project,package=war

    创建maven project时，选择package类型为war

    在webapp文件夹下创建META-INF和WEB-INF文件夹及web.xml文件必须有。否则pom.xml中的<package>war</package>会报错
    
    在pom.xml中添加java ee相关的三个jar包
    
        javax.servlet  : 后端servlet
        
        javax.servlet.jsp :前端jsp
        
        tomcat插件

##### 2.导入web相关依赖包，否则创建jsp文件会报错

    <dependencies></dependencies>
    
--------

#### 六。使用Maven完成SSM练习--Maven项目较多宠大的时候使用

    New -> project -> maven project
    
--------

#### 七。Maven热部署

项目经验:[docker安装tomcat&部署javaweb程序](https://www.cnblogs.com/kaixinyufeng/p/9689982.html)

参考: [Maven+Tomcat实现热部署](https://blog.csdn.net/qq_32625839/article/details/81253564)

    Step1: 修改tomcat的conf/tomcat-users.xml配置文件，添加用户名，密码，权限。
        <role rolename="manager-gui" />
        <role rolename="manager-script" />
        <user username="tomcat" password="tomcat" roles="manager-gui,manager-script" />
        
    Step2: 重新启动tomcat (./startup),输入用户/密码
    
    Step3: Maven配置(pom.xml)
        <build>
            <plugins>
                <!-- 配置Tomcat插件 -->
                <plugin>
                    <groupId>org.apache.tomcat.maven</groupId>
                    <artifactId>tomcat7-maven-plugin</artifactId>
                    <configuration>
                        <!-- 
                            一般eclipse启动项目时候这里配置什么端口，访问项目的时候就是什么端口；用了热部署后，
                            是部署到目标tomcat里，因此这个port算是没用，访问时，是在tomcat的端口
                         -->
                        <port>8081</port>
                        <!-- 部署到ROOT下 -->
                        <path>/</path>
                        <!-- tomcat的地址和端口，manager/text是固定的 -->
                        <url>http://192.168.70.18:8080/manager/text</url>
                        <username>tomcat</username>
                        <password>tomcat</password>
                    </configuration>        
                </plugin>
            </plugins>
        </build>
        
    Step4:使用Maven命令进行部署
    
        clean tomcat7:redeploy #若为第一次部署，则为deploy，否则为redploy (可勾选Skip Tests跳过测试)
        
        clean tomcat7:redeploy -DskipTests #命令跳过测试
    
--------

#### 八。聚合项目演示 -后台

需求描述： 一个完整的项目，分前台(向用户展示数据)和后台(商家维护数据)

    创建父类项目: ego-parent (Maven package类型为pom)
    
    右击ego-parent项目，创建聚合项目(maven modules),web项目，则maven package类型为war
    
因：pojo类，为前后台共用，则可以将pojo类以聚合模块提出。其他项目依赖该modules

    右击父类项目,创建聚合项目(maven modules),java项目(供其他项目调用)，则maven package类型为jar
    
    项目前台/后台应用该pojo模块，则在pom.xml中以<dependency>pojo的坐标</dependency>方式引入即可

数据访问层：mybatis ,以xml文件形式书写sql

后台：模拟上传图片功能,经验：路径存储到数据库中，前端调用时：固定地址+路径即可

--------

#### 九。聚合项目演示 -前台

前台访问图片经验： 请求图片url ,如 :  http://localhost:8080/projectName/pic/abc.png （pic是存储文件夹）

则http://localhost:8080/projectName可以提取出来，配置到属性文件中(方便后续更换环境时，地址修改)

    如my.properties中存储 image.path=http://localhost:8080/projectName

/pic/abc.png可以存储到db中

则在service层，拿到图片路径时，需要修改图片路径

    (1)。在service层，以Spring注解@Value({image.path}) 取出配置文件路径值 (spring中需要引入该属性文件) (@Value在controller/service层均可以使用)
    
        @Value("${image.path}")
        private String path;
    
    (2)。循环遍历图片，修改原路径path+原path

--------

#### 十。实现图片轮播特效

前端技术

--------


参考学习：

[Maven教程](https://www.yiibai.com/maven/)   |   [Maven官网](http://maven.apache.org/)    |   [Maven仓库](https://mvnrepository.com/)

---------

**体胖还需勤跑步，人丑就该多读书!  ----开心玉凤**
