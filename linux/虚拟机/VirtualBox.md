### VirtualBox扩展增强

因centos应用VBoxManage guestproperty enumerate vmName | grep "Net.*V4.*IP"命令

不显示虚拟Ip信息。因此需要安装VirtualBox扩展插件.同时解决共享目录到虚拟机问题。如下所示：

>[Centos Install VirtualBox Guest Additions](https://www.if-not-true-then-false.com/2010/install-virtualbox-guest-additions-on-fedora-centos-red-hat-rhel/)

    Step1: $vagrant ssh node1
    
    Step2: $su -  /  $sudo -i
    
    Step3: $yum update kernel*
    
    Step4: $reboot
    
    Step5: Install Guest Additions...
    
[下载镜像: VboxGuestAdditions_6.0.0_RC1.iso](http://download.virtualbox.org/virtualbox/6.0.0_RC1/)

[Index of /virtualbox所有版本](http://download.virtualbox.org/virtualbox/)

    Step6: $ wget http://download.virtualbox.org/virtualbox/6.0.0_RC1/VBoxGuestAdditions_6.0.0_RC1.iso
    
    Step7: $ mkdir /media/VirtualBoxGuestAdditions
    
    Step8: $ mount -r VBoxGuestAdditions_6.0.0_RC1.iso /media/VirtualBoxGuestAdditions
    
    Step9: Install following packages
    
    $ rpm -Uvh https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
    
    $ yum install gcc kernel-devel kernel-headers dkms make bzip2 perl
    
    Step10: Add KERN_DIR environment variable
    
    $ KERN_DIR=/usr/src/kernels/`uname -r`
    $ export KERN_DIR
    
    Step11: Install Guest Additions
    
    $ cd /media/VirtualBoxGuestAdditions
    
    $ ./VBoxLinuxAdditions.run
    
    此时执行命令：VBoxManage guestproperty enumerate vmName | grep "Net.*V4.*IP" 能够显示虚拟机ip信息
    
    





