package context

import (
	"encoding/hex"
	"time"

	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
)

// Used to forward unregistered UEs to an preselected AMF
// func ForwardToNextAmf(lbConn *LBConn, message *ngapType.NGAPPDU, ue *LbUe, startTime int64) {
// 	lb := LB_Self()
	

// 	// Forwarding the message
// 	var mes []byte
// 	mes, _ = ngap.Encoder(*message)
// 	next.LbConn.Conn.Write(mes)
	
// 	logger.NgapLog.Tracef("Forward to nextAMF:")
// 	logger.NgapLog.Tracef("Packet content:\n%+v", hex.Dump(mes))
// 	logger.NgapLog.Tracef("UeRanID: %d | UeLbID: %d", uint64(ue.UeRanID), uint64(ue.UeLbID))
	
// 	/* Metrics */
// 	// Adding new Trace to the routing table 
// 	now :=  int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Microsecond)
// 	if lb.MetricsLevel > 0 {
// 		AddRouting_Element(lb.MetricsUEs, ue.RanID, ue.UeRanID, ue.AmfID, TypeAmf, ue.UeStateIdent, startTime, now)
// 	}
	
// 	// Selecting AMF that will be used for the next new UE
// 	lb.SelectNextAmf()
// }

// Used to forward registered UE's messages to an AMF
func ForwardToAmf(message *ngapType.NGAPPDU, ue *LbUe, startTime int64) {
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
	if ue.UeAmfID != 0 {
		logger.NgapLog.Tracef("UeRanID: %d | UeLbID: %d | UeAmfID: %d", uint64(ue.UeRanID), uint64(ue.UeLbID), uint64(ue.UeAmfID))
	} else {
		logger.NgapLog.Tracef("UeRanID: %d | UeLbID: %d", uint64(ue.UeRanID), uint64(ue.UeLbID))
	}
	
	/* Metrics */
	// Adding new Trace to the routing table 
	now :=  int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Microsecond)
	lb := LB_Self()
	if lb.MetricsLevel > 0 {
		AddRouting_Element(lb.MetricsGNBs, ue.RanID, ue.UeRanID, ue.AmfID, TypeAmf, ue.UeStateIdent, startTime, now)
	}
}

// Used to forward registered UE's messages to an GNB
func ForwardToGnb(message *ngapType.NGAPPDU, ue *LbUe, startTime int64) {
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
	if ue.UeAmfID != 0 {
		logger.NgapLog.Tracef("UeRanID: %d | UeLbID: %d | UeAmfID: %d", uint64(ue.UeRanID), uint64(ue.UeLbID), uint64(ue.UeAmfID))
	} else {
		logger.NgapLog.Tracef("UeRanID: %d | UeLbID: %d", uint64(ue.UeRanID), uint64(ue.UeLbID))
	}

	/* Metrics */
	// Adding new Trace to the routing table 
	now :=  int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Microsecond)
	lb := LB_Self()

	_, ok := lb.MetricsGNBs.Load(ue.UeRanID)
	if !ok {
		logger.ContextLog.Errorln("failed mgnb add #3")
	} 

	if lb.MetricsLevel > 0 {
		AddRouting_Element(lb.MetricsGNBs, ue.AmfID, ue.UeRanID, ue.RanID, TypeGnb, ue.UeStateIdent, startTime, now)
	}
}
