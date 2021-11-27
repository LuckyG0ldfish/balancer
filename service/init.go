package service

import (
	// "bufio"
	// "fmt"
	"os"
	"strconv"
	// "os/exec"
	"os/signal"
	// "sync"
	"syscall"


	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/LuckyG0ldfish/balancer/context"

	"github.com/LuckyG0ldfish/balancer/factory"

	"github.com/LuckyG0ldfish/balancer/logger"

	"github.com/LuckyG0ldfish/balancer/ngap"

	ngap_service "github.com/LuckyG0ldfish/balancer/ngap/service"

	"github.com/LuckyG0ldfish/balancer/util"

	"github.com/free5gc/path_util"

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

	// lb = NewLB()
	if config.lbcfg != "" {
		if err := factory.InitConfigFactory(config.lbcfg); err != nil {
			return err
		}
	} else {
		DefaultAmfConfigPath := path_util.Free5gcPath("balancer/config/lbcfg.yaml")
		if err := factory.InitConfigFactory(DefaultAmfConfigPath); err != nil {
			return err
		}
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

	if factory.LbConfig.Logger.AMF != nil {
		if factory.LbConfig.Logger.AMF.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.LbConfig.Logger.AMF.DebugLevel); err != nil {
				initLog.Warnf("AMF Log level [%s] is invalid, set to [info] level",
					factory.LbConfig.Logger.AMF.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("AMF Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Warnln("AMF Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.LbConfig.Logger.AMF.ReportCaller)
	}

	
}
	// if factory.AmfConfig.Logger.NGAP != nil {
	// 	if factory.AmfConfig.Logger.NGAP.DebugLevel != "" {
	// 		if /*level*/_, err := logrus.ParseLevel(factory.AmfConfig.Logger.NGAP.DebugLevel); err != nil {
	// 			// ngapLogger.NgapLog.Warnf("NGAP Log level [%s] is invalid, set to [info] level",
	// 				// factory.AmfConfig.Logger.NGAP.DebugLevel)
	// 			// ngapLogger.SetLogLevel(logrus.InfoLevel)
	// 		} else {
	// 			// ngapLogger.SetLogLevel(level)
	// 		}
	// 	} else {
	// 		// ngapLogger.NgapLog.Warnln("NGAP Log level not set. Default set to [info] level")
	// 		// ngapLogger.SetLogLevel(logrus.InfoLevel)
	// 	}
	// 	// ngapLogger.SetReportCaller(factory.AmfConfig.Logger.NGAP.ReportCaller)
	// }

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

	// addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)

	ngapHandler := ngap_service.NGAPHandler{
		HandleMessage:      ngap.Dispatch,
		HandleNotification: ngap.HandleSCTPNotification,
	}

	go Lb.InitAmfs(ngapHandler)

	go ngap_service.Run(self.LbListenAddr, ngapHandler)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		Lb.Terminate()
		os.Exit(0)
	}()
}

func (Lb *Load) InitAmfs(ngapHandler ngap_service.NGAPHandler) {
	self := context.LB_Self()
	for {
		if !self.Running { return }
		if self.NewAmf {
			var ip string 
			var port string  
			if len(self.NewAmfIpList) == len(self.NewAmfPortList) {
				for i := 0; i < len(self.NewAmfIpList); i++  {
					ip = self.NewAmfIpList[i]
					port = self.NewAmfPortList[i]
					logger.NgapLog.Tracef("connecting to: " + ip + ":" + port)
					if a, err := strconv.Atoi(port); err == nil {
						go Lb.StartAmfs(ip, a, ngapHandler)
					} else {
						logger.CfgLog.Errorf("port conversion to int failed")
					}
				}
			} else {
				logger.CfgLog.Errorf("length of IP-List and Port-List aren't identical")
			}
		}
		self.NewAmfPortList = []string{}
		self.NewAmfIpList = []string{}
		self.NewAmf = false
	}
}

func (Lb *Load) StartAmfs(amfIP string, amfPort int, ngapHandler ngap_service.NGAPHandler) {
	self := context.LB_Self()
	amf := context.NewLbAmf()
	self.AddAmfToLB(amf)
	ngap_service.StartAmf(amf, self.LbToAmfAddr, amfIP, amfPort, ngapHandler)
	initLog.Infoln("connected to amf: IP " + amfIP + " Port: " + strconv.Itoa(amfPort))
}

// Used in LB planned removal procedure
func (Lb *Load) Terminate() {
	logger.InitLog.Infof("Terminating LB...")
	lbSelf := context.LB_Self()

	lbSelf.Running = false 
	ngap_service.Stop()

	logger.InitLog.Infof("LB terminated")
}
