### springmvc面试题

参考：[腾讯这套SpringMvc面试题你懂多少（面试必备）](https://juejin.im/post/5cc2de6f5188252d9109875d)

#### 1.什么是 SpringMvc？

    SpringMvc 是 spring 的一个模块，基于 MVC 的一个框架，无需中间整合层来整

#### 2.Spring MVC 的优点

    1）它是基于组件技术的.全部的应用对象,无论控制器和视图,还是业务对象之类的都是组件.并且和Spring提供的其他基础结构紧密集成.
    2）不依赖于Servlet API(目标虽是如此,但是在实现的时候确实是依赖于Servlet的)
    3）可以任意使用各种视图技术,而不仅仅局限于 JSP
    4）支持各种请求资源的映射策略
    5）它应是易于扩展的
    
#### 3.SpringMVC 工作原理

    1）客户端发送请求到 DispatcherServlet
    2）DispatcherServlet 查询 handlerMapping 找到处理请求的 Controller
    3）Controller 调用业务逻辑后，返回 ModelAndView
    4）DispatcherServlet 查询 ModelAndView，找到指定视图
    5）视图将结果返回到客户端
   
#### 4.SpringMVC 流程

    1）用户发送请求至前端控制器 DispatcherServlet。
    2）DispatcherServlet 收到请求调用 HandlerMapping 处理器映射器。
    3）处理器映射器找到具体的处理器(可以根据 xml 配置、注解进行查找)，生成处理器及处理器拦截器(如果有则生成)一并返回给 DispatcherServlet。
    4）DispatcherServlet 调用 HandlerAdapter 处理器适配器。
    5）HandlerAdapter 经过适配调用具体的处理器(Controller，也叫后端控制器)
    6）Controller 执行完成返回 ModelAndView。
    7）HandlerAdapter 将 controller 执行结果 ModelAndView 返回给 DispatcherServlet。8）DispatcherServlet 将 ModelAndView 传给 ViewReslover 视图解析器。
    9）ViewReslover 解析后返回具体 View。
    10）DispatcherServlet 根据 View 进行渲染视图（即将模型数据填充至视图中）。
    11）DispatcherServlet 响应用户。

#### 5.SpringMvc 的控制器是不是单例模式,如果是,有什么问题,怎么解决？

    是单例模式,所以在多线程访问的时候有线程安全问题,不要用同步,会影响性能的,解方案是在控制器里面不能写字段。
    
#### 6.如果你也用过 struts2.简单介绍下 springMVC 和 struts2 的区别有哪些?

    1）springmvc 的入口是一个 servlet 即前端控制器，而 struts2 入口是一个 filter 过虑器
    2）springmvc 是基于方法开发(一个url对应一个方法)，请求参数传递到方法的形参，设计为单例或多例(建议单例)，struts2是基于类开发，传递参数是通过类的属性，只能计为多例。
    3）Struts 采用值栈存储请求和响应的数据，通过 OGNL 存取数据，springmvc 通过参析器是将 request 请求内容解析，并给方法形参赋值，将数据和视图封装成 ModelAnd对象，最后又将 ModelAndView 中的模型数据通过 reques 域传输到页面。Jsp 视图解析认使用 jstl。
     
#### 7.SpingMvc 中的控制器的注解一般用那个,有没有别的注解可以替代

    一般用@Conntroller 注解,表示是表现层,不能用用别的注解代替
    
#### 8.@RequestMapping 注解用在类上面有什么作用？

    是一个用来处理请求地址映射的注解，可用于类或方法上。用于类上，表示类有响应请求的方法都是以该地址作为父路径。
     
#### 9.怎么样把某个请求映射到特定的方法上面

    直接在方法上面加上注解@RequestMapping,并且在这个注解里面写上要拦截的路
    
#### 10.如果在拦截请求中,我想拦截 get 方式提交的方法,怎么配置？

    可以在@RequestMapping 注解里面加上 method=RequestMethod.GET
    
#### 11.怎么样在方法里面得到 Request,或者 Session？

    直接在方法的形参中声明 request,SpringMvc 就自动把 request 对象传
    
#### 12.


