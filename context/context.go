package context

import (
	// "sync"
	"git.cs.nctu.edu.tw/calee/sctp"
)

var (
	lbContext 					= LBContext{}
)

func init() {
	LB_Self().Name = "lb"
	//LB_Self().NetworkName.Full = "free5GC"
}

type LBContext struct{
	Name 			string 
	//NetworkName   factory.NetworkName

	LbRanPool		[]*LbGnb 		// gNBs connected to the LB 
	LbAmfPool		[]*LbAmf			// amfs (each connected to AMF 1:1) connected to LB
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

func (context *LBContext) LbGnbFindByID(ranID int) (*LbGnb, bool) {
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

func (context *LBContext) LbAmfFindByID(amfID int) (*LbAmf, bool) {
	for _, v := range context.LbAmfPool {
		if v.AmfID == amfID {
			return v, true 
		}
	} 
	return nil, false
}

func (context *LBContext) LbAmfFindByUeID(UeID int) *LbAmf{
	for _, Amf := range context.LbAmfPool {
		if check := Amf.ContainsUE(UeID); check {
			return Amf
		}
	}
	return nil 
}

func (context *LBContext) LbGnbFindByUeID(UeID int) *LbGnb{
	for _, Gnb := range context.LbRanPool {
		if check := Gnb.ContainsUE(UeID); check {
			return Gnb
		}
	}
	return nil 
}

func LB_Self() *LBContext {
	return &lbContext
}
