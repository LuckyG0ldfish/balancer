package main

import (
	"fmt"
	"os"
	// "time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/LuckyG0ldfish/balancer/service"
	"github.com/free5gc/version"
)

var LB = &service.Load{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {	
	app := cli.NewApp()
	app.Name = "lb"
	appLog.Infoln(app.Name)
	appLog.Infoln("LB version: ", version.GetVersion())
	app.Usage = "-lbcfg lb configuration file" // -free5gccfg common configuration file
	app.Action = action
	app.Flags = LB.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		appLog.Errorf("LB Run error: %v", err)
		return
	}
	
	// for{
	// 	time.Sleep(1 *time.Hour)
	// }
}

func action(c *cli.Context) error {
	if err := LB.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("failed to initialize")
	}

	LB.Start()

	return nil
}
