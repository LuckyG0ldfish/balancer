package ngap

import (
	// "net"

	"git.cs.nctu.edu.tw/calee/sctp"

	"github.com/LuckyG0ldfish/balancer/context" 
	"github.com/free5gc/amf/logger"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
)

func Dispatch(lbConn *context.LBConn, msg []byte) {
	// var lbConn *context.LBConn
	// lbSelf := context.LB_Self()

	if len(msg) == 0 {
		// ran.Log.Infof("RAN close the connection.")
		// ran.Remove()
		return
	}
	msgCopy := make([]byte, len(msg))
	copy(msgCopy, msg)
	
	// ran, ok := lbSelf.LbGnbFindByConn(lbConn.Conn)
	// if !ok {
	// 	logger.NgapLog.Infof("Create a new NG connection for: %s", lbConn.Conn.RemoteAddr().String())
	// 	ran = lbSelf.AddGnbToLB(lbConn.Conn)
	// }

	pdu, err := ngap.Decoder(msg)
	if err != nil {
		// ran.Log.Errorf("NGAP decode error : %+v", err)
		return
	}

	// lbConn = ran.LbConn

	switch pdu.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := pdu.InitiatingMessage
		if initiatingMessage == nil {
			// lbConn.Log.Errorln("Initiating Message is nil")
			return
		}
		switch initiatingMessage.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			HandleNGSetupRequest(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeInitialUEMessage:
			HandleInitialUEMessage(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkNASTransport:
			HandleUplinkNasTransport(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeNGReset:
			HandleNGReset(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverCancel:
			HandleHandoverCancel(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUEContextReleaseRequest:
			HandleUEContextReleaseRequest(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeNASNonDeliveryIndication:
			HandleNasNonDeliveryIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeLocationReportingFailureIndication:
			HandleLocationReportingFailureIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeErrorIndication:
			HandleErrorIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUERadioCapabilityInfoIndication:
			HandleUERadioCapabilityInfoIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverNotification:
			HandleHandoverNotify(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverPreparation:
			HandleHandoverRequired(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			HandleRanConfigurationUpdate(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeRRCInactiveTransitionReport:
			HandleRRCInactiveTransitionReport(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceNotify:
			HandlePDUSessionResourceNotify(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePathSwitchRequest:
			HandlePathSwitchRequest(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeLocationReport:
			HandleLocationReport(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkUEAssociatedNRPPaTransport:
			HandleUplinkUEAssociatedNRPPATransport(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkRANConfigurationTransfer:
			HandleUplinkRanConfigurationTransfer(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
			HandlePDUSessionResourceModifyIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeCellTrafficTrace:
			HandleCellTrafficTrace(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkRANStatusTransfer:
			HandleUplinkRanStatusTransfer(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkNonUEAssociatedNRPPaTransport:
			HandleUplinkNonUEAssociatedNRPPATransport(lbConn, pdu, msgCopy)
		default:
			// lbConn.Log.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, initiatingMessage.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := pdu.SuccessfulOutcome
		if successfulOutcome == nil {
			// lbConn.Log.Errorln("successful Outcome is nil")
			return
		}
		switch successfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGReset:
			HandleNGResetAcknowledge(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUEContextRelease:
			HandleUEContextReleaseComplete(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			HandlePDUSessionResourceReleaseResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUERadioCapabilityCheck:
			HandleUERadioCapabilityCheckResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			HandleAMFconfigurationUpdateAcknowledge(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeInitialContextSetup:
			HandleInitialContextSetupResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUEContextModification:
			HandleUEContextModificationResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			HandlePDUSessionResourceSetupResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceModify:
			HandlePDUSessionResourceModifyResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverResourceAllocation:
			HandleHandoverRequestAcknowledge(lbConn, pdu, msgCopy)
		default:
			// lbConn.Log.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, successfulOutcome.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := pdu.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			// lbConn.Log.Errorln("unsuccessful Outcome is nil")
			return
		}
		switch unsuccessfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			HandleAMFconfigurationUpdateFailure(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeInitialContextSetup:
			HandleInitialContextSetupFailure(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUEContextModification:
			HandleUEContextModificationFailure(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverResourceAllocation:
			HandleHandoverFailure(lbConn, pdu, msgCopy)
		default:
			// lbConn.Log.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, unsuccessfulOutcome.ProcedureCode.Value)
		}
	}
}

func HandleSCTPNotification(conn *sctp.SCTPConn, notification sctp.Notification) {
	lbSelf := context.LB_Self()

	logger.NgapLog.Infof("Handle SCTP Notification[addr: %+v]", conn.RemoteAddr())

	_, ok := lbSelf.LbAmfFindByConn(conn)
	if !ok {
		logger.NgapLog.Warnf("RAN context has been removed[addr: %+v]", conn.RemoteAddr())
		return
	}

	switch notification.Type() {
	case sctp.SCTP_ASSOC_CHANGE:
		// ran.Log.Infof("SCTP_ASSOC_CHANGE notification")
		event := notification.(*sctp.SCTPAssocChangeEvent)
		switch event.State() {
		case sctp.SCTP_COMM_LOST:
			// ran.Log.Infof("SCTP state is SCTP_COMM_LOST, close the connection")
			// ran.Remove()
		case sctp.SCTP_SHUTDOWN_COMP:
			// ran.Log.Infof("SCTP state is SCTP_SHUTDOWN_COMP, close the connection")
			// ran.Remove()
		default:
			// ran.Log.Warnf("SCTP state[%+v] is not handled", event.State())
		}
	case sctp.SCTP_SHUTDOWN_EVENT:
		// ran.Log.Infof("SCTP_SHUTDOWN_EVENT notification, close the connection")
		// ran.Remove()
	default:
		// ran.Log.Warnf("Non handled notification type: 0x%x", notification.Type())
	}
}
