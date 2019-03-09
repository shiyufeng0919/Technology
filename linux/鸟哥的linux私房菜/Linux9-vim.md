# 第9章 vim编辑器

## vi编辑器

![](resources/images/84.jpg)

###按键说明:

+ 第一部份：一般指令模式可用的按钮说明，游标移动、复制粘贴、搜寻取代等

![](resources/images/85-1.jpg)
![](resources/images/85-2.jpg)
![](resources/images/85-3.jpg)
![](resources/images/85-4.jpg)

+ 第二部份：一般指令模式切换到编辑模式的可用的按钮说明

![](resources/images/86.jpg)

+ 第三部份：一般指令模式切换到指令列模式的可用按钮说明

![](resources/images/87.jpg)

### vim的暂存档、救援恢复与开启时的警告讯息

**使用vim编辑时， vim会在与被编辑的档案的目录下，再建立一个名为.filename.swp的档案**

若因异常情况，文件未来得及保存，则.filename.swp即发恢作用

![](resources/images/88.jpg)
![](resources/images/89.jpg)
![](resources/images/90.jpg)

### vim额外功能

+ 颜色

+ 区块选择

![](resources/images/91.jpg)

+ 多档案编辑

    vim file1 file2 ... fileN #同时编辑多个档案
    
    ![](resources/images/92.jpg)
    
    **示例：$ vim domain.csr domain.key**
    
    + :files #查看当前打开了哪些文件
    
    + 文件第一列输入「4yy」复制四列 (domain.csr)
    
    + :n #切换到下一个文档 (domain.key)->「G」最后一行 「p」粘贴 #「u」还原原本资料

+ 多视窗功能

 ![](resources/images/93.jpg)
 
+ 挑字补全功能

![](resources/images/94.jpg)

+ vim 环境设定与记录： ~/.vimrc, ~/.viminfo

![](resources/images/95-1.jpg)
![](resources/images/95-2.jpg)

**$ vim ~/.vimrc**
![](resources/images/96.jpg)


### vim其他注意事项

+ 中文编码问题

+ DOS 与Linux 的断行字元

+ 语系编码转换


------------------------


[鸟哥的linux私房菜](http://linux.vbird.org/linux_basic/0310vi.php)