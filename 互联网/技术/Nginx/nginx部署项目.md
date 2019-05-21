### 实战一：应用nginx部署一个简单的前端web项目

前端项目框架: React

项目结构: 

    myproject --项目名称
    
      dist   --部署关键
      
         web
        
         index.html
      
      node_modules --依赖包
      
      server
      
      src  --源码
      
      ...
      
部署，拷备dist目录到/Users/shiyufeng/learn/myproject/

> url请求：http://localhost:8080/

修改nginx.conf配置文件

    location / {
                root   /Users/shiyufeng/learn/myproject/dist;
                index  index.html index.htm;
            }

> url请求：http://localhost:8080/admin

若请求路径不是根路径，如/admin,则需要在/Users/shiyufeng/learn/myproject/dist目录下面创建admin目录，并将web及index.html拷备到admin目录下面

    location /admin { #此时请求路径指向admin
      root /Users/shiyufeng/learn/myproject/dist; #root目录指向admin的上层目录
      index index.html;
    }

   **上述可理解为，查找dist/admin目录下的index.html文件**
   
----------------------------------------
   

