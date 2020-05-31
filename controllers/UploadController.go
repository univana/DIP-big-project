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

var oriMapping [256]int           //原始图像灰度直方图均衡化映射表
var oriPicture *models.Picture    //原始图像实例
var oriHistogramData [256]float64 //原始图像灰度直方图数据

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
		oriHistogramData, resHistogramData, resPicture, _ := models.Equalize(originalPicture)

		/* 生成结果图像 */
		resPicture.Path = models.MakePicture(resPicture)

		//设置信息 供前端调用
		c.Data["originalPicture"] = originalPicture
		c.Data["oriHistogramData"] = oriHistogramData
		c.Data["resHistogramData"] = resHistogramData
		c.Data["resPicture"] = resPicture

		//跳转到显示结果页面
		c.TplName = "equalize/display.html"
	case "specificate":
		//规定化-上传原始图像
		//解析原始图片
		oriPicture = models.NewPicture().Analyse(filePath)

		//进行均衡化
		oriHistogramData, _, _, oriMapping = models.Equalize(oriPicture)

		c.TplName = "specificate/match.html"

	case "match":
		//规定化-上传匹配图像
		matchPicture := models.NewPicture().Analyse(filePath)

		//进行均衡化
		matchHistogramData, _, _, matchMapping := models.Equalize(matchPicture)

		//进行规定化
		resHistogramData, resPicture := models.Specificate(oriPicture, oriHistogramData, oriMapping, matchMapping)

		resPicture.Path = models.MakePicture(resPicture)

		c.Data["resPicture"] = resPicture
		c.Data["oriPicture"] = oriPicture
		c.Data["matchPicture"] = matchPicture
		c.Data["oriHistogramData"] = oriHistogramData
		c.Data["matchHistogramData"] = matchHistogramData
		c.Data["resHistogramData"] = resHistogramData
		c.TplName = "specificate/display.html"
	}

}
