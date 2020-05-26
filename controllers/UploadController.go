package controllers

import (
	"fmt"
	"log"
	"path"

	"github.com/astaxie/beego"
	uuid "github.com/iris-contrib/go.uuid"
)

type UploadController struct {
	beego.Controller
}

//Upload :图片上传
func (c *UploadController) Upload() {

	f, info, err := c.GetFile("file")

	if err != nil {
		log.Fatal("getfile err ", err)
	}
	defer f.Close()

	//获取文件后缀
	ext := path.Ext(info.Filename)

	//生成文件名
	uuid, _ := uuid.NewV4()
	fileName := fmt.Sprintf("%s%s", uuid, ext)
	fmt.Println(fileName)
	c.SaveToFile("file", "static/upload/"+fileName)

	//设置上传后文件的路径 供前端调用
	c.Data["uploadFilePath"] = fmt.Sprintf("/static/upload/%s", fileName)

	//跳转到显示结果页面
	c.TplName = "equalize/display.html"
}
