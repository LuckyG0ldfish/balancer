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
	lb_Context = Lb_Context{}
)

type Lb_Context struct {
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

	Next_Regist_Amf 			*Lb_Amf
	Next_Regular_Amf 			*Lb_Amf
	Next_Deregist_Amf 			*Lb_Amf

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
func NewLBContext() (*Lb_Context){
	var new Lb_Context
	return &new
}

func LB_Self() *Lb_Context {
	return &lb_Context
}