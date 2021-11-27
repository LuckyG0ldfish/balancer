package context

import (
	"sync"
	
	"github.com/sirupsen/logrus"
)

const TypeIdRegistAMF 		int	= 0
const TypeIdRegularAMF 		int	= 1
const TypeIdDeregistAMF		int = 2

var nextAmfID int64 = 1

type LbAmf struct {
	AmfID  			int64

	AmfTypeIdent 	int

	Capacity 		int64

	LbConn 			*LBConn
	Ues    			sync.Map

	/* logger */
	Log 			*logrus.Entry
}

func (amf *LbAmf) FindUeByUeRanID(id int64) (*LbUe, bool){
	ue, _ := amf.Ues.Load(id)
	ue2, ok :=  ue.(*LbUe)
	return ue2, ok
}

func NewLbAmf() *LbAmf {
	var amf LbAmf
	amf.AmfID = nextAmfID
	amf.LbConn = NewLBConn()
	amf.LbConn.ID = nextAmfID
	amf.LbConn.TypeID = TypeIdAMFConn
	nextAmfID++
	return &amf
}

func (amf *LbAmf) ContainsUE(id int64) (cont bool) {
	_, cont = amf.Ues.Load(id)
	return
}



