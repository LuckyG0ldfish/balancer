package ngap

import (

	"github.com/ishidawataru/sctp"

	"github.com/LuckyG0ldfish/balancer/context"
	amf_ngap "github.com/LuckyG0ldfish/balancer/ngap/amf_ngap"
	gnb_ngap "github.com/LuckyG0ldfish/balancer/ngap/gnb_ngap"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
)

// Distributes message to the correct handler based on Type of the incoming Connection 
func Dispatch(lbConn *context.LBConn, msg []byte) {
	if lbConn.TypeID == context.TypeIdGNBConn {
		DispatchForMessageToAmf(lbConn, msg) 		
	} else if lbConn.TypeID == context.TypeIdAMFConn {
		DispatchForMessageToGnb(lbConn, msg)		
	} else {
		logger.NgapLog.Errorf("Connection undefiend!")
	}
}

// This handles messages incoming from GNB with the functions of the AMFs handler 
func DispatchForMessageToAmf(lbConn *context.LBConn, msg []byte) {
	if len(msg) == 0 {
		lbConn.Log.Infof("RAN close the connection.")
		// ran.Remove() TODO
		return
	}

	pdu, err := ngap.Decoder(msg)
	if err != nil {
		lbConn.Log.Errorf("NGAP decode error : %+v", err)
		return
	}

	switch pdu.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := pdu.InitiatingMessage
		if initiatingMessage == nil {
			lbConn.Log.Errorln("Initiating Message is nil")
			return
		}
		switch initiatingMessage.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			logger.NgapLog.Tracef("Handling NGSetupRequest")
			go gnb_ngap.HandleNGSetupRequest(lbConn, pdu)
		case ngapType.ProcedureCodeInitialUEMessage:
			logger.NgapLog.Tracef("Handling InitialUEMessage")
			go gnb_ngap.HandleInitialUEMessage(lbConn, pdu)
		case ngapType.ProcedureCodeUplinkNASTransport:
			logger.NgapLog.Tracef("Handling UplinkNasTransport")
			go gnb_ngap.HandleUplinkNasTransport(lbConn, pdu)
		case ngapType.ProcedureCodeNGReset:
			logger.NgapLog.Tracef("Handling NGReset")
			go gnb_ngap.HandleNGReset(lbConn, pdu)
		case ngapType.ProcedureCodeHandoverCancel:
			logger.NgapLog.Tracef("Handling HandoverCancel")
			go gnb_ngap.HandleHandoverCancel(lbConn, pdu)
		case ngapType.ProcedureCodeUEContextReleaseRequest:
			logger.NgapLog.Tracef("Handling UEContextReleaseRequest")
			go gnb_ngap.HandleUEContextReleaseRequest(lbConn, pdu)
		case ngapType.ProcedureCodeNASNonDeliveryIndication:
			logger.NgapLog.Tracef("Handling NasNonDeliveryIndication")
			go gnb_ngap.HandleNasNonDeliveryIndication(lbConn, pdu)
		case ngapType.ProcedureCodeLocationReportingFailureIndication:
			logger.NgapLog.Tracef("Handling LocationReportingFailureIndication")
			go gnb_ngap.HandleLocationReportingFailureIndication(lbConn, pdu)
		case ngapType.ProcedureCodeErrorIndication:
			logger.NgapLog.Tracef("Handling ErrorIndication")
			go gnb_ngap.HandleErrorIndication(lbConn, pdu)
		case ngapType.ProcedureCodeUERadioCapabilityInfoIndication:
			logger.NgapLog.Tracef("Handling UERadioCapabilityInfoIndication")
			go gnb_ngap.HandleUERadioCapabilityInfoIndication(lbConn, pdu)
		case ngapType.ProcedureCodeHandoverNotification:
			logger.NgapLog.Tracef("Handling HandoverNotify")
			go gnb_ngap.HandleHandoverNotify(lbConn, pdu)
		case ngapType.ProcedureCodeHandoverPreparation:
			logger.NgapLog.Tracef("Handling HandoverRequired") //
			go gnb_ngap.HandleHandoverRequired(lbConn, pdu)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			logger.NgapLog.Tracef("Handling RanConfigurationUpdate")
			go gnb_ngap.HandleRanConfigurationUpdate(lbConn, pdu)
		case ngapType.ProcedureCodeRRCInactiveTransitionReport:
			logger.NgapLog.Tracef("Handling RRCInactiveTransitionReport")
			go gnb_ngap.HandleRRCInactiveTransitionReport(lbConn, pdu)
		case ngapType.ProcedureCodePDUSessionResourceNotify:
			logger.NgapLog.Tracef("Handling PDUSessionResourceNotify")
			go gnb_ngap.HandlePDUSessionResourceNotify(lbConn, pdu)
		case ngapType.ProcedureCodePathSwitchRequest:
			logger.NgapLog.Tracef("Handling PathSwitchRequest")
			go gnb_ngap.HandlePathSwitchRequest(lbConn, pdu)
		case ngapType.ProcedureCodeLocationReport:
			logger.NgapLog.Tracef("Handling LocationReport")
			go gnb_ngap.HandleLocationReport(lbConn, pdu)
		case ngapType.ProcedureCodeUplinkUEAssociatedNRPPaTransport:
			logger.NgapLog.Tracef("Handling UplinkUEAssociatedNRPPATransport")
			go gnb_ngap.HandleUplinkUEAssociatedNRPPATransport(lbConn, pdu)
		case ngapType.ProcedureCodeUplinkRANConfigurationTransfer:
			logger.NgapLog.Tracef("Handling UplinkRanConfigurationTransfer")
			go gnb_ngap.HandleUplinkRanConfigurationTransfer(lbConn, pdu)
		case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
			logger.NgapLog.Tracef("Handling PDUSessionResourceModifyIndication")
			go gnb_ngap.HandlePDUSessionResourceModifyIndication(lbConn, pdu)
		case ngapType.ProcedureCodeCellTrafficTrace:
			logger.NgapLog.Tracef("Handling CellTrafficTrace")
			go gnb_ngap.HandleCellTrafficTrace(lbConn, pdu)
		case ngapType.ProcedureCodeUplinkRANStatusTransfer:
			logger.NgapLog.Tracef("Handling UplinkRanStatusTransfer")
			go gnb_ngap.HandleUplinkRanStatusTransfer(lbConn, pdu)
		case ngapType.ProcedureCodeUplinkNonUEAssociatedNRPPaTransport:
			logger.NgapLog.Tracef("Handling UplinkNonUEAssociatedNRPPATransport")
			go gnb_ngap.HandleUplinkNonUEAssociatedNRPPATransport(lbConn, pdu)
		default:
			lbConn.Log.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, initiatingMessage.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := pdu.SuccessfulOutcome
		if successfulOutcome == nil {
			lbConn.Log.Errorln("successful Outcome is nil")
			return
		}
		switch successfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGReset:
			logger.NgapLog.Tracef("Handling NGResetAcknowledge")
			go gnb_ngap.HandleNGResetAcknowledge(lbConn, pdu)
		case ngapType.ProcedureCodeUEContextRelease:
			logger.NgapLog.Tracef("Handling UEContextReleaseComplete")
			go gnb_ngap.HandleUEContextReleaseComplete(lbConn, pdu)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			logger.NgapLog.Tracef("Handling PDUSessionResourceReleaseResponse")
			go gnb_ngap.HandlePDUSessionResourceReleaseResponse(lbConn, pdu)
		case ngapType.ProcedureCodeUERadioCapabilityCheck:
			logger.NgapLog.Tracef("Handling UERadioCapabilityCheckResponse")
			go gnb_ngap.HandleUERadioCapabilityCheckResponse(lbConn, pdu)
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			logger.NgapLog.Tracef("Handling AMFconfigurationUpdateAcknowledge")
			go gnb_ngap.HandleAMFconfigurationUpdateAcknowledge(lbConn, pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			logger.NgapLog.Tracef("Handling InitialContextSetupResponse")
			go gnb_ngap.HandleInitialContextSetupResponse(lbConn, pdu)
		case ngapType.ProcedureCodeUEContextModification:
			logger.NgapLog.Tracef("Handling UEContextModificationResponse")
			go gnb_ngap.HandleUEContextModificationResponse(lbConn, pdu)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			logger.NgapLog.Tracef("Handling PDUSessionResourceSetupResponse")
			go gnb_ngap.HandlePDUSessionResourceSetupResponse(lbConn, pdu)
		case ngapType.ProcedureCodePDUSessionResourceModify:
			logger.NgapLog.Tracef("Handling PDUSessionResourceModifyResponse")
			go gnb_ngap.HandlePDUSessionResourceModifyResponse(lbConn, pdu)
		case ngapType.ProcedureCodeHandoverResourceAllocation:
			logger.NgapLog.Tracef("Handling HandoverRequestAcknowledge")
			go gnb_ngap.HandleHandoverRequestAcknowledge(lbConn, pdu)
		default:
			lbConn.Log.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, successfulOutcome.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := pdu.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			lbConn.Log.Errorf("unsuccessful Outcome is nil")
			return
		}
		switch unsuccessfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			logger.NgapLog.Tracef("Handling AMFconfigurationUpdateFailure")
			go gnb_ngap.HandleAMFconfigurationUpdateFailure(lbConn, pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			logger.NgapLog.Tracef("Handling InitialContextSetupFailure")
			go gnb_ngap.HandleInitialContextSetupFailure(lbConn, pdu)
		case ngapType.ProcedureCodeUEContextModification:
			logger.NgapLog.Tracef("Handling UEContextModificationFailure")
			go gnb_ngap.HandleUEContextModificationFailure(lbConn, pdu)
		case ngapType.ProcedureCodeHandoverResourceAllocation:
			logger.NgapLog.Tracef("Handling HandoverFailure")
			go gnb_ngap.HandleHandoverFailure(lbConn, pdu)
		default:
			lbConn.Log.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, unsuccessfulOutcome.ProcedureCode.Value)
		}
	}
}

// This handles messages incoming from AMF with the functions of the GNBs handler 
func DispatchForMessageToGnb(lbConn *context.LBConn, msg []byte) {
	// Decode
	pdu, err := ngap.Decoder(msg)
	if err != nil {
		logger.NgapLog.Errorf("NGAP decode error: %+v\n", err)
		return
	}

	switch pdu.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := pdu.InitiatingMessage
		if initiatingMessage == nil {
			logger.NgapLog.Errorf("Initiating Message is nil")
			return
		}

		switch initiatingMessage.ProcedureCode.Value {
		//case ngapType.ProcedureCodeNGReset:
		//	handler.HandleNGReset(amf, pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			logger.NgapLog.Tracef("Handling InitialContextSetupRequest")
			go amf_ngap.HandleInitialContextSetupRequest(lbConn, pdu)
		//case ngapType.ProcedureCodeUEContextModification:
		//	handler.HandleUEContextModificationRequest(amf, pdu)
		case ngapType.ProcedureCodeUEContextRelease:
			logger.NgapLog.Tracef("Handling UEContextReleaseCommand")
			go amf_ngap.HandleUEContextReleaseCommand(lbConn, pdu)
		case ngapType.ProcedureCodeDownlinkNASTransport:
			logger.NgapLog.Tracef("Handling DownlinkNASTransport")
			go amf_ngap.HandleDownlinkNASTransport(lbConn, pdu)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			logger.NgapLog.Tracef("Handling PDUSessionResourceSetupRequest")
			go amf_ngap.HandlePDUSessionResourceSetupRequest(lbConn, pdu)
		// TODO: This will be commented for the time being, after adding other procedures will be uncommented.
		//case ngapType.ProcedureCodePDUSessionResourceModify:
		//	handler.HandlePDUSessionResourceModifyRequest(amf, pdu)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			logger.NgapLog.Tracef("Handling PDUSessionResourceReleaseCommand")
			go amf_ngap.HandlePDUSessionResourceReleaseCommand(lbConn, pdu)
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
			go amf_ngap.HandleNGSetupResponse(lbConn, pdu)

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
			logger.NgapLog.Errorf("Unsuccessful Outcome is nil")
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
