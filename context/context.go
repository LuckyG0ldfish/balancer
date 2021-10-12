package context

import (
	"sync"
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

	LbRanPool		sync.Map 		// gNBs connected to the LB 
	LbAmfPool		sync.Map		// amfs (each connected to AMF 1:1) connected to LB
}

// use net.Conn to find RAN context, return *AmfRan and ok bit
func (context *LBContext) LbGnbFindByConn(conn sctp.SCTPConn) (*LbGnb, bool) {
	if value, ok := context.LbRanPool.Load(conn); ok {
		return value.(*LbGnb), ok
	}
	return nil, false
}

func (context *LBContext) LbGnbFindByID(ranID int) (*LbGnb, bool) {
	if value, ok := context.LbRanPool.Load(ranID); ok {
		return value.(*LbGnb), ok
	}
	return nil, false
}

func (context *LBContext) LbAmfFindByConn(conn sctp.SCTPConn) (*LbGnb, bool) {
	if value, ok := context.LbRanPool.Load(conn); ok {
		return value.(*LbGnb), ok
	}
	return nil, false
}

func (context *LBContext) LbAmfFindByID(amfID int) (*LbGnb, bool) {
	if value, ok := context.LbRanPool.Load(amfID); ok {
		return value.(*LbGnb), ok
	}
	return nil, false
}



func LB_Self() *LBContext {
	return &lbContext
}
