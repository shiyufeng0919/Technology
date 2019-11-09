# Hyperledger fabric开发实战   --kaixinyufeng 2019.10.26

## 五。solo多机部署

### 1。生成证书

    [root@node3 bin]# ./cryptogen generate --config=../crypto-config.yaml
    org1.example.com
    org2.example.com

    [root@node3 bin]# ls
    configtxgen  configtxlator  crypto-config  cryptogen  get-docker-images.sh  orderer  peer

    [root@node3 bin]# cd ..

### 2。生成创世块「genesis.block」和通道channel配置文件「syfchannel.tx」

    # 设置 configtx.yaml文件目录
    [root@node3 fabric.solo]# export FABRIC_CFG_PATH=$PWD

    [root@node3 fabric.solo]# cd bin/

    # 生成创世块，用于orderer排序服务启动使用 「genesis.block」
    [root@node3 bin]# ./configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ../channel-artifacts/genesis.block

    # 生成channel配置文件，用于peer创建channel及加入channel使用 「syfchannel.tx」 
    [root@node3 bin]# ./configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ../channel-artifacts/syfchannel.tx -channelID syfchannel
    
    [root@node3 bin]# ls ../channel-artifacts/
    genesis.block  syfchannel.tx

### 3。部署orderer节点

    # (1).启动orderer节点
    [root@node1 deployment]# docker-compose -f docker-orderer.yaml up -d

### 4。部署peer0.org1节点，创建、加入channel，安装、实例化chaincode

    # 替换docker-peer0.org1.yaml文件中_sk值
    [root@node2 bin]# cd /opt/gopath/src/github.com/hyperledger/fabric.solo/bin/crypto-config/peerOrganizations/org1.example.com/ca
    [root@node2 ca]# ls
    3058fe0e17972bf6b3d97f18efeccc260417ab095dc727a04631aae0f630ee13_sk  ca.org1.example.com-cert.pem
    
    # (1).启动peer0.org1节点服务
    [root@node2 deployment]# docker-compose -f docker-peer0.org1.example.com up -d

    # (2).进入fabrictool客户端cli容器
    [root@node2 deployment]# docker exec -it cli bash
    root@474c8fe3061b:/opt/gopath/src/github.com/hyperledger/fabric/peer# ls
    channel-artifacts  crypto

    root@474c8fe3061b:/opt/gopath/src/github.com/hyperledger/fabric/peer# ls channel-artifacts/
    genesis.block  syfchannel.tx

    # (3).创建channel
    root@474c8fe3061b:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel create -o orderer.example.com:7050 -c syfchannel -t 50s -f ./channel-artifacts/syfchannel.tx

    # (4).加入channel
    root@474c8fe3061b:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel join -b syfchannel.block
    
    # (5).安装链码
    root@474c8fe3061b:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode install -n syfchaincode -p github.com/hyperledger/fabric/fabric.solo/chaincode/go/chaincode_example02 -v 1.0

    # (6).实例化链码
    root@474c8fe3061b:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode instantiate -o orderer.example.com:7050 -C syfchannel -n syfchaincode -c '{"Args":["init","A","10","B","10"]}' -P "OR ('Org1MSP.member')" -v 1.0

    # (7).查询链码(A值，B值分别为10)
    root@474c8fe3061b:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode query -C syfchannel -n syfchaincode -c '{"Args":["query","A"]}'

    # (8).调用链码(A值为5，B值为15)
    root@474c8fe3061b:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode invoke -C syfchannel -n syfchaincode -c '{"Args":["invoke","A","B","5"]}'

    # (9).查看已安装链码
    root@6feb9563cdb2:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode -C syfchannel list --installed
    Get installed chaincodes on peer:
    Name: syfchaincode, Version: 1.0, Path: github.com/hyperledger/fabric/fabric.solo/chaincode/go/chaincode_example02, Id: 833eb428faef7cee136f0d86b564fdea6d4a837d5e93ed1aeb81752d97e385b9

    # (10).查看已实例化链码
    root@6feb9563cdb2:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode -C syfchannel list --instantiated
    Get instantiated chaincodes on channel syfchannel:
    Name: syfchaincode, Version: 1.0, Path: github.com/hyperledger/fabric/fabric.solo/chaincode/go/chaincode_example02, Escc: escc, Vscc: vscc

### 5。部署peer1.org1节点，加入channel，安装chaincode
    
    # (1).部署peer1.org1
    [root@node2 deployment]# docker-compose -f docker-peer1.org1.yaml up -d

    # (2).拷备cli_peer0容器中的通道文件syfchannel.block拷备到cli_peer1容器
    # 拷备cli_peer0容器中通道文件syfchannel.block到物理机
    [root@node2 deployment]# docker cp cli_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/syfchannel.block .
    [root@node2 deployment]# ls
    docker-peer0.org1.yaml  docker-peer1.org1.yaml  syfchannel.block

    # 拷备物理机通道文件到cli_peer1容器中
    [root@node2 deployment]# docker cp syfchannel.block cli_peer1:/opt/gopath/src/github.com/hyperledger/fabric/peer
    [root@node2 deployment]# docker exec -it cli_peer1 bash
    root@f70ff562ba0c:/opt/gopath/src/github.com/hyperledger/fabric/peer# ls
    channel-artifacts  crypto  syfchannel.block「peer0.org1创建的channel文件」

### 6。部署peer0.org2节点，加入channel，安装chaincode



