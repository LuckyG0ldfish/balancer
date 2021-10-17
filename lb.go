package main

import (
	// "fmt"
	// "os"

	"github.com/sirupsen/logrus"
	// "github.com/urfave/cli"

	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/LuckyG0ldfish/balancer/service"
	// "github.com/free5gc/version"
)

var LB = &service.LB{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	LB.Initialize()
	LB.Start()
	// app := cli.NewApp()
	// app.Name = "amf"
	// appLog.Infoln(app.Name)
	// appLog.Infoln("AMF version: ", version.GetVersion())
	// app.Usage = "-free5gccfg common configuration file -amfcfg amf configuration file"
	// app.Action = action
	// app.Flags = LB.GetCliCmd()
	// if err := app.Run(os.Args); err != nil {
	// 	appLog.Errorf("AMF Run error: %v", err)
	// 	return
	// }
}

// func action(c *cli.Context) error {
// 	if err := LB.Initialize(c); err != nil {
// 		// logger.CfgLog.Errorf("%+v", err)
// 		return fmt.Errorf("Failed to initialize !!")
// 	}

// 	LB.Start()

// 	return nil
// }
