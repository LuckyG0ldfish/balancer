package logger_helper

type Logger struct {
	LB   			*LogSetting `yaml:"LB"`
	
	// AMF   			*LogSetting `yaml:"AMF"`

	// GNB				*LogSetting `yaml:"GNB"`
	
	// NGAP            *LogSetting `yaml:"NGAP"`
	// AMFHandler 		*LogSetting `yaml:"AMFHandler"`
	// GNBHandler		*LogSetting `yaml:"GNBHandler"`
	
	// NAS             *LogSetting `yaml:"NAS"`
	
	// UELog 			*LogSetting `yaml:"UE"`
}

type LogSetting struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}
