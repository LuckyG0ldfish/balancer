package util

import (
	"github.com/google/uuid"

	"github.com/LuckyG0ldfish/balancer/context"
	"github.com/LuckyG0ldfish/balancer/factory"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/LuckyG0ldfish/balancer/metrics"
)

func InitLbContext(self *context.LBContext) {
	config := factory.LbConfig
	logger.UtilLog.Infof("amfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
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
	self.Table = metrics.NewTable()

	// adding AMFs 
	if configuration.AmfNgapIpList != nil {
		self.NewAmfIpList = configuration.AmfNgapIpList
		self.NewAmf = true
	} else {
		self.NewAmfIpList = []string{"127.0.0.1"} // default localhost
		logger.CfgLog.Warnf("Default AMF-list selected")
		self.NewAmf = false
	}
}
