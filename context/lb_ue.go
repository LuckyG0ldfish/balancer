package context 

import (
	"github.com/sirupsen/logrus"
)

// Type, that stores all relevant information of UEs
type LbUe struct{
	UeStateIdent 	int			// Identifies the state of the UE 

	UeRanID 		int64		// ID given to the UE by GNB/RAN
	UeLbID 			int64		// ID given to the UE by LB
	UeAmfId 		int64		// ID given to the UE by AMF -> set but unused atm 
	
	RanID			int64		// LB-internal ID of GNB that issued the UE 
	RanPointer 		*LbGnb

	AmfID		 	int64		// LB-internal ID of AMF that processes the UE  
	AmfPointer		*LbAmf

	/* logger */
	Log 			*logrus.Entry
}

// Creates, initializes and returns a new *LbUe
func NewUE() (*LbUe){
	var ue LbUe
	ue.UeStateIdent = TypeIdRegist
	return &ue
}

// Removes LbUe from AMF and RAN Context withing LB  
func (ue *LbUe) RemoveUeEntirely() {
	ue.RemoveUeFromAMF()
	ue.RemoveUeFromGNB()
}

// Removes LbUe from AMF Context withing LB 
func (ue *LbUe) RemoveUeFromAMF() {
	if ue.AmfPointer != nil {
		ue.AmfPointer.Ues.Delete(ue.UeLbID) // sync.Map key here is the LB internal UE-ID 
		ue.AmfPointer.Log.Traceln("UE context removed from AMF")
		ue.AmfPointer = nil 
		ue.AmfID = 0 
	}
}

// Removes LbUe from RAN Context withing LB 
func (ue *LbUe) RemoveUeFromGNB() {
	if ue.RanPointer != nil {
		ue.RanPointer.Ues.Delete(ue.UeRanID) // sync.Map key here is the RAN UE-ID
		ue.RanPointer.Log.Traceln("UE context removed from GNB")
		ue.RanPointer = nil 
		ue.RanID = 0 
	}
}