package controller

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/snguovo/web/v2.0/dao"
	"github.com/snguovo/web/v2.0/log"
	"github.com/snguovo/web/v2.0/util"
	"github.com/snguovo/web/v2.0/view"

	"github.com/julienschmidt/httprouter"
)

//Index 首页
func Index(res http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var html string
	var options string
	var js string
	apps, err := dao.QueryAppIDs()
	if err != nil {
		log.Error.Println("QueryAppNames failed:", err)
		return
	}

	for _, k := range apps {
		view.AddOptions(&options, k)
	}

	query := r.URL.Query()
	filepath, ok := query["filepath"]
	if !ok || len(filepath[0]) < 1 {
		js = ""
	} else {
		//读文件
		content, err := ioutil.ReadFile(filepath[0])
		if err != nil {
			log.Error.Println("read file failed, err:", err)
			fmt.Fprintln(res, fmt.Sprintf("<body><script>alert('%s')</script> </body>", err))
			return
		}
		log.Info.Println("read file " + filepath[0] + " success")
		view.AddConfFile(&js, string(content), filepath[0])
	}
	html = fmt.Sprintf(view.INDEX, options, js)
	fmt.Fprintln(res, html)
}

//System 系统信息
func System(res http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query()
	tmp, ok := query["time"]
	t, err := strconv.Atoi(tmp[0])
	if !ok || len(tmp[0]) < 1 || err != nil {
		fmt.Fprintln(res, "<body><script>alert('请选择刷新时间！')</script> </body>")
		return
	}
	log.Info.Println("get flush time ", t, "s")
	tmp, ok = query["content"]
	if !ok || len(tmp[0]) < 1 {
		fmt.Fprintln(res, "<body><script>alert('请选择查看内容！')</script> </body>")
		return
	}
	var html string
	html = view.SYSTEM
	switch tmp[0] {
	case "system":
		view.AddSystem(&html)
		fmt.Fprintln(res, html)
		for {
			select {
			case <-res.(http.CloseNotifier).CloseNotify():
				log.Info.Println("connection closed")
				//cancel()
				return
			default:
				//收集数据
				systemInfo := util.GetSystemStatus()
				//发送
				status := view.UpdateSystem(systemInfo)
				fmt.Fprint(res, status)
				if f, ok := res.(http.Flusher); ok {
					f.Flush()
				} else {
					log.Info.Println("Damn, no flush")
				}
				log.Info.Println("send system info success", systemInfo)
				time.Sleep(time.Second * time.Duration(t))
			}
		}
	case "app":
		view.AddApps(&html, t)
		fmt.Fprintln(res, html)
	case "docker":
		if runtime.GOOS != "linux" {
			fmt.Fprintln(res, "<body><script>alert('docker info only for linux！')</script> </body>")
			return
		}
		view.AddDocker(&html, t)
		fmt.Fprintln(res, html)
	}
}

//AppLog app日志
func AppLog(res http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query()
	app, ok := query["app"]
	if !ok || len(app[0]) < 1 {
		fmt.Fprintln(res, "<body><script>alert('请选择app！')</script> </body>")
		return
	}
	html := view.AppLog(app[0])
	fmt.Fprintln(res, html)
}

//Appaction 控制app
func Appaction(res http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	action := r.FormValue("action")
	if action == "" {
		res.WriteHeader(http.StatusBadRequest)
		log.Error.Println("action is nil")
		return
	}
	appName := r.FormValue("app_name")
	if appName == "" {
		res.WriteHeader(http.StatusBadRequest)
		log.Error.Println("app_name is nil")
		return
	}
	app, err := dao.QueryAppByID(appName)
	if err != nil {
		res.WriteHeader(http.StatusMethodNotAllowed)
		log.Error.Println("QueryAppByName failed:", err)
		return
	}
	switch action {
	case "stop":
		if app.Pid == 0 {
			log.Error.Println(app.Name + " is already stop")
			res.Write([]byte(app.Name + " is already stop"))
			return
		}
		err = util.KillProcessByPid(int(app.Pid))
		if err != nil {
			log.Error.Println("KillProcessByPid failed:", err)
			res.Write([]byte("停止app " + app.Name + " 失败"))
			return
		}
		log.Info.Println("停止app " + app.Name + " 成功")
		res.Write([]byte("停止app " + app.Name + " 成功"))
	case "start":
		if app.Cmdline == "" {
			log.Error.Println(app.Name + " cmd is not defined")
			res.Write([]byte(app.Name + " cmd is not defined"))
			return
		}
		if app.Pid != 0 {
			log.Error.Println(app.Name + " is already running")
			res.Write([]byte(app.Name + " is already running"))
			return
		}
		err := util.RunCommandAtDir(app.Cwd, app.Cmdline)
		if err != nil {
			log.Error.Println("启动app " + app.Name + " 失败")
			res.Write([]byte("启动app " + app.Name + " 失败"))
			return
		}
		log.Info.Println("启动app " + app.Name + " 成功")
		res.Write([]byte("启动app " + app.Name + " 成功"))
	default:
		res.WriteHeader(http.StatusBadRequest)
		log.Error.Println("action is worry")
		return
	}
}

//Updateconf 显示配置文件
func Updateconf(res http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 2. 请求类型是application/json时从r.Body读取数据
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error.Printf("read request.Body failed, err:%v\n", err)
		return
	}
	values, err := url.ParseQuery(string(b))
	config := values.Get("config")
	filepath := values.Get("filepath")
	//log.Println(config)
	err = ioutil.WriteFile(filepath, []byte(config), 0777)
	if err != nil {
		log.Error.Printf("write file %v failed, err:%v\n", filepath, err)
		return
	}
	log.Info.Println("write to file", filepath, "success")
	answer := `<body>
	<form action="/" method="get">
	配置文件更新成功！
	<input type="submit" value="返回">
	</form> 
	</body>`
	res.Write([]byte(answer))
}

//Uploadfile 上传文件
func Uploadfile(res http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filepath := r.FormValue("path")
	if filepath == "" {
		res.WriteHeader(http.StatusBadRequest)
		log.Error.Println("path is nil")
		return
	}
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Error.Println(err)
		return
	}
	defer file.Close()

	log.Info.Println("uploadfile ", handler.Filename, "to path ", filepath)
	os.MkdirAll(filepath, os.ModePerm)
	f, err := os.OpenFile(path.Join(filepath, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error.Println(err)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		log.Error.Println(err)
		return
	}
	res.Write([]byte("上传文件成功"))
}
