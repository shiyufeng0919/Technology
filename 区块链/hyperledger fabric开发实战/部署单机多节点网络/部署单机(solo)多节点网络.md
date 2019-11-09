# Hyperledger fabric开发实战   --kaixinyufeng 2019.10.26

## 四。部署单机(solo)多节点网络

### 1。下载平台二进制文件

>step1:下载平台二进制文件

    $ cd $GOPATH/src/github.com/hyperledger/fabric/

    $ mkdir aberic

    将下载的平台二进制文件bin/拷备到aberic/

[hyperledger fabric下载二进制Linux版本](https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/linux-amd64-1.1.0/)

[hyperledger fabric下载二进制](https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/)

    $ cd bin && chmod +x *  #设置的所有脚本可执行权限

    [root@node3 aberic]# tree -L 2
    .
    ├── bin
    │   ├── configtxgen         #「脚本」:工具=》 执行configtx.yaml文件以创建创世区块和通道认证文件
    │   ├── configtxlator
    │   ├── crypto-config       #执行./cryptogen generate --config=../crypto-config.yaml生成的证书文件目录
    │   ├── cryptogen           #「脚本」: 工具=》 ./cryptogen generate --config=../crypto-config.yaml
    │   ├── get-docker-images.sh
    │   ├── orderer
    │   └── peer
    ├── configtx.yaml          #用于: [创建orderer Genesis block]的配置文件
    └── crypto-config.yaml     #用于: [生成证书]的配置文件

### 2。生成证书文件

>step2:生成证书文件

    $ ./cryptogen generate --config=../crypto-config.yaml    #在bin/下会生成crypto-config目录

    [root@node3 crypto-config]# tree -L 3
        .
        ├── ordererOrganizations
        │   └── example.com
        │       ├── ca
        │       ├── msp
        │       ├── orderers   #orderer目录
        │       ├── tlsca
        │       └── users
        └── peerOrganizations
            ├── org1.example.com
            │   ├── ca
            │   ├── msp
            │   ├── peers     #peer目录
            │   ├── tlsca
            │   └── users
            └── org2.example.com
                ├── ca
                ├── msp
                ├── peers
                ├── tlsca
                └── users

>step3:使用configtxgen工具执行configtx.yaml文件创建orderer Genesis block(创世区块)

    (1).为configtxgen工具指定configtx.yaml文件路径，设置环境变量
        [root@node3 aberic]# export FABRIC_CFG_PATH=$PWD
        [root@node3 aberic]# echo $PWD
        /opt/gopath/src/github.com/hyperledger/fabric/aberic
        [root@node3 aberic]# mkdir channel-artifacts   #生成创世区块及channel认证文件放置于该目录下

    (2).根据configtx.yaml生成创世区块
        [root@node3 bin]# ./configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ../channel-artifacts/genesis.block

    注意：configtx.yaml中关于msp路径配置需根据实际msp目录位置做修改
    结果：
    [root@node3 channel-artifacts]# ls
    genesis.block  #生成创世区块

    (3).根据configtx.yaml生成通道(channel)认证文件
        [root@node3 bin]# ./configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ../channel-artifacts/syfchannel.tx -channelID syfchannel
        结果：
        [root@node3 channel-artifacts]# ls
        genesis.block「创世区块，orderer服务启动需要该配置」   syfchannel.tx「通道认证文件,peer可以执行channel的创建」

### 3。部署Orderer节点

>step4: 单机多节点部署采用solo启动方式，编写orderer节点启动文件docker-orderer.yaml
    
    $ cd $GOPATH/src/github.com/hyperledger/fabric/aberic
    $ mkdir deployment
    $ cd deployment
    $ vim docker-orderer.yaml

### 4。部署peer0.org1节点

>step5:编写peer节点启动文件docker-peer0.org1.yaml

> 修改 _sk值，目的：加载CA并登录CA用户
    $ cd /opt/gopath/src/github.com/hyperledger/fabric/aberic/bin/crypto-config/peerOrganizations/org1.example.com/ca
    $ ls
    [root@node3 ca]# ls
    66840bd901c2c4c4798e845094cd60cc9543a48ed950c1d2ff864120ac2e7417_sk  ca.org1.example.com-cert.pem

### 5。搭建Fabric网络

>step6: 启动orderer节点

    $ docker-compose -f docker-orderer.yaml up -d
    [root@node3 deployment]# docker ps
    CONTAINER ID        IMAGE                        COMMAND             CREATED             STATUS              PORTS                    NAMES
    c9e7d35a1faa        hyperledger/fabric-orderer   "orderer"           31 seconds ago      Up 30 seconds       0.0.0.0:7050->7050/tcp   orderer.example.com

>step7: 启动peer节点

    [root@node3 deployment]# docker-compose -f docker-peer.yaml up -d
    [root@node3 deployment]# docker ps
    CONTAINER ID        IMAGE                        COMMAND                  CREATED             STATUS              PORTS                                        NAMES
    1594f1946c18        hyperledger/fabric-tools     "/bin/bash"              2 minutes ago       Up 2 minutes                                                     cli
    c509d73904a8        hyperledger/fabric-peer      "peer node start"        2 minutes ago       Up 2 minutes        0.0.0.0:7051-7053->7051-7053/tcp             peer0.org1.example.com
    77e63b6c5626        hyperledger/fabric-ca        "sh -c 'fabric-ca-se…"   3 minutes ago       Up 2 minutes        0.0.0.0:7054->7054/tcp                       ca
    3b817b4ebc92        hyperledger/fabric-couchdb   "tini -- /docker-ent…"   3 minutes ago       Up 2 minutes        4369/tcp, 9100/tcp, 0.0.0.0:5984->5984/tcp   couchdb

>step8: channel的创建和加盟

**对peer节点操作需要依赖客户端(fabric-tools/cli)或SDK完成**

    # 1.进入cli客户端
    [root@node3 deployment]# docker exec -it cli bash

    root@1594f1946c18:/opt/gopath/src/github.com/hyperledger/fabric/peer# ls
    channel-artifacts  crypto

    root@1594f1946c18:/opt/gopath/src/github.com/hyperledger/fabric/peer# ls channel-artifacts/
    genesis.block  syfchannel.tx

    root@1594f1946c18:/opt/gopath/src/github.com/hyperledger/fabric/peer# ls crypto/
    ordererOrganizations  peerOrganizations

    # 2.创建channel
    root@1594f1946c18:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel create -o orderer.example.com:7050 -c syfchannel -t 50s -f ./channel-artifacts/syfchannel.tx
    2019-10-26 02:41:26.534 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 02:41:26.547 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 02:41:26.550 UTC [channelCmd] InitCmdFactory -> INFO 003 Endorser and orderer connections initialized
    2019-10-26 02:41:26.889 UTC [cli.common] readBlock -> INFO 004 Received block: 0

    root@1594f1946c18:/opt/gopath/src/github.com/hyperledger/fabric/peer# ls
    channel-artifacts  crypto  syfchannel.block 「生成通道文件」

    # 3.加入channel 「通过syfchannel.block文件加入syfchannel，加入channel即可安装、实例化、操作智能合约」
    root@1594f1946c18:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel join -b syfchannel.block
    2019-10-26 02:47:53.979 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 02:47:53.984 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 02:47:53.987 UTC [channelCmd] InitCmdFactory -> INFO 003 Endorser and orderer connections initialized
    2019-10-26 02:47:54.128 UTC [channelCmd] executeJoin -> INFO 004 Successfully submitted proposal to join channel

**注意：创建channel报错**

    Error: got unexpected status: BAD_REQUEST -- error applying config update to existing channel 'syfchannel': error authorizing update: error validating ReadSet: proposed update requires that key [Group]  /Channel/Application be at version 0, but it is currently at version 1

    问题原因：事先已成功创建channel，加入channel。关闭peer容器后，再启动peer容器，进入cli容器，再创建channel会报错

    解决方案：关闭orderer容器，重启
    
以上，一个最小单位的Fabric网络已成功搭建!!!

### 6。智能合约安装部署、实例化、测试

>step9: 智能合约安装部署

    root@1594f1946c18:/opt/gopath/src/github.com/hyperledger/fabric/peer# cd ..

    root@1594f1946c18:/opt/gopath/src/github.com/hyperledger/fabric# ls
    aberic  peer

    root@1594f1946c18:/opt/gopath/src/github.com/hyperledger/fabric# ls aberic/
    chaincode 「安装的智能合约路径」

    root@1594f1946c18:/opt/gopath/src/github.com/hyperledger/fabric# ls peer/
    channel-artifacts  crypto  syfchannel.block 「通道文件」

    #安装智能合约
    root@1594f1946c18:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode install -n syfchannel -p github.com/hyperledger/fabric/aberic/chaincode/go/chaincode_example02 -v 1.0
    2019-10-26 03:00:12.517 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 03:00:12.521 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 03:00:12.526 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 003 Using default escc
    2019-10-26 03:00:12.526 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 004 Using default vscc
    2019-10-26 03:00:12.898 UTC [chaincodeCmd] install -> INFO 005 Installed remotely response:<status:200 payload:"OK" >

>step10: 实例化智能合约

    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode instantiate -o orderer.example.com:7050 -C syfchannel -n syfchannel -c '{"Args":["init","A","10","B","10"]}' -P "OR ('Org1MSP.member')" -v 1.0
    2019-10-26 05:40:57.561 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 05:40:57.565 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 05:40:57.571 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 003 Using default escc
    2019-10-26 05:40:57.571 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 004 Using default vscc

**说明：**

    -P: 指定背书策略，-P "OR ('Org1MSP.member')"仅有Org1成员才具有背书能力。因此智能合约的invoke操作仅有Org1才会执行成功!!!
    -c: 指定参数，init为初始化方法，key=A，设置值10,key=B,设置值10

**实例化链码成功，即会启动链码容器**

    [root@node3 deployment]# docker ps
    CONTAINER ID        IMAGE                                                                                                               COMMAND                  CREATED             STATUS              PORTS                                        NAMES
    37bdee846a95        deployment-peer0.org1.example.com-syfchannel-1.0-a9dd647e9c1ff9969b84339ebf429274b8165ef3365706caad7ab4fa37fef9de   "chaincode -peer.add…"   12 minutes ago      Up 12 minutes                                                    deployment-peer0.org1.example.com-syfchannel-1.0


**注意：实例化链码时发生错误：**
    Error: could not assemble transaction, err proposal response was not successful, error code 500, msg error starting container: error starting container: API error (404): network aberic_default not found

    错误原因：docker-orderer.yaml与docker-peer.yaml文件在目录/deployment下面。因此:
        CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=aberic_default   #修改此处为deployment_default
        及
        networks:
            default:
                aliases:
                     - aberic    #修改此处为deployment
        
    #关闭容器
    $ docker-compose -f docker-peer.yaml down 
    $ docker-compose -f docker-orderer.yaml down

    #启动容器
    $ docker-compose -f docker-orderer.yaml up -d
    $ docker-compose -f docker-peer.yaml up -d

    #创建账本
    $ docker exec -it cli bash
    $ peer channel create -o orderer.example.com:7050 -c syfchannel -t 50s -f ./channel-artifacts/syfchannel.tx

    #加入账本
    $ peer channel join -b syfchannel.block

    #安装链码
    $ peer chaincode install -n syfchannel -p github.com/hyperledger/fabric/aberic/chaincode/go/chaincode_example02 -v 1.0

    #实例化链码
    $ peer chaincode instantiate -o orderer.example.com:7050 -C syfchannel -n syfchannel -c '{"Args":["init","A","10","B","10"]}' -P "OR ('Org1MSP.member')" -v 1.0

>step11: 查询智能合约(query)

    # 查询A值 (也可查询B值，均为10)
    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode query -C syfchannel -n syfchannel -c '{"Args":["query","A"]}'
    2019-10-26 05:58:16.663 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 05:58:16.668 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    10

>step12: 调用智能合约(invoke)

    # A向B转账5
    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode invoke -C syfchannel -n syfchannel -c '{"Args":["invoke","A","B","5"]}'
    2019-10-26 06:01:50.735 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:01:50.739 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:01:50.748 UTC [chaincodeCmd] InitCmdFactory -> INFO 003 Retrieved channel (syfchannel) orderer endpoint: orderer.example.com:7050
    2019-10-26 06:01:50.772 UTC [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 004 Chaincode invoke successful. result: status:200

    #再次查询A和B值,A为5,B为15
    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode query -C syfchannel -n syfchannel -c '{"Args":["query","B"]}'
    2019-10-26 06:02:33.485 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:02:33.491 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    15
    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode query -C syfchannel -n syfchannel -c '{"Args":["query","A"]}'
    2019-10-26 06:03:03.503 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:03:03.508 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    5

以上，智能合约测试完毕!!!

### 7。部署peer0.org2节点

>step13: 部署peer0.org2节点

    $ vim docker-peer0.org2.yaml

>step14: 修改cli容器全局变量，操作peer0.org2容器

    $ docker exec -it cli bash
    #修改全局变量,cli容器即可对peer0.org2操作
    CORE_PEER_ID=peer0.org2.example.com
    CORE_PEER_ADDRESS=peer0.org2.example.com:7051
    CORE_PEER_CHAINCODELISTENADDRESS=peer0.org2.example.com:7052
    CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org2.example.com:7051
    CORE_PEER_LOCALMSPID=Org2MSP
    CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp

>step15: peer0org2加入channel

    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# ls
    channel-artifacts  crypto  syfchannel.block

    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel join -b syfchannel.block
    2019-10-26 06:36:29.896 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:36:29.908 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:36:29.910 UTC [channelCmd] InitCmdFactory -> INFO 003 Endorser and orderer connections initialized
    2019-10-26 06:36:29.944 UTC [channelCmd] executeJoin -> INFO 004 Successfully submitted proposal to join channel

>step16: peer0org2安装智能合约

    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode install -n syfchannel -p github.com/hyperledger/fabric/aberic/chaincode/go/chaincode_example02 -v 1.0
    2019-10-26 06:38:55.741 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:38:55.747 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:38:55.755 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 003 Using default escc
    2019-10-26 06:38:55.755 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 004 Using default vscc
    2019-10-26 06:38:55.941 UTC [chaincodeCmd] install -> INFO 005 Installed remotely response:<status:200 payload:"OK" >

>step17: peer0org2查询智能合约

    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode query -C syfchannel -n syfchannel -c '{"Args":["query","A"]}'
    2019-10-26 06:39:47.882 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:39:47.886 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    5

>step18: peer0org2调用智能合约

    # peer0org2调用智能合约
    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode invoke -C syfchannel -n syfchannel -c '{"Args":["invoke","B","A","5"]}'
    2019-10-26 06:41:51.896 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:41:51.900 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:41:51.910 UTC [chaincodeCmd] InitCmdFactory -> INFO 003 Retrieved channel (syfchannel) orderer endpoint: orderer.example.com:7050
    2019-10-26 06:41:51.919 UTC [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 004 Chaincode invoke successful. result: status:200

    # 再次查询A资产，发现：A资产值依然为5，并没有变化
    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode query -C syfchannel -n syfchannel -c '{"Args":["query","A"]}'
    2019-10-26 06:41:56.171 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:41:56.175 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    5 「数值未变化」

**注意：实例化链码时，指定org1具备背书能力，而org2不具备背书能力，因此：来自org2对资产变更操作并未经过背书，即不具备任何效力。因org2组织加入了该channel且安装了合法的合约，即可以对区块链中的数据进行检索。**

    $ exit  #退出容器

    $ docker exec -it cli bash  #再次进入cli容器，此时默认是对peer0org1节点操作
    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# echo $CORE_PEER_ADDRESS
    peer0.org1.example.com:7051

    #执行B转账给A
    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode invoke -C syfchannel -n syfchannel -c '{"Args":["invoke","B","A","5"]}'
    2019-10-26 06:48:34.687 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:48:34.691 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:48:34.701 UTC [chaincodeCmd] InitCmdFactory -> INFO 003 Retrieved channel (syfchannel) orderer endpoint: orderer.example.com:7050
    2019-10-26 06:48:34.733 UTC [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 004 Chaincode invoke successful. result: status:200

    # 再查询资产A，发现A值已变化，验证Org1具有背书能力
    root@564d2b6dae7a:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode query -C syfchannel -n syfchannel -c '{"Args":["query","A"]}'
    2019-10-26 06:48:48.721 UTC [main] InitCmd -> WARN 001 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    2019-10-26 06:48:48.726 UTC [main] SetOrdererEnv -> WARN 002 CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable
    10 「数值变化」

--------------------------------------------------

>实例总结：

    [root@node3 fabric1.1.0]# pwd
    /opt/gopath/src/github.com/hyperledger/fabric1.1.0

    [root@node3 fabric1.1.0]# tree -L 2
        .
        ├── bin 「fabric平台二进制文件，官方可下载」
        │   ├── configtxgen    「生成创世块及channel配置文件工具」
        │   ├── configtxlator
        │   ├── crypto-config  「执行：./cryptogen generate --config=../crypto-config.yaml 生成证书文件目录」
        │   ├── cryptogen      「生成证书文件工具」
        │   ├── get-docker-images.sh
        │   ├── orderer
        │   └── peer
        ├── chaincode 「智能合约/链码，拷备自fabric/examples/chaincode/go/chaincode_example02/」
        │   └── go
        ├── channel-artifacts  
        │   ├── genesis.block   「执行:./configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ../channel-artifacts/genesis.block 生成该文件」
        │   └── syfchannel.tx   「执行:./configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ../channel-artifacts/syfchannel.tx -channelID syfchannel 生成该文件」
        ├── configtx.yaml       「生成创世区块(orderer服务启动必须)和通道配置文件(peer创建、加入channel)」
        ├── crypto-config.yaml  「生成证书配置文件」
        └── deployment  「部署orderer和peer0.org1和peer0.org2节点yaml文件」
            ├── docker-orderer.yaml
            ├── docker-peer0.org1.yaml
            └── docker-peer0.org2.yaml

     Step1:生成证书
     $ ./cryptogen generate --config=../crypto-config.yaml  #bin目录生成crypto-config证书目录

     Step2:生成创世块及通道配置文件
     $ export FABRIC_CFG_PATH=$PWD     #指定configtx.yaml文件目录
     $ ./configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ../channel-artifacts/genesis.block #在channel-artifacts目录下生成genesis.block 
     ./configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ../channel-artifacts/syfchannel.tx -channelID syfchannel #在hannel-artifacts目录下生成syfchannel.tx通证认证文件

     Step3:编写部署orderer及peer的yaml启动文件,启动order和peer服务
     $ vim docker-orderer.yaml
     $ vim docker-peer0.org1.yaml
     $ vim docker-peer0.org2.yaml

     $ docker-compose -f docker-orderer.yaml up -d      #启动排序服务节点
     $ docker-compose -f docker-peer0.org1.yaml up -d   #启动peer服务节点

     Step4:创建channel
     $ docker exec -it cli  #进入cli容器，fabric-tool客户端，用于调用fabric API
     $ peer channel create -o orderer.example.com:7050 -c syfchannel -t 50s -f ./channel-artifacts/syfchannel.tx

     Step5:加入channel
     $ peer channel join -b syfchannel.block

     Step6:安装链码
     $ peer chaincode install -n syfchannel -p github.com/hyperledger/fabric/fabric1.1.0/chaincode/go/chaincode_example02 -v 1.0

     Step7:实例化链码(此时会启动链码容器)
     $ peer chaincode instantiate -o orderer.example.com:7050 -C syfchannel -n syfchannel -c '{"Args":["init","A","10","B","10"]}' -P "OR ('Org1MSP.member')" -v 1.0

     Step8:查询链码
     $ peer chaincode query -C syfchannel -n syfchannel -c '{"Args":["query","A"]}'

     Step9:调用链码
     $ peer chaincode invoke -C syfchannel -n syfchannel -c '{"Args":["invoke","A","B","5"]}'

     Step10:动态加入账本
     #启动peer0.org2服务
     $ docker-compose -f docker-peer0.org2.yaml up -d

     #进入cli容器
     $ docker exec -it cli bash

     #修改全局环境变量，操作peer0.org2
     CORE_PEER_ID=peer0.org2.example.com
     CORE_PEER_ADDRESS=peer0.org2.example.com:7051
     CORE_PEER_CHAINCODELISTENADDRESS=peer0.org2.example.com:7052
     CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org2.example.com:7051
     CORE_PEER_LOCALMSPID=Org2MSP
     CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
     CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp

     #加入channel
     $ peer channel join -b syfchannel.block

     Step11:安装链码
     $ peer chaincode install -n syfchannel -p github.com/hyperledger/fabric/fabric1.1.0/chaincode/go/chaincode_example02 -v 1.0

     Step12:查询链码
     $ peer chaincode query -C syfchannel -n syfchannel -c '{"Args":["query","A"]}'

     Step13:调用链码(此处说明：实例化链码时，已指定org1具有背书能力，org2并不具有，因此更改数据无效)
     $ peer chaincode invoke -C syfchannel -n syfchannel -c '{"Args":["invoke","B","A","5"]}'

--------------------------------------------------

[参见:HyperLedger Fabric 1.1 手动部署单机单节点](https://www.cnblogs.com/aberic/p/8618556.html)

[Fabric命令手册](http://cw.hubwiz.com/card/c/fabric-command-manual/1/1/9/)


