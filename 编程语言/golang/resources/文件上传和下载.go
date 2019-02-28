package resources

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

/*
 ####golang的gin框架实现文件下载

 它的核心是调用了http包    http.ServeFile(c.Writer, c.Request, filepath)

 ServeFile  又调用了func serveFile(w ResponseWriter, r *Request, fs FileSystem, name string, redirect bool) 方法

*/

func downCFCARootCert(c *gin.Context) {
	filename:="xx.cer"
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File("msp/crypto-config/ordererOrganizations/example.com/ca/xx.cer")
}