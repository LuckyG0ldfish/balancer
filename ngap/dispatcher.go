package ngap

import (
	// "net"

	"fmt"

	"git.cs.nctu.edu.tw/calee/sctp"

	"github.com/LuckyG0ldfish/balancer/context"
	amf_ngap "github.com/LuckyG0ldfish/balancer/ngap/amf_ngap"
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
			amf_ngap.HandleNGSetupRequest(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeInitialUEMessage:
			fmt.Println("Handling InitialUEMessage")
			amf_ngap.HandleInitialUEMessage(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkNASTransport:
			fmt.Println("Handling UplinkNasTransport")
			amf_ngap.HandleUplinkNasTransport(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeNGReset:
			fmt.Println("Handling NGReset")
			amf_ngap.HandleNGReset(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverCancel:
			fmt.Println("Handling HandoverCancel")
			amf_ngap.HandleHandoverCancel(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUEContextReleaseRequest:
			fmt.Println("Handling UEContextReleaseRequest")
			amf_ngap.HandleUEContextReleaseRequest(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeNASNonDeliveryIndication:
			fmt.Println("Handling NasNonDeliveryIndication")
			amf_ngap.HandleNasNonDeliveryIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeLocationReportingFailureIndication:
			fmt.Println("Handling LocationReportingFailureIndication")
			amf_ngap.HandleLocationReportingFailureIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeErrorIndication:
			fmt.Println("Handling ErrorIndication")
			amf_ngap.HandleErrorIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUERadioCapabilityInfoIndication:
			fmt.Println("Handling UERadioCapabilityInfoIndication")
			amf_ngap.HandleUERadioCapabilityInfoIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverNotification:
			fmt.Println("Handling HandoverNotify")
			amf_ngap.HandleHandoverNotify(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverPreparation:
			fmt.Println("Handling HandoverRequired") //
			amf_ngap.HandleHandoverRequired(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			fmt.Println("Handling RanConfigurationUpdate")
			amf_ngap.HandleRanConfigurationUpdate(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeRRCInactiveTransitionReport:
			fmt.Println("Handling RRCInactiveTransitionReport")
			amf_ngap.HandleRRCInactiveTransitionReport(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceNotify:
			fmt.Println("Handling PDUSessionResourceNotify")
			amf_ngap.HandlePDUSessionResourceNotify(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePathSwitchRequest:
			fmt.Println("Handling PathSwitchRequest")
			amf_ngap.HandlePathSwitchRequest(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeLocationReport:
			fmt.Println("Handling LocationReport")
			amf_ngap.HandleLocationReport(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkUEAssociatedNRPPaTransport:
			fmt.Println("Handling UplinkUEAssociatedNRPPATransport")
			amf_ngap.HandleUplinkUEAssociatedNRPPATransport(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkRANConfigurationTransfer:
			fmt.Println("Handling UplinkRanConfigurationTransfer")
			amf_ngap.HandleUplinkRanConfigurationTransfer(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
			fmt.Println("Handling PDUSessionResourceModifyIndication")
			amf_ngap.HandlePDUSessionResourceModifyIndication(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeCellTrafficTrace:
			fmt.Println("Handling CellTrafficTrace")
			amf_ngap.HandleCellTrafficTrace(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkRANStatusTransfer:
			fmt.Println("Handling UplinkRanStatusTransfer")
			amf_ngap.HandleUplinkRanStatusTransfer(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUplinkNonUEAssociatedNRPPaTransport:
			fmt.Println("Handling UplinkNonUEAssociatedNRPPATransport")
			amf_ngap.HandleUplinkNonUEAssociatedNRPPATransport(lbConn, pdu, msgCopy)
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
			amf_ngap.HandleNGResetAcknowledge(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUEContextRelease:
			fmt.Println("Handling UEContextReleaseComplete")
			amf_ngap.HandleUEContextReleaseComplete(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			fmt.Println("Handling PDUSessionResourceReleaseResponse")
			amf_ngap.HandlePDUSessionResourceReleaseResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUERadioCapabilityCheck:
			fmt.Println("Handling UERadioCapabilityCheckResponse")
			amf_ngap.HandleUERadioCapabilityCheckResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			fmt.Println("Handling AMFconfigurationUpdateAcknowledge")
			amf_ngap.HandleAMFconfigurationUpdateAcknowledge(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeInitialContextSetup:
			fmt.Println("Handling InitialContextSetupResponse")
			amf_ngap.HandleInitialContextSetupResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUEContextModification:
			fmt.Println("Handling UEContextModificationResponse")
			amf_ngap.HandleUEContextModificationResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			fmt.Println("Handling PDUSessionResourceSetupResponse")
			amf_ngap.HandlePDUSessionResourceSetupResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodePDUSessionResourceModify:
			fmt.Println("Handling PDUSessionResourceModifyResponse")
			amf_ngap.HandlePDUSessionResourceModifyResponse(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverResourceAllocation:
			fmt.Println("Handling HandoverRequestAcknowledge")
			amf_ngap.HandleHandoverRequestAcknowledge(lbConn, pdu, msgCopy)
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
			amf_ngap.HandleAMFconfigurationUpdateFailure(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeInitialContextSetup:
			fmt.Println("Handling InitialContextSetupFailure")
			amf_ngap.HandleInitialContextSetupFailure(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeUEContextModification:
			fmt.Println("Handling UEContextModificationFailure")
			amf_ngap.HandleUEContextModificationFailure(lbConn, pdu, msgCopy)
		case ngapType.ProcedureCodeHandoverResourceAllocation:
			fmt.Println("Handling HandoverFailure")
			amf_ngap.HandleHandoverFailure(lbConn, pdu, msgCopy)
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
