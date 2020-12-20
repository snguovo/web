package log

import (
	"io"
	"log"
	"os"
)

//Logger 日志对象
var (
	Warnning *log.Logger
	Info     *log.Logger
	Error    *log.Logger
)

func init() {
	logFile, err := os.OpenFile("./web.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("open log file failed:", err)
	}

	Warnning = log.New(io.MultiWriter(os.Stdout, logFile), "[WARNNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(io.MultiWriter(os.Stdout, logFile), "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stdout, logFile), "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}
