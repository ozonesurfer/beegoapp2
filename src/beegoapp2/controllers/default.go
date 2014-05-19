package controllers

import (
	"beegoapp2/conf"
	"beegoapp2/models"
	"github.com/astaxie/beegae"
)

type MainController struct {
	beegae.Controller
}

func (this *MainController) Get() {
	/*	this.Data["Website"] = "beegae.me"
		this.Data["Email"] = "astaxie@gmail.com"
		this.TplNames = "index.tpl"
	*/
	context := this.AppEngineCtx
	bands, _ := models.GetAllDocs(&context, conf.BAND_TYPE)
	this.Data["Bands"] = bands
	this.TplNames = "home/index.html"
}
