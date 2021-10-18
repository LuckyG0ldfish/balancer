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

func NewLbAmf() *LbAmf {
	var amf LbAmf
	amf.AmfID = nextAmfID
	amf.LbConn = NewLBConn()
	amf.LbConn.ID = nextAmfID
	amf.LbConn.TypeID = TypeIdentAMFConn
	nextAmfID++
	return &amf
}

func (amf *LbAmf) AddAMFUe(id int64) {
	amf.Ues.Store(id, NewUE(id))
}

func (amf *LbAmf) ContainsUE(id int64) (cont bool) {
	_, cont = amf.Ues.Load(id)
	return
}



