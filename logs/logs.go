package logs

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	_ "log"
	"os"
	"sync"
	"time"
)

const (
	//Panic = 0
	//Error = 1
	//Warn  = 2
	//Info  = 3
	//Debug = 4

	pfx     = "[KRAKEND][SAGA-CLIENT]"
	resetC  = "\033[0m"
	redC    = "\033[31m"
	greenC  = "\033[32m"
	yellowC = "\033[33m"
	blueC   = "\033[34m"
	purpleC = "\033[35m"
	grayC   = "\033[37m"
)

var (
	lock         = &sync.Mutex{}
	cLogInstance *logrus.Logger
)

// GetInstance GetLogger Get an instance of logger
func GetInstance(lvl logrus.Level) (logger *logrus.Logger) {
	if cLogInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if cLogInstance == nil {
			fmt.Println("logrus instance created")
			file := fmt.Sprintf("plugins/saga-plugin-%s.log", time.Now().Format("2006-01-02"))
			f, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
			cLogInstance = &logrus.Logger{
				Out:          io.MultiWriter(f),
				Formatter:    &logrus.JSONFormatter{},
				ReportCaller: false,
				Level:        lvl,
			}
		} else {
			fmt.Println("logrus instance already created")
		}
	} else {
		fmt.Println("logrus instance already created")
	}

	return cLogInstance
}

//Logs log to stdout
func Logs(level logrus.Level, m string) {
	d := time.Now().Local().Format("2006/01/02 - 15:04:05.000")

	switch level {
	case logrus.InfoLevel:
		fmt.Println(fmt.Sprintf("%s %s %s▶ INFO%s %s", pfx, d, greenC, resetC, m))
	case logrus.DebugLevel:
		fmt.Println(fmt.Sprintf("%s %s %s▶ DEBUG%s %s", pfx, d, blueC, resetC, m))
	case logrus.WarnLevel:
		fmt.Println(fmt.Sprintf("%s %s %s▶ WARN%s %s", pfx, d, yellowC, resetC, m))
	case logrus.ErrorLevel:
		fmt.Println(fmt.Sprintf("%s %s %s▶ ERROR%s %s", pfx, d, redC, resetC, m))
	case logrus.PanicLevel:
		fmt.Println(fmt.Sprintf("%s %s %s▶ PANIC%s %s", pfx, d, purpleC, resetC, m))
	default:
		fmt.Println(fmt.Sprintf("%s %s %s▶ UNKOWN%s %s", pfx, d, grayC, resetC, m))
	}
}

type Output struct {
	Header struct {
		CorrelationID string `json:"correlation_id"`
		Timestamp     struct {
			HumanReadable string `json:"human_readable"`
			Epoch         int64  `json:"epoch"`
		} `json:"timestamp"`
		LogGenerator         string       `json:"log_generator"`
		LogFunctionGenerator string       `json:"log_function_generator"`
		Severity             logrus.Level `json:"severity"`
	} `json:"header"`
	Body struct {
		Message string `json:"message"`
		Extra   string `json:"extra"`
	} `json:"body"`
}

func GenerateLog(m map[string]interface{}) (msg string) {
	var out Output

	out.Header.Timestamp.HumanReadable = time.Now().Format("2006-01-02 15:04:05.000000")
	out.Header.Timestamp.Epoch = time.Now().Unix()
	out.Header.LogGenerator = "KrakenD"
	out.Header.LogFunctionGenerator = "ProcessRequest"
	out.Header.Severity, _ = m["severity"].(logrus.Level)
	if len(m["utid"].(string)) != 0 {
		out.Header.CorrelationID = m["utid"].(string)
	}

	out.Body.Message = m["message"].(string)
	if len(m["extra"].(string)) != 0 {
		out.Body.Extra = m["extra"].(string)
	}

	b, _ := json.Marshal(out)
	fmt.Println(string(b))
	return string(b)
}
