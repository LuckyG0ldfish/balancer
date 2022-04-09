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
	lB_Context = LB_Context{}
)

// The LB's main context 
type LB_Context struct {
	Name 						string
	NfId               			string

	LbIP 						string // IPv4 adress of the LB 

	LbToAmfPort					int // inital prot to connect to the first AMF 
	LbToAmfAddr					*sctp.SCTPAddr // SCTP adress generated for AMFs

	LbListenPort				int // port to listen for GNBs
	LbListenAddr				*sctp.SCTPAddr // SCTP adress generated for GNBs

	Running 					bool 	// true while the LB is not beeing terminated 

	NewAmf						bool // indicates that a new AMF IP+Port have been added so that the LB can connect to it 
	NewRegistAmfIpList 			[]string // Registration AMFs that have to be registred
	NewRegularAmfIpList 		[]string // Regular AMFs that have to be registred
	NewDeregistAmfIpList 		[]string // Deregistration AMFs that have to be registred
	
	DifferentAmfTypes			int // amount of different AMF types used 
	ContinuesAmfRegistration	bool // true for continues accepting AMFs for registration 
	
	LB_GNBPool 					sync.Map  // GNBs connected to the LB
	RegistAMFPool 				sync.Map  // AMFs connected to LB
	RegularAMFPool 				sync.Map  // AMFs connected to LB
	DeregistAMFPool 			sync.Map  // AMFs connected to LB

	Next_Regist_Amf 			*LB_AMF // AMF that is used for the next UE entering registration 
	Next_Regular_Amf 			*LB_AMF	// if 3 AMF types: AMF that is used for the next UE finishing registration
	Next_Deregist_Amf 			*LB_AMF	// if > 2 AMF types: AMF that is used for the next UE starting deregistration

	IDGen 						*UniqueNumberGen 
	
	RelativeCapacity 			int64 // To build setup response

	/* temp */
	PlmnSupportList 			[]factory.PlmnSupportItem
	ServedGuamiList 			[]models.Guami

	/* logger */
	Log 						*logrus.Entry

	/* metrics */
	MetricsLevel 				int 
	MetricsGNBs					sync.Map
}

// Creates and returns a new *LBContext
func NewLBContext() (*LB_Context){
	var new LB_Context
	return &new
}

// returns the current instance of the LB_Context
func LB_Self() *LB_Context {
	return &lB_Context
}