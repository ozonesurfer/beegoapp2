// album
package controllers

import (
	//	"encoding/json"
	//	"fmt"
	"appengine/datastore"
	"beegoapp2/conf"
	"beegoapp2/models"
	"log"
	"os"
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
	this.Data["Request"] = c
	this.TplNames = "album/index.html"
}

func (this *AlbumAddController) Get() {
	rawId := this.Ctx.Input.Params[":id"]
	c := &this.AppEngineCtx
	genres, _ := models.GetAllDocs(c, conf.GENRE_TYPE)
	this.Data["Title"] = "Add Album"
	this.Data["Id"] = rawId
	this.Data["Genres"] = genres
	this.TplNames = "album/add.html"
}
func (this *AlbumVerifyController) Post() {
	c := &this.AppEngineCtx
	rawId := this.Ctx.Input.Params[":id"]
	log.Println("rawId =", rawId)
	q := strings.Split(rawId, ",")
	x := q[1]
	id_int, e := strconv.ParseInt(x, 10, 64)
	if e != nil {
		log.Println("Parse error:", e)
		os.Exit(1)

	}
	bandId := datastore.NewKey(*c, conf.BAND_TYPE, "", id_int, nil)
	message := "no errors"
	var genreId *datastore.Key

	genreType := this.GetString("genretype")
	switch genreType {
	case "existing":
		rawId := this.GetString("genre_id")
		q := strings.Split(rawId, ",")
		log.Println("rawId =", rawId)
		x := q[1]
		id_int, e := strconv.ParseInt(x, 10, 64)
		if e != nil {
			message = e.Error()
		}
		genreId = datastore.NewKey(*c, conf.GENRE_TYPE, "", id_int, nil)
		log.Printf("genreId =", genreId)
		//	message = "not implemented yet"
		break
	case "new":
		var err error
		genre := models.Genre{this.GetString("genre_name")}
		genreId, err = models.AddGenre(genre, c)
		log.Printf("genreId =", genreId)
		if err != nil {
			message = err.Error()
		}
		break
	}
	if message == "no errors" {

		year, _ := this.GetInt("year")
		album := models.Album{Name: this.GetString("name"),
			Year: int(year), GenreId: genreId}
		err := models.AddAlbum(album, bandId, c)
		if err != nil {
			message = err.Error()
		}

	}
	/*
		t, err := template.ParseFiles("src/datastoremusic/views/album/verify.html")
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("File not found"))
			log.Println("Template parse error:", err)
			return
		}
		t.Execute(rw, struct {
			Id      *datastore.Key
			Title   string
			Message string
		}{Id: bandId, Title: "Verifying Album", Message: message})
	*/
	this.Data["Id"] = bandId
	this.Data["Title"] = "Verifying Album"
	this.Data["Message"] = message

	this.TplNames = "album/verify.html"
}
