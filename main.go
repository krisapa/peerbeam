package main

import (
	"github.com/6b70/peerbeam/cmd"
	log "github.com/sirupsen/logrus"
)

//func configureLogger() {
//	log.SetLevel(log.TraceLevel)
//	log.SetReportCaller(true)
//	callerFormatter := func(f *runtime.Frame) string {
//		s := strings.Split(f.Function, ".")
//		funcName := s[len(s)-1]
//		return fmt.Sprintf(" [%s:%d][%s()]", path.Base(f.File), f.Line, funcName)
//	}
//	log.SetFormatter(&nested.Formatter{
//		TimestampFormat:       "2006-01-02 15:04:05",
//		CallerFirst:           true,
//		CustomCallerFormatter: callerFormatter,
//	})
//}

func main() {
	//configureLogger()
	err := cmd.App()
	if err != nil {
		log.Fatal(err)
	}

}
