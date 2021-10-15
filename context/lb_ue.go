package context 

import (
	// "git.cs.nctu.edu.tw/calee/sctp"
)

type LbUe struct{
	UeID 		int
}

func NewUE(id int) (ue *LbUe){
	ue.UeID = id
	return 
}