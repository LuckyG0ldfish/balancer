package message

import (
	"fmt"
	"github.com/free5gc/ngap/ngapType"
	"github.com/LuckyG0ldfish/balancer/context"
	"github.com/LuckyG0ldfish/balancer/logger"
)

func SendNGSetupRequest(conn *context.LBConn) {

	// lbCtx := context.LB_Self()

	// TODO: Replace ID - is only for testing (use the context ^)
	lbID := 1

	hexGNBID := []byte{0x00, 0x02, byte(lbID)}
	gnbName := fmt.Sprintf("free5gc_%d", lbID)

	sendMsg, err := BuildNGSetupRequest(hexGNBID, 24, gnbName)
	if err != nil {
		logger.NgapLog.Errorf("Build NGSetup Request failed: [%+v]\n", err)
		return
	}
	context.SendByteToConn(conn.Conn, sendMsg)
	// n, err := conn.Conn.Write(sendMsg)
	// if err != nil {
	// 	// NGAPLog.Errorf("Write to SCTP socket failed: [%+v]", err)
	// } else {
	// 	// NGAPLog.Tracef("Wrote %d bytes", n)
	// }
}

func SendNGSetupResponse(conn *context.LBConn) {
	// lb.Log.Info("Send NG-Setup response")
	fmt.Println("Send NG-Setup response")
	// pkt, err := BuildNGSetupResponse()
	//
	pkt, err := BuildNGSetupResponse()

	if err != nil {
		logger.NgapLog.Errorf("Build NGSetupResponse failed : %s", err.Error())
		// fmt.Println("BuildNGSetupResponse failed")
		return
	}
	context.SendByteToConn(conn.Conn, pkt)
}

func SendNGSetupFailure(conn *context.LBConn, cause ngapType.Cause) {
	// ran.Log.Info("Send NG-Setup failure")

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