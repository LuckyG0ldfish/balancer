package context

import (
	// "fmt"
	"sync"
	// "time"

	// "git.cs.nctu.edu.tw/calee/sctp"

	
)

const lbPPID uint32 = 0x3c000000

var nextAmfID int64 = 1

type LbAmf struct {
	AmfID  int64
	LbConn *LBConn
	Ues    sync.Map
}

func (amf *LbAmf) FindUeByUeRanID(id int64) ([]*LbUe, bool){
	//var ue LbUe
	ue, _ := amf.Ues.Load(id)
	ue2, ok :=  ue.([]*LbUe)
	return ue2, ok
}

func NewLbAmf() *LbAmf {
	var amf LbAmf
	amf.AmfID = nextAmfID
	amf.LbConn = NewLBConn()
	amf.LbConn.ID = nextAmfID
	amf.LbConn.TypeID = TypeIdentAMFConn
	nextAmfID++
	return &amf
}

func (amf *LbAmf) ContainsUE(id int64) (cont bool) {
	_, cont = amf.Ues.Load(id)
	return
}



