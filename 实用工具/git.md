# git学习网址

## [git pro中文版](https://gitee.com/progit/)

## [廖雪峰](https://www.liaoxuefeng.com/)

## [Linux终端复用神器-Tmux使用梳理](http://www.cnblogs.com/kevingrace/p/6496899.html)

$ tmux #进入tmux模式  $ exit #退出tmux模式

# git常用命令

## git操作远程仓库相关命令

+ 检出仓库：$ git clone git://github.com/xx.git

+ 查看远程仓库：$ git remote -v

+ 添加远程仓库：$ git remote add [name][url]

+ 删除远程仓库：$ git remote rm [name]

+ 修改远程仓库：$ git remote set-url --push[name][newUrl]

+ 拉取远程仓库：$ git pull [remoteName][localBranchName]

+ 推送远程仓库：$ git push [remoteName][localBranchName]

## git分支(branch)操作相关命令

+ 查看本地分支：$git branch

+ 查看远程分支：$ git branch -r

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
