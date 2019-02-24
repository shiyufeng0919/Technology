###1.CentOS linux系统默认的shell是什么?  

答:bash

$ echo $SHELL  或 $ grep root /etc/passwd 

###2.请问echo $user的返回结果是什么？

$ cat test.sh #输出user=`whoami`
$ sh test.sh #执行脚本 (sh和bash执行脚本会新启动一个子shell,执行完成后退回父shell,因此子shell中执行返回的值无法保留)
$ echo $user #输出内容是? (3)

(1)当前用户 (2)oldboy (3)空(无内容输出)



