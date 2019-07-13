#### <dependencies> maven项目的依赖关系

    即把另一个项目的jar包引入到当前项目。自动下载另一个项目依赖的其他项目
    
    //子类能够继承父类的dependencies
    <dependencies>
        <dependency>
            <groupId>junit</groupId>
            <artifactId>junit</artifactId>
            <version>4.12</version>
            <scope>test</scope> #范围：provided编译期生效，test测试期生效
        </dependency>
        
        <dependency>
        </dependency>
    </dependencies>
    
#### <dependencyManagement>管理依赖

    写在父项目中，作用：声明可能使用到的所有jar,子项目中只需要有<groupId>和<artifactId>,<version>继承父项目

    //若不希望子项目继承父项目dependencies，可以引入dependencyManagement，此时需手动引入;也可将公共应用dependencyManagement管理，所有子项目均可继承，不必再重写
    <dependencyManagement>
        <dependencies>
        </dependencies>
    </dependencyManagement>
 
#### <properties>属性配置

    将项目版本更好管理
    
    <properties>
        <xxx>1.0.0</xxx>  #可在下述以${xxx}引入该值，如spring版本一致，可配置版本号属性。更方便管理
    </properties>
    
#### <parent> maven项目的继承关系(父项目和子项目)--可引入其他项目
    
    <parent>
        <groupId>com.xx</groupId>
        <artifactId>A</artifactId>
        <version>1.0.0</version>
    </parent>

#### <modules> maven项目的聚合关系 --必定为同一个项目

    父项目将子项目包含进来,优点：父项目包含了哪些子项目。告诉maven，子项目需要被父项目引进来

    <modules>
      <module>A</module>
    </modules>

#### <build> 添加插件

    如tomcat配置，则不会应用本地tomcat,会应用该tomcat插件，可控制tomcat的端口号和项目发布到tomcat的名称.分布式思想

    <build>
      <plugins>
        <plugin>
          坐标
          <configuration>
            <port>tomcat端口号</port>
            <path>项目发布到tomcat名称</path> //此处若为/，则相当于把项目发布名称为ROOT(浏览器访问:http://localhost:8080/ 为tomcat的ROOT),如果是/xx，则项目发布名称为xx
          </configuration>
        </plugin>
      </plugins>
    </build>
    
**启动插件：Run As->maven build->Goals:clean tomcat7:run**

    tomcat加7，否则默认为tomcat6