package context

import (
	"sync"

	//"github.com/ishidawataru/sctp"
)

type LbGnb struct{
	GnbID 		int64
	LbConn 		*LBConn
	Ues 		sync.Map
}

func NewLbGnb(id int64) (gnb *LbGnb){
	gnb.GnbID = id
	gnb.LbConn = NewLBConn()
	gnb.LbConn.ID = id
	gnb.LbConn.TypeID = TypeIdentAMFConn
	return gnb
}

func (gnb *LbGnb) AddAMFUe(id int64) {
	gnb.Ues.Store(id, NewUE(id))
}

func (gnb *LbGnb) ContainsUE(id int64) (cont bool) {
	_, cont = gnb.Ues.Load(id)
	return 
}