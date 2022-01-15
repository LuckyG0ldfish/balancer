package context

import (
	"sync"

	"github.com/ishidawataru/sctp"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/sirupsen/logrus"
)

var nextGnbID int64 = 1

// Type, that stores all relevant information of connected GNBs 
type LbGnb struct{
	GnbID 		int64			// INTERNAL ID for this GNB
	
	LbConn 		*LBConn			// stores all the connection related information 
	Ues 		sync.Map		// "list" of all UE that are processed by this GNB

	/* logger */
	Log 		*logrus.Entry
}

// use sctp.SCTPConn to find RAN context, return *AmfRan and true if found
func CreateAndAddNewGnbToLB(conn *sctp.SCTPConn) *LbGnb{
	self := LB_Self()
	gnb := newLbGnb(conn)
	self.LbRanPool.Store(conn, gnb)
	return gnb
}

// Creates, initializes and returns a new *LbGnb
func newLbGnb(conn *sctp.SCTPConn) *LbGnb{
	var gnb LbGnb
	gnb.GnbID = nextGnbID
	gnb.LbConn = newLBConn(nextGnbID, TypeIdGNBConn)
	gnb.Log = logger.GNBLog
	gnb.LbConn.RanPointer = &gnb
	gnb.LbConn.Conn = conn
	nextGnbID++
	return &gnb
}

// Use a UE-ID to find UE context, return *LbUe and true if found
func (gnb *LbGnb) FindUeByUeRanID(id int64) (*LbUe, bool){
	ue, ok := gnb.Ues.Load(id)
	if !ok {
		gnb.Log.Errorf("UE is not registered to this RAN")
		return nil, false 
	}
	ue2, ok :=  ue.(*LbUe)
	if !ok {
		gnb.Log.Errorf("couldn't be converted")
		return nil, false 
	}
	return ue2, ok
}

// takes UeID and returns true if UE exists in the GNBs list 
func (gnb *LbGnb) ContainsUE(id int64) (cont bool) {
	_, cont = gnb.Ues.Load(id)
	return 
}

// Removes GNB-Context and closes the Connection
func (gnb *LbGnb) RemoveGnbContext() {
	lb := LB_Self()
	lb.LbRanPool.Delete(gnb.LbConn.Conn)
	gnb.LbConn.Conn.Close()
	gnb = nil 
}