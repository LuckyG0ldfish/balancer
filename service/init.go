package service

import (
	// "bufio"
	// "fmt"
	// "os"
	// "os/exec"
	// "os/signal"
	// "sync"
	// "syscall"

	// "github.com/gin-contrib/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	// "github.com/free5gc/amf/communication"
	// "github.com/free5gc/amf/consumer"
	"github.com/LuckyG0ldfish/balancer/context"
	// "github.com/free5gc/amf/eventexposure"
	// "github.com/free5gc/amf/factory"
	// "github.com/free5gc/amf/httpcallback"
	// "github.com/free5gc/amf/location"
	"github.com/LuckyG0ldfish/balancer/logger"
	// "github.com/free5gc/amf/mt"
	"github.com/LuckyG0ldfish/balancer/ngap"
	// ngap_message "github.com/free5gc/amf/ngap/message"
	ngap_service "github.com/LuckyG0ldfish/balancer/ngap/service"
	// "github.com/free5gc/amf/oam"
	// "github.com/free5gc/amf/producer/callback"
	// "github.com/free5gc/amf/util"
	// aperLogger "github.com/free5gc/aper/logger"
	// fsmLogger "github.com/free5gc/fsm/logger"
	// "github.com/free5gc/http2_util"
	// "github.com/free5gc/logger_util"
	// nasLogger "github.com/free5gc/nas/logger"
	// ngapLogger "github.com/free5gc/ngap/logger"
	// openApiLogger "github.com/free5gc/openapi/logger"
	// "github.com/free5gc/openapi/models"
	// "github.com/free5gc/path_util"
	// pathUtilLogger "github.com/free5gc/path_util/logger"
)

type LB struct{
	LbContext context.LBContext
}

type (
	// Config information.
	Config struct {
		lbcfg string
	}
)

const lbIP string = "127.0.0.1"
const amfIP string = "127.0.0.1"

const lbPort int = 48484
const amfPort int = 38412

var config Config

var amfCLi = []cli.Flag{
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

func (*LB) GetCliCmd() (flags []cli.Flag) {
	return amfCLi
}

func NewLB() (lb *LB){
	lb.LbContext = NewLB().LbContext
	return 
}

func (lb *LB) Initialize()  { // c *cli.Context) error {
	// config = Config{
	// 	lbcfg: c.String("lbcfg"),
	// }

	Lb := NewLB()

	Lb.Start()

	// if config.lbcfg != "" {
	// 	if err := factory.InitConfigFactory(config.lbcfg); err != nil {
	// 		return err
	// 	}
	// } else {
	// 	DefaultAmfConfigPath := path_util.Free5gcPath("balancer/config/lbcfg.yaml")
	// 	if err := factory.InitConfigFactory(DefaultAmfConfigPath); err != nil {
	// 		return err
	// 	}
	// }

	// lb.setLogLevel()
	lb.LbContext.init()
	lb.LbContext.LbIP = lbIP
	lb.LbContext.LbPort = lbPort

	// if err := factory.CheckConfigVersion(); err != nil {
	// 	return err
	// }
}

// func (amf *LB) setLogLevel() {
// 	if factory.AmfConfig.Logger == nil {
// 		initLog.Warnln("AMF config without log level setting!!!")
// 		return
// 	}

// 	if factory.AmfConfig.Logger.AMF != nil {
// 		if factory.AmfConfig.Logger.AMF.DebugLevel != "" {
// 			if level, err := logrus.ParseLevel(factory.AmfConfig.Logger.AMF.DebugLevel); err != nil {
// 				initLog.Warnf("AMF Log level [%s] is invalid, set to [info] level",
// 					factory.AmfConfig.Logger.AMF.DebugLevel)
// 				logger.SetLogLevel(logrus.InfoLevel)
// 			} else {
// 				initLog.Infof("AMF Log level is set to [%s] level", level)
// 				logger.SetLogLevel(level)
// 			}
// 		} else {
// 			initLog.Warnln("AMF Log level not set. Default set to [info] level")
// 			logger.SetLogLevel(logrus.InfoLevel)
// 		}
// 		logger.SetReportCaller(factory.AmfConfig.Logger.AMF.ReportCaller)
// 	}

// 	if factory.AmfConfig.Logger.NAS != nil {
// 		if factory.AmfConfig.Logger.NAS.DebugLevel != "" {
// 			if level, err := logrus.ParseLevel(factory.AmfConfig.Logger.NAS.DebugLevel); err != nil {
// 				nasLogger.NasLog.Warnf("NAS Log level [%s] is invalid, set to [info] level",
// 					factory.AmfConfig.Logger.NAS.DebugLevel)
// 				logger.SetLogLevel(logrus.InfoLevel)
// 			} else {
// 				nasLogger.SetLogLevel(level)
// 			}
// 		} else {
// 			nasLogger.NasLog.Warnln("NAS Log level not set. Default set to [info] level")
// 			nasLogger.SetLogLevel(logrus.InfoLevel)
// 		}
// 		nasLogger.SetReportCaller(factory.AmfConfig.Logger.NAS.ReportCaller)
// 	}

// 	if factory.AmfConfig.Logger.NGAP != nil {
// 		if factory.AmfConfig.Logger.NGAP.DebugLevel != "" {
// 			if level, err := logrus.ParseLevel(factory.AmfConfig.Logger.NGAP.DebugLevel); err != nil {
// 				ngapLogger.NgapLog.Warnf("NGAP Log level [%s] is invalid, set to [info] level",
// 					factory.AmfConfig.Logger.NGAP.DebugLevel)
// 				ngapLogger.SetLogLevel(logrus.InfoLevel)
// 			} else {
// 				ngapLogger.SetLogLevel(level)
// 			}
// 		} else {
// 			ngapLogger.NgapLog.Warnln("NGAP Log level not set. Default set to [info] level")
// 			ngapLogger.SetLogLevel(logrus.InfoLevel)
// 		}
// 		ngapLogger.SetReportCaller(factory.AmfConfig.Logger.NGAP.ReportCaller)
// 	}

// 	if factory.AmfConfig.Logger.FSM != nil {
// 		if factory.AmfConfig.Logger.FSM.DebugLevel != "" {
// 			if level, err := logrus.ParseLevel(factory.AmfConfig.Logger.FSM.DebugLevel); err != nil {
// 				fsmLogger.FsmLog.Warnf("FSM Log level [%s] is invalid, set to [info] level",
// 					factory.AmfConfig.Logger.FSM.DebugLevel)
// 				fsmLogger.SetLogLevel(logrus.InfoLevel)
// 			} else {
// 				fsmLogger.SetLogLevel(level)
// 			}
// 		} else {
// 			fsmLogger.FsmLog.Warnln("FSM Log level not set. Default set to [info] level")
// 			fsmLogger.SetLogLevel(logrus.InfoLevel)
// 		}
// 		fsmLogger.SetReportCaller(factory.AmfConfig.Logger.FSM.ReportCaller)
// 	}

// 	if factory.AmfConfig.Logger.Aper != nil {
// 		if factory.AmfConfig.Logger.Aper.DebugLevel != "" {
// 			if level, err := logrus.ParseLevel(factory.AmfConfig.Logger.Aper.DebugLevel); err != nil {
// 				aperLogger.AperLog.Warnf("Aper Log level [%s] is invalid, set to [info] level",
// 					factory.AmfConfig.Logger.Aper.DebugLevel)
// 				aperLogger.SetLogLevel(logrus.InfoLevel)
// 			} else {
// 				aperLogger.SetLogLevel(level)
// 			}
// 		} else {
// 			aperLogger.AperLog.Warnln("Aper Log level not set. Default set to [info] level")
// 			aperLogger.SetLogLevel(logrus.InfoLevel)
// 		}
// 		aperLogger.SetReportCaller(factory.AmfConfig.Logger.Aper.ReportCaller)
// 	}

// 	if factory.AmfConfig.Logger.PathUtil != nil {
// 		if factory.AmfConfig.Logger.PathUtil.DebugLevel != "" {
// 			if level, err := logrus.ParseLevel(factory.AmfConfig.Logger.PathUtil.DebugLevel); err != nil {
// 				pathUtilLogger.PathLog.Warnf("PathUtil Log level [%s] is invalid, set to [info] level",
// 					factory.AmfConfig.Logger.PathUtil.DebugLevel)
// 				pathUtilLogger.SetLogLevel(logrus.InfoLevel)
// 			} else {
// 				pathUtilLogger.SetLogLevel(level)
// 			}
// 		} else {
// 			pathUtilLogger.PathLog.Warnln("PathUtil Log level not set. Default set to [info] level")
// 			pathUtilLogger.SetLogLevel(logrus.InfoLevel)
// 		}
// 		pathUtilLogger.SetReportCaller(factory.AmfConfig.Logger.PathUtil.ReportCaller)
// 	}

// 	if factory.AmfConfig.Logger.OpenApi != nil {
// 		if factory.AmfConfig.Logger.OpenApi.DebugLevel != "" {
// 			if level, err := logrus.ParseLevel(factory.AmfConfig.Logger.OpenApi.DebugLevel); err != nil {
// 				openApiLogger.OpenApiLog.Warnf("OpenAPI Log level [%s] is invalid, set to [info] level",
// 					factory.AmfConfig.Logger.OpenApi.DebugLevel)
// 				openApiLogger.SetLogLevel(logrus.InfoLevel)
// 			} else {
// 				openApiLogger.SetLogLevel(level)
// 			}
// 		} else {
// 			openApiLogger.OpenApiLog.Warnln("OpenAPI Log level not set. Default set to [info] level")
// 			openApiLogger.SetLogLevel(logrus.InfoLevel)
// 		}
// 		openApiLogger.SetReportCaller(factory.AmfConfig.Logger.OpenApi.ReportCaller)
// 	}
// }

// func (amf *LB) FilterCli(c *cli.Context) (args []string) {
// 	for _, flag := range amf.GetCliCmd() {
// 		name := flag.GetName()
// 		value := fmt.Sprint(c.Generic(name))
// 		if value == "" {
// 			continue
// 		}

// 		args = append(args, "--"+name, value)
// 	}
// 	return args
// }

func (lb *LB) Start() {
	// initLog.Infoln("Server started")

	
	ngapHandler := ngap_service.NGAPHandler{
		HandleMessage:      ngap.Dispatch,
		HandleNotification: ngap.HandleSCTPNotification,
	}
	ngap_service.Run(lb.LbContext.LbIP, lb.LbContext.LbPortlb, ngapHandler)

}
	// router := logger_util.NewGinWithLogrus(logger.GinLog)
	// router.Use(cors.New(cors.Config{
	// 	AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
	// 	AllowHeaders: []string{
	// 		"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host",
	// 		"Token", "X-Requested-With",
	// 	},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	AllowAllOrigins:  true,
	// 	MaxAge:           86400,
	// }))

	// httpcallback.AddService(router)
	// oam.AddService(router)
	// for _, serviceName := range factory.AmfConfig.Configuration.ServiceNameList {
	// 	switch models.ServiceName(serviceName) {
	// 	case models.ServiceName_NAMF_COMM:
	// 		communication.AddService(router)
	// 	case models.ServiceName_NAMF_EVTS:
	// 		eventexposure.AddService(router)
	// 	case models.ServiceName_NAMF_MT:
	// 		mt.AddService(router)
	// 	case models.ServiceName_NAMF_LOC:
	// 		location.AddService(router)
	// 	}
	// }

	// self := context.AMF_Self()
	// util.InitAmfContext(self)

	// addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)

	// ngapHandler := ngap_service.NGAPHandler{
	// 	HandleMessage:      ngap.Dispatch,
	// 	HandleNotification: ngap.HandleSCTPNotification,
	// }
	// ngap_service.Run(self.NgapIpList, 38412, ngapHandler)

	// // Register to NRF
	// var profile models.NfProfile
	// if profileTmp, err := consumer.BuildNFInstance(self); err != nil {
	// 	initLog.Error("Build AMF Profile Error")
	// } else {
	// 	profile = profileTmp
	// }

	// if _, nfId, err := consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile); err != nil {
	// 	initLog.Warnf("Send Register NF Instance failed: %+v", err)
	// } else {
	// 	self.NfId = nfId
	// }

	// signalChannel := make(chan os.Signal, 1)
	// signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-signalChannel
	// 	lb.Terminate()
	// 	os.Exit(0)
	// }()

	// server, err := http2_util.NewServer(addr, util.AmfLogPath, router)

	// if server == nil {
	// 	initLog.Errorf("Initialize HTTP server failed: %+v", err)
	// 	return
	// }

	// if err != nil {
	// 	initLog.Warnf("Initialize HTTP server: %+v", err)
	// }

	// serverScheme := factory.AmfConfig.Configuration.Sbi.Scheme
	// if serverScheme == "http" {
	// 	err = server.ListenAndServe()
	// } else if serverScheme == "https" {
	// 	err = server.ListenAndServeTLS(util.AmfPemPath, util.AmfKeyPath)
	// }

	// if err != nil {
	// 	initLog.Fatalf("HTTP server setup failed: %+v", err)
	// }
//}

// func (amf *LB) Exec(c *cli.Context) error {
// 	// AMF.Initialize(cfgPath, c)

// 	initLog.Traceln("args:", c.String("amfcfg"))
// 	args := amf.FilterCli(c)
// 	initLog.Traceln("filter: ", args)
// 	command := exec.Command("./amf", args...)

// 	stdout, err := command.StdoutPipe()
// 	if err != nil {
// 		initLog.Fatalln(err)
// 	}
// 	wg := sync.WaitGroup{}
// 	wg.Add(3)
// 	go func() {
// 		in := bufio.NewScanner(stdout)
// 		for in.Scan() {
// 			fmt.Println(in.Text())
// 		}
// 		wg.Done()
// 	}()

// 	stderr, err := command.StderrPipe()
// 	if err != nil {
// 		initLog.Fatalln(err)
// 	}
// 	go func() {
// 		in := bufio.NewScanner(stderr)
// 		for in.Scan() {
// 			fmt.Println(in.Text())
// 		}
// 		wg.Done()
// 	}()

// 	go func() {
// 		if err = command.Start(); err != nil {
// 			initLog.Errorf("AMF Start error: %+v", err)
// 		}
// 		wg.Done()
// 	}()

// 	wg.Wait()

// 	return err
// }

// // Used in AMF planned removal procedure
// func (amf *LB) Terminate() {
// 	logger.InitLog.Infof("Terminating AMF...")
// 	amfSelf := context.AMF_Self()

// 	// TODO: forward registered UE contexts to target AMF in the same AMF set if there is one

// 	// deregister with NRF
// 	problemDetails, err := consumer.SendDeregisterNFInstance()
// 	if problemDetails != nil {
// 		logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
// 	} else if err != nil {
// 		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
// 	} else {
// 		logger.InitLog.Infof("[AMF] Deregister from NRF successfully")
// 	}

// 	// send AMF status indication to ran to notify ran that this AMF will be unavailable
// 	logger.InitLog.Infof("Send AMF Status Indication to Notify RANs due to AMF terminating")
// 	unavailableGuamiList := ngap_message.BuildUnavailableGUAMIList(amfSelf.ServedGuamiList)
// 	amfSelf.AmfRanPool.Range(func(key, value interface{}) bool {
// 		ran := value.(*context.AmfRan)
// 		ngap_message.SendAMFStatusIndication(ran, unavailableGuamiList)
// 		return true
// 	})

// 	ngap_service.Stop()

// 	callback.SendAmfStatusChangeNotify((string)(models.StatusChange_UNAVAILABLE), amfSelf.ServedGuamiList)
// 	logger.InitLog.Infof("AMF terminated")
// }