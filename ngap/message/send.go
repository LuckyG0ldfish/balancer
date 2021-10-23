package message

import (
	"fmt"
	"github.com/free5gc/ngap/ngapType"
	"github.com/LuckyG0ldfish/balancer/context"
)



func SendNGSetupResponse(conn *context.LBConn) {
	// lb.Log.Info("Send NG-Setup response")
	fmt.Println("Send NG-Setup response")

	pkt, err := BuildNGSetupResponse()
	if err != nil {
		// lb.Log.Errorf("Build NGSetupResponse failed : %s", err.Error())
		return
	}
	context.SendByteToConn(conn.Conn, pkt)
}

func SendNGSetupFailure(conn *context.LBConn, cause ngapType.Cause) {
	// ran.Log.Info("Send NG-Setup failure")

	if cause.Present == ngapType.CausePresentNothing {
		// lb.Log.Errorf("Cause present is nil")
		return
	}

	pkt, err := BuildNGSetupFailure(cause)
	if err != nil {
		// lb.Log.Errorf("Build NGSetupFailure failed : %s", err.Error())
		return
	}
	context.SendByteToConn(conn.Conn, pkt)
}