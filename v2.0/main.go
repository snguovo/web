package main

import (
	"net/http"

	"github.com/snguovo/web/v2.0/conf"
	"github.com/snguovo/web/v2.0/controller"
	"github.com/snguovo/web/v2.0/log"
	"github.com/snguovo/web/v2.0/service"

	"github.com/julienschmidt/httprouter"

	"golang.org/x/net/websocket"
)

func main() {
	config := conf.NewConfig("./appconf.json")

	go service.NewAppCron(config.Data.AppManagerEnable, 2)

	if config.Data.WebEnable {
		router := httprouter.New()
		router.GET("/", controller.Index)

		router.GET("/log", controller.AppLog)
		router.Handler("GET", "/log/upper", websocket.Handler(controller.UpdateAppLog))

		router.GET("/system", controller.System)
		router.Handler("GET", "/system/upapp", websocket.Handler(controller.UpdateApp))
		router.Handler("GET", "/system/updocker", websocket.Handler(controller.UpdateDocker))
		router.POST("/system/appaction", controller.Appaction)

		router.POST("/updateconf", controller.Updateconf)
		router.POST("/uploadfile", controller.Uploadfile)
		log.Info.Println("web server is running,listen localhost:8080")
		err := http.ListenAndServe(":8080", router)
		if err != nil {
			log.Error.Fatalln(err)
		}
	} else {
		select {}
	}
}
