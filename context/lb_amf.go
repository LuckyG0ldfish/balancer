package context

import (
	"sync"
	"github.com/sirupsen/logrus"
)

const TypeIdRegistAMF 		int	= 0
const TypeIdRegularAMF 		int	= 1
const TypeIdDeregistAMF		int = 2

var nextAmfID int64 = 1

// Type, that stores all relevant information of connected AMFs 
type LbAmf struct {
	AmfID  			int64 			// INTERNAL ID for this AMF 

	AmfTypeIdent 	int 			// Identifies the type of AMF 

	Capacity 		int64 			// AMFs Relative Cap. -> extracted out of NGSetup

	LbConn 			*LBConn 		// Stores all the connection related information 
	Ues    			sync.Map 		// "List" of all UE that are processed by this AMF 

	/* logger */
	Log 			*logrus.Entry
}

// Use a UE-ID to find UE context, return *LbUe and true if found
func (amf *LbAmf) FindUeByUeID(id int64) (*LbUe, bool){
	ue, ok := amf.Ues.Load(id)
	if !ok {
		amf.Log.Errorf("UE not found")
		return nil, false
	}
	ue2, ok :=  ue.(*LbUe)
	return ue2, ok
}

// Creates, initializes and returns a new *LbAmf
func NewLbAmf() *LbAmf {
	var amf LbAmf
	amf.AmfID = nextAmfID
	amf.LbConn = newLBConn(nextAmfID, TypeIdAMFConn)
	nextAmfID++
	return &amf
}

// takes UeID and returns true if UE exists in the AMFs list 
func (amf *LbAmf) ContainsUE(id int64) (cont bool) {
	_, cont = amf.Ues.Load(id)
	return
}



