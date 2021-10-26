package gnb_ngap

import (
	//"fmt"
	// "context"

	// "github.com/free5gc/aper"
	// "github.com/free5gc/fsm"
	"fmt"

	"github.com/LuckyG0ldfish/balancer/context"
	"github.com/free5gc/ngap/ngapType"
	"github.com/sirupsen/logrus"
	// "github.com/ishidawataru/sctp"
	// "gitlab.lrz.de/lkn_free5gc/gnbsim/context"
	// "gitlab.lrz.de/lkn_free5gc/gnbsim/gmm"
	// "gitlab.lrz.de/lkn_free5gc/gnbsim/logger"
	// gnb_nas "gitlab.lrz.de/lkn_free5gc/gnbsim/nas"
	// ngap_message "gitlab.lrz.de/lkn_free5gc/gnbsim/ngap/message"
	// RAN "gitlab.lrz.de/lkn_free5gc/gnbsim/util/ran_helper"
	// "time"
)

var NGAPLog *logrus.Entry

// func init() {
// 	NGAPLog = logger.NGAPLog
// }

func HandleNGSetupResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	// NGAPLog.Infoln("[gNB] Handle NG Setup Response")

	// var amfName *ngapType.AMFName
	// var servedGUAMIList *ngapType.ServedGUAMIList
	// var plmnSupportList *ngapType.PLMNSupportList

	// var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	// if message == nil {
	// 	NGAPLog.Error("NGAP Message is nil")
	// 	return
	// }

	// successfulOutcome := message.SuccessfulOutcome
	// if successfulOutcome == nil {
	// 	NGAPLog.Error("Successful Outcome is nil")
	// 	return
	// }

	// ngSetupResponse := successfulOutcome.Value.NGSetupResponse
	// if ngSetupResponse == nil {
	// 	NGAPLog.Error("ngSetupResponse is nil")
	// 	return
	// }

	// for _, ie := range ngSetupResponse.ProtocolIEs.List {
	// 	switch ie.Id.Value {
	// 	case ngapType.ProtocolIEIDAMFName:
	// 		NGAPLog.Traceln("[NGAP] Decode IE AMFName")
	// 		amfName = ie.Value.AMFName
	// 		if amfName == nil {
	// 			NGAPLog.Errorf("AMFName is nil")
	// 			item := buildCriticalityDiagnosticsIEItem(
	// 				ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
	// 			iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	// 		}
	// 	case ngapType.ProtocolIEIDServedGUAMIList:
	// 		NGAPLog.Traceln("[NGAP] Decode IE ServedGUAMIList")
	// 		servedGUAMIList = ie.Value.ServedGUAMIList
	// 		if servedGUAMIList == nil {
	// 			NGAPLog.Errorf("ServedGUAMIList is nil")
	// 			item := buildCriticalityDiagnosticsIEItem(
	// 				ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
	// 			iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	// 		}
	// 	case ngapType.ProtocolIEIDRelativeAMFCapacity:
	// 		NGAPLog.Traceln("[NGAP] Decode IE RelativeAMFCapacity")
	// 		//relativeAMFCapacity = ie.Value.RelativeAMFCapacity
	// 	case ngapType.ProtocolIEIDPLMNSupportList:
	// 		NGAPLog.Traceln("[NGAP] Decode IE PLMNSupportList")
	// 		plmnSupportList = ie.Value.PLMNSupportList
	// 		if plmnSupportList == nil {
	// 			NGAPLog.Errorf("PLMNSupportList is nil")
	// 			item := buildCriticalityDiagnosticsIEItem(
	// 				ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
	// 			iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	// 		}
	// 	case ngapType.ProtocolIEIDCriticalityDiagnostics:
	// 		NGAPLog.Traceln("[NGAP] Decode IE CriticalityDiagnostics")
	// 		//criticalityDiagnostics = ie.Value.CriticalityDiagnostics
	// 	}
	// }

	// if len(iesCriticalityDiagnostics.List) != 0 {
	// 	NGAPLog.Traceln("[NGAP] Sending error indication to AMF, because some mandatory IEs were not included")

	// 	cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentAbstractSyntaxErrorReject)

	// 	procedureCode := ngapType.ProcedureCodeNGSetup
	// 	triggeringMessage := ngapType.TriggeringMessagePresentSuccessfulOutcome
	// 	procedureCriticality := ngapType.CriticalityPresentReject

	// 	criticalityDiagnostics := buildCriticalityDiagnostics(
	// 		&procedureCode, &triggeringMessage, &procedureCriticality, &iesCriticalityDiagnostics)

	// 	ngap_message.SendErrorIndicationWithSctpConn(lbConn, nil, nil, cause, &criticalityDiagnostics)

	// 	return
	// }

	// // amfInfo := n3iwfSelf.NewN3iwfAmf(sctpAddr, conn)

}

// TODO 
func HandleInitialContextSetupRequest(lbConn *context.LBConn, message *ngapType.NGAPPDU) {

	LB := context.LB_Self()
	// TODO: add NAS Registration Complete Message (or this maybe will be done by the nas_handler).

	NGAPLog.Infoln("[gNB] Handle Initial Context Setup Request")

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	// var ueSecurityCapabilities *ngapType.UESecurityCapabilities
	// var securityKey *ngapType.SecurityKey
	// var traceActivation *ngapType.TraceActivation
	// var nasPDU *ngapType.NASPDU
	// var emergencyFallbackIndicator *ngapType.EmergencyFallbackIndicator
	// var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	// var ueCtx *RAN.UeContext
	// var emulatorCtx = context.EmulatorSelf()

	if message == nil {
		NGAPLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		NGAPLog.Error("Initiating Message is nil")
		return
	}

	initialContextSetupRequest := initiatingMessage.Value.InitialContextSetupRequest
	if initialContextSetupRequest == nil {
		NGAPLog.Error("InitialContextSetupRequest is nil")
		return
	}

	for _, ie := range initialContextSetupRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			NGAPLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			aMFUENGAPID = ie.Value.AMFUENGAPID
			if aMFUENGAPID == nil {
				NGAPLog.Errorf("AMFUENGAPID is nil")
				// item := buildCriticalityDiagnosticsIEItem(
				// 	ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				// iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			NGAPLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				NGAPLog.Errorf("RANUENGAPID is nil")
				// item := buildCriticalityDiagnosticsIEItem(
				// 	ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				// iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		// case ngapType.ProtocolIEIDUESecurityCapabilities:
		// 	NGAPLog.Traceln("[NGAP] Decode IE UESecurityCapabilities")
		// 	ueSecurityCapabilities = ie.Value.UESecurityCapabilities
		// 	if ueSecurityCapabilities == nil {
		// 		NGAPLog.Errorf("UESecurityCapabilities is nil")
		// 		item := buildCriticalityDiagnosticsIEItem(
		// 			ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
		// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
		// 	}
		// case ngapType.ProtocolIEIDSecurityKey:
		// 	NGAPLog.Traceln("[NGAP] Decode IE SecurityKey")
		// 	securityKey = ie.Value.SecurityKey
		// 	if securityKey == nil {
		// 		NGAPLog.Errorf("SecurityKey is nil")
		// 		item := buildCriticalityDiagnosticsIEItem(
		// 			ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
		// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
		// 	}
		// case ngapType.ProtocolIEIDTraceActivation:
		// 	NGAPLog.Traceln("[NGAP] Decode IE TraceActivation")
		// 	traceActivation = ie.Value.TraceActivation
		// 	if traceActivation != nil {
		// 		NGAPLog.Warnln("Not Supported IE [TraceActivation]")
		// 	}
		// case ngapType.ProtocolIEIDNASPDU:
		// 	NGAPLog.Traceln("[NGAP] Decode IE NAS PDU")
		// 	nasPDU = ie.Value.NASPDU
		// case ngapType.ProtocolIEIDEmergencyFallbackIndicator:
		// 	NGAPLog.Traceln("[NGAP] Decode IE EmergencyFallbackIndicator")
		// 	emergencyFallbackIndicator = ie.Value.EmergencyFallbackIndicator
		// 	if emergencyFallbackIndicator != nil {
		// 		NGAPLog.Warnln("Not Supported IE [EmergencyFallbackIndicator]")
		// 	}
		}
	}

	if lbConn.TypeID == context.TypeIdentGNBConn {
		gnb, _ := LB.LbGnbFindByConn(lbConn.Conn)
		UEs, ok := gnb.FindUeByUeRanID(rANUENGAPID.Value)
		var ue *context.LbUe
		if !ok {
			fmt.Println("Ue not of type UE/not found")
		} else {
			ue2, empty := context.FindUeInSlice(UEs, aMFUENGAPID.Value)
			switch empty{
			case 0: 
				fmt.Println("no UE Found")
				return
			case 1: 
				fmt.Println("UE Found")
				ue = ue2
			case 2: 
				fmt.Println("UE Found") // 
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToAmf(lbConn, m2, ue)
		return
	}
	// if (amfUeNgapID != nil) && (ranUeNgapID != nil) {
	// 	// Find UE context
	// 	var ok bool
	// 	ueCtx, ok = emulatorCtx.UePoolLoad(ranUeNgapID.Value)
	// 	if !ok {
	// 		NGAPLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", ranUeNgapID.Value)
	// 		// TODO: build cause and handle error
	// 		// Cause: Unknown local UE NGAP ID
	// 		return
	// 	} else {
	// 		if ueCtx.AmfUeNgapId != amfUeNgapID.Value {
	// 			// TODO: build cause and handle error
	// 			// Cause: Inconsistent remote UE NGAP ID
	// 			return
	// 		}
	// 	}
	// }

	// ueCtx.AmfUeNgapId = amfUeNgapID.Value
	// ueCtx.RanUeNgapId = ranUeNgapID.Value

	// ngap_message.SendInitialContextSetupResponse(lbConn, ueCtx.AmfUeNgapId, ueCtx.RanUeNgapId)

	// // TODO: The flow is as follows:
	// // send NAS Registration Complete Msg
	// // send NAS Deregistration Request (UE Originating)
	// // TODO: Here I need a handler that will send the Registration Accept to the NAS handler, which in turn will put the UE to the Registered State.
	// if nasPDU != nil {
	// 	// TODO: Handle NAS Packet
	// 	go gnb_nas.HandleNAS(ueCtx, ngapType.ProcedureCodeDownlinkNASTransport, nasPDU.Value)
	// }
}

func HandleUEContextReleaseCommand(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	// NGAPLog.Infoln("[gNB] Handle UE Context Release Command")

	// var ueNgapIDs *ngapType.UENGAPIDs
	// var cause *ngapType.Cause
	// var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	// var ueCtx *RAN.UeContext
	// var lbCtx = context.LB_Self()

	// if message == nil {
	// 	NGAPLog.Error("NGAP Message is nil")
	// 	return
	// }

	// initiatingMessage := message.InitiatingMessage
	// if initiatingMessage == nil {
	// 	NGAPLog.Error("Initiating Message is nil")
	// 	return
	// }

	// ueContextReleaseCommand := initiatingMessage.Value.UEContextReleaseCommand
	// if ueContextReleaseCommand == nil {
	// 	NGAPLog.Error("UEContextReleaseCommand is nil")
	// 	return
	// }

	// for _, ie := range ueContextReleaseCommand.ProtocolIEs.List {
	// 	switch ie.Id.Value {
	// 	case ngapType.ProtocolIEIDUENGAPIDs:
	// 		NGAPLog.Traceln("[NGAP] Decode IE UENGAPIDs")
	// 		ueNgapIDs = ie.Value.UENGAPIDs
	// 		if ueNgapIDs == nil {
	// 			NGAPLog.Errorf("UENGAPIDs is nil")
	// 			item := buildCriticalityDiagnosticsIEItem(
	// 				ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
	// 			iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	// 		}
	// 	case ngapType.ProtocolIEIDCause:
	// 		NGAPLog.Traceln("[NGAP] Decode IE Cause")
	// 		cause = ie.Value.Cause
	// 	}
	// }

	// if len(iesCriticalityDiagnostics.List) > 0 {
	// 	// TODO: send error indication
	// 	return
	// }

	// switch ueNgapIDs.Present {
	// case ngapType.UENGAPIDsPresentUENGAPIDPair:
	// 	ueCtx, _ = lbCtx.UePoolLoad(ueNgapIDs.UENGAPIDPair.RANUENGAPID.Value)
	// case ngapType.UENGAPIDsPresentAMFUENGAPID:
	// 	// TODO: find UE according to specific AMF
	// 	// The implementation here may have error when N3IWF need to
	// 	// connect multiple AMFs.
	// 	// Use UEpool in AMF context can solve this problem
	// 	// ueCtx = amf.FindUeByAmfUeNgapID(ueNgapIDs.AMFUENGAPID.Value)
	// }

	// if ueCtx == nil {
	// 	// TODO: send error indication(unknown local ngap ue id)
	// 	return
	// }

	// if cause != nil {
	// 	printAndGetCause(cause)
	// }

	// ngap_message.SendUEContextReleaseComplete(lbConn, ueCtx.AmfUeNgapId, ueCtx.RanUeNgapId)

	// ueCtx.DeregisrtationFinished = true
	// ueCtx.TimestampT8 = int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)


	// if err := gmm.GmmFSM.SendEvent(ueCtx.State, gmm.DeregistrationAcceptEvent, fsm.ArgsType{
	// 	gmm.ArgRanUe:         ueCtx,
	// }); err != nil {
	// 	logger.GmmLog.Errorln(err)
	// }
}

// TODO 
func HandleDownlinkNASTransport(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	NGAPLog.Infoln("[gNB] Handle Downlink NAS Transport")

	var amfUeNgapID *ngapType.AMFUENGAPID
	var ranUeNgapID *ngapType.RANUENGAPID
	// var nasPDU *ngapType.NASPDU
	// var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	// var ueCtx *RAN.UeContext
	// var emulatorCtx = context.EmulatorSelf()

	if message == nil {
		NGAPLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		NGAPLog.Error("Initiating Message is nil")
		return
	}

	downlinkNASTransport := initiatingMessage.Value.DownlinkNASTransport
	if downlinkNASTransport == nil {
		NGAPLog.Error("DownlinkNASTransport is nil")
		return
	}

	for _, ie := range downlinkNASTransport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			NGAPLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			amfUeNgapID = ie.Value.AMFUENGAPID
			if amfUeNgapID == nil {
				NGAPLog.Errorf("AMFUENGAPID is nil")
				// item := buildCriticalityDiagnosticsIEItem(
				// 	ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				// iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			NGAPLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				NGAPLog.Errorf("RANUENGAPID is nil")
				// item := buildCriticalityDiagnosticsIEItem(
				// 	ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				// iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		// case ngapType.ProtocolIEIDNASPDU:
		// 	NGAPLog.Traceln("[NGAP] Decode IE NASPDU")
		// 	nasPDU = ie.Value.NASPDU
		// 	if nasPDU == nil {
		// 		NGAPLog.Errorf("NASPDU is nil")
		// 		item := buildCriticalityDiagnosticsIEItem(
		// 			ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
		// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
		// 	}
		}
	}


	// if ranUeNgapID != nil {
	// 	var ok bool
	// 	ueCtx, ok = emulatorCtx.UePoolLoad(ranUeNgapID.Value)
	// 	//fmt.Println(ueCtx)
	// 	if !ok {
	// 		NGAPLog.Warnf("No UE Context[RanUeNgapID:%d]\n", ranUeNgapID.Value)
	// 		return
	// 	}
	// }

	// if amfUeNgapID != nil {
	// 	if ueCtx.AmfUeNgapId == 0 {
	// 		NGAPLog.Tracef("Create new logical UE-associated NG-connection")
	// 		ueCtx.AmfUeNgapId = amfUeNgapID.Value
	// 	} else {
	// 		if ueCtx.AmfUeNgapId != amfUeNgapID.Value {
	// 			NGAPLog.Warn("AMFUENGAPID unmatched")
	// 			return
	// 		}
	// 	}
	// }

	// if nasPDU != nil {
	// 	// TODO: Handle NAS Packet
	// 	go gnb_nas.HandleNAS(ueCtx, ngapType.ProcedureCodeDownlinkNASTransport, nasPDU.Value)
	// }
}

// TODO 
func HandlePDUSessionResourceSetupRequest(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	NGAPLog.Infoln("[gNB] Handle PDU Session Resource Setup Request")

	var amfUeNgapID *ngapType.AMFUENGAPID
	var ranUeNgapID *ngapType.RANUENGAPID
	// var nasPDU *ngapType.NASPDU
	// var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList
	// var pduSessionResourceSetupListSUReq *ngapType.PDUSessionResourceSetupListSUReq

	// var ueCtx *RAN.UeContext
	// var emulatorCtx = context.EmulatorSelf()

	if message == nil {
		NGAPLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		NGAPLog.Error("Initiating Message is nil")
		return
	}

	pduSessionResourceSetupRequest := initiatingMessage.Value.PDUSessionResourceSetupRequest
	if pduSessionResourceSetupRequest == nil {
		NGAPLog.Error("PDUSessionResourceSetupRequest is nil")
		return
	}

	for _, ie := range pduSessionResourceSetupRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			NGAPLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			amfUeNgapID = ie.Value.AMFUENGAPID
			if amfUeNgapID == nil {
				NGAPLog.Errorf("AMFUENGAPID is nil")
				// item := buildCriticalityDiagnosticsIEItem(
				// 	ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				// iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			NGAPLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				NGAPLog.Errorf("RANUENGAPID is nil")
				// item := buildCriticalityDiagnosticsIEItem(
				// 	ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				// iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		// case ngapType.ProtocolIEIDNASPDU:
		// 	NGAPLog.Traceln("[NGAP] Decode IE NASPDU")
		// 	nasPDU = ie.Value.NASPDU
		// 	if nasPDU == nil {
		// 		NGAPLog.Errorf("NASPDU is nil")
		// 		item := buildCriticalityDiagnosticsIEItem(
		// 			ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
		// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
		// 	}
		// case ngapType.ProtocolIEIDPDUSessionResourceSetupListSUReq:
		// 	NGAPLog.Traceln("[NGAP] Decode IE PDUSessionResourceSetupRequestList")
		// 	pduSessionResourceSetupListSUReq = ie.Value.PDUSessionResourceSetupListSUReq
		// 	if pduSessionResourceSetupListSUReq == nil {
		// 		NGAPLog.Errorf("PDUSessionResourceSetupRequestList is nil")
		// 		item := buildCriticalityDiagnosticsIEItem(
		// 			ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
		// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
		// 	}
		}
	}

	// if len(iesCriticalityDiagnostics.List) > 0 {
	// 	// TODO: Send error indication to AMF
	// 	NGAPLog.Errorln("Sending error indication to AMF")
	// 	return
	// }

	// if (amfUeNgapID != nil) && (ranUeNgapID != nil) {
	// 	// Find UE context
	// 	var ok bool
	// 	ueCtx, ok = emulatorCtx.UePoolLoad(ranUeNgapID.Value)
	// 	if !ok {
	// 		NGAPLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", ranUeNgapID.Value)
	// 		// TODO: build cause and handle error
	// 		// Cause: Unknown local UE NGAP ID
	// 		return
	// 	} else {
	// 		if ueCtx.AmfUeNgapId != amfUeNgapID.Value {
	// 			// TODO: build cause and handle error
	// 			NGAPLog.Warn("AMFUENGAPID unmatched")
	// 			return
	// 		}
	// 	}
	// }
	// ngap_message.SendPDUSessionResourceSetupResponse(lbConn, ueCtx.AmfUeNgapId, ueCtx.RanUeNgapId, ueCtx.PduSessionID, emulatorCtx.GNBIPAddress)

	// if pduSessionResourceSetupListSUReq != nil {
	// 	for _, item := range pduSessionResourceSetupListSUReq.List {
	// 		//pduSessionID := item.PDUSessionID.Value
	// 		pduSessionNasPdu := item.PDUSessionNASPDU.Value
	// 		//snssai := item.SNSSAI
	// 		// TODO: procedure Code may need to be changed...
	// 		dummyProcedureCode := int64(16)
	// 		go gnb_nas.HandleNAS(ueCtx, dummyProcedureCode, pduSessionNasPdu)
	// 	}
	// }
}

// TODO
func HandlePDUSessionResourceReleaseCommand(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	NGAPLog.Infoln("[gNB] Handle PDU Session Resource Release Command")
	var amfUeNgapID *ngapType.AMFUENGAPID
	var ranUeNgapID *ngapType.RANUENGAPID
	// var nasPDU *ngapType.NASPDU
	// var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList
	// var pDUSessionResourceToReleaseListRelCmd *ngapType.PDUSessionResourceToReleaseListRelCmd

	// var ueCtx *RAN.UeContext
	// var emulatorCtx = context.EmulatorSelf()

	if message == nil {
		NGAPLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		NGAPLog.Error("Initiating Message is nil")
		return
	}

	pDUSessionResourceReleaseCommand := initiatingMessage.Value.PDUSessionResourceReleaseCommand
	if pDUSessionResourceReleaseCommand == nil {
		NGAPLog.Error("pDUSessionResourceReleaseCommand is nil")
		return
	}

	for _, ie := range pDUSessionResourceReleaseCommand.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			NGAPLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			amfUeNgapID = ie.Value.AMFUENGAPID
			if amfUeNgapID == nil {
				NGAPLog.Error("AMFUENGAPID is nil")
				// item := buildCriticalityDiagnosticsIEItem(
				// 	ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				// iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			NGAPLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			ranUeNgapID = ie.Value.RANUENGAPID
			if ranUeNgapID == nil {
				NGAPLog.Error("RANUENGAPID is nil")
				// item := buildCriticalityDiagnosticsIEItem(
				// 	ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				// iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		// case ngapType.ProtocolIEIDRANPagingPriority:
		// 	NGAPLog.Traceln("[NGAP] Decode IE RANPagingPriority")
		// 	// rANPagingPriority = ie.Value.RANPagingPriority
		// case ngapType.ProtocolIEIDNASPDU:
		// 	NGAPLog.Traceln("[NGAP] Decode IE NASPDU")
		// 	nasPDU = ie.Value.NASPDU
		// 	if nasPDU == nil {
		// 		NGAPLog.Errorf("NASPDU is nil")
		// 		item := buildCriticalityDiagnosticsIEItem(
		// 			ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
		// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
		// 	}
		// case ngapType.ProtocolIEIDPDUSessionResourceToReleaseListRelCmd:
		// 	NGAPLog.Traceln("[NGAP] Decode IE PDUSessionResourceToReleaseListRelCmd")
		// 	pDUSessionResourceToReleaseListRelCmd = ie.Value.PDUSessionResourceToReleaseListRelCmd
		// 	if pDUSessionResourceToReleaseListRelCmd == nil {
		// 		NGAPLog.Error("PDUSessionResourceToReleaseListRelCmd is nil")
		// 		item := buildCriticalityDiagnosticsIEItem(
		// 			ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
		// 		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
		// 	}
		}
	}

	// if len(iesCriticalityDiagnostics.List) > 0 {
	// 	procudureCode := ngapType.ProcedureCodePDUSessionResourceRelease
	// 	trigger := ngapType.TriggeringMessagePresentInitiatingMessage
	// 	criticality := ngapType.CriticalityPresentReject
	// 	criticalityDiagnostics := buildCriticalityDiagnostics(
	// 		&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
	// 	ngap_message.SendErrorIndicationWithSctpConn(lbConn, nil, nil, nil, &criticalityDiagnostics)
	// 	return
	// }

	// ueCtx, ok := emulatorCtx.UePoolLoad(ranUeNgapID.Value)
	// if !ok {
	// 	NGAPLog.Errorf("Unknown local UE NGAP ID. RanUENGAPID: %d", ranUeNgapID.Value)
	// 	cause := buildCause(ngapType.CausePresentRadioNetwork, ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID)
	// 	ngap_message.SendErrorIndicationWithSctpConn(lbConn, nil, nil, cause, nil)
	// 	return
	// }

	// if ueCtx.AmfUeNgapId != amfUeNgapID.Value {
	// 	NGAPLog.Errorf("Inconsistent remote UE NGAP ID, AMFUENGAPID: %d, ue.AmfUeNgapId: %d",
	// 		amfUeNgapID.Value, ueCtx.AmfUeNgapId)
	// 	cause := buildCause(ngapType.CausePresentRadioNetwork,
	// 		ngapType.CauseRadioNetworkPresentInconsistentRemoteUENGAPID)
	// 	ngap_message.SendErrorIndicationWithSctpConn(lbConn, nil, &ranUeNgapID.Value, cause, nil)
	// 	return
	// }

	// for _, item := range pDUSessionResourceToReleaseListRelCmd.List {
	// 	pduSessionId := item.PDUSessionID.Value
	// 	transfer := ngapType.PDUSessionResourceReleaseCommandTransfer{}
	// 	err := aper.UnmarshalWithParams(item.PDUSessionResourceReleaseCommandTransfer, &transfer, "valueExt")
	// 	if err != nil {
	// 		NGAPLog.Warnf("[PDUSessionID: %d] PDUSessionResourceReleaseCommandTransfer Decode Error: %+v\n", pduSessionId, err)
	// 	} else {
	// 		printAndGetCause(&transfer.Cause)
	// 	}
	// 	NGAPLog.Tracef("Release PDU Session Id[%d] due to PDU Session Resource Release Command", pduSessionId)
	// 	//delete(ueCtx.PduSessionList, pduSessionId)
	// }

	// if nasPDU != nil {
	// 	dummyProcedureCore := ngapType.ProcedureCodeDownlinkNASTransport
	// 	go gnb_nas.HandleNAS(ueCtx, dummyProcedureCore, nasPDU.Value)
	// }
	// ngap_message.SendPDUSessionResourceReleaseResponse(lbConn, ueCtx.AmfUeNgapId, ueCtx.RanUeNgapId, ueCtx.PduSessionID)
}

// func buildCriticalityDiagnosticsIEItem(ieCriticality aper.Enumerated, ieID int64, typeOfErr aper.Enumerated) (
// 	item ngapType.CriticalityDiagnosticsIEItem) {

// 	item = ngapType.CriticalityDiagnosticsIEItem{
// 		IECriticality: ngapType.Criticality{
// 			Value: ieCriticality,
// 		},
// 		IEID: ngapType.ProtocolIEID{
// 			Value: ieID,
// 		},
// 		TypeOfError: ngapType.TypeOfError{
// 			Value: typeOfErr,
// 		},
// 	}

// 	return item
// }

// func buildCause(present int, value aper.Enumerated) (cause *ngapType.Cause) {
// 	cause = new(ngapType.Cause)
// 	cause.Present = present

// 	switch present {
// 	case ngapType.CausePresentRadioNetwork:
// 		cause.RadioNetwork = new(ngapType.CauseRadioNetwork)
// 		cause.RadioNetwork.Value = value
// 	case ngapType.CausePresentTransport:
// 		cause.Transport = new(ngapType.CauseTransport)
// 		cause.Transport.Value = value
// 	case ngapType.CausePresentNas:
// 		cause.Nas = new(ngapType.CauseNas)
// 		cause.Nas.Value = value
// 	case ngapType.CausePresentProtocol:
// 		cause.Protocol = new(ngapType.CauseProtocol)
// 		cause.Protocol.Value = value
// 	case ngapType.CausePresentMisc:
// 		cause.Misc = new(ngapType.CauseMisc)
// 		cause.Misc.Value = value
// 	case ngapType.CausePresentNothing:
// 	}

// 	return
// }

// func buildCriticalityDiagnostics(
// 	procedureCode *int64,
// 	triggeringMessage *aper.Enumerated,
// 	procedureCriticality *aper.Enumerated,
// 	iesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList) (
// 	criticalityDiagnostics ngapType.CriticalityDiagnostics) {

// 	if procedureCode != nil {
// 		criticalityDiagnostics.ProcedureCode = new(ngapType.ProcedureCode)
// 		criticalityDiagnostics.ProcedureCode.Value = *procedureCode
// 	}

// 	if triggeringMessage != nil {
// 		criticalityDiagnostics.TriggeringMessage = new(ngapType.TriggeringMessage)
// 		criticalityDiagnostics.TriggeringMessage.Value = *triggeringMessage
// 	}

// 	if procedureCriticality != nil {
// 		criticalityDiagnostics.ProcedureCriticality = new(ngapType.Criticality)
// 		criticalityDiagnostics.ProcedureCriticality.Value = *procedureCriticality
// 	}

// 	if iesCriticalityDiagnostics != nil {
// 		criticalityDiagnostics.IEsCriticalityDiagnostics = iesCriticalityDiagnostics
// 	}

// 	return criticalityDiagnostics
// }


// func printAndGetCause(cause *ngapType.Cause) (present int, value aper.Enumerated) {

// 	present = cause.Present
// 	switch cause.Present {
// 	case ngapType.CausePresentRadioNetwork:
// 		NGAPLog.Warnf("Cause RadioNetwork[%d]", cause.RadioNetwork.Value)
// 		value = cause.RadioNetwork.Value
// 	case ngapType.CausePresentTransport:
// 		NGAPLog.Warnf("Cause Transport[%d]", cause.Transport.Value)
// 		value = cause.Transport.Value
// 	case ngapType.CausePresentProtocol:
// 		NGAPLog.Warnf("Cause Protocol[%d]", cause.Protocol.Value)
// 		value = cause.Protocol.Value
// 	case ngapType.CausePresentNas:
// 		NGAPLog.Warnf("Cause Nas[%d]", cause.Nas.Value)
// 		value = cause.Nas.Value
// 	case ngapType.CausePresentMisc:
// 		NGAPLog.Warnf("Cause Misc[%d]", cause.Misc.Value)
// 		value = cause.Misc.Value
// 	default:
// 		NGAPLog.Errorf("Invalid Cause group[%d]", cause.Present)
// 	}
// 	return
// }