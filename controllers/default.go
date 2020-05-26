package controllers

import (
	"fmt"
	"log"
	"path"

	"github.com/astaxie/beego"
	uuid "github.com/iris-contrib/go.uuid"
)

type MainController struct {
	beego.Controller
}

//check :错误检测函数
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *MainController) Get() {

	c.TplName = "index.html"
}

func (c *MainController) Equalize() {

	c.TplName = "equalize.html"
}

//Upload :图片上传
func (c *MainController) Upload() {
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
	c.SaveToFile("file", "upload/"+fileName)

	//设置上传后文件的路径 供前端调用
	c.Data["uploadFilePath"] = fmt.Sprintf("/upload/%s", fileName)

	c.TplName = "equalize.html"
}
