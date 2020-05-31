package routers

import (
	"myApp/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "*:Get")
	beego.Router("/equalize", &controllers.MainController{}, "*:Equalize")
	beego.Router("/specificate", &controllers.MainController{}, "*:Specificate")
	beego.Router("/specificate/display", &controllers.MainController{}, "*:EqualizeDisplay")

	//图片上传
	beego.Router("/api/upload/:key", &controllers.UploadController{}, "post:Upload")
}
