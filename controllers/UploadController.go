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

	//key: 上传后申请的操作 均衡化/规格化
	key := c.Ctx.Input.Param(":key")
	fmt.Println(key)

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

	switch key {
	case "equalize":
		//均衡化
		//解析原始图片
		originalPicture := models.NewPicture().Analyse(filePath)

		//进行均衡化
		source := models.Equalize(originalPicture)

		//制作新图像
		models.MakePicture(source, originalPicture.Bounds, originalPicture.Ext)

		//设置原始图片实例 供前端调用
		c.Data["OriginalPicture"] = originalPicture

		//跳转到显示结果页面
		c.TplName = "equalize/display.html"
	case "specify":
		//规定化

	}

}
