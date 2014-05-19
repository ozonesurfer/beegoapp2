package routers

import (
	"beegoapp2/controllers"
	"github.com/astaxie/beegae"
)

func init() {
	beegae.Router("/", &controllers.MainController{})
	beegae.Router("/home/index", &controllers.MainController{})
	beegae.Router("/band/add", &controllers.BandAddController{})
	beegae.Router("/band/verify", &controllers.BandVerifyController{})
	beegae.Router("/album/index/:id", &controllers.AlbumIndexController{})
}
