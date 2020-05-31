package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {

	c.TplName = "index.html"
}

func (c *MainController) Equalize() {

	c.TplName = "equalize/index.html"
}

func (c *MainController) EqualizeDisplay() {
	c.TplName = "equalize/display.html"
}

func (c *MainController) Specificate() {
	c.TplName = "specificate/index.html"
}

func (c *MainController) SpecificateDisplay() {
	c.TplName = "specificate/display.html"
}
