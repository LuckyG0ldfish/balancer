/*
 * LB Configuration Factory
 */

 package factory

 import (
	Logger "github.com/LuckyG0ldfish/balancer/util/logger_helper"
	"github.com/free5gc/openapi/models"
 )
 
 const (
	 AMF_EXPECTED_CONFIG_VERSION = "1.0.2"
 )
 
 type Config struct {
	 Info          *Info               `yaml:"info"`
	 Configuration *Configuration      `yaml:"configuration"`
	 Logger        *Logger.Logger 		`yaml:"logger"`
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
	 DifferentAmfTypes 				 int 						`yaml:"differentAmfTypes,omitempty"`
	 RegistrationAmfNgapIpList		 []string					`yaml:"registrationAmfNgapIpList,omitempty"`
	 DeregistrationAmfNgapIpList	 []string					`yaml:"deregistrationAmfNgapIpList,omitempty"`
	 RegularAmfNgapIpList			 []string					`yaml:"regularAmfNgapIpList,omitempty"`
	 ServedGumaiList                 []models.Guami            	`yaml:"servedGuamiList,omitempty"`
	 PlmnSupportList                 []PlmnSupportItem         	`yaml:"plmnSupportList,omitempty"`
	 MetricsLevel 					 int 						`yaml:"metricsLevel,omitempty"`
	 ContinuesAmfRegistration 		 bool 						`yaml:"continuesAmfRegistration,omitempty"`
 }

 type PlmnSupportItem struct {
	PlmnId     models.PlmnId   `yaml:"plmnId"`
	SNssaiList []models.Snssai `yaml:"snssaiList,omitempty"`
}

 func (c *Config) GetVersion() string {
	 if c.Info != nil && c.Info.Version != "" {
		 return c.Info.Version
	 }
	 return ""
 }
 