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

// func init() {
// 	LB_Self().Name = "lb"
// 	//LB_Self().NetworkName.Full = "free5GC"
// }

type LBContext struct {
	Name 				string
	// NetworkName   		factory.NetworkName
	NfId               	string

	LbIP 				string
	LbPort				int

	Running 			bool

	NewAmf				bool
	NewAmfIp 			string 
	NewAmfPort 			int 
	
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
	// if mes, err := ngap.Encoder(*message); err == nil {
	// }
	ue.AmfID = lb.Next_Amf.AmfID
	_, ok := lb.Next_Amf.Ues.Load(ue.UeLbID)
	if lb.Next_Amf == nil {
		logger.NgapLog.Errorf("No Connected AMF")
		return 
	}
	if ok {
		logger.NgapLog.Errorf("UE already exists")
		return 
	} 
	lb.Next_Amf.Ues.Store(ue.UeLbID, ue)

	var mes []byte
	mes, _  = ngap.Encoder(*message)
	lb.Next_Amf.LbConn.Conn.Write(mes)
	lb.Next_Amf.Capacity -= 1
	lb.SelectNextAmf()
	logger.NgapLog.Debugf("forward to nextAMF:")
	logger.NgapLog.Debugf("Packet content:\n%+v", hex.Dump(mes))
	logger.NgapLog.Tracef("UeLbID: " + strconv.FormatInt(ue.UeLbID, 10) + " | UeRanID: " + strconv.FormatInt(ue.UeRanID, 10))
}

func (lb *LBContext) ForwardToAmf(lbConn *LBConn, message *ngapType.NGAPPDU, ue *LbUe) {
	amf, ok := lb.LbAmfFindByID(ue.AmfID)
	// if mes, err := ngap.Encoder(*message); err == nil {	
	// }
	if ok {
		var mes []byte
		mes, _  = ngap.Encoder(*message)
		amf.LbConn.Conn.Write(mes)
		logger.NgapLog.Debugf("forward to AMF:")
		logger.NgapLog.Debugf("Packet content:\n%+v", hex.Dump(mes))
		logger.NgapLog.Tracef("UeLbID: " + strconv.FormatInt(ue.UeLbID, 10) + " | UeRanID: " + strconv.FormatInt(ue.UeRanID, 10))
	} else {
		logger.NgapLog.Errorf("AMF not found")
	}
}

func (lb *LBContext) ForwardToGnb(lbConn *LBConn, message *ngapType.NGAPPDU, ue *LbUe) { //*ngapType.NGAPPDU
	gnb, ok := lb.LbGnbFindByID(ue.RanID)
	// if mes, err := ngap.Encoder(*message); err == nil {	
	// }
	if ok {
		var mes []byte
		mes, _  = ngap.Encoder(*message)
		gnb.LbConn.Conn.Write(mes)
		logger.NgapLog.Debugf("forward to GNB:")
		logger.NgapLog.Debugf("Packet content:\n%+v", hex.Dump(mes))
		logger.NgapLog.Tracef("UeLbID: " + strconv.FormatInt(ue.UeLbID, 10) + " | UeRanID: " + strconv.FormatInt(ue.UeRanID, 10))
	} else {
		logger.NgapLog.Errorf("GNB not found")
	}
}

// use net.Conn to find RAN context, return *AmfRan and ok bit
func (context *LBContext) LbGnbFindByConn(conn *sctp.SCTPConn) (*LbGnb, bool) {
	for _, v := range context.LbRanPool {
		if v.LbConn.Conn == conn {
			return v, true
		}
	}
	return nil, false
}

func (context *LBContext) LbGnbFindByID(ranID int64) (*LbGnb, bool) {
	for _, v := range context.LbRanPool {
		if v.GnbID == ranID {
			return v, true
		}
	}
	return nil, false
}

func (context *LBContext) LbAmfFindByConn(conn *sctp.SCTPConn) (*LbAmf, bool) {
	for _, v := range context.LbAmfPool {
		if v.LbConn.Conn == conn {
			return v, true
		}
	}
	return nil, false
}

func (context *LBContext) LbAmfFindByID(amfID int64) (*LbAmf, bool) {
	for _, v := range context.LbAmfPool {
		if v.AmfID == amfID {
			return v, true
		}
	}
	return nil, false
}

func (context *LBContext) LbAmfFindByUeID(UeID int64) (*LbAmf, bool) {
	for _, Amf := range context.LbAmfPool {
		if check := Amf.ContainsUE(UeID); check {
			return Amf, true
		}
	}
	return nil, false
}

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

func (context *LBContext) AddGnbToLB(conn *sctp.SCTPConn) *LbGnb{
	gnb := NewLbGnb()
	gnb.LbConn.Conn = conn
	context.LbRanPool = append(context.LbRanPool, gnb)
	return gnb
}

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
