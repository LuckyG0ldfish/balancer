package context

import (
	// "sync"

	"fmt"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/amf/factory" // TODO
	"github.com/free5gc/openapi/models"
)

var (
	lbContext = LBContext{}
)

// func init() {
// 	LB_Self().Name = "lb"
// 	//LB_Self().NetworkName.Full = "free5GC"
// }

type LBContext struct {
	Name string
	//NetworkName   factory.NetworkName
	LbIP 	string
	LbPort	int


	LbRanPool []*LbGnb // gNBs connected to the LB
	LbAmfPool []*LbAmf // amfs (each connected to AMF 1:1) connected to LB
	
	Next_Amf *LbAmf

	IDGen 	*UeIdGen
	
	RelativeCapacity 	int64 // To build setup response
	PlmnSupportList 	[]factory.PlmnSupportItem
	ServedGuamiList 	[]models.Guami
}

func NewLBContext() (LbContext *LBContext){
	return 
}

func (lb *LBContext) ForwardToNextAmf(lbConn *LBConn, message *ngapType.NGAPPDU, ue *LbUe) { 
	// if mes, err := ngap.Encoder(*message); err == nil {
	// }
	fmt.Println("forward to nextAMF")

	ue.AmfID = lb.Next_Amf.AmfID
	_, ok := lb.Next_Amf.Ues.Load(ue.UeLbID)
	if ok {
		fmt.Println("UE already exists")
		return 
	} 
	lb.Next_Amf.Ues.Store(ue.UeLbID, ue)
	var mes []byte
	mes, _  = ngap.Encoder(*message)
	lb.Next_Amf.LbConn.Conn.Write(mes)
	fmt.Println(mes)
}

func (lb *LBContext) ForwardToAmf(lbConn *LBConn, message *ngapType.NGAPPDU, ue *LbUe) {
	amf, ok := lb.LbGnbFindByID(ue.AmfID)
	// if mes, err := ngap.Encoder(*message); err == nil {	
	// }
	if ok {
		fmt.Println("forward to AMF:")
		var mes []byte
		mes, _  = ngap.Encoder(*message)
		amf.LbConn.Conn.Write(mes)
		fmt.Println(mes)
		fmt.Println(ue.UeLbID)
	} else {
		fmt.Println("AMF not found")
	}
}

func (lb *LBContext) ForwardToGnb(lbConn *LBConn, message *ngapType.NGAPPDU, ue *LbUe) { //*ngapType.NGAPPDU
	gnb, ok := lb.LbGnbFindByID(ue.RanID)
	// if mes, err := ngap.Encoder(*message); err == nil {	
	// }
	if ok {
		fmt.Println("forward to GNB:")
		var mes []byte
		mes, _  = ngap.Encoder(*message)
		gnb.LbConn.Conn.Write(mes)
		fmt.Println(mes)
		fmt.Println(ue.UeLbID)
	} else {
		fmt.Println("GNB not found")
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

func LB_Self() *LBContext {
	return &lbContext
}
