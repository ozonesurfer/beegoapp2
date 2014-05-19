// album
package controllers

import (
	//	"encoding/json"
	//	"fmt"
	"appengine/datastore"
	"beegoapp2/conf"
	"beegoapp2/models"
	"log"
	"strconv"
	"strings"
)

type AlbumIndexController MainController
type AlbumAddController MainController
type AlbumVerifyController MainController

/*
func main() {
	fmt.Println("Hello World!")
}
*/

func (this *AlbumIndexController) Get() {
	c := &this.AppEngineCtx
	rawId := this.Ctx.Input.Params[":id"]
	title := "Invalid band"
	q := strings.Split(rawId, ",")
	log.Println("rawId =", rawId)
	x := q[1]
	id_int, e := strconv.ParseInt(x, 10, 64)
	if e != nil {
		log.Panic(e.Error())
	}
	bandId := datastore.NewKey(*c, conf.BAND_TYPE, "", id_int, nil)
	log.Println("/Album/Index received", rawId)
	band, err := models.GetBand(bandId, c)
	if err != nil {
		log.Println("GetBand error:", err)
	} else {
		title = "Albums by " + band.Name
		this.Data["Band"] = band
	}
	this.Data["Title"] = title
	this.Data["Id"] = rawId
	this.TplNames = "album/index.html"
}
