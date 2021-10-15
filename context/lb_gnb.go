package context

import (
	"sync"

	//"github.com/ishidawataru/sctp"
)

var nextGnbID int64 = 1

type LbGnb struct{
	GnbID 		int64
	LbConn 		*LBConn
	Ues 		sync.Map
}

func NewLbGnb() (gnb *LbGnb){
	gnb.GnbID = nextGnbID
	gnb.LbConn = NewLBConn()
	gnb.LbConn.ID = nextGnbID
	gnb.LbConn.TypeID = TypeIdentGNBConn
	nextGnbID++
	return gnb
}

func (gnb *LbGnb) AddAMFUe(id int64) {
	gnb.Ues.Store(id, NewUE(id))
}

func (gnb *LbGnb) ContainsUE(id int64) (cont bool) {
	_, cont = gnb.Ues.Load(id)
	return 
}