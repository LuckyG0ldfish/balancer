package context 

import (
	"github.com/sirupsen/logrus"
)

const StateIdRegistration 		int	= 0
const StateIdRegular 			int	= 1
const StateIdDeregistration		int = 2

type LbUe struct{
	UeStateIdent 	int

	UeRanID 		int64
	UeLbID 			int64
	UeAmfId 		int64
	
	RanID			int64
	AmfID		 	int64

	/* logger */
	Log 			*logrus.Entry
}

func NewUE() (*LbUe){
	var ue LbUe
	ue.UeStateIdent = StateIdRegistration
	return &ue
}