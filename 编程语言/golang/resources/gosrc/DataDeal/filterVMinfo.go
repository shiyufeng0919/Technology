package DataDeal

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"net"
	"strings"
)

var newVmList []interface{}

func filterVmInfo(vmInfo, physicalIp string) ([]map[string][]Vmlist, error) {

	vmInfo="k8s-node1"

	physicalIp="127.0.0.1"

	vmIp := ""

	vmName := ""

	var err error

	if vmInfo != "" { //根据虚拟机IP搜索

		//ipReg := `^((0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])\.){3}(0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])$`
		//
		//match, _ := regexp.MatchString(ipReg, vmInfo)
		//
		//if !match {
		//
		//	vmName = vmInfo
		//
		//} else {
		//
		//	vmIp = vmInfo
		//}

		address := net.ParseIP(vmInfo) //校验是否为合法IP

		if address != nil {

			vmName = vmInfo //合法IP

		} else {

			vmIp = vmInfo //虚拟机名称
		}
	}

	filterRs, err := filterVMData(physicalIp, vmName, vmIp)

	if err != nil {
		return nil, err
	}

	return filterRs, nil
}

func filterVMData(physicalIp, vmName, vmIp string) ([]map[string][]Vmlist, error) {

	vmMap := []map[string][]Vmlist{} //源数据

	newVmMap := []map[string][]Vmlist{} //存储过滤后的数据

	if newVmList == nil {
		logs.Error("待过滤的虚拟机列表数据为空")
		return nil, errors.New("虚拟机列表数据为空")
	}

	strJson, err := json.Marshal(newVmList)

	if err != nil {
		logs.Error("NewVmList - > json 异常:", err.Error())
		return nil, err
	}

	err = json.Unmarshal(strJson, &vmMap)

	if err != nil {
		logs.Error("json - > map异常:", err.Error())
		return nil, err
	}

	//循环map,处理数据

	for _, v := range vmMap {

		for k1, v1Arr := range v {

			var strArr []Vmlist

			if physicalIp != "" {

				if k1 == physicalIp {

					if vmName != "" {
						//数组
						for _, v2Struct := range v1Arr {

							if strings.Contains(v2Struct.Node, vmName) { //包含虚拟机名称

								strArr = append(strArr, v2Struct) //组装新数组

							}

						}
					}

					if vmIp != "" {

						for _, v2Struct := range v1Arr {

							if len(v2Struct.VMIp) > 0 {

								for _, vkv := range v2Struct.VMIp {

									if strings.Trim(vkv, " ") == vmIp {
										strArr = append(strArr, v2Struct)
									}

								}
							}
						}
					}

				}

			} else {

				if vmName != "" {
					//数组
					for _, v2Struct := range v1Arr {

						if strings.Contains(v2Struct.Node, vmName) { //包含虚拟机名称

							strArr = append(strArr, v2Struct) //组装新数组
						}
					}
				}

				if vmIp != "" {

					for _, v2Struct := range v1Arr {

						if len(v2Struct.VMIp) > 0 {

							for _, vkv := range v2Struct.VMIp {

								if strings.Trim(vkv, " ") == vmIp {
									strArr = append(strArr, v2Struct)
								}

							}
						}
					}
				}

			}

			if len(strArr) == 0 {

				delete(v, k1) //没有匹配的数据则删除key

			} else {

				v[k1] = strArr //有匹配的数据，则赋予满足条件的数据

			}
		}

		if v != nil && len(v) > 0 {
			newVmMap = append(newVmMap, v)
		}

	}

	logs.Info("newVmMap:", newVmMap)

	return newVmMap, nil
}
