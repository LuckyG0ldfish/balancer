package context

import (
	"sync"

	"github.com/LuckyG0ldfish/balancer/logger"
)

// TODO:
func (context *LB_Context) SelectNextRegistAmf() bool{
	if context.Next_Regist_Amf == nil {
		logger.NgapLog.Errorf("No Amf found")
		return false 
	}
	if context.Next_Regist_Amf.NumberOfConnectedUEs % 3 != 0 {
		return true 
	}
	nextAmf := context.findNextAMF(TypeIdRegist)
	context.Next_Regist_Amf = nextAmf
	logger.ContextLog.Tracef("NextRegistAMF = AMFID: %d", nextAmf.AmfID)
	return true 
}

// TODO:
func (context *LB_Context) SelectNextRegularAmf() bool{
	if context.Next_Regular_Amf == nil {
		logger.NgapLog.Errorf("No Amf found")
		return false 
	}
	if context.Next_Regular_Amf.NumberOfConnectedUEs % 3 != 0 {
		return true 
	}
	amfWithMaxCap := context.findNextAMF(TypeIdRegular)
	context.Next_Regular_Amf = amfWithMaxCap
	logger.ContextLog.Tracef("NextRegularAMF = AMFID: %d", amfWithMaxCap.AmfID)
	return true 
}

// TODO:
func (context *LB_Context) SelectNextDeregistAmf() bool{
	if context.Next_Deregist_Amf == nil {
		logger.ContextLog.Errorf("No Amf found")
		return false 
	}
	if context.Next_Deregist_Amf.NumberOfConnectedUEs % 3 != 0 {
		return true 
	}
	amfWithMaxCap := context.findNextAMF(TypeIdDeregist)
	context.Next_Deregist_Amf = amfWithMaxCap
	logger.ContextLog.Tracef("NextDeregistAMF = AMFID: %d", amfWithMaxCap.AmfID)
	return true 
}

func (context *LB_Context) findNextAMF(state int) *LB_AMF{
	lb := LB_Self()
	var amfWithMaxCap *LB_AMF
	var amfUsage float32
	var pool *sync.Map

	switch state {
	case TypeIdRegist:
		amfWithMaxCap = lb.Next_Regist_Amf
		pool = &lb.RegistAMFPool
	case TypeIdRegular:
		amfWithMaxCap = lb.Next_Regular_Amf
		pool = &lb.RegularAMFPool
	case TypeIdDeregist:
		amfWithMaxCap = lb.Next_Deregist_Amf
		pool = &lb.DeregistAMFPool
	} 
	amfUsage = amfWithMaxCap.calculateAMFUsage()
	
	pool.Range(func(key, value interface{}) bool{
		amfTemp, ok := value.(*LB_AMF)
		if !ok {
			logger.NgapLog.Errorf("couldn't be converted")
		}
		
		tempUsage := amfTemp.calculateAMFUsage()
		
		// chooses the AMF with the lowest Usage
		if  amfUsage > tempUsage {
			amfWithMaxCap = amfTemp
			amfUsage = tempUsage
		} 
		
		return true
	})
	return amfWithMaxCap
}