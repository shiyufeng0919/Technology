Nginx命令参数

    nginx -t 测试配置是否正确
    nginx -s reload 加载最新配置
    nginx -s stop 立即停止
    nginx -s quit 优雅停止
    nginx -s reopen 重新打开日志
    kill -USR2 cat /usr/local/nginx/logs/nginx.pid 快速重启