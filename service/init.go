package service

import (
	"os"
	"os/signal"
	"syscall"


	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/LuckyG0ldfish/balancer/context"

	"github.com/LuckyG0ldfish/balancer/factory"

	"github.com/LuckyG0ldfish/balancer/logger"

	"github.com/LuckyG0ldfish/balancer/ngap"

	ngap_service "github.com/LuckyG0ldfish/balancer/ngap/service"

	util "github.com/LuckyG0ldfish/balancer/util/context_helper"
)

type Load struct{
	LbContext 	*context.LBContext
}

type (
	// Config information.
	Config struct {
		lbcfg string
	}
)

var config Config

var lbCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "lbcfg",
		Usage: "lb config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*Load) GetCliCmd() (flags []cli.Flag) {
	return lbCLi
}

func (Lb *Load) Initialize(c *cli.Context)  error{ // c *cli.Context) error {
	config = Config{
		lbcfg: c.String("lbcfg"),
	}

	if config.lbcfg != "" {
		if err := factory.InitConfigFactory(config.lbcfg); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Config is empty")
	}
	Lb.setLogLevel()
	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}
	return nil 
}

func (lb *Load) setLogLevel() {
	if factory.LbConfig.Logger == nil {
		initLog.Warnln("AMF config without log level setting!!!")
		return
	}

	if factory.LbConfig.Logger.LB != nil {
		if factory.LbConfig.Logger.LB.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.LbConfig.Logger.LB.DebugLevel); err != nil {
				initLog.Warnf("AMF Log level [%s] is invalid, set to [info] level",
					factory.LbConfig.Logger.LB.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("AMF Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Warnln("AMF Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.LbConfig.Logger.LB.ReportCaller)
	}

	// if factory.LbConfig.Logger.NGAP != nil {
	// 	if factory.LbConfig.Logger.NGAP.DebugLevel != "" {
	// 		if level, err := logrus.ParseLevel(factory.LbConfig.Logger.NGAP.DebugLevel); err != nil {
	// 			initLog.Warnf("AMF Log level [%s] is invalid, set to [info] level",
	// 				factory.LbConfig.Logger.NGAP.DebugLevel)
	// 			logger.SetLogLevel(logrus.InfoLevel)
	// 		} else {
	// 			initLog.Infof("AMF Log level is set to [%s] level", level)
	// 			logger.SetLogLevel(level)
	// 		}
	// 	} else {
	// 		initLog.Warnln("AMF Log level not set. Default set to [info] level")
	// 		logger.SetLogLevel(logrus.InfoLevel)
	// 	}
	// 	logger.SetReportCaller(factory.LbConfig.Logger.NGAP.ReportCaller)
	// }

	
}

func (amf *Load) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range amf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (Lb *Load) Start() {
	initLog.Infoln("Server started")

	self := context.LB_Self()
	util.InitLbContext(self)

	ngapHandler := ngap_service.NGAPHandler{
		HandleMessage:      ngap.Dispatch,
		HandleNotification: ngap.HandleSCTPNotification,
	}

	// Starting NGAP Services
	go ngap_service.Run(self.LbListenAddr, ngapHandler)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		Lb.Terminate()
		os.Exit(0)
	}()
}

// Used in LB planned removal procedure
func (Lb *Load) Terminate() {
	logger.InitLog.Infof("Terminating LB...")
	lbSelf := context.LB_Self()

	lbSelf.Running = false 
	ngap_service.Stop()

	/* Metrics */
	if lbSelf.Metrics {
	lbSelf.Table.Print()
	}

	logger.InitLog.Infof("LB terminated")
}
