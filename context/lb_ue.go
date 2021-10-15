package context 

import (
	// "github.com/ishidawataru/sctp"
)

type LbUe struct{
	UeID 		int64
}

func NewUE(id int64) (ue *LbUe){
	ue.UeID = id
	return 
}