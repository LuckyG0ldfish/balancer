package nas

import (
	"fmt"

	"github.com/LuckyG0ldfish/balancer/context"
	"github.com/LuckyG0ldfish/balancer/logger"

	"github.com/free5gc/nas"
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

	err := IdentMsgType(ue, nasPdu)
	if err != nil { 
		logger.NASLog.Errorln(err)
		return
	}
}

/*
payload either a security protected 5GS NAS message or a plain 5GS NAS message which
format is followed TS 24.501 9.1.1
*/
func IdentMsgType(ue *context.LbUe, payload []byte) (error) {
	if ue == nil {
		return fmt.Errorf("amfUe is nil")
	}
	if payload == nil {
		return fmt.Errorf("Nas payload is empty")
	}

	// self := context.LB_Self()

	msg := new(nas.Message)
	msg.SecurityHeaderType = nas.GetSecurityHeaderType(payload) & 0x0f
	if msg.SecurityHeaderType == nas.SecurityHeaderTypePlainNas {
			if err := msg.PlainNasDecode(&payload); err != nil {
				return err
			}

			if msg.GmmMessage == nil {
				return fmt.Errorf("Gmm Message is nil")
			}

			switch msg.GmmHeader.GetMessageType() {
			case nas.MsgTypeRegistrationComplete:
				ue.UeStateIdent = context.TypeIdRegular
				// next := self.Next_Regular_Amf
				// ue.RemoveUeFromAMF()
				// ue.AddUeToAmf(next)
				// go context.SelectNextRegularAmf()
				return nil
			case nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration:
				ue.UeStateIdent = context.TypeIdDeregist
				// next := self.Next_Deregist_Amf
				// ue.RemoveUeFromAMF()
				// ue.AddUeToAmf(next)
				// go context.SelectNextDeregistAmf()
				return nil
			case nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration:
				ue.UeStateIdent = context.TypeIdDeregist
				// next := self.Next_Deregist_Amf
				// ue.RemoveUeFromAMF()
				// ue.AddUeToAmf(next)
				// go context.SelectNextDeregistAmf()
				return nil
			default:
				return fmt.Errorf(
					"UE can not send plain nas for non-emergency service when there is a valid security context")
		}
	}
	return fmt.Errorf("Nas payload is not in plain")
}
