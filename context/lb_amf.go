package context

import (
	"sync"

	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/sirupsen/logrus"
)

var nextAmfID int64 = 1

// Type, that stores all relevant information of connected AMFs 
type LB_AMF struct {
	AmfID  					int64 			// INTERNAL ID for this AMF 

	AmfTypeIdent 			int 			// Identifies the type of AMF 

	LB_Conn 				*LB_Conn 		// Stores all the connection related information 

	RelativeCapacity 		int64 			// AMFs Relative Cap. -> extracted out of NGSetup
	NumberOfConnectedUEs 	int64 			// Amount of UEs that are connected to this AMF 
	
	Ues    					sync.Map 		// "List" of all UE that are processed by this AMF 

	/* logger */
	Log 					*logrus.Entry
}

// Use a UELB-ID to find UE context, return *LbUe and true if found
func (amf *LB_AMF) FindUeByRAN_UE_ID(id int64) (*LB_UE, bool){
	ue, ok := amf.Ues.Load(id)
	if !ok {
		amf.Log.Errorf("UE is not registered to this AMF %d", id)
		return nil, false 
	}
	ue2, ok :=  ue.(*LB_UE)
	if !ok {
		amf.Log.Errorf("couldn't be converted")
		return nil, false 
	}
	return ue2, ok
}

// Use a UEAMF-ID to find UE context, return *LbUe and true if found TODO
func (amf *LB_AMF) FindUeByAMF_UE_ID(id int64) (*LB_UE, bool){
	var ue *LB_UE
	var ok bool = false 
	amf.Ues.Range(func(key, value interface{}) bool{
		ueTemp, okTemp := value.(*LB_UE)
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

// Creates and adds an Amf to the LB_Context
func CreateAndAddAmfToLB(amfType int) *LB_AMF{
	self := LB_Self()
	amf := newLbAmf(amfType)
	if amf.AmfTypeIdent == TypeIdRegist {
		self.RegistAMFPool.Store(amf.AmfID, amf)
	} else if amf.AmfTypeIdent == TypeIdDeregist {
		self.DeregistAMFPool.Store(amf.AmfID, amf)
	} else {
		self.RegularAMFPool.Store(amf.AmfID, amf)
	}
	return amf
}

// Creates, initializes and returns a new *LbAmf
func newLbAmf(amfType int) *LB_AMF {
	var amf LB_AMF
	amf.AmfID = nextAmfID
	amf.LB_Conn = newLBConn(nextAmfID, TypeIdAMFConn)
	amf.LB_Conn.AmfPointer = &amf
	amf.Log = logger.AMFLog
	amf.RelativeCapacity = 0 
	amf.NumberOfConnectedUEs = 0 
	amf.AmfTypeIdent = amfType
	nextAmfID++
	return &amf
}

// takes UeID and returns true if UE exists in the AMFs list 
func (amf *LB_AMF) ContainsUE(id int64) (cont bool) {
	_, cont = amf.Ues.Load(id)
	return
}

// calculates a number reflecting the AMF-Usage that is comparable within the loadbalancer
func (amf *LB_AMF) calculateAMFUsage() float32{
	if amf.RelativeCapacity == 0 {
		return 10000.0 	// equivalent to between 10000 and 255000 
						// depending on relative capacity
	}
	return float32(amf.NumberOfConnectedUEs) / float32(amf.RelativeCapacity)
}

// Removes AMF-Context and closes the Connection  
func (amf *LB_AMF) RemoveAmfContext() {
	lb := LB_Self()

	if amf.AmfTypeIdent == TypeIdRegist {
		lb.RegistAMFPool.Delete(amf.LB_Conn.Conn)
	} else if amf.AmfTypeIdent == TypeIdDeregist {
		lb.DeregistAMFPool.Delete(amf.LB_Conn.Conn)
	} else {
		lb.RegularAMFPool.Delete(amf.LB_Conn.Conn)
	}
	amf.LB_Conn.Conn.Close()
	amf = nil 
}


