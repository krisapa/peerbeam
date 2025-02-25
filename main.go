package main

import (
	"github.com/6b70/peerbeam/cmd"
	log "github.com/sirupsen/logrus"
	"os"
)

func configureLogger() {
	log.SetLevel(log.TraceLevel)
	file, err := os.OpenFile("log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("无法创建日志文件: ", err)
	}
	log.SetOutput(file)
	//log.SetLevel(log.TraceLevel)
	//log.SetReportCaller(true)
	//callerFormatter := func(f *runtime.Frame) string {
	//	s := strings.Split(f.Function, ".")
	//	funcName := s[len(s)-1]
	//	return fmt.Sprintf(" [%s:%d][%s()]", path.Base(f.File), f.Line, funcName)
	//}
	//log.SetFormatter(&nested.Formatter{
	//	TimestampFormat:       "2006-01-02 15:04:05",
	//	CallerFirst:           true,
	//	CustomCallerFormatter: callerFormatter,
	//})
}

func main() {
	configureLogger()
	err := cmd.App()
	if err != nil {
		log.Fatal(err)
	}

}
