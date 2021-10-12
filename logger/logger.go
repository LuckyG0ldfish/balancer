package logger

import (
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	"github.com/free5gc/logger_util"
)

var logFileName = "loadbalancer.log"

var log *logrus.Logger
var AppLog *logrus.Entry
var InitLog *logrus.Entry
var NASLog *logrus.Entry
var NGAPLog *logrus.Entry
var LoadbalancerLog *logrus.Entry
var GmmLog *logrus.Entry

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339Nano,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"category"},
	}

	selfLogHook, err := logger_util.NewFileHook(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		log.Hooks.Add(selfLogHook)
	}

	AppLog = log.WithFields(logrus.Fields{"category": "App"})
	InitLog = log.WithFields(logrus.Fields{"category": "Init"})
	LoadbalancerLog = log.WithFields(logrus.Fields{"category": "LB"})
	NASLog = log.WithFields(logrus.Fields{"category": "NAS"})
	NGAPLog = log.WithFields(logrus.Fields{"category": "NGAP"})
	GmmLog = log.WithFields(logrus.Fields{"category": "GMM"})

}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}

func TryAddLogFileHook(fileName string) {
	selfLogHook, err := logger_util.NewFileHook(fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err == nil {
		log.Hooks.Add(selfLogHook)
	}
}
