package DataDeal

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

/*
前端Post请求，传递参数：json格式

 [
        {
            "10.11.81.11": [
                {
                    "Spec": {
                        "Boot": null,
                        "CPU": {
                            "Count": 2
                        },
                        "Disks": [
                            {
                                "Controller": {
                                    "Device": 0,
                                    "Name": "IDE",
                                    "Port": 0,
                                    "Type": ""
                                },
                                "Format": "",
                                "Path": "/Users/shiyufeng/VirtualBox VMs/node1/centos-7-1-1.x86_64.vmdk",
                                "SizeMB": 0,
                                "Type": "",
                                "UUID": "e5fc0283-2c00-4dc3-93c3-a38974cb70bb"
                            }
                        ],
                        "Group": "",
                        "Memory": {
                            "SizeMB": 2048
                        },
                        "NICs": [
                            {
                                "BootPrio": 0,
                                "CableConnected": true,
                                "Index": 0,
                                "MAC": "525400261060",
                                "Mode": "nat",
                                "NetworkName": "",
                                "PromiscuousMode": "",
                                "Speedkbps": 0,
                                "Type": "82540EM"
                            },
                            {
                                "BootPrio": 0,
                                "CableConnected": true,
                                "Index": 0,
                                "MAC": "0800279E102F",
                                "Mode": "bridged",
                                "NetworkName": "",
                                "PromiscuousMode": "",
                                "Speedkbps": 0,
                                "Type": "82540EM"
                            }
                        ],
                        "Name": "node1",
                        "OSType": {
                            "Bit64": false,
                            "Description": "",
                            "FamilyDescription": "",
                            "FamilyID": "",
                            "ID": ""
                        },
                        "StorageControllers": [
                            {
                                "Bootable": "on",
                                "Instance": 0,
                                "Name": "IDE",
                                "PortCount": 2,
                                "Type": ""
                            }
                        ]
                    },
                    "UUID": "f6fa9977-fe1e-44c9-ae14-35450b52898f",
                    "node": "node1",
                    "status": "start",
                    "uuid": "f6fa9977-fe1e-44c9-ae14-35450b52898f"
                },
                {
                    "Spec": {
                        "Boot": null,
                        "CPU": {
                            "Count": 2
                        },
                        "Disks": [
                            {
                                "Controller": {
                                    "Device": 0,
                                    "Name": "IDE",
                                    "Port": 0,
                                    "Type": ""
                                },
                                "Format": "",
                                "Path": "/Users/shiyufeng/VirtualBox VMs/node2/centos-7-1-1.x86_64.vmdk",
                                "SizeMB": 0,
                                "Type": "",
                                "UUID": "ef408b4b-88eb-4950-972a-118170939620"
                            }
                        ],
                        "Group": "",
                        "Memory": {
                            "SizeMB": 2048
                        },
                        "NICs": [
                            {
                                "BootPrio": 0,
                                "CableConnected": true,
                                "Index": 0,
                                "MAC": "525400261060",
                                "Mode": "nat",
                                "NetworkName": "",
                                "PromiscuousMode": "",
                                "Speedkbps": 0,
                                "Type": "82540EM"
                            },
                            {
                                "BootPrio": 0,
                                "CableConnected": true,
                                "Index": 0,
                                "MAC": "080027D031BA",
                                "Mode": "bridged",
                                "NetworkName": "",
                                "PromiscuousMode": "",
                                "Speedkbps": 0,
                                "Type": "82540EM"
                            }
                        ],
                        "Name": "node2",
                        "OSType": {
                            "Bit64": false,
                            "Description": "",
                            "FamilyDescription": "",
                            "FamilyID": "",
                            "ID": ""
                        },
                        "StorageControllers": [
                            {
                                "Bootable": "on",
                                "Instance": 0,
                                "Name": "IDE",
                                "PortCount": 2,
                                "Type": ""
                            }
                        ]
                    },
                    "UUID": "49f147c1-757d-48cd-938b-658389d30209",
                    "node": "node2",
                    "status": "start",
                    "uuid": "49f147c1-757d-48cd-938b-658389d30209"
                }
            ]
        },
        {
            "127.0.0.1": [
                {
                    "Spec": {
                        "Boot": null,
                        "CPU": {
                            "Count": 2
                        },
                        "Disks": [
                            {
                                "Controller": {
                                    "Device": 0,
                                    "Name": "IDE",
                                    "Port": 0,
                                    "Type": ""
                                },
                                "Format": "",
                                "Path": "/Users/shiyufeng/VirtualBox VMs/node1/centos-7-1-1.x86_64.vmdk",
                                "SizeMB": 0,
                                "Type": "",
                                "UUID": "e5fc0283-2c00-4dc3-93c3-a38974cb70bb"
                            }
                        ],
                        "Group": "",
                        "Memory": {
                            "SizeMB": 2048
                        },
                        "NICs": [
                            {
                                "BootPrio": 0,
                                "CableConnected": true,
                                "Index": 0,
                                "MAC": "525400261060",
                                "Mode": "nat",
                                "NetworkName": "",
                                "PromiscuousMode": "",
                                "Speedkbps": 0,
                                "Type": "82540EM"
                            },
                            {
                                "BootPrio": 0,
                                "CableConnected": true,
                                "Index": 0,
                                "MAC": "0800279E102F",
                                "Mode": "bridged",
                                "NetworkName": "",
                                "PromiscuousMode": "",
                                "Speedkbps": 0,
                                "Type": "82540EM"
                            }
                        ],
                        "Name": "node1",
                        "OSType": {
                            "Bit64": false,
                            "Description": "",
                            "FamilyDescription": "",
                            "FamilyID": "",
                            "ID": ""
                        },
                        "StorageControllers": [
                            {
                                "Bootable": "on",
                                "Instance": 0,
                                "Name": "IDE",
                                "PortCount": 2,
                                "Type": ""
                            }
                        ]
                    },
                    "UUID": "f6fa9977-fe1e-44c9-ae14-35450b52898f",
                    "node": "node1",
                    "status": "stop",
                    "uuid": "f6fa9977-fe1e-44c9-ae14-35450b52898f"
                },
                {
                    "Spec": {
                        "Boot": null,
                        "CPU": {
                            "Count": 2
                        },
                        "Disks": [
                            {
                                "Controller": {
                                    "Device": 0,
                                    "Name": "IDE",
                                    "Port": 0,
                                    "Type": ""
                                },
                                "Format": "",
                                "Path": "/Users/shiyufeng/VirtualBox VMs/node2/centos-7-1-1.x86_64.vmdk",
                                "SizeMB": 0,
                                "Type": "",
                                "UUID": "ef408b4b-88eb-4950-972a-118170939620"
                            }
                        ],
                        "Group": "",
                        "Memory": {
                            "SizeMB": 2048
                        },
                        "NICs": [
                            {
                                "BootPrio": 0,
                                "CableConnected": true,
                                "Index": 0,
                                "MAC": "525400261060",
                                "Mode": "nat",
                                "NetworkName": "",
                                "PromiscuousMode": "",
                                "Speedkbps": 0,
                                "Type": "82540EM"
                            },
                            {
                                "BootPrio": 0,
                                "CableConnected": true,
                                "Index": 0,
                                "MAC": "080027D031BA",
                                "Mode": "bridged",
                                "NetworkName": "",
                                "PromiscuousMode": "",
                                "Speedkbps": 0,
                                "Type": "82540EM"
                            }
                        ],
                        "Name": "node2",
                        "OSType": {
                            "Bit64": false,
                            "Description": "",
                            "FamilyDescription": "",
                            "FamilyID": "",
                            "ID": ""
                        },
                        "StorageControllers": [
                            {
                                "Bootable": "on",
                                "Instance": 0,
                                "Name": "IDE",
                                "PortCount": 2,
                                "Type": ""
                            }
                        ]
                    },
                    "UUID": "49f147c1-757d-48cd-938b-658389d30209",
                    "node": "node2",
                    "status": "stop",
                    "uuid": "49f147c1-757d-48cd-938b-658389d30209"
                }
            ]
        }
    ]
*/

//处理参数

//1.
func WebStartAndStop(ctx *gin.Context) {

	logs.Info("--------web前端接口：启停虚拟机Start-----------")
	defer logs.Info("--------web前端接口：启停虚拟机End-----------")

	vmMap := []map[string][]VMList{}

	err := ctx.ShouldBindJSON(&vmMap)

	if err != nil {
		logs.Error("read request err:%s", err.Error())
		return
	}

	for _, v := range vmMap {

		for key, val := range v {

			if err := CallVmClientStopAndStartVM(key, val); err != nil {
				return
			}

		}

	}
}

//2.调用虚拟机客户端启动/停止虚拟机
func CallVmClientStopAndStartVM(ip string, vminfo []VMList) error {

	//建立http连接，停止/启动虚拟机
	url := GetVMClientUrl(ip, "operate")

	logs.Info("调用VM客户端的url为：", url)

	body, err := ConnectVmClientHttpPost(url, vminfo)

	if err != nil {
		logs.Error("连接vm-client错误:", err.Error())
		return err
	}

	logs.Info("响应string(body):", string(body))

	return nil
}

//3.获取vm客户端url
func GetVMClientUrl(ip, callType string) string {

	protocol := viper.GetString("vmclient.protocol")

	port := viper.GetString("vmclient.port")

	url := ""

	switch callType {
	case "list":
		url = protocol + "://" + ip + ":" + port + "/vmclient/GetVMListApi"
	case "operate":
		url = protocol + "://" + ip + ":" + port + "/vmclient/GetVMOperateApi"
	}

	return url
}

//4.POST请求
func ConnectVmClientHttpPost(url string, params []VMList) ([]byte, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	jsonParam, err := json.Marshal(params)

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonParam))

	if err != nil {
		logs.Error("发送http post连接发生错误###", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := client.Do(req)

	if err != nil {
		logs.Info("发送请求失败###", err)
		return nil, err
	}

	response, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logs.Error("返回响应错误:", err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	return response, nil
}
