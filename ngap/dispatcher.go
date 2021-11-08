package ngap

import (
	// "net"

	// "fmt"

	"git.cs.nctu.edu.tw/calee/sctp"

	"github.com/LuckyG0ldfish/balancer/context"
	amf_ngap "github.com/LuckyG0ldfish/balancer/ngap/amf_ngap"
	gnb_ngap "github.com/LuckyG0ldfish/balancer/ngap/gnb_ngap"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
)

func Dispatch(lbConn *context.LBConn, msg []byte) {
	if lbConn.TypeID == context.TypeIdGNBConn {
		DispatchForMessageToAmf(lbConn, msg)
	} else {
		DispatchForMessageToGnb(lbConn, msg)
	}
}

func DispatchForMessageToAmf(lbConn *context.LBConn, msg []byte) {
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
			logger.NgapLog.Tracef("Handling NGSetupRequest")
			amf_ngap.HandleNGSetupRequest(lbConn, pdu)
		case ngapType.ProcedureCodeInitialUEMessage:
			logger.NgapLog.Tracef("Handling InitialUEMessage")
			amf_ngap.HandleInitialUEMessage(lbConn, pdu)
		case ngapType.ProcedureCodeUplinkNASTransport:
			logger.NgapLog.Tracef("Handling UplinkNasTransport")
			amf_ngap.HandleUplinkNasTransport(lbConn, pdu)
		case ngapType.ProcedureCodeNGReset:
			logger.NgapLog.Tracef("Handling NGReset")
			amf_ngap.HandleNGReset(lbConn, pdu)
		case ngapType.ProcedureCodeHandoverCancel:
			logger.NgapLog.Tracef("Handling HandoverCancel")
			amf_ngap.HandleHandoverCancel(lbConn, pdu)
		case ngapType.ProcedureCodeUEContextReleaseRequest:
			logger.NgapLog.Tracef("Handling UEContextReleaseRequest")
			amf_ngap.HandleUEContextReleaseRequest(lbConn, pdu)
		case ngapType.ProcedureCodeNASNonDeliveryIndication:
			logger.NgapLog.Tracef("Handling NasNonDeliveryIndication")
			amf_ngap.HandleNasNonDeliveryIndication(lbConn, pdu)
		case ngapType.ProcedureCodeLocationReportingFailureIndication:
			logger.NgapLog.Tracef("Handling LocationReportingFailureIndication")
			amf_ngap.HandleLocationReportingFailureIndication(lbConn, pdu)
		case ngapType.ProcedureCodeErrorIndication:
			logger.NgapLog.Tracef("Handling ErrorIndication")
			amf_ngap.HandleErrorIndication(lbConn, pdu)
		case ngapType.ProcedureCodeUERadioCapabilityInfoIndication:
			logger.NgapLog.Tracef("Handling UERadioCapabilityInfoIndication")
			amf_ngap.HandleUERadioCapabilityInfoIndication(lbConn, pdu)
		case ngapType.ProcedureCodeHandoverNotification:
			logger.NgapLog.Tracef("Handling HandoverNotify")
			amf_ngap.HandleHandoverNotify(lbConn, pdu)
		case ngapType.ProcedureCodeHandoverPreparation:
			logger.NgapLog.Tracef("Handling HandoverRequired") //
			amf_ngap.HandleHandoverRequired(lbConn, pdu)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			logger.NgapLog.Tracef("Handling RanConfigurationUpdate")
			amf_ngap.HandleRanConfigurationUpdate(lbConn, pdu)
		case ngapType.ProcedureCodeRRCInactiveTransitionReport:
			logger.NgapLog.Tracef("Handling RRCInactiveTransitionReport")
			amf_ngap.HandleRRCInactiveTransitionReport(lbConn, pdu)
		case ngapType.ProcedureCodePDUSessionResourceNotify:
			logger.NgapLog.Tracef("Handling PDUSessionResourceNotify")
			amf_ngap.HandlePDUSessionResourceNotify(lbConn, pdu)
		case ngapType.ProcedureCodePathSwitchRequest:
			logger.NgapLog.Tracef("Handling PathSwitchRequest")
			amf_ngap.HandlePathSwitchRequest(lbConn, pdu)
		case ngapType.ProcedureCodeLocationReport:
			logger.NgapLog.Tracef("Handling LocationReport")
			amf_ngap.HandleLocationReport(lbConn, pdu)
		case ngapType.ProcedureCodeUplinkUEAssociatedNRPPaTransport:
			logger.NgapLog.Tracef("Handling UplinkUEAssociatedNRPPATransport")
			amf_ngap.HandleUplinkUEAssociatedNRPPATransport(lbConn, pdu)
		case ngapType.ProcedureCodeUplinkRANConfigurationTransfer:
			logger.NgapLog.Tracef("Handling UplinkRanConfigurationTransfer")
			amf_ngap.HandleUplinkRanConfigurationTransfer(lbConn, pdu)
		case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
			logger.NgapLog.Tracef("Handling PDUSessionResourceModifyIndication")
			amf_ngap.HandlePDUSessionResourceModifyIndication(lbConn, pdu)
		case ngapType.ProcedureCodeCellTrafficTrace:
			logger.NgapLog.Tracef("Handling CellTrafficTrace")
			amf_ngap.HandleCellTrafficTrace(lbConn, pdu)
		case ngapType.ProcedureCodeUplinkRANStatusTransfer:
			logger.NgapLog.Tracef("Handling UplinkRanStatusTransfer")
			amf_ngap.HandleUplinkRanStatusTransfer(lbConn, pdu)
		case ngapType.ProcedureCodeUplinkNonUEAssociatedNRPPaTransport:
			logger.NgapLog.Tracef("Handling UplinkNonUEAssociatedNRPPATransport")
			amf_ngap.HandleUplinkNonUEAssociatedNRPPATransport(lbConn, pdu)
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
			logger.NgapLog.Tracef("Handling NGResetAcknowledge")
			amf_ngap.HandleNGResetAcknowledge(lbConn, pdu)
		case ngapType.ProcedureCodeUEContextRelease:
			logger.NgapLog.Tracef("Handling UEContextReleaseComplete")
			amf_ngap.HandleUEContextReleaseComplete(lbConn, pdu)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			logger.NgapLog.Tracef("Handling PDUSessionResourceReleaseResponse")
			amf_ngap.HandlePDUSessionResourceReleaseResponse(lbConn, pdu)
		case ngapType.ProcedureCodeUERadioCapabilityCheck:
			logger.NgapLog.Tracef("Handling UERadioCapabilityCheckResponse")
			amf_ngap.HandleUERadioCapabilityCheckResponse(lbConn, pdu)
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			logger.NgapLog.Tracef("Handling AMFconfigurationUpdateAcknowledge")
			amf_ngap.HandleAMFconfigurationUpdateAcknowledge(lbConn, pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			logger.NgapLog.Tracef("Handling InitialContextSetupResponse")
			amf_ngap.HandleInitialContextSetupResponse(lbConn, pdu)
		case ngapType.ProcedureCodeUEContextModification:
			logger.NgapLog.Tracef("Handling UEContextModificationResponse")
			amf_ngap.HandleUEContextModificationResponse(lbConn, pdu)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			logger.NgapLog.Tracef("Handling PDUSessionResourceSetupResponse")
			amf_ngap.HandlePDUSessionResourceSetupResponse(lbConn, pdu)
		case ngapType.ProcedureCodePDUSessionResourceModify:
			logger.NgapLog.Tracef("Handling PDUSessionResourceModifyResponse")
			amf_ngap.HandlePDUSessionResourceModifyResponse(lbConn, pdu)
		case ngapType.ProcedureCodeHandoverResourceAllocation:
			logger.NgapLog.Tracef("Handling HandoverRequestAcknowledge")
			amf_ngap.HandleHandoverRequestAcknowledge(lbConn, pdu)
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
			logger.NgapLog.Tracef("Handling AMFconfigurationUpdateFailure")
			amf_ngap.HandleAMFconfigurationUpdateFailure(lbConn, pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			logger.NgapLog.Tracef("Handling InitialContextSetupFailure")
			amf_ngap.HandleInitialContextSetupFailure(lbConn, pdu)
		case ngapType.ProcedureCodeUEContextModification:
			logger.NgapLog.Tracef("Handling UEContextModificationFailure")
			amf_ngap.HandleUEContextModificationFailure(lbConn, pdu)
		case ngapType.ProcedureCodeHandoverResourceAllocation:
			logger.NgapLog.Tracef("Handling HandoverFailure")
			amf_ngap.HandleHandoverFailure(lbConn, pdu)
		default:
			// lbConn.Log.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, unsuccessfulOutcome.ProcedureCode.Value)
		}
	}
}

func DispatchForMessageToGnb(lbConn *context.LBConn, msg []byte) {
	// AMF SCTP address
	// AMF context
	// lbConn := context.LB_Self()
	// Decode
	pdu, err := ngap.Decoder(msg)
	if err != nil {
		// NGAPLog.Errorf("NGAP decode error: %+v\n", err)
		return
	}

	switch pdu.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := pdu.InitiatingMessage
		if initiatingMessage == nil {
			// NGAPLog.Errorln("Initiating Message is nil")
			return
		}

		switch initiatingMessage.ProcedureCode.Value {
		//case ngapType.ProcedureCodeNGReset:
		//	handler.HandleNGReset(amf, pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			logger.NgapLog.Tracef("Handling InitialContextSetupRequest")
			gnb_ngap.HandleInitialContextSetupRequest(lbConn, pdu)
		//case ngapType.ProcedureCodeUEContextModification:
		//	handler.HandleUEContextModificationRequest(amf, pdu)
		case ngapType.ProcedureCodeUEContextRelease:
			logger.NgapLog.Tracef("Handling UEContextReleaseCommand")
			gnb_ngap.HandleUEContextReleaseCommand(lbConn, pdu)
		case ngapType.ProcedureCodeDownlinkNASTransport:
			logger.NgapLog.Tracef("Handling DownlinkNASTransport")
			gnb_ngap.HandleDownlinkNASTransport(lbConn, pdu)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			logger.NgapLog.Tracef("Handling PDUSessionResourceSetupRequest")
			gnb_ngap.HandlePDUSessionResourceSetupRequest(lbConn, pdu)
		// TODO: This will be commented for the time being, after adding other procedures will be uncommented.
		//case ngapType.ProcedureCodePDUSessionResourceModify:
		//	handler.HandlePDUSessionResourceModifyRequest(amf, pdu)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			logger.NgapLog.Tracef("Handling PDUSessionResourceReleaseCommand")
			gnb_ngap.HandlePDUSessionResourceReleaseCommand(lbConn, pdu)
		//case ngapType.ProcedureCodeErrorIndication:
		//	handler.HandleErrorIndication(amf, pdu)
		//case ngapType.ProcedureCodeUERadioCapabilityCheck:
		//	handler.HandleUERadioCapabilityCheckRequest(amf, pdu)
		//case ngapType.ProcedureCodeAMFConfigurationUpdate:
		//	handler.HandleAMFConfigurationUpdate(amf, pdu)
		//case ngapType.ProcedureCodeDownlinkRANConfigurationTransfer:
		//	handler.HandleDownlinkRANConfigurationTransfer(pdu)
		//case ngapType.ProcedureCodeDownlinkRANStatusTransfer:
		//	handler.HandleDownlinkRANStatusTransfer(pdu)
		//case ngapType.ProcedureCodeAMFStatusIndication:
		//	handler.HandleAMFStatusIndication(pdu)
		//case ngapType.ProcedureCodeLocationReportingControl:
		//	handler.HandleLocationReportingControl(pdu)
		//case ngapType.ProcedureCodeUETNLABindingRelease:
		//	handler.HandleUETNLAReleaseRequest(pdu)
		//case ngapType.ProcedureCodeOverloadStart:
		//	handler.HandleOverloadStart(amf, pdu)
		//case ngapType.ProcedureCodeOverloadStop:
		//	handler.HandleOverloadStop(amf, pdu)
		default:
			// NGAPLog.Warnf("Not implemented NGAP message(initiatingMessage), procedureCode:%d]\n",
				// initiatingMessage.ProcedureCode.Value)
				logger.NgapLog.Warnf("Not implemented NGAP message(initiatingMessage), procedureCode:%d]\n", initiatingMessage.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := pdu.SuccessfulOutcome
		if successfulOutcome == nil {
			// NGAPLog.Errorln("Successful Outcome is nil")
			logger.NgapLog.Tracef("Successful Outcome is nil")
			return
		}

		switch successfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			logger.NgapLog.Tracef("Handling NGSetupResponse")
			LB := context.LB_Self() 
			LB.NGSetupRes = pdu
			gnb_ngap.HandleNGSetupResponse(lbConn, pdu)

		//case ngapType.ProcedureCodeNGReset:
		//	handler.HandleNGResetAcknowledge(amf, pdu)
		//case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
		//	handler.HandlePDUSessionResourceModifyConfirm(amf, pdu)
		//case ngapType.ProcedureCodeRANConfigurationUpdate:
		//	handler.HandleRANConfigurationUpdateAcknowledge(amf, pdu)
		default:
			// NGAPLog.Warnf("Not implemented NGAP message(successfulOutcome), procedureCode:%d]\n",
				// successfulOutcome.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := pdu.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			// NGAPLog.Errorln("Unsuccessful Outcome is nil")
			return
		}

		switch unsuccessfulOutcome.ProcedureCode.Value {
		//case ngapType.ProcedureCodeNGSetup:
		//	handler.HandleNGSetupFailure(sctpAddr, conn, pdu)
		//case ngapType.ProcedureCodeRANConfigurationUpdate:
		//	handler.HandleRANConfigurationUpdateFailure(amf, pdu)
		default:
			// NGAPLog.Warnf("Not implemented NGAP message(unsuccessfulOutcome), procedureCode:%d]\n",
				// unsuccessfulOutcome.ProcedureCode.Value)
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
