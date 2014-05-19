package beegoapp2

import (
	_ "beegoapp2/routers"
	"github.com/astaxie/beegae"
)

func init() {
	beegae.Run()
}
