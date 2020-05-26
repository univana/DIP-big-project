package routers

import (
	"myApp/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "*:Get")
	beego.Router("/equalize", &controllers.MainController{}, "*:Equalize")

	//图片上传
	beego.Router("/api/upload", &controllers.MainController{}, "post:Upload")
}
