package context

import (
	"time"
	"sync"
	"fmt"

	"git.cs.nctu.edu.tw/calee/sctp"
)

const lbPPID uint32 = 0x3c000000

var nextAmfID int64 = 1

type LbAmf struct{
	AmfID 		int64
	LbConn 		*LBConn
	Ues			sync.Map
}

func NewLbAmf() ( *LbAmf){
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

func (amf *LbAmf) Start(lbaddr *sctp.SCTPAddr, amfIP string, amfPort int) {
	for{
		fmt.Println("connecting")
		err := amf.ConnectToAmf(lbaddr, amfIP, amfPort)
		if err == nil {
			// amf.up = true
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func (amf *LbAmf) ConnectToAmf(lbaddr *sctp.SCTPAddr, amfIP string, amfPort int) error{
	amfAddr, _ := GenSCTPAddr(amfIP, amfPort)
	conn, err := sctp.DialSCTP("sctp", lbaddr, amfAddr)
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
	fmt.Println("connected")
	return  nil
}