package service

import (
	"time"

	"github.com/snguovo/web/v2.0/dao"
	"github.com/snguovo/web/v2.0/log"
	"github.com/snguovo/web/v2.0/util"
)

//AppChangeCh 配文件改变通道
var AppChangeCh = make(chan struct{}, 10)

//NewAppCron 轮询App状态定时任务,（输入参数：App使能服务 轮询间隔单位：秒）
func NewAppCron(appManagerEnable bool, interval int) {
	ticker1 := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		log.Info.Println("App Cron jobs is running")
		apps, err := dao.QueryAppConfs()
		if err != nil {
			log.Error.Println("QueryAppNames failed:", err)
		}
		for range ticker1.C {
			select {
			case <-AppChangeCh:
				//数据库查询app数
				apps = make([]*dao.App, 0)
				apps, err = dao.QueryAppConfs()
				if err != nil {
					log.Error.Println("QueryAppNames failed:", err)
					continue
				}
			default:
			}
			var appsName []string
			for _, app := range apps {
				appsName = append(appsName, app.Name)
			}
			//获得具体pid信息
			processes, err := util.GetProcessByName(appsName)
			if err != nil {
				log.Error.Println("GetProcessByName failed:", err)
				continue
			}
			for _, process := range processes {
				for _, app := range apps {
					if app.Name == process.Name {
						app.Cmdline = process.Cmdline
						app.Cwd = process.Cwd
						app.Pid = process.Pid
						app.MemPercent = process.MemPercent
						app.CPUPercent = process.CPUPercent
						app.Running = process.Running
						//log.Info.Printf("update app name=%v cmd_line=%v,work_dir=%v,pid=%v,mem_percent=%v,cpu_percent=%v,running=%v",
						//app.Name, app.Cmdline, app.Cwd, app.Pid, app.MemPercent, app.CPUPercent, app.Running)
					}
				}
			}
			//存到数据库
			err = dao.UpdateApps(apps)
			if err != nil {
				log.Error.Println("Addapps failed:", err)
				continue
			}

		}
	}()
	if appManagerEnable {
		time.Sleep(2 * time.Second)
		ticker2 := time.NewTicker(2 * time.Second)
		//如果App enable但不在线，启动app
		go func() {
			log.Info.Println("App enable jobs is running")
			for range ticker2.C {
				//拿到enable not running app
				apps, err := dao.QueryNotRunApps()
				if err != nil {
					log.Error.Println("QueryNotRunApps failed:", err)
					continue
				}
				//启动app
				for _, app := range apps {
					if app.Cmdline == "" {
						log.Warnning.Printf("app %s cmdline is not exist\n", app.Name+":"+app.Version)
						continue
					}
					log.Info.Printf("find app %s not running,run cmd %s at %s\n", app.Name+":"+app.Version, app.Cmdline, app.Cwd)
					err = util.RunCommandAtDir(app.Cwd, app.Cmdline)
					if err != nil {
						log.Error.Printf("run cmd %s at %s failed:%S \n", app.Cmdline, app.Cwd, err)
						continue
					}
				}
				if apps != nil {
					time.Sleep(time.Duration(interval) * time.Second)
				}
			}
		}()
	}
}
