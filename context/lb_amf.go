package context

import (
	"sync"

	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/sirupsen/logrus"
)

var nextAmfID int64 = 1

// Type, that stores all relevant information of connected AMFs 
type Lb_Amf struct {
	AmfID  					int64 			// INTERNAL ID for this AMF 

	AmfTypeIdent 			int 			// Identifies the type of AMF 

	Lb_Conn 				*Lb_Conn 		// Stores all the connection related information 

	RelativeCapacity 		int64 			// AMFs Relative Cap. -> extracted out of NGSetup
	NumberOfConnectedUEs 	int64 			// Amount of UEs that are connected to this AMF 
	
	Ues    					sync.Map 		// "List" of all UE that are processed by this AMF 

	/* logger */
	Log 					*logrus.Entry
}

// Use a UELB-ID to find UE context, return *LbUe and true if found
func (amf *Lb_Amf) FindUeByRAN_UE_ID(id int64) (*LbUe, bool){
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
func (amf *Lb_Amf) FindUeByUeAmfID(id int64) (*LbUe, bool){
	var ue *LbUe
	var ok bool = false 
	amf.Ues.Range(func(key, value interface{}) bool{
		ueTemp, okTemp := value.(*LbUe)
		if !okTemp {
			logger.NgapLog.Errorf("couldn't be converted")
		}
		if ueTemp.AMF_UE_ID == id{
			ue = ueTemp
			ok = true 
		}
		return true
	})
	return ue, ok
}

// 
func CreateAndAddAmfToLB(amfType int) *Lb_Amf{
	self := LB_Self()
	amf := newLbAmf(amfType)
	if amf.AmfTypeIdent == TypeIdRegist {
		self.LbRegistAmfPool.Store(amf.AmfID, amf)
	} else if amf.AmfTypeIdent == TypeIdDeregist {
		self.LbDeregistAmfPool.Store(amf.AmfID, amf)
	} else {
		self.LbRegularAmfPool.Store(amf.AmfID, amf)
	}
	return amf
}

// Creates, initializes and returns a new *LbAmf
func newLbAmf(amfType int) *Lb_Amf {
	var amf Lb_Amf
	amf.AmfID = nextAmfID
	amf.Lb_Conn = newLBConn(nextAmfID, TypeIdAMFConn)
	amf.Lb_Conn.AmfPointer = &amf
	amf.Log = logger.AMFLog
	amf.RelativeCapacity = 0 
	amf.NumberOfConnectedUEs = 0 
	amf.AmfTypeIdent = amfType
	nextAmfID++
	return &amf
}

// takes UeID and returns true if UE exists in the AMFs list 
func (amf *Lb_Amf) ContainsUE(id int64) (cont bool) {
	_, cont = amf.Ues.Load(id)
	return
}

// calculates a number reflecting the AMF-Usage that is comparable within the loadbalancer
func (amf *Lb_Amf) calculateAMFUsage() float32{
	if amf.RelativeCapacity == 0 {
		return 1000.0
	}
	return float32(amf.NumberOfConnectedUEs) / float32(amf.RelativeCapacity)
}

// Removes AMF-Context and closes the Connection  
func (amf *Lb_Amf) RemoveAmfContext() {
	lb := LB_Self()

	if amf.AmfTypeIdent == TypeIdRegist {
		lb.LbRegistAmfPool.Delete(amf.Lb_Conn.Conn)
	} else if amf.AmfTypeIdent == TypeIdDeregist {
		lb.LbDeregistAmfPool.Delete(amf.Lb_Conn.Conn)
	} else {
		lb.LbRegularAmfPool.Delete(amf.Lb_Conn.Conn)
	}
	amf.Lb_Conn.Conn.Close()
	amf = nil 
}


