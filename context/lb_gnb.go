package context

import (
	"github.com/ishidawataru/sctp"
)

type LbGnb struct{
	gnbID 		int 
	Conn 		*sctp.SCTPConn
}

func NewLbGnb(id int) (amf *LbAmf){
	amf.amfID = id
	amf.Conn = nil 
	return amf
}