### git学习网址

#### [git pro中文版](https://gitee.com/progit/)

#### [廖雪峰](https://www.liaoxuefeng.com/)

#### [Linux终端复用神器-Tmux使用梳理](http://www.cnblogs.com/kevingrace/p/6496899.html)

$ tmux #进入tmux模式  $ exit #退出tmux模式

#### [漫话：如何给女朋友解释什么是Git和GitHub](https://juejin.im/post/5ca16da36fb9a05e42555bf5)

##### Git:

    git其实就是一款分布式版本控制软件

##### Github

    GitHub是一个面向开源及私有软件项目的托管平台，因为只支持git 作为唯一的版本库格式进行托管，故名GitHub。

    在GitHub上面，你可以提交你自己写的代码（发微博）、start关注（粉）某人、关注（赞）某个项目、添加评论、Fork（转发）一个项目来自己修改

### git常用命令

##### git操作远程仓库相关命令

+ 检出仓库：$ git clone git://github.com/xx.git

+ 查看远程仓库：$ git remote -v

+ 添加远程仓库：$ git remote add [name][url]

+ 删除远程仓库：$ git remote rm [name]

    git push origin --delete branchA branchB #同时删除远程branchA和branchB分支

+ 修改远程仓库：$ git remote set-url --push[name][newUrl]

+ 拉取远程仓库：$ git pull [remoteName][localBranchName]

+ 推送远程仓库：$ git push [remoteName][localBranchName]

## git分支(branch)操作相关命令

+ 查看本地分支：$git branch

+ 查看远程分支：$ git branch -r

+ 查看所有分支：$ git branch -a

+ 创建本地分支：$ git branch [new-branchname]  ##注意新分支创建后不会自动切换为当前分支

+ 切换分支：$ git checkout [branchname]

+ 创建新分支并切换到新分支：$ git checkout -b [new-branchname]

+ 删除分支：$ git branch -d [branch-name]  #-d:只能删除已经参与了合并的分支，对于没有合并的分支是无法删除的。

+ 删除分支(强制)：$ git branch -D [branch-name] #-D:选项强制删除

## git版本(tag)操作相关命令

+ 查看版本：$ git tag

+ 创建版本：$ git tag [tag-name]

+ 删除版本：$ git tag -d [tag-name]

+ 查看远程版本：$ git tag -r

+ 创建远程版本(本地版本push到远程)：$ git push origin [remote-tag-name]

+ 删除远程版本：$ git push origin :refs/tags/[remote-tag-name]
