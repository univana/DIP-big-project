package controllers

import (
	"fmt"
	"log"
	"myApp/models"
	"path"
	"strings"

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
	ext := strings.ToLower(path.Ext(info.Filename))

	//保存图片
	uuid, _ := uuid.NewV4()
	fileName := fmt.Sprintf("%s%s", uuid, ext)
	filePath := fmt.Sprintf("static/upload/%s", fileName)
	c.SaveToFile("file", filePath)

	//解析图片
	picture := models.NewPicture().Analyse(filePath)

	//设置上传后文件的路径 供前端调用
	c.Data["Picture"] = picture

	//跳转到显示结果页面
	c.TplName = "equalize/display.html"
}
