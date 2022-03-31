package context

import (
	"sync"

	"github.com/ishidawataru/sctp"
	"github.com/sirupsen/logrus"
	
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

	NewAmf						bool // indicates that a new AMF IP+Port have been added so that the LB can connect to it 
	NewRegistAmfIpList 			[]string 
	NewRegularAmfIpList 		[]string 
	NewDeregistAmfIpList 		[]string 
	
	DifferentAmfTypes			int 
	ContinuesAmfRegistration	bool // true for continues accepting AMFs for registration 
	
	LbRanPool 			sync.Map //[]*LbGnb // gNBs connected to the LB
	LbRegistAmfPool 	sync.Map //[]*LbAmf // amfs (each connected to AMF 1:1) connected to LB
	LbRegularAmfPool 	sync.Map //[]*LbAmf // amfs (each connected to AMF 1:1) connected to LB
	LbDeregistAmfPool 	sync.Map //[]*LbAmf // amfs (each connected to AMF 1:1) connected to LB

	Next_Regist_Amf 			*LbAmf
	Next_Regular_Amf 			*LbAmf
	Next_Deregist_Amf 			*LbAmf

	IDGen 				*UniqueNumberGen
	
	RelativeCapacity 	int64 // To build setup response

	/* temp */
	PlmnSupportList 	[]factory.PlmnSupportItem
	ServedGuamiList 	[]models.Guami

	/* logger */
	Log 				*logrus.Entry

	/* metrics */
	MetricsLevel 		int 
	MetricsGNBs			sync.Map
}

// Creates and returns a new *LBContext
func NewLBContext() (*LBContext){
	var new LBContext
	return &new
}

func LB_Self() *LBContext {
	return &lbContext
}