## Centos下NFS服务安装，配置，使用

### 一。什么是NFS

>NFS是Network File System的缩写，即网络文件系统.

>客户端通过挂载的方式将NFS服务器端共享的数据目录挂载到本地目录下

### 二。NFS为什么需要RPC

    因为NFS支持的功能很多，不同功能会使用不同程序来启动，因此，NFS对应的功能所对应的端口无法固定。
    端口不固定造成客户端与服务端之间的通信障碍，所以需要RPC来从中帮忙。
    
    NFS启动时会随机取用若干端口，然后主动向RPC服务注册取用相关端口和功能信息，RPC使用固定端口111来监听来自NFS客户端的请求，
    并将正确的NFS服务端口信息返回给客户端，这样客户端与服务端就可以进行数据传输了。

### 三。NFS工作流程

    1、由程序在NFS客户端发起存取文件的请求，客户端本地的RPC(rpcbind)服务会通过网络向NFS服务端的RPC的111端口发出文件存取功能的请求。
    
    2、NFS服务端的RPC找到对应已注册的NFS端口，通知客户端RPC服务。
    
    3、客户端获取正确的端口，并与NFS daemon联机存取数据。
    
    4、存取数据成功后，返回前端访问程序，完成一次存取操作。
    
    所以无论客户端，服务端，需要使用NFS，必须安装RPC服务。

>NFS的RPC服务，在Centos5下名为portmap,Centos6下名称为rpcbind。

### 四。服务端安装NFS服务

####1。安装NFS

    $ yum -y install nfs-utils rpcbind
    
#### 2。服务端：在NFS服务端上创建共享目录/export/nfs并设置权限(node1节点)

    $ mkdir -p /export/nfs
    $ chmod 666 /export/nfs/
    
    ##编辑export文件,并添加内容，NFS服务装在192.168.1.31(node1)上,应用192.168.1.32-33(node2-node3)作客户端
    vim /etc/exports 
    [root@node1 ~]# cat /etc/exports
    /export/nfs 192.168.1.0/24(rw,no_root_squash,no_all_squash,sync)
    
    **上述表示将/export/nfs这个目录共享给：192.168.1.0代表192.168.1.*这些客户机**
    
    ---
    
    ##上述配置内容说明
    NFS共享的目录 NFS客户端地址1(参数1,参数2,...) 客户端地址2(参数1,参数2,...)
    
    >NFS共享的目录:要用绝对路径，可被nfsnobody读写。
    >NFS客户端地址:
        指定IP: 192.168.0.1
        指定子网所有主机: 192.168.0.0/24
        指定域名的主机: test.com
        指定域名所有主机: *.test.com
        所有主机: *
    >参数:
        ro：目录只读
        rw：目录读写
        sync：将数据同步写入内存缓冲区与磁盘中，效率低，但可以保证数据的一致性
        async：将数据先保存在内存缓冲区中，必要时才写入磁盘
        all_squash：将远程访问的所有普通用户及所属组都映射为匿名用户或用户组(nfsnobody)
        no_all_squash：与all_squash取反(默认设置)
        root_squash：将root用户及所属组都映射为匿名用户或用户组(默认设置)
        no_root_squash：与rootsquash取反
        anonuid=xxx：将远程访问的所有用户都映射为匿名用户，并指定该用户为本地用户(UID=xxx)
        anongid=xxx：将远程访问的所有用户组都映射为匿名用户组账户
       
    ---
    
    ##配置生效 
    $ exportfs -r
    
    ##或重新加载nfs配置
    $ exportfs -rv
    
    ---
    
    ##启动rpc
    [root@node1 ~]# service rpcbind start
    Redirecting to /bin/systemctl start rpcbind.service
    
    ##启动nfs
    [root@node1 ~]# service nfs start
    Redirecting to /bin/systemctl start nfs.service
    
    ---
    
    ##查看 RPC 服务的注册状况
    [root@node1 ~]# rpcinfo -p
        program vers proto   port  service
            100000    4   tcp    111  portmapper
            100000    3   tcp    111  portmapper
            100000    2   tcp    111  portmapper
            100000    4   udp    111  portmapper
            100000    3   udp    111  portmapper
            100000    2   udp    111  portmapper
            100005    1   udp  20048  mountd
            100005    1   tcp  20048  mountd
            100024    1   udp  46656  status
            100024    1   tcp  60667  status
            100005    2   udp  20048  mountd
            100005    2   tcp  20048  mountd
            100005    3   udp  20048  mountd
            100005    3   tcp  20048  mountd
            100003    3   tcp   2049  nfs
            100003    4   tcp   2049  nfs
            100227    3   tcp   2049  nfs_acl
            100003    3   udp   2049  nfs
            100003    4   udp   2049  nfs
            100227    3   udp   2049  nfs_acl
            100021    1   udp  55226  nlockmgr
            100021    3   udp  55226  nlockmgr
            100021    4   udp  55226  nlockmgr
            100021    1   tcp  46760  nlockmgr
            100021    3   tcp  46760  nlockmgr
            100021    4   tcp  46760  nlockmgr

    ---
    
    ###查看服务端抛出的共享目录信息
    [root@node1 ~]# showmount -e
    Export list for node1:
    /export/nfs 192.168.1.0/24
    
    ###查看rpcbind状态
    [root@node1 ~]# systemctl status rpcbind
    ● rpcbind.service - RPC bind service
       Loaded: loaded (/usr/lib/systemd/system/rpcbind.service; enabled; vendor preset: enabled)
       Active: active (running) since Mon 2019-05-13 06:53:04 UTC; 13min ago
     Main PID: 10137 (rpcbind)
        Tasks: 1
       Memory: 612.0K
       CGroup: /system.slice/rpcbind.service
               └─10137 /sbin/rpcbind -w
    
    May 13 06:53:04 node1 systemd[1]: Stopped RPC bind service.
    May 13 06:53:04 node1 systemd[1]: Starting RPC bind service...
    May 13 06:53:04 node1 systemd[1]: Started RPC bind service.
    
### 五。客户端配置(node2节点)

    ##1.安装nfs-utils客户端
    $ yum -y install nfs-utils
    
    ##2.创建挂载目录(挂载到服务端的/export/nfs目录下)
    $ mkdir /testnfs
    
    ##查看服务器抛出的共享目录信息
    [root@node2 ~]# showmount -e 192.168.1.31
    Export list for 192.168.1.31:
    /export/nfs 192.168.1.31/24
    
    选项与参数：
    -a ：显示目前主机与客户端的 NFS 联机分享的状态；
    -e ：显示某部主机的 /etc/exports 所分享的目录数据。
    
    ##挂载
    [root@node2 ~]# mount -t nfs 192.168.1.31:/export/nfs /testnfs/
    
    ##查看挂载结果
    [root@node2 ~]# df -h
    Filesystem                Size  Used Avail Use% Mounted on
    /dev/sda1                  40G  4.9G   36G  13% /
    ......
    192.168.1.31:/export/nfs   40G  4.6G   36G  12% /testnfs

### 五。挂载测试

    ##node2节点(NFS客户端)
    [root@node2 testnfs]# echo "test nfs" > test-nfs.txt
    
    ##切换到node1节点查看(NFS服务端)
    [root@node1 k8s]# ls /export/nfs/
    test-nfs.txt
    
    上述：已成功挂载
    
    ##node1节点(NFS服务端)
    [root@node1 nfs]# echo "server nfs" > server.txt
    
    ##node2节点(NFS客户端)查看
    [root@node2 testnfs]# ls
    server.txt  test-nfs.txt
    
### 六。卸载已挂载的NFS
    
    ###卸载已挂载的目录
    [root@node2 /]# umount /testnfs/
    
    ###查看挂载目录
    [root@node2 /]# df -h
    
    
>yum卸载软件： yum remove nfs-utils rpcbind


--------------------------------------------    

参考：

1.[centos安装NFS服务](https://www.cnblogs.com/hwp0710/p/7942222.html)

2.[centos7下NFS使用与配置](https://www.cnblogs.com/jkko123/p/6361476.html?utm_source=itdadao&utm_medium=referral)