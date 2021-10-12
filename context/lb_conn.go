package context

import(
	"github.com/ishidawataru/sctp"
)

const TypeIndetNotThere 	int	= 0
const TypeIdentAMFConn 		int	= 1
const TypeIdentGNBConn		int = 2

type LBConn struct{
	TypeID 		int 
	ID 			int 
	Conn 		*sctp.SCTPConn
}

func NewLBConn() (lbConn *LBConn){
	lbConn.TypeID = 0 
	lbConn.ID = 0 
	lbConn.Conn = nil 
	return 
}