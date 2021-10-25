package context 

import (
	// "github.com/ishidawataru/sctp"
)

type LbUe struct{
	UeRanID 		int64
	UeLbID 			int64
	UeAmfId 		int64
	RanID			int64
	AmfID		 	int64
}

func NewUE() (*LbUe){
	var ue LbUe
	return &ue
}