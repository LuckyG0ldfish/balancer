package context

import (
	"sync"

	"git.cs.nctu.edu.tw/calee/sctp"
)

type LbGnb struct{
	GnbID 		int
	LbConn 		*LBConn
	Ues 		sync.Map
}

func NewLbGnb(id int) (gnb *LbGnb){
	gnb.GnbID = id
	gnb.LbConn = NewLBConn()
	gnb.LbConn.ID = id
	gnb.LbConn.TypeID = TypeIdentAMFConn
	return gnb
}

func (gnb *LbGnb) AddAMFUe(id int) {
	gnb.Ues.Store(id, NewUE(id))
}

func (gnb *LbGnb) ContainsUE(id int) (cont bool) {
	_, cont = gnb.Ues.Load(id)
	return 
}