package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/snguovo/web/v2.0/dao"
	"github.com/snguovo/web/v2.0/log"
	"github.com/snguovo/web/v2.0/service"
	"github.com/snguovo/web/v2.0/util"

	"golang.org/x/net/websocket"
)

//Time 返回flush time
type Time struct {
	Time int64 `json:"time"`
}

//UpdateDocker websocket上传docker状态
func UpdateDocker(ws *websocket.Conn) {
	var reply Time
	if err := websocket.JSON.Receive(ws, &reply); err != nil {
		log.Error.Println(err)
		return
	}
	log.Info.Println("get flush time ", reply.Time, "s")
	for {
		contaniers, err := util.GetContainerInfo()
		if err != nil {
			log.Error.Println("Get Container Info failed:", err)
			return
		}

		if contaniers != nil {
			if err = websocket.JSON.Send(ws, contaniers); err != nil {
				log.Warnning.Println(err)
				return
			}
		}
		time.Sleep(time.Second * time.Duration(reply.Time))
	}
}

//UpdateApp websocket上传app状态
func UpdateApp(ws *websocket.Conn) {
	var reply Time
	if err := websocket.JSON.Receive(ws, &reply); err != nil {
		log.Error.Println(err)
		return
	}
	log.Info.Println("get flush time ", reply.Time, "s")
	ctx, cancel := context.WithCancel(context.Background())
	appchan := service.NewAppTask(ctx, reply.Time)
	for {
		apps := <-appchan
		if apps != nil {
			if err := websocket.JSON.Send(ws, apps); err != nil {
				log.Warnning.Println(err)
				cancel()
				return
			}
		}
		time.Sleep(time.Second * time.Duration(reply.Time))
	}
}

//UpdateAppLog websocket上传App日志
func UpdateAppLog(ws *websocket.Conn) {
	var reply string
	if err := websocket.Message.Receive(ws, &reply); err != nil {
		log.Error.Println(err)
		return
	}
	fmt.Println("upper app", reply)
	app, err := dao.QueryAppByID(reply)
	if err != nil {
		log.Error.Println("QueryAppNames falied:", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	lineChan, err := util.InitTailTask(ctx, app.Logpath)
	if err != nil {
		log.Error.Println("init taillog failed:", err)
		websocket.Message.Send(ws, fmt.Sprintf("%v <br>", err))
		cancel()
		return
	}
	var tmp string
	for {
		select {
		case tmp = <-lineChan:
			if err = websocket.Message.Send(ws, tmp+" <br>"); err != nil {
				log.Warnning.Println(err)
				cancel()
				return
			}
			fmt.Println("get new log:", tmp)
		case <-time.After(5 * time.Minute):
			websocket.Message.Send(ws, "5 min no input,close <br>")
			log.Info.Println("5 min no input,close ")
			cancel()
			return
		}
	}
}
