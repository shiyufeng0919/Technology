package DataDeal

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/viper"
	vbg "github.com/uruddarraju/virtualbox-go"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

type Vmlist struct {
	Node   string `json:"node"`
	Status string `json:"status"`
	UUID   string `json:"uuid"`
	VMIp   []string `json:"vmIp"`
	*vbg.VirtualMachine
}

//1.
func ManageVMLists() ([]Vmlist, error) {

	logs.Info("-------客户端获取虚拟机明细信息Start--------")
	defer logs.Info("-------客户端获取虚拟机明细信息End--------")

	vmNameListArr, err := ExecCmdGetVMLists()
	if err != nil {
		return nil, err
	}

	for key, vm := range vmNameListArr {

		machine, err := GetVMInfos(vm.Node)

		vmNameListArr[key].VirtualMachine = machine

		if err != nil {
			logs.Error("根据虚拟机名称获取虚拟机信息错误:", err.Error())
			return nil, err
		}
	}

	logs.Info("vmNameListArr:", vmNameListArr)

	/*
		返回结果：

		[{node1  f6fa9977-fe1e-44c9-ae14-35450b52111f 0xc000328000} {node2  49f147c1-757d-48cd-938b-658389d31119 0xc0003281e0}

	*/

	return vmNameListArr, nil
}

//2.执行命令获取虚拟机
func ExecCmdGetVMLists() ([]Vmlist, error) {

	logs.Info("-----执行命令VBoxManage list vms Start------")
	defer logs.Info("-----执行命令VBoxManage list vms END------")

	val, err := ExecCmds("list", "vms")

	if err != nil {
		logs.Error("执行命令，获取vm list异常,", err.Error())
		return nil, err
	}

	logs.Info("执行VBoxManage list vms返回结果:", string(val))

	valRunning, err := ExecCmds("list", "runningvms")

	if err != nil {
		logs.Error("处理vmRunningList数据为数组错误:", err.Error())
		return nil, err
	}

	logs.Info("执行VBoxManage list runningvms返回结果:", string(valRunning))

	vmlistArr, err := DealVmListDatas(val, valRunning)

	if err != nil {
		logs.Error("处理vmlist数据为数组错误:", err.Error())
		return nil, err
	}

	return vmlistArr, nil

}

//3.
func ExecCmds(arg ...string) ([]byte, error) {
	cmd := exec.Command("VBoxManage", arg...)
	opBytes, err := cmd.Output()
	if err != nil {
		logs.Error("执行命令错误:", err.Error())
	}
	return opBytes, nil
}

//4.处理执行命令后返回的数据
func DealVmListDatas(val, valRunning []byte) ([]Vmlist, error) {

	logs.Info("-----DealVmListData Start----")
	logs.Info("-----DealVmListData END----")

	str1 := string(val)

	str2 := string(valRunning)

	var result []string //虚拟机列表

	var runResult []string //运行的虚拟机列表

	result = strings.Split(str1, "\n")

	fmt.Printf("itemList:%+v\n\n", result)

	runResult = strings.Split(str2, "\n")

	fmt.Printf("itemListRunning:%+v\n\n", runResult)

	var itemList = make([]Vmlist, len(result)-1)

	for i, v := range result {

		if i != len(result)-1 {

			val := strings.Split(strings.Trim(v, " "), " ")
			itemList[i].Node = val[0][1 : len(val[0])-1]
			itemList[i].UUID = val[1][1 : len(val[1])-1]

			for _, runv := range runResult {

				if runv != "" {
					if v == runv {
						itemList[i].Status = "start"
						break
					} else {
						itemList[i].Status = "stop"
						break
					}
				}
			}
		}
	}

	if len(itemList) == 0 {
		logs.Error("没有获取到虚拟机")
		return nil, errors.New("没有获取到虚拟机")
	}

	logs.Info("itemList值:", itemList)

	return itemList, nil
}

//5.
func GetVMInfos(name string) (machine *vbg.VirtualMachine, err error) {
	logs.Info("-----获取虚拟机名为%s的信息 Start----", name)
	defer logs.Info("-----获取虚拟机名为%s的信息 END------", name)
	vb := vbg.NewVBox(vbg.Config{})
	return vb.VMInfo(name)
}

type VMClientListResponse struct {
	Success  bool                     `json:"success"`
	Result   []map[string]interface{} `json:"result"`
	Errors   []interface{}            `json:"errors"`
	Messages []interface{}            `json:"messages"`
}

//1. 构造物理IP与虚拟机间关系提供给前端展示
func CallVmClientGetVmLists() {

	//获取 ip地址，建立http连接
	clientIpArr := viper.GetStringSlice("vmclient.ip")

	logs.Info("clientIpArr:", clientIpArr)

	var vmArr []interface{}

	for _, ip := range clientIpArr {

		url := GetVMClientUrl(ip, "list")

		logs.Info("调用VM客户端的url为：", url)

		body, err := ConnectHttpGetCommons(url) //上述ManageVMList返回的结果(本示例分虚拟机客户端，部署在每一台物理机，及web-server两部分，所以此处为http连接虚拟机客户端)

		if err != nil {
			logs.Error("连接vm-client错误:", err.Error())
			return
		}

		logs.Info("响应string(body):", string(body))

		vmRs := VMClientListResponse{}

		if err := json.Unmarshal(body, &vmRs); err != nil {
			logs.Error("json->struct异常:", err.Error())
			return
		}

		vmMap := make(map[string]interface{})

		if vmRs.Success && len(vmRs.Result) > 0 {

			vmMap[ip] = vmRs.Result

		} else {
			logs.Info("调取vm-client端返回结果有异常,返回success=%s,errors=%s,message=%s", vmRs.Success, vmRs.Errors, vmRs.Messages)
			return
		}

		vmArr = append(vmArr, vmMap)
	}

	logs.Info("调用vm的客户端获取虚拟机列表结果:", vmArr)

	/*
		返回结果即OperateVM.go中的参数
	*/
}

//普通http GET请求
func ConnectHttpGetCommons(url string) ([]byte, error) {

	//通过设置tls.Config的InsecureSkipVerify为true，client将不再对服务端的证书进行校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	//发送http请求
	resp, err := client.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	logs.Info("调用VM客户端响应结果:%s", string(body))

	return body, nil
}
