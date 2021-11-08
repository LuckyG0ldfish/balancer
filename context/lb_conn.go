package context

import(
	"git.cs.nctu.edu.tw/calee/sctp"
)

const TypeIdNotThere 	int	= 0
const TypeIdAMFConn 	int	= 1
const TypeIdGNBConn		int = 2

type LBConn struct{
	TypeID 		int 
	ID 			int64 
	Conn 		*sctp.SCTPConn
}

func NewLBConn() (*LBConn){
	var lbConn LBConn
	lbConn.TypeID = 0 
	lbConn.ID = 0 
	return &lbConn
}
