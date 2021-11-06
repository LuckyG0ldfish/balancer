package context 

import (
	// "github.com/ishidawataru/sctp"
	"github.com/sirupsen/logrus"
)

type LbUe struct{
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
	return &ue
}