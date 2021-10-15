package context

import (
	"time"
	"sync"

	"github.com/ishidawataru/sctp"
)

const lbPPID uint32 = 0x3c000000

type LbAmf struct{
	AmfID 		int
	LbConn 		*LBConn
	Ues			sync.Map
}

func NewLbAmf(id int) (amf *LbAmf){
	amf.AmfID = id
	amf.LbConn = NewLBConn()
	amf.LbConn.ID = id
	amf.LbConn.TypeID = TypeIdentAMFConn
	return amf
}

func (amf *LbAmf) AddAMFUe(id int) {
	amf.Ues.Store(id, NewUE(id))
}

func (amf *LbAmf) ContainsUE(id int) (cont bool) {
	_, cont = amf.Ues.Load(id)
	return 
}

func (amf *LbAmf) Start(lbaddr sctp.SCTPAddr, amfIP string, amfPort int) {
	for{
		err := amf.ConnectToAmf(lbaddr, amfIP, amfPort)
		if err == nil {
			// amf.up = true
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func (amf *LbAmf) ConnectToAmf(lbaddr sctp.SCTPAddr, amfIP string, amfPort int) error{
	amfAddr, _ := GenSCTPAddr(amfIP, amfPort)
	conn, err := sctp.DialSCTP("sctp", &lbaddr, amfAddr)
	if err != nil {
		return  err
	}
	info, err := conn.GetDefaultSentParam()
	if err != nil {
		return  err
	}
	info.PPID = lbPPID
	err = conn.SetDefaultSentParam(info)
	if err != nil {
		return  err
	}
	//setting this connection as the amf SCTPConn
	amf.LbConn.Conn = conn
	return  nil
}