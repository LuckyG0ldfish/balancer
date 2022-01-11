package context

import (
	"github.com/LuckyG0ldfish/balancer/logger"
)

// TODO:
func (context *LBContext) SelectNextAmf() bool{
	if context.Next_Regist_Amf == nil {
		logger.ContextLog.Errorf("No Amf found")
		return false
	}

	var amfWithMaxCap *LbAmf = context.Next_Regist_Amf
	var amfUsage float64 = context.Next_Regist_Amf.calculateAMFUsage()
	
	context.LbAmfPool.Range(func(key, value interface{}) bool{
		amfTemp, ok := value.(*LbAmf)
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
	context.Next_Regist_Amf = amfWithMaxCap
	logger.ContextLog.Tracef("NextAMF = AMFID: %d", amfWithMaxCap.AmfID)
	return true 
}

// TODO:
func SelectNextRegistAmf() bool{
	context := LB_Self()
	if context.Next_Regular_Amf == nil {
		logger.NgapLog.Errorf("No Amf found")
		return false 
	}

	nextAmf := findNextAMF(TypeIdRegist)
	context.Next_Regist_Amf = nextAmf
	logger.ContextLog.Tracef("NextAMF = AMFID: %d", nextAmf.AmfID)
	return true 
}

// TODO:
func SelectNextRegularAmf() bool{
	context := LB_Self()
	if context.Next_Regular_Amf == nil {
		logger.NgapLog.Errorf("No Amf found")
		return false 
	}

	amfWithMaxCap := findNextAMF(TypeIdRegular)
	context.Next_Regist_Amf = amfWithMaxCap
	logger.ContextLog.Tracef("NextRegularAMF = AMFID: %d", amfWithMaxCap.AmfID)
	return true 
}

// TODO:
func SelectNextDeregistAmf() bool{
	context := LB_Self()
	if context.Next_Deregist_Amf == nil {
		logger.ContextLog.Errorf("No Amf found")
		return false 
	}
	amfWithMaxCap := findNextAMF(TypeIdDeregist)
	context.Next_Regist_Amf = amfWithMaxCap
	logger.ContextLog.Tracef("NextDeregistAMF = AMFID: %d", amfWithMaxCap.AmfID)
	return true 
}

func findNextAMF(state int) *LbAmf{
	lb := LB_Self()
	var amfWithMaxCap *LbAmf
	var amfUsage float64
	
	switch state {
	case TypeIdRegist:
		amfWithMaxCap = lb.Next_Regist_Amf
	case TypeIdRegular:
		amfWithMaxCap = lb.Next_Regular_Amf
	case TypeIdDeregist:
		amfWithMaxCap = lb.Next_Deregist_Amf
	} 
	amfUsage = amfWithMaxCap.calculateAMFUsage()
	
	lb.LbAmfPool.Range(func(key, value interface{}) bool{
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