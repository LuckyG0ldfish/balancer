package context

import (
	"sync"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/ngap/ngapType"

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

	Running 			bool 	// true while the LB is not beeing terminated 

	NewAmf				bool // indicates that a new AMF IP+Port have been added so that the LB can connect to it 
	NewAmfIpList 		[]string 
	NewAmfPortList		[]string 
	
	LbRanPool 			sync.Map //[]*LbGnb // gNBs connected to the LB
	LbAmfPool 			sync.Map //[]*LbAmf // amfs (each connected to AMF 1:1) connected to LB

	Next_Amf 			*LbAmf

	IDGen 				*UniqueNumberGen
	
	RelativeCapacity 	int64 // To build setup response

	/* temp */
	NGSetupRes 			*ngapType.NGAPPDU
	PlmnSupportList 	*ngapType.PLMNSupportList
	ServedGuamiList 	*ngapType.ServedGUAMIList

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

// TODO: 
func (context *LBContext) SelectNextAmf() bool{
	
	var amfWithMaxCap *LbAmf = context.Next_Amf
	var amfUsage float64 = context.Next_Amf.calculateAMFUsage()
	
	context.LbAmfPool.Range(func(key, value interface{}) bool{
		amfTemp, ok := value.(*LbAmf)
		if !ok {
			logger.NgapLog.Errorf("couldn't be converted")
		}
		tempUsage := amfTemp.calculateAMFUsage()
		
		// chooses the AMF with the lowest Usage
		if  amfUsage > tempUsage {
			amfWithMaxCap = amfTemp
			amfUsage = tempUsage
		} 
		return true
	})
	if amfWithMaxCap == nil {
		logger.ContextLog.Errorf("No Amf found")
		return true 
	}
	context.Next_Amf = amfWithMaxCap
	logger.ContextLog.Tracef("NextAMF = AMFID: %d", amfWithMaxCap.AmfID)
	return true 
}



func LB_Self() *LBContext {
	return &lbContext
}
