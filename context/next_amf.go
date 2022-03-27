package context

import (
	"sync"

	"github.com/LuckyG0ldfish/balancer/logger"
)

// // TODO:
// func (context *LBContext) SelectNextAmf() bool{
// 	if context.Next_Regist_Amf == nil {
// 		logger.ContextLog.Errorf("No Amf found")
// 		return false
// 	}

// 	var amfWithMaxCap *LbAmf = context.Next_Regist_Amf
// 	var amfUsage float64 = context.Next_Regist_Amf.calculateAMFUsage()
	
// 	context.LbRegistAmfPool.Range(func(key, value interface{}) bool{
// 		amfTemp, ok := value.(*LbAmf)
// 		if !ok {
// 			logger.NgapLog.Errorf("couldn't be converted")
// 		}
// 		tempUsage := amfTemp.calculateAMFUsage()
		
// 		// chooses the AMF with the lowest Usage
// 		if  amfUsage > tempUsage {
// 			amfWithMaxCap = amfTemp
// 			amfUsage = tempUsage
// 		} 
// 		return true
// 	})
// 	context.Next_Regist_Amf = amfWithMaxCap
// 	logger.ContextLog.Tracef("NextAMF = AMFID: %d", amfWithMaxCap.AmfID)
// 	return true 
// }

// TODO:
func (context *LBContext) SelectNextRegistAmf() bool{
	if context.Next_Regular_Amf == nil {
		logger.NgapLog.Errorf("No Amf found")
		return false 
	}
	if context.Next_Regist_Amf.NumberOfConnectedUEs % 5 != 0 {
		return true 
	}
	nextAmf := context.findNextAMF(TypeIdRegist)
	context.Next_Regist_Amf = nextAmf
	logger.ContextLog.Tracef("NextAMF = AMFID: %d", nextAmf.AmfID)
	return true 
}

// TODO:
func (context *LBContext) SelectNextRegularAmf() bool{
	if context.Next_Regular_Amf == nil {
		logger.NgapLog.Errorf("No Amf found")
		return false 
	}
	if context.Next_Regular_Amf.NumberOfConnectedUEs % 5 != 0 {
		return true 
	}
	amfWithMaxCap := context.findNextAMF(TypeIdRegular)
	context.Next_Regular_Amf = amfWithMaxCap
	logger.ContextLog.Tracef("NextRegularAMF = AMFID: %d", amfWithMaxCap.AmfID)
	return true 
}

// TODO:
func (context *LBContext) SelectNextDeregistAmf() bool{
	if context.Next_Deregist_Amf == nil {
		logger.ContextLog.Errorf("No Amf found")
		return false 
	}
	if context.Next_Deregist_Amf.NumberOfConnectedUEs % 5 != 0 {
		return true 
	}
	amfWithMaxCap := context.findNextAMF(TypeIdDeregist)
	context.Next_Deregist_Amf = amfWithMaxCap
	logger.ContextLog.Tracef("NextDeregistAMF = AMFID: %d", amfWithMaxCap.AmfID)
	return true 
}

func (context *LBContext) findNextAMF(state int) *LbAmf{
	lb := LB_Self()
	var amfWithMaxCap *LbAmf
	var amfUsage float64
	var pool *sync.Map

	switch state {
	case TypeIdRegist:
		amfWithMaxCap = lb.Next_Regist_Amf
		pool = &lb.LbRegistAmfPool
	case TypeIdRegular:
		amfWithMaxCap = lb.Next_Regular_Amf
		pool = &lb.LbRegularAmfPool
	case TypeIdDeregist:
		amfWithMaxCap = lb.Next_Deregist_Amf
		pool = &lb.LbDeregistAmfPool
	} 
	amfUsage = amfWithMaxCap.calculateAMFUsage()
	
	pool.Range(func(key, value interface{}) bool{
		amfTemp, ok := value.(*LbAmf)
		if !ok {
			logger.NgapLog.Errorf("couldn't be converted")
		}
		if amfTemp.AmfTypeIdent == state{
			tempUsage := amfTemp.calculateAMFUsage()
		
			// chooses the AMF with the lowest Usage
			if  amfUsage > tempUsage {
				amfWithMaxCap = amfTemp
				amfUsage = tempUsage
			} 
		}
		return true
	})
	return amfWithMaxCap
}