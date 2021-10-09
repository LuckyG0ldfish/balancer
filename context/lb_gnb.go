package context

import (
	"github.com/ishidawataru/sctp"
)

type LbGnb struct{
	gnbID 		int 
	Conn 		*sctp.SCTPConn
}

