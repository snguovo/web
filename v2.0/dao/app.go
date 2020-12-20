package dao

import (
	"github.com/snguovo/web/v2.0/log"
)

//App app结构体定义
type App struct {
	ID         string  `json:"-"`
	Name       string  `json:"name"`
	Version    string  `json:"version"`
	Logpath    string  `json:"logpath"`
	Confpath   string  `json:"confpath"`
	Enable     bool    `json:"enable"`
	Cmdline    string  `json:"cmd_line"`
	Cwd        string  `json:"work_dir"` //current working directory
	Pid        int32   `json:"pid"`
	MemPercent float32 `json:"mem_percent"`
	CPUPercent float64 `json:"cpu_percent"`
	Running    bool    `json:"running"`
}

//AddApps 添加app
func AddApps(apps []*App) (err error) {
	tx, err := DB.Begin()
	if err != nil {
		log.Error.Println(err)
		return
	}
	stmt, err := tx.Prepare("insert into app(id,name,version,logpath,confpath,enable,cmd_line,work_dir) values(?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Error.Println(err)
		tx.Rollback()
		return
	}
	defer stmt.Close()
	for _, app := range apps {
		_, err = stmt.Exec(app.Name+":"+app.Version, app.Name, app.Version, app.Logpath, app.Confpath, app.Enable, app.Cmdline, app.Cwd)
		if err != nil {
			log.Error.Println(err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

//DelAllApps 删除所有app
func DelAllApps() (err error) {
	_, err = DB.Exec("delete from app")
	if err != nil {
		log.Error.Println(err)
	}
	return
}

//UpdateApps 更新App状态
func UpdateApps(apps []*App) (err error) {
	//更新数据
	tx, err := DB.Begin()
	if err != nil {
		log.Error.Println(err)
		return
	}
	stmt, err := tx.Prepare("update app set cmd_line=?,work_dir=?,pid=?,mem_percent=?,cpu_percent=?,running=? where id=?;")
	if err != nil {
		log.Error.Println(err)
		tx.Rollback()
		return
	}
	defer stmt.Close()
	for _, app := range apps {
		_, err = stmt.Exec(app.Cmdline, app.Cwd, app.Pid,
			app.MemPercent, app.CPUPercent, app.Running, app.ID)
		if err != nil {
			log.Error.Println(err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

//QueryApps 查询所有app
func QueryApps() (apps []*App, err error) {
	rows, err := DB.Query("select * from app")
	if err != nil {
		log.Error.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var app App
		err = rows.Scan(&app.ID, &app.Name, &app.Version, &app.Logpath, &app.Confpath,
			&app.Enable, &app.Cmdline, &app.Cwd, &app.Pid,
			&app.MemPercent, &app.CPUPercent, &app.Running)
		if err != nil {
			log.Error.Println(err)
			return
		}
		apps = append(apps, &app)
	}
	err = rows.Err()
	if err != nil {
		log.Error.Println(err)
	}
	return
}

//QueryAppByID 通过App名称查询
func QueryAppByID(id string) (*App, error) {
	stmt, err := DB.Prepare("select * from app where id = ?")
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}
	defer stmt.Close()
	var app App
	err = stmt.QueryRow(id).Scan(&app.ID, &app.Name, &app.Version, &app.Logpath, &app.Confpath,
		&app.Enable, &app.Cmdline, &app.Cwd, &app.Pid,
		&app.MemPercent, &app.CPUPercent, &app.Running)
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}
	return &app, nil
}

//QueryAppIDs 查询所有App名称
func QueryAppIDs() (appIDs []string, err error) {
	rows, err := DB.Query("select id from app")
	if err != nil {
		log.Error.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var appID string
		err = rows.Scan(&appID)
		if err != nil {
			log.Error.Println(err)
			return
		}
		appIDs = append(appIDs, appID)
	}
	err = rows.Err()
	if err != nil {
		log.Error.Println(err)
	}
	return
}

//QueryNotRunApps 查询没有启动App
func QueryNotRunApps() (apps []*App, err error) {
	rows, err := DB.Query("select * from app where enable = true and running = false")
	if err != nil {
		log.Error.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var app App
		err = rows.Scan(&app.ID, &app.Name, &app.Version, &app.Logpath, &app.Confpath,
			&app.Enable, &app.Cmdline, &app.Cwd, &app.Pid,
			&app.MemPercent, &app.CPUPercent, &app.Running)
		if err != nil {
			log.Error.Println(err)
			return
		}
		apps = append(apps, &app)
	}
	err = rows.Err()
	if err != nil {
		log.Error.Println(err)
	}
	//log.Debug.Println("QueryNotRunApps:", apps)
	return
}

//QueryAppConfs 查询所有app配置
func QueryAppConfs() (apps []*App, err error) {
	rows, err := DB.Query("select id,name,version,logpath,confpath,enable,cmd_line,work_dir from app")
	if err != nil {
		log.Error.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var app App
		err = rows.Scan(&app.ID, &app.Name, &app.Version, &app.Logpath, &app.Confpath,
			&app.Enable, &app.Cmdline, &app.Cwd)
		if err != nil {
			log.Error.Println(err)
			return
		}
		apps = append(apps, &app)
	}
	err = rows.Err()
	if err != nil {
		log.Error.Println(err)
	}
	return
}
