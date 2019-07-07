package DataDeal

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

//管道符应用(https://www.golangtc.com/t/515bdcdb320b5211ba00002f)
func GuanDaofu(vmName string) []string{

	cmd1 := exec.Command("VBoxManage", "guestproperty", "enumerate", vmName)

	cmd2 := exec.Command("grep", "Net.*V4.*IP")

	r, w := io.Pipe()

	cmd1.Stdout = w

	cmd2.Stdin = r

	var out bytes.Buffer
	cmd2.Stdout = &out
	cmd1.Start()
	cmd2.Start()
	cmd1.Wait()
	w.Close()
	cmd2.Wait()

	str := out.String()

	logs.Info("获取虚拟机ip地址:", str)

	//只摘取IP地址，构造IP数组
	var result []string

	result = strings.Split(str, "\n")

	var vmIpArr []string

	for k, v := range result {

		if k != len(result)-1 {

			val := strings.Split(v, ",")

			str := strings.Replace(val[1], "value: ", "", -1)

			vmIpArr = append(vmIpArr, str)
		}
	}

	return vmIpArr
}

//获取本机IP地址
func GetLocalHostIp() {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
			}

		}
	}
}