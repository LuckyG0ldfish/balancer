package context

import (
	"sync"

	//"github.com/ishidawataru/sctp"
	"github.com/sirupsen/logrus"
)

var nextGnbID int64 = 1

type LbGnb struct{
	GnbID 		int64
	LbConn 		*LBConn
	Ues 		sync.Map

	/* logger */
	Log 		*logrus.Entry
}

func NewLbGnb() (*LbGnb){
	var gnb LbGnb
	gnb.GnbID = nextGnbID
	gnb.LbConn = NewLBConn()
	gnb.LbConn.ID = nextGnbID
	gnb.LbConn.TypeID = TypeIdGNBConn
	nextGnbID++
	return &gnb
}

func (gnb *LbGnb) FindUeByUeRanID(id int64) (*LbUe, bool){
	//var ue LbUe
	ue, _ := gnb.Ues.Load(id)
	ue2, ok :=  ue.(*LbUe)
	return ue2, ok
}

func (gnb *LbGnb) ContainsUE(id int64) (cont bool) {
	_, cont = gnb.Ues.Load(id)
	return 
}