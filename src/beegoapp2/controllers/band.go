// band
package controllers

import (
	"beegoapp2/conf"
	"beegoapp2/models"
	//	"fmt"
	"appengine/datastore"
	"log"
	"strconv"
	"strings"
)

type BandAddController MainController
type BandVerifyController MainController

/*func main() {
	fmt.Println("Hello World!")
}
*/
func (this *BandAddController) Get() {
	context := &this.AppEngineCtx
	locations, _ := models.GetAllDocs(context, conf.LOCATION_TYPE)
	this.Data["Locations"] = locations
	this.TplNames = "band/add.html"
}

func (this *BandVerifyController) Post() {
	c := &this.AppEngineCtx
	name := this.GetString("name")
	message := "no errors"
	var locationId *datastore.Key
	var albums []models.Album
	locationType := this.GetString("loctype")
	switch locationType {
	case "existing":
		rawId := this.GetString("location_id")
		q := strings.Split(rawId, ",")
		log.Println("rawId =", rawId)
		x := q[1]
		id_int, e := strconv.ParseInt(x, 10, 64)
		if e != nil {
			message = e.Error()
		}
		locationId = datastore.NewKey(*c, conf.LOCATION_TYPE, "", id_int, nil)
		log.Printf("locationId =", locationId)
		//	message = "not implemented yet"
		break
	case "new":
		var err error
		location := models.Location{this.GetString("city"),
			this.GetString("state"), this.GetString("country")}
		locationId, err = models.AddLocation(location, c)
		log.Printf("locationId =", locationId)
		if err != nil {
			message = err.Error()
		}
		break
	}
	if message == "no errors" {
		band := models.Band{Name: name, LocationId: locationId, Albums: albums}
		_, err := models.AddBand(band, c)
		if err != nil {
			message = "Band add: " + err.Error()
		}
	}
	/*	t, err := template.ParseFiles("src/datastoremusic/views/band/verify.html")
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("File not found"))
			fmt.Println("Template parse error:", err)
			return
		} */
	//http.Redirect(rw, rq, "/home/index", http.StatusFound)
	this.Data["Message"] = message
	this.TplNames = "band/verify.html"
	//	t.Execute(rw, message)
}
