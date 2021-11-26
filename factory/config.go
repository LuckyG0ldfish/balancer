/*
 * AMF Configuration Factory
 */

 package factory

 import (
	 "github.com/free5gc/logger_util" 
 )
 
 const (
	 AMF_EXPECTED_CONFIG_VERSION = "1.0.2"
 )
 
 type Config struct {
	 Info          *Info               `yaml:"info"`
	 Configuration *Configuration      `yaml:"configuration"`
	 Logger        *logger_util.Logger `yaml:"logger"`
 }
 
 type Info struct {
	 Version     string `yaml:"version,omitempty"`
	 Description string `yaml:"description,omitempty"`
 }
 
 type Configuration struct {
	 LbName                          string                    	`yaml:"lbName,omitempty"`
	 NgapIp                      	 string                  	`yaml:"ngapIp,omitempty"`
	 NgapListenPort					 int						`yaml:"ngapListenPort,omitempty"`
	 NgapLbToAmfPort				 int						`yaml:"ngapLbToAmfPort,omitempty"`
	 AmfNgapIpList					 []string					`yaml:"amfNgapIpList,omitempty"`
	 AmfNgapPortList				 []string					`yaml:"amfNgapPortList,omitempty"`
	//  Sbi                             *Sbi                      `yaml:"sbi,omitempty"`
	//  NetworkFeatureSupport5GS        *NetworkFeatureSupport5GS `yaml:"networkFeatureSupport5GS,omitempty"`
	//  ServiceNameList                 []string                  `yaml:"serviceNameList,omitempty"`
	//  ServedGumaiList                 []models.Guami            `yaml:"servedGuamiList,omitempty"`
	//  SupportTAIList                  []models.Tai              `yaml:"supportTaiList,omitempty"`
	//  PlmnSupportList                 []PlmnSupportItem         `yaml:"plmnSupportList,omitempty"`
	//  SupportDnnList                  []string                  `yaml:"supportDnnList,omitempty"`
 }
 
 func (c *Config) GetVersion() string {
	 if c.Info != nil && c.Info.Version != "" {
		 return c.Info.Version
	 }
	 return ""
 }
 