package context

import (
	"sync"

	"github.com/ishidawataru/sctp"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/sirupsen/logrus"
)

var nextGnbID int64 = 1

// Type, that stores all relevant information of connected GNBs 
type LB_GNB struct{
	GnbID 		int64			// INTERNAL ID for this GNB
	
	LB_Conn 	*LB_Conn		// stores all the connection related information 
	Ues 		sync.Map		// "list" of all UE that are processed by this GNB

	/* logger */
	Log 		*logrus.Entry
}

// use sctp.SCTPConn to find RAN context, return *AmfRan and true if found
func CreateAndAddNewGnbToLB(conn *sctp.SCTPConn) *LB_GNB{
	self := LB_Self()
	gnb := newLbGnb(conn)
	self.LB_GNBPool.Store(conn, gnb)
	return gnb
}

// Creates, initializes and returns a new *LbGnb
func newLbGnb(conn *sctp.SCTPConn) *LB_GNB{
	var gnb LB_GNB
	gnb.GnbID = nextGnbID
	gnb.LB_Conn = newLBConn(nextGnbID, TypeIdGNBConn)
	gnb.Log = logger.GNBLog
	gnb.LB_Conn.GnbPointer = &gnb
	gnb.LB_Conn.Conn = conn
	nextGnbID++
	return &gnb
}

// Use a UE-ID to find UE context, return *LbUe and true if found
func (gnb *LB_GNB) FindUeByRAN_UE_ID(id int64) (*LB_UE, bool){
	ue, ok := gnb.Ues.Load(id)
	if !ok {
		gnb.Log.Errorf("UE is not registered to this RAN")
		return nil, false 
	}
	ue2, ok := ue.(*LB_UE)
	if !ok {
		gnb.Log.Errorf("couldn't be converted")
		return nil, false 
	}
	return ue2, ok
}

// takes UeID and returns true if UE exists in the GNBs list 
func (gnb *LB_GNB) ContainsUE(id int64) (cont bool) {
	_, cont = gnb.Ues.Load(id)
	return 
}

// Removes GNB-Context and closes the Connection
func (gnb *LB_GNB) RemoveGnbContext() {
	lb := LB_Self()
	lb.LB_GNBPool.Delete(gnb.LB_Conn.Conn)
	gnb.LB_Conn.Conn.Close()
	gnb = nil 
}