package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

const (
	Panic = 0
	Error = 1
	Warn  = 2
	Info  = 3
	Debug = 4

	pfx     = "[KRAKEND][SAGA-CLIENT]"
	resetC  = "\033[0m"
	redC    = "\033[31m"
	greenC  = "\033[32m"
	yellowC = "\033[33m"
	blueC   = "\033[34m"
	purpleC = "\033[35m"
	grayC   = "\033[37m"
)

var log = logrus.New()

func init() {
	file := fmt.Sprintf("plugins/saga-plugin-%s.log", time.Now().Format("2006-01-02"))
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(io.MultiWriter(f))
}

//Log log to stdout
func Log(level int, m string) {
	fm := "2006/01/02 - 15:04:05.000"
	d := time.Now().Local().Format(fm)

	switch level {
	case Info:
		fmt.Println(fmt.Sprintf("%s %s %s▶ INFO%s %s", pfx, d, greenC, resetC, m))
	case Debug:
		fmt.Println(fmt.Sprintf("%s %s %s▶ DEBUG%s %s", pfx, d, blueC, resetC, m))
	case Warn:
		fmt.Println(fmt.Sprintf("%s %s %s▶ WARN%s %s", pfx, d, yellowC, resetC, m))
	case Error:
		fmt.Println(fmt.Sprintf("%s %s %s▶ ERROR%s %s", pfx, d, redC, resetC, m))
	case Panic:
		fmt.Println(fmt.Sprintf("%s %s %s▶ PANIC%s %s", pfx, d, purpleC, resetC, m))
	default:
		fmt.Println(fmt.Sprintf("%s %s %s▶ UNKOWN%s %s", pfx, d, grayC, resetC, m))
	}
}

//LogF log required data to file
func LogF(level int, m string, f map[string]interface{}) {
	switch level {
	case Info:
		log.WithFields(f).Info(m)
	case Debug:
		log.WithFields(f).Debug(m)
	case Warn:
		log.WithFields(f).Warn(m)
	case Error:
		log.WithFields(f).Error(m)
	case Panic:
		log.WithFields(f).Panic(m)
	default:
		log.WithFields(f).Println(m)
	}
}
