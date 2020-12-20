package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/snguovo/web/v2.0/dao"
	"github.com/snguovo/web/v2.0/log"
	"github.com/snguovo/web/v2.0/service"
)

//Config 配置文件结构体
type Config struct {
	Applist          []*dao.App `json:"applist"`
	WebEnable        bool       `json:"web_enable"`
	AppManagerEnable bool       `json:"appmanager_enable"`
}

//ConfigObj 配置结构体
type ConfigObj struct {
	Filename       string
	LastModifyTime int64
	Lock           *sync.RWMutex
	Data           *Config
}

//NewConfig 初始化配置
func NewConfig(filename string) *ConfigObj {
	conf := &ConfigObj{
		Filename: filename,
		Lock:     &sync.RWMutex{},
		Data:     new(Config),
	}
	if ok := conf.parse(); !ok {
		log.Error.Println("config file parse failed")
	}
	go conf.reload()
	return conf
}

//解析函数
func (c *ConfigObj) parse() bool {
	//记录最后修改时间
	fileInfo, err := os.Stat(c.Filename)
	if err != nil {
		log.Error.Println("read config ModifyTime failed:", err)
		return false
	}
	c.LastModifyTime = fileInfo.ModTime().UnixNano()
	//Logger.Println("LastModifyTime", c.LastModifyTime)

	//读取文件内容
	file, err := ioutil.ReadFile(c.Filename)
	if err != nil {
		log.Error.Println("read config failed:", err)
		return false
	}

	//解json
	c.Lock.RLock()
	if err = json.Unmarshal(file, &c.Data); err != nil {
		log.Error.Println("Unmarshal json failed:", err)
		return false
	}
	c.Lock.RUnlock()

	//更新到数据库
	err = dao.DelAllApps()
	if err != nil {
		log.Error.Println("delete app table failed:", err)
		return false
	}
	err = dao.AddApps(c.Data.Applist)
	if err != nil {
		log.Error.Println("add apps to db failed:", err)
		return false
	}
	log.Info.Println("read config file success")
	return true
}

//重载函数
func (c *ConfigObj) reload() {
	ticker := time.NewTicker(time.Second * 2)
	for range ticker.C {
		fileInfo, _ := os.Stat(c.Filename)
		currModifyTime := fileInfo.ModTime().UnixNano()
		//Logger.Println("currModifyTime", currModifyTime)
		if currModifyTime != c.LastModifyTime {
			if c.parse() {
				log.Info.Println("reload config success!")
				service.AppChangeCh <- struct{}{}
			}
		}

	}
}
