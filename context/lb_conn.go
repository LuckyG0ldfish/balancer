package context

import(
	"git.cs.nctu.edu.tw/calee/sctp"
)

const TypeIdNotThere 	int	= 0
const TypeIdAMFConn 	int	= 1
const TypeIdGNBConn		int = 2

// universal type for all connections of the LB 
type LBConn struct{
	ID 			int64 				// internal AMF/GNB ID that is connected with Conn 
	TypeID 		int 				// type identifier of the connected AMF/GNB 
	Conn 		*sctp.SCTPConn		// actual connection to AMF/GNB 
}

// creates, initializes and returns a new *LbConn
// (only used when initializing a AMF/GNB)
func newLBConn(id int64, typeID int) (*LBConn){
	var lbConn LBConn
	lbConn.ID = id
	lbConn.TypeID = typeID 
	return &lbConn
}
