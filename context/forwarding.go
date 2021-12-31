package context

import (
	"encoding/hex"
	"strconv"

	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
)

// Used to forward unregistered UEs to an preselected AMF
func ForwardToNextAmf(lbConn *LBConn, message *ngapType.NGAPPDU, ue *LbUe) {
	lb := LB_Self()
	if lb.Next_Amf == nil {
		logger.NgapLog.Errorf("No Connected AMF / No AMf set as next AMF")
		return
	}

	// Temporarily stores the pointer to the chosen AMF so no
	// parallelized process will change it during runtime
	next := lb.Next_Amf

	// Checks whether an UE with this UeLbID already exists
	// and otherwise adds it
	ue.AmfID = next.AmfID
	ue.AmfPointer = next
	_, ok := next.Ues.Load(ue.UeLbID)
	if ok {
		logger.NgapLog.Errorf("UE already exists")
		return
	}
	next.Ues.Store(ue.UeLbID, ue)

	// Forwarding the message
	var mes []byte
	mes, _ = ngap.Encoder(*message)
	next.LbConn.Conn.Write(mes)
	next.NumberOfConnectedUEs += 1
	logger.NgapLog.Debugf("Forward to nextAMF:")
	logger.NgapLog.Debugf("Packet content:\n%+v", hex.Dump(mes))
	logger.NgapLog.Tracef("UeLbID: " + strconv.FormatInt(ue.UeLbID, 10) + " | UeRanID: " + strconv.FormatInt(ue.UeRanID, 10))

	// Selecting AMF that will be used for the next new UE
	lb.SelectNextAmf()
}

// Used to forward registered UE's messages to an AMF
func ForwardToAmf(message *ngapType.NGAPPDU, ue *LbUe) {
	// finding the correct AMF by the in UE stored AMF-Pointer
	amf := ue.AmfPointer

	// Encoding
	var mes []byte
	mes, err := ngap.Encoder(*message)
	if err != nil {
		logger.NgapLog.Errorf("Message encoding failed")
		return
	}

	// Forwarding
	SendByteToConn(amf.LbConn.Conn, mes)
	logger.NgapLog.Debugf("Message forwarded to AMF")
	logger.NgapLog.Tracef("Packet content:\n%+v", hex.Dump(mes))
	logger.NgapLog.Tracef("UeLbID: " + strconv.FormatInt(ue.UeLbID, 10) + " | UeRanID: " + strconv.FormatInt(ue.UeRanID, 10))
}

// Used to forward registered UE's messages to an GNB
func ForwardToGnb(message *ngapType.NGAPPDU, ue *LbUe) {
	// finding the correct GNB by the in UE stored AMF-Pointer
	gnb := ue.RanPointer

	// Encoding
	var mes []byte
	mes, err := ngap.Encoder(*message)
	if err != nil {
		logger.NgapLog.Errorf("Message encoding failed")
		return
	}

	// Forwarding
	SendByteToConn(gnb.LbConn.Conn, mes)
	logger.NgapLog.Debugf("Message forwarded to GNB")
	logger.NgapLog.Tracef("Packet content:\n%+v", hex.Dump(mes))
	logger.NgapLog.Tracef("UeLbID: " + strconv.FormatInt(ue.UeLbID, 10) + " | UeRanID: " + strconv.FormatInt(ue.UeRanID, 10))

}
