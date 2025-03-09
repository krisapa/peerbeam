package main

import (
	"fmt"
	"github.com/6b70/peerbeam/cmd"
	log "github.com/sirupsen/logrus"
	"io"
	"path/filepath"
	"runtime"
)

func configureLogger() {
	log.SetReportCaller(true) // 关键配置：启用调用者信息
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // 时间格式
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := filepath.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})
	log.SetLevel(log.TraceLevel)
	//err := os.MkdirAll("log", os.ModePerm)
	//if err != nil {
	//	log.Fatal("failed to create log directory: ", err)
	//}
	//file, err := os.OpenFile("log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatal("无法创建日志文件: ", err)
	//}
	//log.SetOutput(file)
	log.SetOutput(io.Discard)
}

func main() {
	configureLogger()
	err := cmd.App()
	if err != nil {
		log.Fatal(err)
	}

}
