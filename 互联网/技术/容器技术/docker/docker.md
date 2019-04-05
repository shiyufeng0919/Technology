# docker-compose

### [docker-compose安装](https://docs.docker.com/compose/install/)

删除所有镜像: docker rmi `docker images -q`

删除所有容器：docker rm `docker ps -a -q`

按条件删除镜像：
//没有打标签
docker rmi `docker images -q | awk '/^<none>/ { print $3 }'`
//按关键字删除
docker rmi --force `docker images | grep doss-api | awk '{print $3}'`    //其中doss-api为关键字
