package util

import (
	"context"

	"github.com/snguovo/web/v2.0/log"

	"github.com/hpcloud/tail"
)

//InitTailTask 初始化taillog任务
func InitTailTask(ctx context.Context, fileName string) (<-chan string, error) {
	config := tail.Config{
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件那个地方开始读
		MustExist: true,                                 //文件不存在报错
		Poll:      true,
		Logger:    log.Info,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		log.Error.Println("tail file failed, err:", err)
		tails.Cleanup()
		return nil, err
	}
	lineChan := make(chan string, 10)
	go func() {
		defer close(lineChan)
		for {
			select {
			case <-ctx.Done():
				log.Info.Printf("tail task:%s close...\n", fileName)
				err = tails.Stop()
				if err != nil {
					log.Error.Println("stop taillog failed:", err)
				}
				tails.Cleanup()
				return
			case line := <-tails.Lines: // 从tailObj的通道中一行一行的读取日志数据
				lineChan <- line.Text
			}
		}
	}()
	return lineChan, nil
}
