package context 

import (
	"github.com/sirupsen/logrus"
)

const StateIdRegistration 		int	= 0
const StateIdRegular 			int	= 1
const StateIdDeregistration		int = 2

// Type, that stores all relevant information of UEs
type LbUe struct{
	UeStateIdent 	int			// Identifies the state of the UE 

	UeRanID 		int64		// ID given to the UE by GNB/RAN
	UeLbID 			int64		// ID given to the UE by LB
	UeAmfId 		int64		// ID given to the UE by AMF
	
	RanID			int64		// LB-internal ID of GNB that issued the UE 
	AmfID		 	int64		// LB-internal ID of AMF that processes the UE  

	/* logger */
	Log 			*logrus.Entry
}

// Creates, initializes and returns a new *LbUe
func NewUE() (*LbUe){
	var ue LbUe
	ue.UeStateIdent = StateIdRegistration
	return &ue
}