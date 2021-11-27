package context

import (
	// "sync"

	"encoding/hex"
	"strconv"

	"git.cs.nctu.edu.tw/calee/sctp"
	// "github.com/free5gc/amf/factory" // TODO
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"

	// "github.com/free5gc/openapi/models"
	"github.com/sirupsen/logrus"
)

var (
	lbContext = LBContext{}
)

type LBContext struct {
	Name 				string
	// NetworkName   		factory.NetworkName
	NfId               	string

	LbIP 				string

	LbToAmfPort			int 
	LbToAmfAddr			*sctp.SCTPAddr

	LbListenPort		int
	LbListenAddr		*sctp.SCTPAddr

	Running 			bool

	NewAmf				bool // indicates that a new AMF IP+Port have been added so that the LB can connect to it 
	NewAmfIpList 		[]string 
	NewAmfPortList		[]string 
	
	LbRanPool 			[]*LbGnb // gNBs connected to the LB
	LbAmfPool 			[]*LbAmf // amfs (each connected to AMF 1:1) connected to LB

	Next_Amf 			*LbAmf

	IDGen 				*UeIdGen
	
	RelativeCapacity 	int64 // To build setup response

	/* temp */
	NGSetupRes 			*ngapType.NGAPPDU
	PlmnSupportList 	*ngapType.PLMNSupportList
	ServedGuamiList 	*ngapType.ServedGUAMIList

	/* logger */
	Log 	*logrus.Entry
}

func NewLBContext() (LbContext *LBContext){
	return 
}

func (lb *LBContext) ForwardToNextAmf(lbConn *LBConn, message *ngapType.NGAPPDU, ue *LbUe) { 
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
	_, ok := next.Ues.Load(ue.UeLbID)
	if ok {
		logger.NgapLog.Errorf("UE already exists")
		return 
	} 
	next.Ues.Store(ue.UeLbID, ue)

	// Forwarding the message
	var mes []byte
	mes, _  = ngap.Encoder(*message)
	next.LbConn.Conn.Write(mes)
	next.Capacity -= 1
	logger.NgapLog.Debugf("Forward to nextAMF:")
	logger.NgapLog.Debugf("Packet content:\n%+v", hex.Dump(mes))
	logger.NgapLog.Tracef("UeLbID: " + strconv.FormatInt(ue.UeLbID, 10) + " | UeRanID: " + strconv.FormatInt(ue.UeRanID, 10))
	
	// Felecting AMF that will be used for the next new UE 
	lb.SelectNextAmf()
}

func (lb *LBContext) ForwardToAmf(lbConn *LBConn, message *ngapType.NGAPPDU, ue *LbUe) {
	
	
	// finding the correct AMF by the in UE stored AMF-ID 
	amf, ok := lb.LbAmfFindByID(ue.AmfID)
	if ok {
		var mes []byte
		mes, _  = ngap.Encoder(*message)
		amf.LbConn.Conn.Write(mes)
		logger.NgapLog.Debugf("Message forwarded to AMF")
		logger.NgapLog.Tracef("Packet content:\n%+v", hex.Dump(mes))
		logger.NgapLog.Tracef("UeLbID: " + strconv.FormatInt(ue.UeLbID, 10) + " | UeRanID: " + strconv.FormatInt(ue.UeRanID, 10))
	} else {
		logger.NgapLog.Errorf("AMF not found")
	}
}

func (lb *LBContext) ForwardToGnb(lbConn *LBConn, message *ngapType.NGAPPDU, ue *LbUe) { //*ngapType.NGAPPDU
	
	// finding the correct GNB by the in UE stored AMF-ID 
	gnb, ok := lb.LbGnbFindByID(ue.RanID)
	if ok {
		var mes []byte
		mes, _  = ngap.Encoder(*message)
		gnb.LbConn.Conn.Write(mes)
		logger.NgapLog.Debugf("Message forwarded to GNB")
		logger.NgapLog.Tracef("Packet content:\n%+v", hex.Dump(mes))
		logger.NgapLog.Tracef("UeLbID: " + strconv.FormatInt(ue.UeLbID, 10) + " | UeRanID: " + strconv.FormatInt(ue.UeRanID, 10))
	} else {
		logger.NgapLog.Errorf("GNB not found")
	}
}

// use sctp.SCTPConn to find RAN context, return *LbRan and true if found
func (context *LBContext) LbGnbFindByConn(conn *sctp.SCTPConn) (*LbGnb, bool) {
	for _, v := range context.LbRanPool {
		if v.LbConn.Conn == conn {
			return v, true
		}
	}
	return nil, false
}

// use ID to find RAN context, return *LbRan and true if found
func (context *LBContext) LbGnbFindByID(ranID int64) (*LbGnb, bool) {
	for _, v := range context.LbRanPool {
		if v.GnbID == ranID {
			return v, true
		}
	}
	return nil, false
}

// use sctp.SCTPConn to find AMF context, return *LBAmf and true if found
func (context *LBContext) LbAmfFindByConn(conn *sctp.SCTPConn) (*LbAmf, bool) {
	for _, v := range context.LbAmfPool {
		if v.LbConn.Conn == conn {
			return v, true
		}
	}
	return nil, false
}

// use ID to find AMF context, return *LbAmf and true if found
func (context *LBContext) LbAmfFindByID(amfID int64) (*LbAmf, bool) {
	for _, v := range context.LbAmfPool {
		if v.AmfID == amfID {
			return v, true
		}
	}
	return nil, false
}

// use UeID to find UE context, return *AmfRan and true if found
func (context *LBContext) LbAmfFindByUeID(UeID int64) (*LbAmf, bool) {
	for _, Amf := range context.LbAmfPool {
		if check := Amf.ContainsUE(UeID); check {
			return Amf, true
		}
	}
	return nil, false
}

// use sctp.SCTPConn to find RAN context, return *AmfRan and true if found
func (context *LBContext) LbGnbFindByUeID(UeID int64) (*LbGnb, *LbUe, bool) {
	for _, Gnb := range context.LbRanPool {
		if ue, check := Gnb.Ues.Load(UeID); check {
			UE, ok := ue.(*LbUe)
			if !ok {
				return nil, nil, false
			}
			return Gnb, UE, true
		}
	}
	return nil, nil, false
}

// use sctp.SCTPConn to find RAN context, return *AmfRan and true if found
func (context *LBContext) AddNewGnbToLB(conn *sctp.SCTPConn) *LbGnb{
	gnb := NewLbGnb()
	gnb.LbConn.Conn = conn
	context.LbRanPool = append(context.LbRanPool, gnb)
	return gnb
}

// 
func (context *LBContext) AddAmfToLB(amf *LbAmf) *LbAmf{
	context.LbAmfPool = append(context.LbAmfPool, amf)
	return amf
}

// TODO
func (context *LBContext) SelectNextAmf() bool{
	if context.LbAmfPool == nil {
		logger.ContextLog.Errorf("No AMF found")
		return false
	}
	var i int64
	var amf *LbAmf
	for in, v := range context.LbAmfPool {
		if in == 0 {
			i = v.Capacity
			amf = v
		} else if v.Capacity < i {
			i = v.Capacity
			amf = v
		}
	}
	context.Next_Amf = amf
	logger.ContextLog.Tracef("NextAMF = AMFID: %d", amf.AmfID)
	return true 
}

func LB_Self() *LBContext {
	return &lbContext
}
