package nas

import (
	"fmt"
	"reflect"

	"github.com/LuckyG0ldfish/balancer/context"
	// "github.com/LuckyG0ldfish/balancer/gmm"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/openapi/models"
)

const MsgTypeMsgTypeRegistrationComplete int = 1 
const MsgTypeDeregistrationRequestUEOriginatingDeregistration int = 2 
const MsgTypeDeregistrationAcceptUETerminatedDeregistration int = 3
const MsgTypeOther int = 4

func HandleNAS(ue *context.LbUe, nasPdu []byte) {
	if ue == nil {
		logger.NASLog.Error("RanUe is nil")
		return
	}

	if nasPdu == nil {
		logger.NASLog.Error("nasPdu is nil")
		return
	}

	// self := context.LB_Self()
	var accessType models.AccessType // TODO
	msgType, err := IdentMsgType(ue, accessType, nasPdu)
	if err != nil { 
		logger.NASLog.Errorln(err)
		return
	}

	switch msgType {
	case nas.MsgTypeRegistrationComplete:
		ue.UeStateIdent = context.TypeIdRegular
		logger.NASLog.Traceln("MsgTypeRegistrationComplete")
		// next := self.Next_Regular_Amf
		// ue.RemoveUeFromAMF()
		// ue.AddUeToAmf(next)
		// go context.SelectNextRegularAmf()
		return
	case nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration:
		ue.UeStateIdent = context.TypeIdDeregist
		logger.NASLog.Traceln("MsgTypeDeregistrationRequestUEOriginatingDeregistration")
		// next := self.Next_Deregist_Amf
		// ue.RemoveUeFromAMF()
		// ue.AddUeToAmf(next)
		// go context.SelectNextDeregistAmf()
		return
	case nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration:
		ue.UeStateIdent = context.TypeIdDeregist
		logger.NASLog.Traceln("MsgTypeDeregistrationAcceptUETerminatedDeregistration")
		// next := self.Next_Deregist_Amf
		// ue.RemoveUeFromAMF()
		// ue.AddUeToAmf(next)
		// go context.SelectNextDeregistAmf()
		return	
	} 
	return 
}

/*
payload either a security protected 5GS NAS message or a plain 5GS NAS message which
format is followed TS 24.501 9.1.1
*/
func IdentMsgType(ue *context.LbUe, accessType models.AccessType, payload []byte) (uint8, error) {
	if ue == nil {
		return 0, fmt.Errorf("amfUe is nil")
	}
	if payload == nil {
		return 0, fmt.Errorf("nas payload is empty")
	}

	

	msg := new(nas.Message)
	msg.SecurityHeaderType = nas.GetSecurityHeaderType(payload) & 0x0f
	if msg.SecurityHeaderType == nas.SecurityHeaderTypePlainNas {
		logger.NASLog.Traceln("SecurityHeaderType == nas.SecurityHeaderTypePlainNas")
		if ue.RRCECause != "0" { // ue.SecurityContextAvailable && 
			if err := msg.PlainNasDecode(&payload); err != nil {
				return 0, fmt.Errorf("failed to decode")
			}

			if msg.GmmMessage == nil {
				return 0, fmt.Errorf("gmm Message is nil")
			}

			test := msg.GmmHeader.GetMessageType()

			switch test {
			case nas.MsgTypeAuthenticationResponse:
				logger.NASLog.Traceln("MsgTypeAuthenticationResponse")
				// gmm.HandleAuthenticationResponse()
				return nas.MsgTypeAuthenticationResponse, nil
			case nas.MsgTypeRegistrationRequest:
				logger.NASLog.Traceln("MsgTypeRegistrationRequest")
				// gmm.HandleRegistrationRequest()
				return nas.MsgTypeRegistrationRequest, nil
			default:
				logger.NASLog.Tracef("%d", test)
				return 0, nil
				// fmt.Errorf("UE can not send plain nas for non-emergency service when there is a valid security context")
			// }
			} 
		} else {
			ue.MacFailed = false
			err := msg.PlainNasDecode(&payload)
			return 0, err
		}
	
	} else { 
		// Security protected NAS message
		securityHeader := payload[0:6]
		logger.NASLog.Traceln("securityHeader is ", securityHeader)
		sequenceNumber := payload[6]
		logger.NASLog.Traceln("sequenceNumber", sequenceNumber)
	
		receivedMac32 := securityHeader[2:]
		// remove security Header except for sequece Number
		payload = payload[6:]
	
		// a security protected NAS message must be integrity protected, and ciphering is optional
		ciphered := false
		switch msg.SecurityHeaderType {
		case nas.SecurityHeaderTypeIntegrityProtected:
			logger.NASLog.Debugln("Security header type: Integrity Protected")
		case nas.SecurityHeaderTypeIntegrityProtectedAndCiphered:
			logger.NASLog.Debugln("Security header type: Integrity Protected And Ciphered")
			ciphered = true
		case nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext:
			logger.NASLog.Debugln("Security header type: Integrity Protected And Ciphered With New 5G Security Context")
			ciphered = true
				ue.ULCount.Set(0, 0)
		default:
			return 0, fmt.Errorf("Wrong security header type: 0x%0x", msg.SecurityHeader.SecurityHeaderType)
		}
	

		// TODO
		if ue.ULCount.SQN() > sequenceNumber {
			logger.NASLog.Debugf("set ULCount overflow")
			ue.ULCount.SetOverflow(ue.ULCount.Overflow() + 1)
		}
		ue.ULCount.SetSQN(sequenceNumber)
	
		logger.NASLog.Debugf("Calculate NAS MAC (algorithm: %+v, ULCount: 0x%0x)", ue.IntegrityAlg, ue.ULCount.Get())
		logger.NASLog.Tracef("NAS integrity key0x: %0x", ue.KnasInt)
		mac32, err := security.NASMacCalculate(ue.IntegrityAlg, ue.KnasInt, ue.ULCount.Get(),
			GetBearerType(accessType), security.DirectionUplink, payload)
		if err != nil {
			return 0, fmt.Errorf("MAC calcuate error: %+v", err)
		}
	
		if !reflect.DeepEqual(mac32, receivedMac32) {
			logger.NASLog.Tracef("NAS MAC verification failed(received: 0x%08x, expected: 0x%08x)", receivedMac32, mac32)
			ue.MacFailed = true
		} else {
			logger.NASLog.Tracef("cmac value: 0x%08x", mac32)
			ue.MacFailed = false
		}
	
		if ciphered {
			logger.NASLog.Debugf("Decrypt NAS message (algorithm: %+v, ULCount: 0x%0x)", ue.CipheringAlg, ue.ULCount.Get())
			logger.NASLog.Tracef("NAS ciphering key: %0x", ue.KnasEnc)
			// decrypt payload without sequence number (payload[1])
			if err = security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.ULCount.Get(), GetBearerType(accessType),
				security.DirectionUplink, payload[1:]); err != nil {
				return 0, fmt.Errorf("Encrypt error: %+v", err)
			}
		}
	
		// remove sequece Number
		payload = payload[1:]
		err = msg.PlainNasDecode(&payload)
		return msg.GmmHeader.GetMessageType(), err
		}
	// }
	// return fmt.Errorf("nas payload is not in plain")
}

func GetBearerType(accessType models.AccessType) uint8 {
	if accessType == models.AccessType__3_GPP_ACCESS {
		return security.Bearer3GPP
	} else if accessType == models.AccessType_NON_3_GPP_ACCESS {
		return security.BearerNon3GPP
	} else {
		return security.OnlyOneBearer
	}
}

