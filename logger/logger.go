package logger

import (
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	log           *logrus.Logger
	AppLog        *logrus.Entry
	InitLog       *logrus.Entry
	CfgLog        *logrus.Entry
	ContextLog    *logrus.Entry
	NgapLog       *logrus.Entry
	AMFHandlerLog *logrus.Entry
	GNBHandlerLog *logrus.Entry
	UtilLog       *logrus.Entry
	LbConnLog     *logrus.Entry
	AMFLog        *logrus.Entry
	GNBLog        *logrus.Entry
	UELog         *logrus.Entry
	NASLog		  *logrus.Entry
)

const (
	FieldRanAddr    string = "ran_addr"
	FieldAmfAddr    string = "amf_addr"
	FieldLbUeNgapID string = "lb_ue_ngap_id"
	FieldSupi       string = "supi"
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category", FieldRanAddr, FieldAmfAddr, FieldLbUeNgapID, FieldSupi},
	}

	AppLog 			= log.WithFields(logrus.Fields{"component": "LB", "category": "App"})
	InitLog 		= log.WithFields(logrus.Fields{"component": "LB", "category": "Init"})
	CfgLog 			= log.WithFields(logrus.Fields{"component": "LB", "category": "CFG"})
	ContextLog 		= log.WithFields(logrus.Fields{"component": "LB", "category": "Context"})
	NgapLog 		= log.WithFields(logrus.Fields{"component": "LB", "category": "NGAP"})
	AMFHandlerLog 	= log.WithFields(logrus.Fields{"component": "LB", "category": "AMFHandler"})
	GNBHandlerLog 	= log.WithFields(logrus.Fields{"component": "LB", "category": "GNBHandler"})
	UtilLog 		= log.WithFields(logrus.Fields{"component": "LB", "category": "Util"})
	LbConnLog 		= log.WithFields(logrus.Fields{"component": "LB", "category": "LBConn"})
	AMFLog 			= log.WithFields(logrus.Fields{"component": "LB", "category": "AMF"})
	GNBLog 			= log.WithFields(logrus.Fields{"component": "LB", "category": "GNB"})
	UELog 			= log.WithFields(logrus.Fields{"component": "LB", "category": "UE"})
	NASLog			= log.WithFields(logrus.Fields{"component": "LB", "category": "NAS"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(set bool) {
	log.SetReportCaller(set)
}
