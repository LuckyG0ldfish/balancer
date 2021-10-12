package context

import (
	"github.com/ishidawataru/sctp"
)

type LbGnb struct{
	LbConn 		*LBConn
}

func NewLbGnb(id int) (gnb *LbGnb){
	gnb.LbConn = NewLBConn()
	gnb.LbConn.ID = id
	gnb.LbConn.TypeID = TypeIdentAMFConn
	return gnb
}