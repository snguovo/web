package util

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/snguovo/web/v2.0/log"
)

// 字节的单位转换 保留两位小数
func formatFileSize(fileSize uint64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

//保留两位小数
func formatFloat(num float64) float64 {
	v, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return v
}

//FormatCmd 去除wincmd中的"
func FormatCmd(cmd string) string {
	return strings.ReplaceAll(strings.ReplaceAll(cmd, "\\", "/"), `"`, "")
}

//RunCommandAtDir 在指定路径执行命令
func RunCommandAtDir(dir string, cmd string) error {
	log.Info.Println("Running cmd: sh -c " + cmd + " &")
	cmdLine := exec.Command("sh", "-c", cmd+" &")
	if dir != "" {
		cmdLine.Dir = dir
	}
	err := cmdLine.Start()
	if err != nil {
		log.Error.Println("run cmd failed:", err)
		return err
	}
	go func() {
		cmdLine.Wait()
	}()
	return nil

}
