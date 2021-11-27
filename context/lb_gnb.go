package context

import (
	"sync"

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

// creates, initializes and returns a new *LbGnb
func NewLbGnb() (*LbGnb){
	var gnb LbGnb
	gnb.GnbID = nextGnbID
	gnb.LbConn = newLBConn(nextGnbID, TypeIdGNBConn)
	nextGnbID++
	return &gnb
}

func (gnb *LbGnb) FindUeByUeRanID(id int64) (*LbUe, bool){
	ue, _ := gnb.Ues.Load(id)
	ue2, ok :=  ue.(*LbUe)
	return ue2, ok
}

// takes UeID and returns true if UE exists in the GNBs list 
func (gnb *LbGnb) ContainsUE(id int64) (cont bool) {
	_, cont = gnb.Ues.Load(id)
	return 
}