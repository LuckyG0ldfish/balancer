package message

import (
	"fmt"

	"github.com/LuckyG0ldfish/balancer/context"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/ngap/ngapType"
)

// Builds and sends NGSetupRequest
func SendNGSetupRequest(conn *context.Lb_Conn) {
	lbID := 1

	hexGNBID := []byte{0x00, 0x02, byte(lbID)}
	gnbName := fmt.Sprintf("f5gcLB_%d", lbID)

	sendMsg, err := BuildNGSetupRequest(hexGNBID, 24, gnbName)
	if err != nil {
		logger.NgapLog.Errorf("Build NGSetup Request failed: [%+v]\n", err)
		return
	}
	context.SendByteToConn(conn.Conn, sendMsg)
}

// Builds and sends NGSetupResponse
func SendNGSetupResponse(conn *context.Lb_Conn) {
	conn.Log.Info("Send NG-Setup response")
	pkt, err := BuildNGSetupResponse()
	if err != nil {
		logger.NgapLog.Errorf("Build NGSetupResponse failed : %s", err.Error())
		return
	}
	context.SendByteToConn(conn.Conn, pkt)
}

// Builds and sends NGSetupFailure
func SendNGSetupFailure(conn *context.Lb_Conn, cause ngapType.Cause) {
	conn.Log.Infoln("Send NG-Setup failure")

	if cause.Present == ngapType.CausePresentNothing {
		logger.NgapLog.Errorf("Cause present is nil")
		return
	}

	pkt, err := BuildNGSetupFailure(cause)
	if err != nil {
		logger.NgapLog.Errorf("Build NGSetupFailure failed : %s", err.Error())
		return
	}
	context.SendByteToConn(conn.Conn, pkt)
}