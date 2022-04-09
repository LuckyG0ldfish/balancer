package context

import (
	"sync"

	"github.com/ishidawataru/sctp"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/sirupsen/logrus"
)

var nextGnbID int64 = 1

// Type, that stores all relevant information of connected GNBs 
type Lb_Gnb struct{
	GnbID 		int64			// INTERNAL ID for this GNB
	
	Lb_Conn 	*Lb_Conn		// stores all the connection related information 
	Ues 		sync.Map		// "list" of all UE that are processed by this GNB

	/* logger */
	Log 		*logrus.Entry
}

// use sctp.SCTPConn to find RAN context, return *AmfRan and true if found
func CreateAndAddNewGnbToLB(conn *sctp.SCTPConn) *Lb_Gnb{
	self := LB_Self()
	gnb := newLbGnb(conn)
	self.LbRanPool.Store(conn, gnb)
	return gnb
}

// Creates, initializes and returns a new *LbGnb
func newLbGnb(conn *sctp.SCTPConn) *Lb_Gnb{
	var gnb Lb_Gnb
	gnb.GnbID = nextGnbID
	gnb.Lb_Conn = newLBConn(nextGnbID, TypeIdGNBConn)
	gnb.Log = logger.GNBLog
	gnb.Lb_Conn.GnbPointer = &gnb
	gnb.Lb_Conn.Conn = conn
	nextGnbID++
	return &gnb
}

// Use a UE-ID to find UE context, return *LbUe and true if found
func (gnb *Lb_Gnb) FindUeByRAN_UE_ID(id int64) (*LbUe, bool){
	ue, ok := gnb.Ues.Load(id)
	if !ok {
		gnb.Log.Errorf("UE is not registered to this RAN")
		return nil, false 
	}
	ue2, ok := ue.(*LbUe)
	if !ok {
		gnb.Log.Errorf("couldn't be converted")
		return nil, false 
	}
	return ue2, ok
}

// takes UeID and returns true if UE exists in the GNBs list 
func (gnb *Lb_Gnb) ContainsUE(id int64) (cont bool) {
	_, cont = gnb.Ues.Load(id)
	return 
}

// Removes GNB-Context and closes the Connection
func (gnb *Lb_Gnb) RemoveGnbContext() {
	lb := LB_Self()
	lb.LbRanPool.Delete(gnb.Lb_Conn.Conn)
	gnb.Lb_Conn.Conn.Close()
	gnb = nil 
}