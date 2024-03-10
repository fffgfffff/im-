package ctrl

import (
	"fmt"
	"ginChat/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	dir, _ = os.Getwd()
)

func init() {
	if err := os.MkdirAll(dir+"mnt", os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}
func Upload(c *gin.Context) {
	UploadLocal(c.Writer, c.Request)
	//UploadOss(c.Writer, c.Request)
}

// UploadLocal 1.存储位置 ./mnt，需要确保已经创建好
// 2.url格式 /mnt/xxxx.png 需要确保网络能访问/mnt/
func UploadLocal(writer http.ResponseWriter, request *http.Request) {
	//todo 获取上传的源文件
	sfile, fhead, err := request.FormFile("file")
	if err != nil {
		utils.RespFail(writer, err.Error())
		return
	}
	//todo 创建一个新文件
	//如果前端文件名称包含后缀.txt
	suffix := filepath.Ext(fhead.Filename)
	//如果前端指定fileType
	filetype := request.FormValue("filetype")
	if len(filetype) > 0 {
		suffix = filetype
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dst, err := os.Create(dir + "/asset/mnt/" + fileName)
	if err != nil {
		utils.RespFail(writer, err.Error())
		return
	}
	//todo 将源文件内容copy到新文件
	_, err = io.Copy(dst, sfile)
	if err != nil {
		utils.RespFail(writer, err.Error())
		return
	}
	//todo 将新文件路径转换成url地址，响应到前端
	utils.RespOk(writer, "/asset/mnt/"+fileName, "")
	return
}

func UploadOss(writer http.ResponseWriter, request *http.Request) {
	mfile, fhead, err := request.FormFile("file")
	if err != nil {
		utils.RespFail(writer, err.Error())
		return
	}
	//如果前端文件名称包含后缀.txt
	suffix := filepath.Ext(fhead.Filename)
	//如果前端指定fileType
	filetype := request.FormValue("filetype")
	if len(filetype) > 0 {
		suffix = filetype
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	client, err := oss.New(viper.GetString("oss.EndPoint"),
		viper.GetString("oss.AccessKeyId"), viper.GetString("oss.AccessKeySecret"))
	if err != nil {
		utils.RespFail(writer, err.Error())
		return
	}
	//获取存储空间
	bucket, err := client.Bucket(viper.GetString("oss.Bucket"))
	if err != nil {
		utils.RespFail(writer, err.Error())
		return
	}
	//上传文件
	err = bucket.PutObject(fileName, mfile)
	if err != nil {
		utils.RespFail(writer, err.Error())
		return
	}
	//获取url
	url := "http://" + viper.GetString("oss.Bucket") + "." +
		viper.GetString("oss.EndPoint") + "/" + fileName
	utils.RespOk(writer, url, "上传成功")
}
