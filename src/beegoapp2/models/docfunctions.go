package models

import (
	"appengine"
	"appengine/datastore"
	"beegoapp2/conf"
	"log"
	"net/http"
)

type Doc struct {
	Id      *datastore.Key
	Request *appengine.Context
	Value   interface{}
}

type Band struct {
	Name       string
	LocationId *datastore.Key
	Albums     []Album
}

type Album struct {
	Name    string
	GenreId *datastore.Key
	Year    int
	//	Request *http.Request
}

type Genre struct {
	Name string
}

type Location struct {
	City, State, Country string
}

func GetAllDocs(c *appengine.Context, docType string) ([]Doc, error) {
	//	c := appengine.NewContext(rq)
	var docs []Doc
	q := datastore.NewQuery(docType).KeysOnly()
	keys, err := q.GetAll(*c, nil)
	if err != nil {
		log.Println("GetAllDocs error: " + err.Error())
		return nil, err
	}
	//	var obj interface{}
	for _, key := range keys {
		switch docType {
		case conf.BAND_TYPE:
			obj := Band{}
			err2 := datastore.Get(*c, key, &obj)
			if err2 != nil {
				log.Println("GetAllDocs error: " + err2.Error())
				return nil, err2
			}
			doc := Doc{Id: key, Value: obj, Request: c}
			docs = append(docs, doc)
			break
		case conf.LOCATION_TYPE:
			obj := Location{}
			err2 := datastore.Get(*c, key, &obj)
			if err2 != nil {
				log.Println("GetAllDocs error: " + err2.Error())
				return nil, err2
			}
			doc := Doc{Id: key, Value: obj, Request: c}
			docs = append(docs, doc)
			break
		case conf.GENRE_TYPE:
			obj := Genre{}
			err2 := datastore.Get(*c, key, &obj)
			if err2 != nil {
				log.Println("GetAllDocs error: " + err2.Error())
				return nil, err2
			}
			doc := Doc{Id: key, Value: obj, Request: c}
			docs = append(docs, doc)
			break
		}

	}
	return docs, nil
}

func (this Doc) GetLocation() string {
	var location Location
	//	c := appengine.NewContext(this.Request)
	c := this.Request
	band := this.Value.(Band)
	datastore.Get(*c, band.LocationId, &location)
	doc := Doc{Id: band.LocationId, Value: location}
	return doc.LocToString()
}
func (this Doc) LocToString() string {
	var city, state, country string
	location := this.Value.(Location)
	if location.City != "" {
		city = location.City
	} else {
		city = "(city)"
	}
	if location.State != "" {
		state = location.State
	} else {
		state = "(state/province)"
	}
	country = location.Country
	locString := city + ", " + state + " " + country
	return locString
}

func AddBand(value Band, c *appengine.Context) (*datastore.Key, error) {
	//	c := appengine.NewContext(rq)
	key := datastore.NewIncompleteKey(*c, conf.BAND_TYPE, nil)
	_, err := datastore.Put(*c, key, &value)

	return key, err
}

func AddLocation(value Location, c *appengine.Context) (*datastore.Key, error) {
	//	c := appengine.NewContext(rq)
	k := datastore.NewIncompleteKey(*c, conf.LOCATION_TYPE, nil)
	key, err := datastore.Put(*c, k, &value)

	return key, err
}

func FindLocation(location Location, c *appengine.Context) (*datastore.Key, error) {
	//	c := appengine.NewContext(rq)
	q := datastore.NewQuery(conf.LOCATION_TYPE).Filter("City =", location.City).
		Filter("State =", location.State).Filter("Country =", location.Country).
		KeysOnly()
	keys, err := q.GetAll(*c, nil)
	if err != nil {
		return nil, err
	}
	var k *datastore.Key
	for _, key := range keys {
		k = key
		break
	}
	return k, nil
}

func GetBand(bandId *datastore.Key, c *appengine.Context) (Band, error) {
	//	c := appengine.NewContext(rq)
	var band Band
	err := datastore.Get(*c, bandId, &band)
	return band, err
}

func AddGenre(genre Genre, c *appengine.Context) (*datastore.Key, error) {
	//	c := appengine.NewContext(rq)
	k := datastore.NewIncompleteKey(*c, conf.GENRE_TYPE, nil)
	key, err := datastore.Put(*c, k, &genre)

	return key, err
}

func AddAlbum(album Album, key *datastore.Key, c *appengine.Context) error {
	band, err := GetBand(key, c)
	if err != nil {
		return err
	}
	band.Albums = append(band.Albums, album)
	_, err = datastore.Put(*c, key, &band)
	return err
}

func (this Album) GetGenreName(c *appengine.Context) string {
	//	c := appengine.NewContext(rq)
	var genre Genre
	datastore.Get(*c, this.GenreId, &genre)
	return genre.Name
}

func GetGenreName(rq *http.Request, genreId *datastore.Key) string {
	c := appengine.NewContext(rq)
	var genre Genre
	datastore.Get(c, genreId, &genre)
	return genre.Name
}

func GetBandsByGenre(rq *http.Request, genreId *datastore.Key) ([]*Doc, error) {
	c := appengine.NewContext(rq)
	log.Println("Recieved key", genreId)
	q := datastore.NewQuery(conf.BAND_TYPE).Filter("Albums.GenreId =", genreId)
	var bands []Band
	keys, err := q.GetAll(c, &bands)
	if err != nil {
		log.Println("Key retrieval error:", err)
		return nil, err
	} else {
		log.Println("Key retrieval successful")
		log.Println("Found", len(keys), "keys")
	}
	var docs []*Doc
	for i := range keys {
		doc := new(Doc)
		doc.Id = keys[i]
		doc.Value = bands[i]
		docs = append(docs, doc)
	}
	return docs, nil
}
