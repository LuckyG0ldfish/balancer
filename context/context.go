package context

import (
	"sync"

	"github.com/ishidawataru/sctp"
	"github.com/sirupsen/logrus"
	
	// "github.com/LuckyG0ldfish/balancer/logger"
	"github.com/LuckyG0ldfish/balancer/factory"
	
	"github.com/free5gc/openapi/models"
	// "github.com/free5gc/ngap/ngapType"
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

	Running 			bool 	// true while the LB is not beeing terminated 

	NewAmf				bool // indicates that a new AMF IP+Port have been added so that the LB can connect to it 
	NewAmfIpList 		[]string 
	
	LbRanPool 			sync.Map //[]*LbGnb // gNBs connected to the LB
	LbAmfPool 			sync.Map //[]*LbAmf // amfs (each connected to AMF 1:1) connected to LB

	Next_Regist_Amf 			*LbAmf
	Next_Regular_Amf 			*LbAmf
	Next_Deregist_Amf 			*LbAmf

	IDGen 				*UniqueNumberGen
	
	RelativeCapacity 	int64 // To build setup response

	/* temp */
	PlmnSupportList 	[]factory.PlmnSupportItem
	ServedGuamiList 	[]models.Guami

	/* logger */
	Log 	*logrus.Entry
}

// Creates and returns a new *LBContext
func NewLBContext() (*LBContext){
	var new LBContext
	return &new
}

// use sctp.SCTPConn to find RAN context, return *LbRan and true if found
func (context *LBContext) LbGnbFindByConn(conn *sctp.SCTPConn) (*LbGnb, bool) {
	gnbTemp, ok := context.LbRanPool.Load(conn)
	if !ok {
		return nil, false
	}
	gnb, ok := gnbTemp.(*LbGnb)
	if !ok {
		return nil, false
	}
	return gnb, ok
}

// use sctp.SCTPConn to find Amf context, return *LbAmf and true if found
func (context *LBContext) LbAmfFindByConn(conn *sctp.SCTPConn) (*LbAmf, bool) {
	amfTemp, ok := context.LbAmfPool.Load(conn)
	if !ok {
		return nil, false
	}
	amf, ok := amfTemp.(*LbAmf)
	if !ok {
		return nil, false
	}
	return amf, ok
}

func LB_Self() *LBContext {
	return &lbContext
}
