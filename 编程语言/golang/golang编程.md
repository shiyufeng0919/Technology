# golang实现文件的上传和下载

**[文件上传和下载](resources/文件上传和下载.go)**

# [JWT实现权限验证(Json web Token)](https://www.cnblogs.com/kaixinyufeng/p/9651304.html)

**Token是令牌，是用户身份的验证方式**

**简单的Token组成:**
+ uid(用户唯一的身份标识)
+ time(当前时间的时间戳)
+ sign(签名，由token的前几位+盐以哈希算法压缩成一定长的十六进制字符串，可以防止恶意第三方拼接token请求服务器)。还可以把不变的参数也放进token，避免多次查库

**JWT Token包含三个部分**
+ header: 告诉我们使用的算法和 token 类型 
+ Payload: 必须使用 sub key 来指定用户 ID, 还可以包括其他信息比如 email, username 等. 
+ Signature: 用来保证 JWT 的真实性. 可以使用不同算法 

