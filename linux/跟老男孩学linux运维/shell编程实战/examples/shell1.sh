#! /bin/bash
#清除日志脚本,版本2
LOG_DIR=/var/log
ROOT_UID=0  #<==$UID为0的用户，即root用户
#脚本需要使用root用户权限来运行，因此，对当前用户进行判断，对不合要求的用户给出友好提示，并终止程序运行
if [ "$UID" -ne "$ROOT_UID" ]    #<==如果当前用户不是root,则不允许执行脚本。注意if后面空格，否则报错
   then
     echo "Must be root to run this script." #<==给出提示后退出
     exit 1 #<==退出脚本
fi 
#如果切换到指定目录不成功，则给出提示，并终止程序运行
cd $LOG_DIR  || {
  echo "cannot change to necessary directory."
  exit 1
}
#经过上述两个判断后，此处的用户权限和路径应该就是对的了，只有清空成功，才打印成功提示
cat /dev/null>messges && {
  echo "Logs cleaned up."
  exit 0 #退出前返回0表示成功，返回1表示失败
}
  echo "Logs cleaned up fail."
  exit 1

##思路
##1。必须是root用户才可执行脚本，否则给出友好提示并退出
##2。成功切换目录，否则给出友好提示并终止脚本运行
##3。清理日志，若成功则给出正确提示
##4。通关或失败，分别给出相应提示
