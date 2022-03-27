package context

import (
	"sync"

	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/sirupsen/logrus"
)

var nextAmfID int64 = 1

// Type, that stores all relevant information of connected AMFs 
type LbAmf struct {
	AmfID  					int64 			// INTERNAL ID for this AMF 

	AmfTypeIdent 			int 			// Identifies the type of AMF 

	LbConn 					*LBConn 		// Stores all the connection related information 

	RelativeCapacity 		int64 			// AMFs Relative Cap. -> extracted out of NGSetup
	NumberOfConnectedUEs 	int64 			// Amount of UEs that are connected to this AMF 
	
	Ues    					sync.Map 		// "List" of all UE that are processed by this AMF 

	/* logger */
	Log 					*logrus.Entry
}

// Use a UELB-ID to find UE context, return *LbUe and true if found
func (amf *LbAmf) FindUeByUeID(id int64) (*LbUe, bool){
	ue, ok := amf.Ues.Load(id)
	if !ok {
		amf.Log.Errorf("UE is not registered to this AMF %d", id)
		return nil, false 
	}
	ue2, ok :=  ue.(*LbUe)
	if !ok {
		amf.Log.Errorf("couldn't be converted")
		return nil, false 
	}
	return ue2, ok
}

// Use a UEAMF-ID to find UE context, return *LbUe and true if found TODO
func (amf *LbAmf) FindUeByUeAmfID(id int64) (*LbUe, bool){
	var ue *LbUe
	var ok bool = false 
	amf.Ues.Range(func(key, value interface{}) bool{
		ueTemp, okTemp := value.(*LbUe)
		if !okTemp {
			logger.NgapLog.Errorf("couldn't be converted")
		}
		if ueTemp.UeAmfID == id{
			ue = ueTemp
			ok = true 
		}
		return true
	})
	return ue, ok
}

// 
func CreateAndAddAmfToLB(amfType int) *LbAmf{
	self := LB_Self()
	amf := newLbAmf(amfType)
	self.LbRegistAmfPool.Store(amf.LbConn.Conn, amf)
	// self.Table.addAmfCounter(amf)
	return amf
}

// Creates, initializes and returns a new *LbAmf
func newLbAmf(amfType int) *LbAmf {
	var amf LbAmf
	amf.AmfID = nextAmfID
	amf.LbConn = newLBConn(nextAmfID, TypeIdAMFConn)
	amf.LbConn.AmfPointer = &amf
	amf.Log = logger.AMFLog
	amf.RelativeCapacity = 0 
	amf.NumberOfConnectedUEs = 0 
	amf.AmfTypeIdent = amfType
	nextAmfID++
	return &amf
}

// takes UeID and returns true if UE exists in the AMFs list 
func (amf *LbAmf) ContainsUE(id int64) (cont bool) {
	_, cont = amf.Ues.Load(id)
	return
}

// calculates a number reflecting the AMF-Usage that is comparable within the loadbalancer
func (amf *LbAmf) calculateAMFUsage() float64{
	return float64(amf.NumberOfConnectedUEs) / float64(amf.RelativeCapacity)
}

// Removes AMF-Context and closes the Connection  
func (amf *LbAmf) RemoveAmfContext() {
	lb := LB_Self()
	lb.LbRegistAmfPool.Delete(amf.LbConn.Conn)
	amf.LbConn.Conn.Close()
	amf = nil 
}


