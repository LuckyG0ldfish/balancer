package context

import (
	"github.com/ishidawataru/sctp"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/sirupsen/logrus"
)

const TypeIdNotThere 	int	= 0
const TypeIdAMFConn 	int	= 1
const TypeIdGNBConn		int = 2

// universal type for all connections of the LB 
type Lb_Conn struct{
	ID 			int64 				// internal AMF/GNB ID that is connected with Conn 
	TypeID 		int 				// type identifier of the connected AMF/GNB 
	
	Conn 		*sctp.SCTPConn		// actual connection to AMF/GNB 
	Closed 		bool 				// determines whether the connection is closed 

	GnbPointer 	*Lb_Gnb
	AmfPointer 	*Lb_Amf
	/* logger */
	Log 			*logrus.Entry
}

// creates, initializes and returns a new *LbConn
// (only used when initializing a AMF/GNB)
func newLBConn(id int64, typeID int) (*Lb_Conn){
	var lbConn Lb_Conn
	lbConn.ID = id
	lbConn.TypeID = typeID 
	lbConn.Log = logger.LbConnLog
	lbConn.Closed = false 
	return &lbConn
}
