package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

const (
	PANIC = 0
	ERROR = 1
	WARN  = 2
	INFO  = 3
	DEBUG = 4

	Prefix = "[KRAKEND][SAGA-CLIENT]"
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Gray   = "\033[37m"
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
	case INFO:
		fmt.Println(fmt.Sprintf("%s %s %s▶ INFO%s %s", Prefix, d, Green, Reset, m))
	case DEBUG:
		fmt.Println(fmt.Sprintf("%s %s %s▶ DEBUG%s %s", Prefix, d, Blue, Reset, m))
	case WARN:
		fmt.Println(fmt.Sprintf("%s %s %s▶ WARN%s %s", Prefix, d, Yellow, Reset, m))
	case ERROR:
		fmt.Println(fmt.Sprintf("%s %s %s▶ ERROR%s %s", Prefix, d, Red, Reset, m))
	case PANIC:
		fmt.Println(fmt.Sprintf("%s %s %s▶ PANIC%s %s", Prefix, d, Purple, Reset, m))
	default:
		fmt.Println(fmt.Sprintf("%s %s %s▶ UNKOWN%s %s", Prefix, d, Gray, Reset, m))
	}
}

//Log2File log required data to file
func Log2File(level int, m string, f map[string]interface{}) {
	switch level {
	case INFO:
		log.WithFields(f).Info(m)
	case DEBUG:
		log.WithFields(f).Debug(m)
	case WARN:
		log.WithFields(f).Warn(m)
	case ERROR:
		log.WithFields(f).Error(m)
	case PANIC:
		log.WithFields(f).Panic(m)
	default:
		log.WithFields(f).Println(m)
	}
}
