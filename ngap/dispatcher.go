package ngap

import (
	// "net"

	"fmt"

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
			fmt.Println("Handling NGSetupRequest")
			HandleNGSetupRequest(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeInitialUEMessage:
			fmt.Println("Handling InitialUEMessage")
			HandleInitialUEMessage(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkNASTransport:
			fmt.Println("Handling UplinkNasTransport")
			HandleUplinkNasTransport(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeNGReset:
			fmt.Println("Handling NGReset")
			HandleNGReset(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverCancel:
			fmt.Println("Handling HandoverCancel")
			HandleHandoverCancel(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUEContextReleaseRequest:
			fmt.Println("Handling UEContextReleaseRequest")
			HandleUEContextReleaseRequest(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeNASNonDeliveryIndication:
			fmt.Println("Handling NasNonDeliveryIndication")
			HandleNasNonDeliveryIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeLocationReportingFailureIndication:
			fmt.Println("Handling LocationReportingFailureIndication")
			HandleLocationReportingFailureIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeErrorIndication:
			fmt.Println("Handling ErrorIndication")
			HandleErrorIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUERadioCapabilityInfoIndication:
			fmt.Println("Handling UERadioCapabilityInfoIndication")
			HandleUERadioCapabilityInfoIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverNotification:
			fmt.Println("Handling HandoverNotify")
			HandleHandoverNotify(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverPreparation:
			fmt.Println("Handling HandoverRequired") //
			HandleHandoverRequired(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			fmt.Println("Handling RanConfigurationUpdate")
			HandleRanConfigurationUpdate(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeRRCInactiveTransitionReport:
			fmt.Println("Handling RRCInactiveTransitionReport")
			HandleRRCInactiveTransitionReport(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceNotify:
			fmt.Println("Handling PDUSessionResourceNotify")
			HandlePDUSessionResourceNotify(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePathSwitchRequest:
			fmt.Println("Handling PathSwitchRequest")
			HandlePathSwitchRequest(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeLocationReport:
			fmt.Println("Handling LocationReport")
			HandleLocationReport(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkUEAssociatedNRPPaTransport:
			fmt.Println("Handling UplinkUEAssociatedNRPPATransport")
			HandleUplinkUEAssociatedNRPPATransport(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkRANConfigurationTransfer:
			fmt.Println("Handling UplinkRanConfigurationTransfer")
			HandleUplinkRanConfigurationTransfer(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
			fmt.Println("Handling PDUSessionResourceModifyIndication")
			HandlePDUSessionResourceModifyIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeCellTrafficTrace:
			fmt.Println("Handling CellTrafficTrace")
			HandleCellTrafficTrace(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkRANStatusTransfer:
			fmt.Println("Handling UplinkRanStatusTransfer")
			HandleUplinkRanStatusTransfer(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkNonUEAssociatedNRPPaTransport:
			fmt.Println("Handling UplinkNonUEAssociatedNRPPATransport")
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
			fmt.Println("Handling NGResetAcknowledge")
			HandleNGResetAcknowledge(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUEContextRelease:
			fmt.Println("Handling UEContextReleaseComplete")
			HandleUEContextReleaseComplete(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			fmt.Println("Handling PDUSessionResourceReleaseResponse")
			HandlePDUSessionResourceReleaseResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUERadioCapabilityCheck:
			fmt.Println("Handling UERadioCapabilityCheckResponse")
			HandleUERadioCapabilityCheckResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			fmt.Println("Handling AMFconfigurationUpdateAcknowledge")
			HandleAMFconfigurationUpdateAcknowledge(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeInitialContextSetup:
			fmt.Println("Handling InitialContextSetupResponse")
			HandleInitialContextSetupResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUEContextModification:
			fmt.Println("Handling UEContextModificationResponse")
			HandleUEContextModificationResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			fmt.Println("Handling PDUSessionResourceSetupResponse")
			HandlePDUSessionResourceSetupResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceModify:
			fmt.Println("Handling PDUSessionResourceModifyResponse")
			HandlePDUSessionResourceModifyResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverResourceAllocation:
			fmt.Println("Handling HandoverRequestAcknowledge")
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
			fmt.Println("Handling AMFconfigurationUpdateFailure")
			HandleAMFconfigurationUpdateFailure(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeInitialContextSetup:
			fmt.Println("Handling InitialContextSetupFailure")
			HandleInitialContextSetupFailure(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUEContextModification:
			fmt.Println("Handling UEContextModificationFailure")
			HandleUEContextModificationFailure(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverResourceAllocation:
			fmt.Println("Handling HandoverFailure")
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
