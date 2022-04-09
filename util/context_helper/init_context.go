package context_helper

import (
	"github.com/google/uuid"

	"github.com/LuckyG0ldfish/balancer/context"
	"github.com/LuckyG0ldfish/balancer/factory"
	"github.com/LuckyG0ldfish/balancer/logger"
)

func InitLbContext(self *context.LB_Context) {
	config := factory.LbConfig
	logger.UtilLog.Infof("lbconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	
	// LB Settings 
	self.NfId = uuid.New().String()
	if configuration.LbName != "" {
		self.Name = configuration.LbName
	}
	if configuration.NgapIp != "" {
		self.LbIP = configuration.NgapIp
	} else {
		self.LbIP = "127.0.0.1" // default localhost
		logger.CfgLog.Warnf("Default IP selected")
	}
	if configuration.NgapListenPort != 0 {
		self.LbListenPort = configuration.NgapListenPort
	} else {
		self.LbListenPort = 48484 // default Port
		logger.CfgLog.Warnf("Default Listen-Port selected")
	}
	if configuration.NgapLbToAmfPort != 0 {
		self.LbToAmfPort = configuration.NgapLbToAmfPort
	} else {
		self.LbToAmfPort = 23232 // default Port
		logger.CfgLog.Warnf("Default LbToAmf-Port selected")
	}
	self.MetricsLevel = configuration.MetricsLevel

	addr, err := context.GenSCTPAddr(self.LbIP, self.LbListenPort)
	if err == nil {
		self.LbListenAddr = addr
		logger.CfgLog.Tracef("LbAddr set")
	} else {
		logger.CfgLog.Warnf("LbAddr couldn't be set")
	}
	addr, err = context.GenSCTPAddr(self.LbIP, self.LbToAmfPort)
	if err == nil {
		self.LbToAmfAddr = addr
		logger.CfgLog.Tracef("LbAddr set")
	} else {
		logger.CfgLog.Warnf("LbAddr couldn't be set")
	}
	self.ServedGuamiList = configuration.ServedGumaiList
	self.PlmnSupportList = configuration.PlmnSupportList
	
	self.Running = true
	self.IDGen = context.NewUniqueNumberGen(1) // internal LbUe.ID for the first UE 

	self.DifferentAmfTypes = configuration.DifferentAmfTypes
	self.ContinuesAmfRegistration = configuration.ContinuesAmfRegistration

	// adding AMFs 
	if configuration.RegistrationAmfNgapIpList != nil {
		self.NewRegistAmfIpList = configuration.RegistrationAmfNgapIpList
		self.NewAmf = true
	} else {
		self.NewRegistAmfIpList = []string{"127.0.0.1"} // default localhost
		logger.CfgLog.Warnf("Default Registration-AMF-list selected")
	}
	if self.DifferentAmfTypes == 3 {
		if configuration.RegularAmfNgapIpList != nil {
			self.NewRegularAmfIpList = configuration.RegularAmfNgapIpList
			self.NewAmf = true
		} else {
			self.NewRegularAmfIpList = []string{"127.0.0.1"} // default localhost
			logger.CfgLog.Warnf("Default Regular-AMF-list selected")
		}
	}
	if self.DifferentAmfTypes >= 2 {
		if configuration.DeregistrationAmfNgapIpList != nil {
			self.NewDeregistAmfIpList = configuration.DeregistrationAmfNgapIpList
			self.NewAmf = true
		} else {
			self.NewDeregistAmfIpList = []string{"127.0.0.1"} // default localhost
			logger.CfgLog.Warnf("Default Deregistration-AMF-list selected")
		}
	}
}
