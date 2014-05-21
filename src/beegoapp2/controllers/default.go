package controllers

import (
	"appengine/datastore"
	"beegoapp2/conf"
	"beegoapp2/models"
	"github.com/astaxie/beegae"
	"log"
	"strconv"
	"strings"
)

type MainController struct {
	beegae.Controller
}

type GenreListController MainController
type ByGenreController MainController

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

func (this *GenreListController) Get() {
	context := this.AppEngineCtx
	genres, _ := models.GetAllDocs(&context, conf.GENRE_TYPE)
	this.Data["Genres"] = genres
	this.TplNames = "home/genrelist.html"
}

func (this *ByGenreController) Get() {
	c := &this.AppEngineCtx
	rawId := this.Ctx.Input.Params[":id"]
	q := strings.Split(rawId, ",")
	log.Println("rawId =", rawId)
	x := q[1]
	id_int, e := strconv.ParseInt(x, 10, 64)
	if e != nil {
		log.Fatal(e.Error())
	}
	genreId := datastore.NewKey(*c, conf.GENRE_TYPE, "", id_int, nil)
	genreName := models.GetGenreName(c, genreId)
	this.Data["Title"] = genreName + " Albums"
	bands, _ := models.GetBandsByGenre(c, genreId)
	this.Data["Model"] = bands
	this.Data["Request"] = c
	this.TplNames = "home/bygenre.html"
}
