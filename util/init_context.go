package util

import (
	// "os"

	"github.com/google/uuid"

	"github.com/LuckyG0ldfish/balancer/context"
	"github.com/LuckyG0ldfish/balancer/factory"
	"github.com/LuckyG0ldfish/balancer/logger"
	// "github.com/free5gc/nas/security"
	// "github.com/free5gc/openapi/models"
)

func InitLbContext(context *context.LBContext) {
	config := factory.AmfConfig
	logger.UtilLog.Infof("amfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	context.NfId = uuid.New().String()
	if configuration.LbName != "" {
		context.Name = configuration.LbName
	}
	if configuration.NgapIp != "" {
		context.LbIP = configuration.NgapIp
	} else {
		context.LbIP = "127.0.0.1" // default localhost
	}
	if configuration.NgapPort != 0 {
		context.LbPort = configuration.NgapPort
	} else {
		context.LbPort = 48484 // default Port
	}

	// adding AMFs 
	if configuration.AmfNgapIp != "" {
		context.NewAmfIp = configuration.AmfNgapIp
	} else {
		context.NewAmfIp = "127.0.0.1" // default localhost
	}
	if configuration.AmfNgapPort != 0 {
		context.NewAmfPort = configuration.AmfNgapPort
	} else {
		context.NewAmfPort = 38412 // default Port for AMF
	}
	context.NewAmf = true 

	context.Running = true 
	// if 
	// sbi := configuration.Sbi
	// if sbi.Scheme != "" {
	// 	context.UriScheme = models.UriScheme(sbi.Scheme)
	// } else {
	// 	logger.UtilLog.Warnln("SBI Scheme has not been set. Using http as default")
	// 	context.UriScheme = "http"
	// }
	// context.RegisterIPv4 = factory.AMF_DEFAULT_IPV4 // default localhost
	// context.SBIPort = factory.AMF_DEFAULT_PORT_INT  // default port
	// if sbi != nil {
	// 	if sbi.RegisterIPv4 != "" {
	// 		context.RegisterIPv4 = sbi.RegisterIPv4
	// 	}
	// 	if sbi.Port != 0 {
	// 		context.SBIPort = sbi.Port
	// 	}
	// 	context.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
	// 	if context.BindingIPv4 != "" {
	// 		logger.UtilLog.Info("Parsing ServerIPv4 address from ENV Variable.")
	// 	} else {
	// 		context.BindingIPv4 = sbi.BindingIPv4
	// 		if context.BindingIPv4 == "" {
	// 			logger.UtilLog.Warn("Error parsing ServerIPv4 address from string. Using the 0.0.0.0 as default.")
	// 			context.BindingIPv4 = "0.0.0.0"
	// 		}
	// 	}
	// }
	// serviceNameList := configuration.ServiceNameList
	// context.InitNFService(serviceNameList, config.Info.Version)
	// context.ServedGuamiList = configuration.ServedGumaiList
	// context.SupportTaiLists = configuration.SupportTAIList
	// for i := range context.SupportTaiLists {
	// 	context.SupportTaiLists[i].Tac = TACConfigToModels(context.SupportTaiLists[i].Tac)
	// }
	// context.PlmnSupportList = configuration.PlmnSupportList
	// context.SupportDnnLists = configuration.SupportDnnList
	// if configuration.NrfUri != "" {
	// 	context.NrfUri = configuration.NrfUri
	// } else {
	// 	logger.UtilLog.Warn("NRF Uri is empty! Using localhost as NRF IPv4 address.")
	// 	context.NrfUri = factory.AMF_DEFAULT_NRFURI
	// }
	// security := configuration.Security
	// if security != nil {
	// 	context.SecurityAlgorithm.IntegrityOrder = getIntAlgOrder(security.IntegrityOrder)
	// 	context.SecurityAlgorithm.CipheringOrder = getEncAlgOrder(security.CipheringOrder)
	// }
	// context.NetworkName = configuration.NetworkName
	// context.T3502Value = configuration.T3502Value
	// context.T3512Value = configuration.T3512Value
	// context.Non3gppDeregistrationTimerValue = configuration.Non3gppDeregistrationTimerValue
	// context.T3513Cfg = configuration.T3513
	// context.T3522Cfg = configuration.T3522
	// context.T3550Cfg = configuration.T3550
	// context.T3560Cfg = configuration.T3560
	// context.T3565Cfg = configuration.T3565
	// context.Locality = configuration.Locality
}

// func getIntAlgOrder(integrityOrder []string) (intOrder []uint8) {
// 	for _, intAlg := range integrityOrder {
// 		switch intAlg {
// 		case "NIA0":
// 			intOrder = append(intOrder, security.AlgIntegrity128NIA0)
// 		case "NIA1":
// 			intOrder = append(intOrder, security.AlgIntegrity128NIA1)
// 		case "NIA2":
// 			intOrder = append(intOrder, security.AlgIntegrity128NIA2)
// 		case "NIA3":
// 			intOrder = append(intOrder, security.AlgIntegrity128NIA3)
// 		default:
// 			logger.UtilLog.Errorf("Unsupported algorithm: %s", intAlg)
// 		}
// 	}
// 	return
// }

// func getEncAlgOrder(cipheringOrder []string) (encOrder []uint8) {
// 	for _, encAlg := range cipheringOrder {
// 		switch encAlg {
// 		case "NEA0":
// 			encOrder = append(encOrder, security.AlgCiphering128NEA0)
// 		case "NEA1":
// 			encOrder = append(encOrder, security.AlgCiphering128NEA1)
// 		case "NEA2":
// 			encOrder = append(encOrder, security.AlgCiphering128NEA2)
// 		case "NEA3":
// 			encOrder = append(encOrder, security.AlgCiphering128NEA3)
// 		default:
// 			logger.UtilLog.Errorf("Unsupported algorithm: %s", encAlg)
// 		}
// 	}
// 	return
// }