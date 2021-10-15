package context 

import (
	// "github.com/ishidawataru/sctp"
)

type LbUe struct{
	UeID 		int
}

func NewUE(id int) (ue *LbUe){
	ue.UeID = id
	return 
}