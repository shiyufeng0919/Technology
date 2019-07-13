#### Spring注解

>[@Value注解](https://blog.csdn.net/woheniccc/article/details/79804600)

    从属性配置文件中取值,如application.properties中配置server.port=8080，则在controller/service层均可以将该值取出
    
    @Value("${server.port}")
    private volatile String ports;