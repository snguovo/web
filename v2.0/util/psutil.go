package util

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/snguovo/web/v2.0/log"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

//CPU CPU信息
type CPU struct {
	CPUCount   int     `json:"cpu_count"`
	CPUPercent float64 `json:"cpu_percent"`
}

//Disk 存储信息
type Disk struct {
	DiskTotal       string  `json:"disk_total"`
	DiskFree        string  `json:"disk_free"`
	DiskUsed        string  `json:"disk_used"`
	DiskUsedPercent float64 `json:"disk_used_percent"`
}

//Host 系统信息
type Host struct {
	Os            string `json:"os"`
	Platform      string `json:"platform"`
	KernelArch    string `json:"arch"`
	KernelVersion string `json:"version"`
}

//Mem 内存信息
type Mem struct {
	MemTotal       string  `json:"mem_total"`
	MemAvailable   string  `json:"mem_available"`
	MemUsed        string  `json:"mem_used"`
	MemUsedPercent float64 `json:"mem_used_percent"`
}

//Docker docker信息
type Docker struct {
	ContainerID string  `json:"container_id"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Status      string  `json:"status"`
	Running     bool    `json:"running"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemUsage    string  `json:"mem_usage"`
}

//Process 进程信息
type Process struct {
	Pid  int32  `json:"pid"`
	Name string `json:"name"`
	//CreateTime string  `json:"createTime"`
	Cmdline    string  `json:"cmd_line"`
	Cwd        string  `json:"work_dir"` //current working directory
	MemPercent float32 `json:"mem_percent"`
	CPUPercent float64 `json:"cpu_percent"`
	Running    bool    `json:"running"`
}

//SystemInfo 系统信息
type SystemInfo struct {
	IP   string    `json:"ip"`
	Mac  string    `json:"mac"`
	Time time.Time `json:"time"`
	Zone string    `json:"zone"`
	*CPU
	*Disk
	*Host
	*Mem
}

//getCPUInfo 查询CPU信息
func getCPUInfo() *CPU {
	newCPU := new(CPU)
	if num, err := cpu.Counts(true); err != nil {
		log.Info.Printf("get cpu count failed,err:%s\n", err)
	} else {
		newCPU.CPUCount = num
	}
	if percent, err := cpu.Percent(time.Second, false); err != nil {
		log.Info.Printf("get cpu percent failed,err:%s\n", err)
	} else {
		//fmt.Println(percent)
		newCPU.CPUPercent = formatFloat(percent[0])
	}
	return newCPU
}
func getDiskInfo() *Disk {
	newDisk := new(Disk)
	if stat, err := disk.Usage("/"); err != nil {
		log.Info.Printf("get disk stat failed,err:%s\n", err)
	} else {
		//fmt.Println(stat.String())
		newDisk.DiskTotal = formatFileSize(stat.Total)
		newDisk.DiskFree = formatFileSize(stat.Free)
		newDisk.DiskUsed = formatFileSize(stat.Used)
		newDisk.DiskUsedPercent = formatFloat(stat.UsedPercent)
	}
	return newDisk
}

func getHostInfo() *Host {
	newHost := new(Host)
	if stat, err := host.Info(); err != nil {
		log.Info.Printf("get host stat failed,err:%s\n", err)
	} else {
		//fmt.Println(stat.String())
		newHost.Os = stat.OS
		newHost.KernelArch = stat.KernelArch
		newHost.KernelVersion = stat.KernelVersion
		newHost.Platform = stat.Platform
	}
	return newHost
}
func getMemInfo() *Mem {
	newMem := new(Mem)
	if memInfo, err := mem.VirtualMemory(); err != nil {
		log.Info.Printf("get mem failed,err:%s\n", err)
	} else {
		//fmt.Println(memInfo)
		newMem.MemTotal = formatFileSize(memInfo.Total)
		newMem.MemAvailable = formatFileSize(memInfo.Available)
		newMem.MemUsed = formatFileSize(memInfo.Used)
		newMem.MemUsedPercent = formatFloat(memInfo.UsedPercent)
	}
	return newMem
}
func getIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Info.Printf("get ip failed,err:%s\n", err)
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	//fmt.Println(localAddr.String())
	return localAddr.IP.String()
}

func getMac(ip string) (mac string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Info.Println("Error:" + err.Error())
		return
	}
	for _, inter := range interfaces {
		addresses, _ := inter.Addrs()
		for _, add := range addresses {
			tmp := strings.Split(add.String(), "/")
			//fmt.Println(tmp[0])
			if tmp[0] == ip {
				//fmt.Println(inter.HardwareAddr)
				mac = inter.HardwareAddr.String()
				return
			}
		}
		// fmt.Println(inter.Name)
		// fmt.Println(inter.Index)
		// fmt.Println(inter.HardwareAddr)
	}
	return
}

//GetContainerInfo 等同于docker ps
func GetContainerInfo() (contaniers []*Docker, err error) {
	if runtime.GOOS != "linux" {
		log.Info.Println("docker info only for linux,skip")
		err = fmt.Errorf("docker info only for linux,skip")
		return
	}
	stats, err := docker.GetDockerStat()
	if err != nil {
		log.Info.Println("get contanier list failed,err:", err)
		return
	}
	for _, stat := range stats {
		var new = Docker{
			ContainerID: stat.ContainerID,
			Name:        stat.Name,
			Image:       stat.Image,
			Status:      stat.Status,
			Running:     stat.Running,
		}
		if stat.Running {
			cpuUsage, err := docker.CgroupCPUUsageDocker(stat.ContainerID)
			if err != nil {
				log.Info.Println("get contanier cpuUsage failed,err:", err)
			} else {
				new.CPUUsage = cpuUsage
			}
			memStat, err := docker.CgroupMemDocker(stat.ContainerID)
			if err != nil {
				log.Info.Println("get contanier memUsage failed,err:", err)
			} else {
				new.MemUsage = formatFileSize(memStat.MemUsageInBytes)
			}
		}
		contaniers = append(contaniers, &new)
	}
	return
}

//GetProcessByName 通过名称获取进程信息
func GetProcessByName(names []string) ([]*Process, error) {
	pids, err := process.Pids()
	if err != nil {
		log.Error.Println("get process failed,err:", err)
		return nil, err
	}
	var (
		flag    int = 0
		proName string
	)
	pros := make([]*Process, 0, len(names))
LOOP:
	for _, pid := range pids {
		pro, err := process.NewProcess(pid)
		if err != nil {
			//log.Info.Println("get process name failed,err:", err)
			continue
		}
		proName, err = pro.Name()
		if err != nil {
			//log.Info.Println("get process name failed,err:", err)
			continue
		}
		for _, name := range names {
			if proName == name {
				new := Process{
					Pid:  pro.Pid,
					Name: proName,
				}

				if cpu, err := pro.CPUPercent(); err != nil {
					log.Warnning.Println("get process cpu Percent failed,err:", err)
				} else {
					new.CPUPercent = formatFloat(cpu)
				}

				if cmd, err := pro.Cmdline(); err != nil {
					log.Warnning.Println("get process Cmdline failed,err:", err)
				} else {
					new.Cmdline = cmd
				}

				if runtime.GOOS == "linux" {
					if cwd, err := pro.Cwd(); err != nil {
						log.Warnning.Println("get process Cwd failed,err:", err)

					} else {
						new.Cwd = cwd
					}

				}
				// if createtime, err := pro.CreateTime(); err != nil {
				// 	log.Warnning.Println("get process createtime Percent failed,err:", err)
				// } else {
				// 	new.CreateTime = time.Unix(createtime/1000, 0).Format("2006-01-02 15:04:05")
				// }
				if mem, err := pro.MemoryPercent(); err != nil {
					log.Warnning.Println("get process mem Percent failed,err:", err)

				} else {
					new.MemPercent = float32(formatFloat(float64(mem)))
				}

				if running, err := pro.IsRunning(); err != nil {
					log.Warnning.Println("get process IsRunning failed,err:", err)

				} else {
					new.Running = running
				}

				pros = append(pros, &new)

				flag++
				if flag == len(names) {
					break LOOP
				}
			}
		}
	}
	return pros, nil
}

//KillProcessByPid 杀死进程
func KillProcessByPid(pid int) (err error) {
	process, err := os.FindProcess(pid)
	if err != nil {
		log.Error.Println("find pid failed:", err)
		return
	}
	return process.Kill()
}

//GetSystemStatus 获取系统全部信息
func GetSystemStatus() *SystemInfo {
	ip := getIP()
	sys := &SystemInfo{
		CPU:  getCPUInfo(),
		Disk: getDiskInfo(),
		Host: getHostInfo(),
		Mem:  getMemInfo(),
		Time: time.Now(),
		Zone: time.Now().Location().String(),
		IP:   ip,
		Mac:  getMac(ip),
	}
	return sys
}
