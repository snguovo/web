package service

import (
	"context"
	"time"

	"github.com/snguovo/web/v2.0/dao"
	"github.com/snguovo/web/v2.0/log"
)

//NewAppTask 创建新协程获取app信息
func NewAppTask(ctx context.Context, flushtime int64) <-chan []*dao.App {
	ch := make(chan []*dao.App, 100)
	//runAppChan 定时任务管道发送App信息
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				log.Info.Println("app task is closed...")
				return
			default:
				apps, err := dao.QueryApps()
				if err != nil {
					log.Error.Println("QueryApps failed:", err)
					continue
				}
				ch <- apps
				time.Sleep(time.Second * time.Duration(flushtime))
			}
		}
	}()
	return ch
}
