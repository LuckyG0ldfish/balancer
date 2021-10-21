package context

import (
	// "sync"

	"fmt"

	"git.cs.nctu.edu.tw/calee/sctp"
	// "github.com/free5gc/ngap"
	// "github.com/free5gc/ngap/ngapType"
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
	
}

func NewLBContext() (LbContext *LBContext){
	return 
}

func (lb *LBContext) ForwardToNextAmf(lbConn *LBConn, message []byte, ue *LbUe) {
	// if mes, err := ngap.Encoder(*message); err == nil {
		
	// }
	fmt.Println("forward to nextAMF")
	fmt.Println(message)
	ue.AmfID = lb.Next_Amf.AmfID
	lb.Next_Amf.Ues.Store(ue.UeRanID, ue)
	lb.Next_Amf.LbConn.Conn.Write(message)
}

func (lb *LBContext) ForwardToAmf(lbConn *LBConn, message []byte, aMFUENGAPID int64) {
	amf, ok := lb.LbAmfFindByUeID(aMFUENGAPID)
	// if mes, err := ngap.Encoder(*message); err == nil {
		
	// }
	if ok {
		fmt.Println("forward to AMF:")
		fmt.Println(message)
		amf.LbConn.Conn.Write(message)
	} else {
		fmt.Println("Amf not found")
	}
}

func (lb *LBContext) ForwardToGnb(lbConn *LBConn, message []byte, rANUENGAPID int64) { //*ngapType.NGAPPDU
	gnb, ok := lb.LbGnbFindByUeID(rANUENGAPID)
	// if mes, err := ngap.Encoder(*message); err == nil {
		
	// }
	if ok {
		fmt.Println("forward to GNB:")
		fmt.Println(message)
		gnb.LbConn.Conn.Write(message)
	} else {
		fmt.Println("Gnb not found")
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

func (context *LBContext) LbGnbFindByUeID(UeID int64) (*LbGnb, bool) {
	for _, Gnb := range context.LbRanPool {
		if check := Gnb.ContainsUE(UeID); check {
			return Gnb, true
		}
	}
	return nil, false
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
