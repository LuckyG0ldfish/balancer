package gnb_ngap

// This handles messages incoming from GNB with the functions of the AMFs handler

import (
	"strconv"

	"github.com/LuckyG0ldfish/balancer/context"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/LuckyG0ldfish/balancer/nas"
	"github.com/free5gc/aper"
	"github.com/free5gc/ngap/ngapType"

	ngap_message "github.com/LuckyG0ldfish/balancer/ngap/message"
)

var LB context.LBContext

//TODO
func HandleNGSetupRequest(LbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var globalRANNodeID *ngapType.GlobalRANNodeID
	var rANNodeName *ngapType.RANNodeName
	var supportedTAList *ngapType.SupportedTAList
	var pagingDRX *ngapType.PagingDRX

	LB = *context.LB_Self()
	var cause ngapType.Cause

	if LbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		LbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		LbConn.Log.Errorf("Initiating Message is nil")
		return
	}
	nGSetupRequest := initiatingMessage.Value.NGSetupRequest
	if nGSetupRequest == nil {
		LbConn.Log.Errorf("NGSetupRequest is nil")
		return
	}
	LbConn.Log.Infoln("Handle NG Setup request")

	for i := 0; i < len(nGSetupRequest.ProtocolIEs.List); i++ {
		ie := nGSetupRequest.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDGlobalRANNodeID:
			globalRANNodeID = ie.Value.GlobalRANNodeID
			LbConn.Log.Traceln("Decode IE GlobalRANNodeID")
			if globalRANNodeID == nil {
				LbConn.Log.Errorf("GlobalRANNodeID is nil")
				return
			}
		case ngapType.ProtocolIEIDSupportedTAList:
			supportedTAList = ie.Value.SupportedTAList
			LbConn.Log.Traceln("Decode IE SupportedTAList")
			if supportedTAList == nil {
				LbConn.Log.Errorf("SupportedTAList is nil")
				return
			}
		case ngapType.ProtocolIEIDRANNodeName:
			rANNodeName = ie.Value.RANNodeName
			LbConn.Log.Traceln("Decode IE RANNodeName")
			if rANNodeName == nil {
				LbConn.Log.Errorf("RANNodeName is nil")
				return
			}
		case ngapType.ProtocolIEIDDefaultPagingDRX:
			pagingDRX = ie.Value.DefaultPagingDRX
			LbConn.Log.Traceln("Decode IE DefaultPagingDRX")
			if pagingDRX == nil {
				LbConn.Log.Errorf("DefaultPagingDRX is nil")
				return
			}
		}
	}

	if cause.Present == ngapType.CausePresentNothing {
		ngap_message.SendNGSetupResponse(LbConn)
	} else {
		ngap_message.SendNGSetupFailure(LbConn, cause)
	}
}

func HandleUplinkNasTransport(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var nASPDU *ngapType.NASPDU
	var ue *context.LbUe

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}

	uplinkNasTransport := initiatingMessage.Value.UplinkNASTransport
	if uplinkNasTransport == nil {
		lbConn.Log.Errorf("UplinkNasTransport is nil")
		return
	}

	lbConn.Log.Infoln("Handle Uplink Nas Transport")

	for i := 0; i < len(uplinkNasTransport.ProtocolIEs.List); i++ {
		ie := uplinkNasTransport.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				var ok bool
				ue, ok = gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
			}
		case ngapType.ProtocolIEIDNASPDU:
			nASPDU = ie.Value.NASPDU
			
		}
	}
	if nASPDU != nil {
		nas.HandleNAS(ue, nASPDU.Value)
	}
	if ue != nil {
		context.ForwardToAmf(message, ue)
	}
}

// TODO
func HandleNGReset(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var cause *ngapType.Cause
	var resetType *ngapType.ResetType

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}
	nGReset := initiatingMessage.Value.NGReset
	if nGReset == nil {
		lbConn.Log.Errorf("NGReset is nil")
		return
	}

	lbConn.Log.Infoln("Handle NG Reset")

	for _, ie := range nGReset.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			lbConn.Log.Traceln("Decode IE Cause")
			if cause == nil {
				lbConn.Log.Errorf("Cause is nil")
				return
			}
		case ngapType.ProtocolIEIDResetType:
			resetType = ie.Value.ResetType
			lbConn.Log.Traceln("Decode IE ResetType")
			if resetType == nil {
				lbConn.Log.Errorf("ResetType is nil")
				return
			}
		}
	}

	printAndGetCause(lbConn, cause)

	switch resetType.Present {
	case ngapType.ResetTypePresentNGInterface:
		lbConn.Log.Traceln("ResetType Present: NG Interface")
		// lbConn.RemoveAllUeInRan()
		// ngap_message.SendNGResetAcknowledge(lbConn, nil, nil)
	case ngapType.ResetTypePresentPartOfNGInterface:
		lbConn.Log.Traceln("ResetType Present: Part of NG Interface")
		partOfNGInterface := resetType.PartOfNGInterface
		if partOfNGInterface == nil {
			lbConn.Log.Errorf("PartOfNGInterface is nil")
			return
		}
	}
}

// TODO
func HandleNGResetAcknowledge(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var uEAssociatedLogicalNGConnectionList *ngapType.UEAssociatedLogicalNGConnectionList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Errorf("SuccessfulOutcome is nil")
		return
	}
	nGResetAcknowledge := successfulOutcome.Value.NGResetAcknowledge
	if nGResetAcknowledge == nil {
		lbConn.Log.Errorf("NGResetAcknowledge is nil")
		return
	}

	lbConn.Log.Infoln("Handle NG Reset Acknowledge")

	for _, ie := range nGResetAcknowledge.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDUEAssociatedLogicalNGConnectionList:
			uEAssociatedLogicalNGConnectionList = ie.Value.UEAssociatedLogicalNGConnectionList
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
		}
	}

	if uEAssociatedLogicalNGConnectionList != nil {
		lbConn.Log.Traceln("%d UE association(s) has been reset", len(uEAssociatedLogicalNGConnectionList.List))
		for _, item := range uEAssociatedLogicalNGConnectionList.List {
			if item.AMFUENGAPID != nil && item.RANUENGAPID != nil {
				// lbConn.Log.Traceln("%d: AmfUeNgapID[%d] RanUeNgapID[%d]", i+1, item.AMFUENGAPID.Value, item.RANUENGAPID.Value)
			} else if item.AMFUENGAPID != nil {
				// lbConn.Log.Traceln("%d: AmfUeNgapID[%d] RanUeNgapID[-1]", i+1, item.AMFUENGAPID.Value)
			} else if item.RANUENGAPID != nil {
				// lbConn.Log.Traceln("%d: AmfUeNgapID[-1] RanUeNgapID[%d]", i+1, item.RANUENGAPID.Value)
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

func HandleUEContextReleaseComplete(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Errorf("SuccessfulOutcome is nil")
		return
	}
	uEContextReleaseComplete := successfulOutcome.Value.UEContextReleaseComplete
	if uEContextReleaseComplete == nil {
		lbConn.Log.Errorf("NGResetAcknowledge is nil")
		return
	}

	lbConn.Log.Infoln("Handle UE Context Release Complete")

	for _, ie := range uEContextReleaseComplete.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
				ue.RemoveUeEntirely()
			}
		}
	}
}

func HandlePDUSessionResourceReleaseResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Errorf("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceReleaseResponse := successfulOutcome.Value.PDUSessionResourceReleaseResponse
	if pDUSessionResourceReleaseResponse == nil {
		lbConn.Log.Errorf("PDUSessionResourceReleaseResponse is nil")
		return
	}

	lbConn.Log.Infoln("Handle PDU Session Resource Release Response")

	for _, ie := range pDUSessionResourceReleaseResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandleUERadioCapabilityCheckResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Errorf("SuccessfulOutcome is nil")
		return
	}

	uERadioCapabilityCheckResponse := successfulOutcome.Value.UERadioCapabilityCheckResponse
	if uERadioCapabilityCheckResponse == nil {
		lbConn.Log.Errorf("UERadioCapabilityCheckResponse is nil")
		return
	}
	lbConn.Log.Infoln("Handle UE Radio Capability Check Response")

	for i := 0; i < len(uERadioCapabilityCheckResponse.ProtocolIEs.List); i++ {
		ie := uERadioCapabilityCheckResponse.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandleLocationReportingFailureIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}
	locationReportingFailureIndication := initiatingMessage.Value.LocationReportingFailureIndication
	if locationReportingFailureIndication == nil {
		lbConn.Log.Errorf("LocationReportingFailureIndication is nil")
		return
	}

	lbConn.Log.Infoln("Handle Location Reporting Failure Indication")

	for i := 0; i < len(locationReportingFailureIndication.ProtocolIEs.List); i++ {
		ie := locationReportingFailureIndication.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

//
func HandleInitialUEMessage(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var rANUENGAPID *ngapType.RANUENGAPID
	var nASPDU *ngapType.NASPDU
	var rRCEstablishmentCause *ngapType.RRCEstablishmentCause

	LB = *context.LB_Self()

	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}
	initialUEMessage := initiatingMessage.Value.InitialUEMessage
	if initialUEMessage == nil {
		lbConn.Log.Errorf("InitialUEMessage is nil")
		return
	}

	lbConn.Log.Infoln("Handle Initial UE Message")

	UeLbID := LB.IDGen.NextNumber()
	var rANUENGAPIDInt int64

	for _, ie := range initialUEMessage.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt = ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("InitialUEMessage: rANUENGAPID == nil")
				return
			}
			ie.Value.RANUENGAPID.Value = UeLbID
		case ngapType.ProtocolIEIDNASPDU: // reject
			nASPDU = ie.Value.NASPDU
			lbConn.Log.Traceln("Decode IE NasPdu")
			if nASPDU == nil {
				lbConn.Log.Errorf("InitialUEMessage: nASPDU == nil")
			}
		case ngapType.ProtocolIEIDRRCEstablishmentCause: // ignore
			rRCEstablishmentCause = ie.Value.RRCEstablishmentCause
			lbConn.Log.Traceln("Decode IE RRCEstablishmentCause")
		}
	}

	if lbConn.TypeID == context.TypeIdGNBConn {
		gnb := lbConn.RanPointer
		ue := context.NewUE()
		ue.UeRanID = rANUENGAPIDInt
		ue.UeLbID = UeLbID
		ue.RanID = gnb.GnbID
		if rRCEstablishmentCause != nil {
			logger.NgapLog.Tracef("[Initial UE Message] RRC Establishment Cause[%d]", rRCEstablishmentCause.Value)
			ue.RRCECause = strconv.Itoa(int(rRCEstablishmentCause.Value))
		} else {
			ue.RRCECause = "0" // TODO: RRCEstablishmentCause 0 is for emergency service
		}
		gnb.Ues.Store(rANUENGAPIDInt, ue)
		ue.RanPointer = gnb
		context.ForwardToNextAmf(lbConn, message, ue)
		lbConn.Log.Traceln("UeRanID: " + strconv.FormatInt(rANUENGAPIDInt, 10))
		return
	}
}

func HandlePDUSessionResourceSetupResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Errorf("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceSetupResponse
	if pDUSessionResourceSetupResponse == nil {
		lbConn.Log.Errorf("PDUSessionResourceSetupResponse is nil")
		return
	}

	lbConn.Log.Infoln("Handle PDU Session Resource Setup Response")

	for _, ie := range pDUSessionResourceSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandlePDUSessionResourceModifyResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Errorf("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceModifyResponse := successfulOutcome.Value.PDUSessionResourceModifyResponse
	if pDUSessionResourceModifyResponse == nil {
		lbConn.Log.Errorf("PDUSessionResourceModifyResponse is nil")
		return
	}

	lbConn.Log.Infoln("Handle PDU Session Resource Modify Response")

	for _, ie := range pDUSessionResourceModifyResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandlePDUSessionResourceNotify(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("InitiatingMessage is nil")
		return
	}
	PDUSessionResourceNotify := initiatingMessage.Value.PDUSessionResourceNotify
	if PDUSessionResourceNotify == nil {
		lbConn.Log.Errorf("PDUSessionResourceNotify is nil")
		return
	}

	lbConn.Log.Infoln("Handle PDU Session Resource Notify")

	for _, ie := range PDUSessionResourceNotify.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandlePDUSessionResourceModifyIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // reject
	if initiatingMessage == nil {
		lbConn.Log.Errorf("InitiatingMessage is nil")
		// cause := ngapType.Cause{
		// 	Present: ngapType.CausePresentProtocol,
		// 	Protocol: &ngapType.CauseProtocol{
		// 		Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorReject,
		// 	},
		// }
		// ngap_message.SendErrorIndication(lbConn, nil, nil, &cause, nil)
		return
	}
	pDUSessionResourceModifyIndication := initiatingMessage.Value.PDUSessionResourceModifyIndication
	if pDUSessionResourceModifyIndication == nil {
		lbConn.Log.Errorf("PDUSessionResourceModifyIndication is nil")
		// cause := ngapType.Cause{
		// 	Present: ngapType.CausePresentProtocol,
		// 	Protocol: &ngapType.CauseProtocol{
		// 		Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorReject,
		// 	},
		// }
		// ngap_message.SendErrorIndication(lbConn, nil, nil, &cause, nil)
		return
	}

	lbConn.Log.Infoln("Handle PDU Session Resource Modify Indication")

	for _, ie := range pDUSessionResourceModifyIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandleInitialContextSetupResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Errorf("SuccessfulOutcome is nil")
		return
	}
	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	if initialContextSetupResponse == nil {
		lbConn.Log.Errorf("InitialContextSetupResponse is nil")
		return
	}

	lbConn.Log.Infoln("Handle Initial Context Setup Response")

	for _, ie := range initialContextSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandleInitialContextSetupFailure(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		lbConn.Log.Errorf("UnsuccessfulOutcome is nil")
		return
	}
	initialContextSetupFailure := unsuccessfulOutcome.Value.InitialContextSetupFailure
	if initialContextSetupFailure == nil {
		lbConn.Log.Errorf("InitialContextSetupFailure is nil")
		return
	}

	lbConn.Log.Infoln("Handle Initial Context Setup Failure")

	for _, ie := range initialContextSetupFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandleUEContextReleaseRequest(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("InitiatingMessage is nil")
		return
	}
	uEContextReleaseRequest := initiatingMessage.Value.UEContextReleaseRequest
	if uEContextReleaseRequest == nil {
		lbConn.Log.Errorf("UEContextReleaseRequest is nil")
		return
	}

	lbConn.Log.Infoln("UE Context Release Request")

	for _, ie := range uEContextReleaseRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandleUEContextModificationResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Errorf("SuccessfulOutcome is nil")
		return
	}
	uEContextModificationResponse := successfulOutcome.Value.UEContextModificationResponse
	if uEContextModificationResponse == nil {
		lbConn.Log.Errorf("UEContextModificationResponse is nil")
		return
	}

	lbConn.Log.Infoln("Handle UE Context Modification Response")

	for _, ie := range uEContextModificationResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandleUEContextModificationFailure(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		lbConn.Log.Errorf("UnsuccessfulOutcome is nil")
		return
	}
	uEContextModificationFailure := unsuccessfulOutcome.Value.UEContextModificationFailure
	if uEContextModificationFailure == nil {
		lbConn.Log.Errorf("UEContextModificationFailure is nil")
		return
	}

	lbConn.Log.Infoln("Handle UE Context Modification Failure")

	for _, ie := range uEContextModificationFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandleRRCInactiveTransitionReport(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}

	rRCInactiveTransitionReport := initiatingMessage.Value.RRCInactiveTransitionReport
	if rRCInactiveTransitionReport == nil {
		lbConn.Log.Errorf("RRCInactiveTransitionReport is nil")
		return
	}

	lbConn.Log.Infoln("Handle RRC Inactive Transition Report")

	for i := 0; i < len(rRCInactiveTransitionReport.ProtocolIEs.List); i++ {
		ie := rRCInactiveTransitionReport.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandleHandoverNotify(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}
	HandoverNotify := initiatingMessage.Value.HandoverNotify
	if HandoverNotify == nil {
		lbConn.Log.Errorf("HandoverNotify is nil")
		return
	}

	lbConn.Log.Infoln("Handle Handover notification")

	for i := 0; i < len(HandoverNotify.ProtocolIEs.List); i++ {
		ie := HandoverNotify.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

// TODO
// TS 23.502 4.9.1
func HandlePathSwitchRequest(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var rANUENGAPID *ngapType.RANUENGAPID
	var sourceAMFUENGAPID *ngapType.AMFUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("InitiatingMessage is nil")
		return
	}
	pathSwitchRequest := initiatingMessage.Value.PathSwitchRequest
	if pathSwitchRequest == nil {
		lbConn.Log.Errorf("PathSwitchRequest is nil")
		return
	}

	lbConn.Log.Infoln("Handle Path Switch Request")

	UeLbID := LB.IDGen.NextNumber()

	for _, ie := range pathSwitchRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Traceln("Decode IE RanUeNgapID")

			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
				return
			}
			ie.Value.RANUENGAPID.Value = UeLbID
		case ngapType.ProtocolIEIDSourceAMFUENGAPID: // reject
			sourceAMFUENGAPID = ie.Value.SourceAMFUENGAPID
			lbConn.Log.Traceln("Decode IE SourceAmfUeNgapID")
			if sourceAMFUENGAPID == nil {
				lbConn.Log.Errorf("SourceAmfUeNgapID is nil")
				return
			}

		}
	}

	//TODO

	if lbConn.TypeID == context.TypeIdGNBConn {
		gnb := lbConn.RanPointer
		ue := context.NewUE()
		ue.UeRanID = rANUENGAPID.Value
		ue.UeLbID = UeLbID
		ue.RanID = gnb.GnbID
		gnb.Ues.Store(rANUENGAPID.Value, ue)
		ue.RanPointer = gnb
		context.ForwardToNextAmf(lbConn, message, ue)
	}
}

func HandleHandoverRequestAcknowledge(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Errorf("SuccessfulOutcome is nil")
		return
	}
	handoverRequestAcknowledge := successfulOutcome.Value.HandoverRequestAcknowledge // reject
	if handoverRequestAcknowledge == nil {
		lbConn.Log.Errorf("HandoverRequestAcknowledge is nil")
		return
	}

	lbConn.Log.Infoln("Handle Handover Request Acknowledge")

	for _, ie := range handoverRequestAcknowledge.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

//TODO
func HandleHandoverFailure(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	// var aMFUENGAPID *ngapType.AMFUENGAPID
	// var cause *ngapType.Cause
	// // var targetUe *context.RanUe
	// var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	LB = *context.LB_Self()

	// if lbConn == nil {
	// 	// logger.NgapLog.Error("ran is nil")
	// 	return
	// }
	// if message == nil {
	// 	// lbConn.Log.Error("NGAP Message is nil")
	// 	return
	// }

	// unsuccessfulOutcome := message.UnsuccessfulOutcome // reject
	// if unsuccessfulOutcome == nil {
	// 	// lbConn.Log.Error("Unsuccessful Message is nil")
	// 	return
	// }

	// handoverFailure := unsuccessfulOutcome.Value.HandoverFailure
	// if handoverFailure == nil {
	// 	// lbConn.Log.Error("HandoverFailure is nil")
	// 	return
	// }

	// for _, ie := range handoverFailure.ProtocolIEs.List {
	// 	switch ie.Id.Value {
	// 	case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
	// 		aMFUENGAPID = ie.Value.AMFUENGAPID
	// 		// lbConn.Log.Trace("Decode IE AmfUeNgapID")
	// 	case ngapType.ProtocolIEIDCause: // ignore
	// 		cause = ie.Value.Cause
	// 		// lbConn.Log.Trace("Decode IE Cause")
	// 	case ngapType.ProtocolIEIDCriticalityDiagnostics: // ignore
	// 		criticalityDiagnostics = ie.Value.CriticalityDiagnostics
	// 		// lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
	// 	}
	// }

	// // causePresent := ngapType.CausePresentRadioNetwork
	// // causeValue := ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem
	// // if cause != nil {
	// // 	causePresent, causeValue = printAndGetCause(lbConn, cause)
	// // }

	// if criticalityDiagnostics != nil {
	// 	printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	// }

	//TODO

	// targetUe = context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)

	// if targetUe == nil {
	// 	lbConn.Log.Errorf("No UE Context[AmfUENGAPID: %d]", aMFUENGAPID.Value)
	// 	cause := ngapType.Cause{
	// 		Present: ngapType.CausePresentRadioNetwork,
	// 		RadioNetwork: &ngapType.CauseRadioNetwork{
	// 			Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
	// 		},
	// 	}
	// 	ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, nil, &cause, nil)
	// 	return
	// }

	// sourceUe := targetUe.SourceUe
	// if sourceUe == nil {
	// 	// TODO: handle N2 Handover between AMF
	// 	lbConn.Log.Error("N2 Handover between AMF has not been implemented yet")
	// } else {
	// 	amfUe := targetUe.AmfUe
	// 	if amfUe != nil {
	// 		amfUe.SmContextList.Range(func(key, value interface{}) bool {
	// 			pduSessionID := key.(int32)
	// 			smContext := value.(*context.SmContext)
	// 			causeAll := context.CauseAll{
	// 				NgapCause: &models.NgApCause{
	// 					Group: int32(causePresent),
	// 					Value: int32(causeValue),
	// 				},
	// 			}
	// 			_, _, _, err := consumer.SendUpdateSmContextN2HandoverCanceled(amfUe, smContext, causeAll)
	// 			if err != nil {
	// 				lbConn.Log.Errorf("Send UpdateSmContextN2HandoverCanceled Error for PduSessionId[%d]", pduSessionID)
	// 			}
	// 			return true
	// 		})
	// 	}
	// 	ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, criticalityDiagnostics)
	// }

	// ngap_message.SendUEContextReleaseCommand(targetUe, context.UeContextReleaseHandover, causePresent, causeValue)
}

//TODO
func HandleHandoverRequired(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}
	HandoverRequired := initiatingMessage.Value.HandoverRequired
	if HandoverRequired == nil {
		lbConn.Log.Errorf("HandoverRequired is nil")
		return
	}

	lbConn.Log.Infoln("Handle HandoverRequired\n")

	for i := 0; i < len(HandoverRequired.ProtocolIEs.List); i++ {
		ie := HandoverRequired.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

//TODO
func HandleHandoverCancel(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	// var aMFUENGAPID *ngapType.AMFUENGAPID
	// var rANUENGAPID *ngapType.RANUENGAPID
	// var cause *ngapType.Cause

	// LB = *context.LB_Self()

	// if lbConn == nil {
	// 	logger.NgapLog.Errorf("ran is nil")
	// 	return
	// }
	// if message == nil {
	// 	lbConn.Log.Errorf("NGAP Message is nil")
	// 	return
	// }

	// initiatingMessage := message.InitiatingMessage
	// if initiatingMessage == nil {
	// 	lbConn.Log.Errorf("Initiating Message is nil")
	// 	return
	// }
	// HandoverCancel := initiatingMessage.Value.HandoverCancel
	// if HandoverCancel == nil {
	// 	lbConn.Log.Errorf("Handover Cancel is nil")
	// 	return
	// }

	// lbConn.Log.Infoln("Handle Handover Cancel")

	// for i := 0; i < len(HandoverCancel.ProtocolIEs.List); i++ {
	// 	ie := HandoverCancel.ProtocolIEs.List[i]
	// 	switch ie.Id.Value {
	// 	case ngapType.ProtocolIEIDAMFUENGAPID:
	// 		aMFUENGAPID = ie.Value.AMFUENGAPID
	// 		lbConn.Log.Traceln("Decode IE AmfUeNgapID")
	// 		if aMFUENGAPID == nil {
	// 			lbConn.Log.Errorf("AMFUENGAPID is nil")
	// 			return
	// 		}
	// 	case ngapType.ProtocolIEIDRANUENGAPID:
	// 		rANUENGAPID = ie.Value.RANUENGAPID
	// 		lbConn.Log.Traceln("Decode IE RanUeNgapID")
	// 		if rANUENGAPID == nil {
	// 			lbConn.Log.Errorf("RANUENGAPID is nil")
	// 			return
	// 		}
	// 	case ngapType.ProtocolIEIDCause:
	// 		cause = ie.Value.Cause
	// 		lbConn.Log.Traceln("Decode IE Cause")
	// 		if cause == nil {
	// 			// lbConn.Log.Errorf(cause, "cause is nil")
	// 			return
	// 		}
	// 	}
	// }

	// if lbConn.TypeID == context.TypeIdAMFConn {
	// 	amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
	// 	lbConn.Log.Traceln("AMF Found")
	// 	UE, _ := amf.FindUeByUeID(rANUENGAPID.Value)
	// 	context.ForwardToGnb(message, UE)
	// 	return
	// }
	// if lbConn.TypeID == context.TypeIdGNBConn {
	// 	gnb, _ := LB.LbGnbFindByConn(lbConn.Conn)
	// 	UE, ok := gnb.FindUeByUeRanID(rANUENGAPID.Value)
	// 	if !ok {
	// 		lbConn.Log.Errorf("UE not found")
	// 		return
	// 	}
	// 	context.ForwardToAmf(message, UE)
	// 	return
	// }
}

func HandleUplinkRanStatusTransfer(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // ignore
	if initiatingMessage == nil {
		lbConn.Log.Errorf("InitiatingMessage is nil")
		return
	}
	uplinkRanStatusTransfer := initiatingMessage.Value.UplinkRANStatusTransfer
	if uplinkRanStatusTransfer == nil {
		lbConn.Log.Errorf("UplinkRanStatusTransfer is nil")
		return
	}

	lbConn.Log.Infoln("Handle Uplink Ran Status Transfer")

	for _, ie := range uplinkRanStatusTransfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandleNasNonDeliveryIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("InitiatingMessage is nil")
		return
	}
	nASNonDeliveryIndication := initiatingMessage.Value.NASNonDeliveryIndication
	if nASNonDeliveryIndication == nil {
		lbConn.Log.Errorf("NASNonDeliveryIndication is nil")
		return
	}

	lbConn.Log.Infoln("Handle Nas Non Delivery Indication")

	for _, ie := range nASNonDeliveryIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

//Todo
func HandleRanConfigurationUpdate(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var rANNodeName *ngapType.RANNodeName
	var supportedTAList *ngapType.SupportedTAList
	var pagingDRX *ngapType.PagingDRX

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}

	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}
	rANConfigurationUpdate := initiatingMessage.Value.RANConfigurationUpdate
	if rANConfigurationUpdate == nil {
		lbConn.Log.Errorf("RAN Configuration is nil")
		return
	}
	lbConn.Log.Infoln("Handle Ran Configuration Update")

	for i := 0; i < len(rANConfigurationUpdate.ProtocolIEs.List); i++ {
		ie := rANConfigurationUpdate.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANNodeName:
			rANNodeName = ie.Value.RANNodeName
			if rANNodeName == nil {
				lbConn.Log.Errorf("RAN Node Name is nil")
				return
			}
			lbConn.Log.Traceln("Decode IE RANNodeName = [%s]", rANNodeName.Value)
		case ngapType.ProtocolIEIDSupportedTAList:
			supportedTAList = ie.Value.SupportedTAList
			lbConn.Log.Traceln("Decode IE SupportedTAList")
			if supportedTAList == nil {
				lbConn.Log.Errorf("Supported TA List is nil")
				return
			}
		case ngapType.ProtocolIEIDDefaultPagingDRX:
			pagingDRX = ie.Value.DefaultPagingDRX
			if pagingDRX == nil {
				lbConn.Log.Errorf("PagingDRX is nil")
				return
			}
			lbConn.Log.Traceln("Decode IE PagingDRX = [%d]", pagingDRX.Value)
		}
	}
}

//TODO
func HandleUplinkRanConfigurationTransfer(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var sONConfigurationTransferUL *ngapType.SONConfigurationTransfer

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("InitiatingMessage is nil")
		return
	}
	uplinkRANConfigurationTransfer := initiatingMessage.Value.UplinkRANConfigurationTransfer
	if uplinkRANConfigurationTransfer == nil {
		lbConn.Log.Errorf("ErrorIndication is nil")
		return
	}

	for _, ie := range uplinkRANConfigurationTransfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDSONConfigurationTransferUL: // optional, ignore
			sONConfigurationTransferUL = ie.Value.SONConfigurationTransferUL
			lbConn.Log.Traceln("Decode IE SONConfigurationTransferUL")
			if sONConfigurationTransferUL == nil {
				lbConn.Log.Warnf("sONConfigurationTransferUL is nil")
			}
		}
	}
}

func HandleUplinkUEAssociatedNRPPATransport(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("InitiatingMessage is nil")
		return
	}
	uplinkUEAssociatedNRPPaTransport := initiatingMessage.Value.UplinkUEAssociatedNRPPaTransport
	if uplinkUEAssociatedNRPPaTransport == nil {
		lbConn.Log.Errorf("uplinkUEAssociatedNRPPaTransport is nil")
		return
	}

	lbConn.Log.Infoln("Handle Uplink UE Associated NRPPA Transpor")

	for _, ie := range uplinkUEAssociatedNRPPaTransport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

// TODO
func HandleUplinkNonUEAssociatedNRPPATransport(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var routingID *ngapType.RoutingID
	var nRPPaPDU *ngapType.NRPPaPDU

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}
	uplinkNonUEAssociatedNRPPATransport := initiatingMessage.Value.UplinkNonUEAssociatedNRPPaTransport
	if uplinkNonUEAssociatedNRPPATransport == nil {
		lbConn.Log.Errorf("Uplink Non UE Associated NRPPA Transport is nil")
		return
	}

	lbConn.Log.Infoln("Handle Uplink Non UE Associated NRPPA Transport")

	for i := 0; i < len(uplinkNonUEAssociatedNRPPATransport.ProtocolIEs.List); i++ {
		ie := uplinkNonUEAssociatedNRPPATransport.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRoutingID:
			routingID = ie.Value.RoutingID
			lbConn.Log.Traceln("Decode IE RoutingID")

		case ngapType.ProtocolIEIDNRPPaPDU:
			nRPPaPDU = ie.Value.NRPPaPDU
			lbConn.Log.Traceln("Decode IE NRPPaPDU")
		}
	}

	if routingID == nil {
		lbConn.Log.Errorf("RoutingID is nil")
		return
	}
	// Forward routingID to LMF
	// Described in (23.502 4.13.5.6)

	if nRPPaPDU == nil {
		lbConn.Log.Errorf("NRPPaPDU is nil")
		return
	}
	// TODO: Forward NRPPaPDU to LMF
}

func HandleLocationReport(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("InitiatingMessage is nil")
		return
	}
	locationReport := initiatingMessage.Value.LocationReport
	if locationReport == nil {
		lbConn.Log.Errorf("LocationReport is nil")
		return
	}

	lbConn.Log.Infoln("Handle Location Report")

	for _, ie := range locationReport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandleUERadioCapabilityInfoIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}
	uERadioCapabilityInfoIndication := initiatingMessage.Value.UERadioCapabilityInfoIndication
	if uERadioCapabilityInfoIndication == nil {
		lbConn.Log.Errorf("UERadioCapabilityInfoIndication is nil")
		return
	}

	lbConn.Log.Infoln("Handle UE Radio Capability Info Indication")

	for i := 0; i < len(uERadioCapabilityInfoIndication.ProtocolIEs.List); i++ {
		ie := uERadioCapabilityInfoIndication.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

//TODO
func HandleAMFconfigurationUpdateFailure(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		lbConn.Log.Errorf("Unsuccessful Message is nil")
		return
	}

	AMFconfigurationUpdateFailure := unsuccessfulOutcome.Value.AMFConfigurationUpdateFailure
	if AMFconfigurationUpdateFailure == nil {
		lbConn.Log.Errorf("AMFConfigurationUpdateFailure is nil")
		return
	}

	lbConn.Log.Infoln("Handle AMF Confioguration Update Failure")

	for _, ie := range AMFconfigurationUpdateFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			lbConn.Log.Traceln("Decode IE Cause")
			if cause == nil {
				lbConn.Log.Errorf("Cause is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Traceln("Decode IE CriticalityDiagnostics")
		}
	}

	//	TODO: Time To Wait

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

//TODO
func HandleAMFconfigurationUpdateAcknowledge(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFTNLAssociationSetupList *ngapType.AMFTNLAssociationSetupList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	var aMFTNLAssociationFailedToSetupList *ngapType.TNLAssociationList

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Errorf("SuccessfulOutcome is nil")
		return
	}
	aMFConfigurationUpdateAcknowledge := successfulOutcome.Value.AMFConfigurationUpdateAcknowledge
	if aMFConfigurationUpdateAcknowledge == nil {
		lbConn.Log.Errorf("AMFConfigurationUpdateAcknowledge is nil")
		return
	}

	lbConn.Log.Infoln("Handle AMF Configuration Update Acknowledge")

	for i := 0; i < len(aMFConfigurationUpdateAcknowledge.ProtocolIEs.List); i++ {
		ie := aMFConfigurationUpdateAcknowledge.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFTNLAssociationSetupList:
			aMFTNLAssociationSetupList = ie.Value.AMFTNLAssociationSetupList
			lbConn.Log.Traceln("Decode IE AMFTNLAssociationSetupList")
			if aMFTNLAssociationSetupList == nil {
				lbConn.Log.Errorf("AMFTNLAssociationSetupList is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Traceln("Decode IE Criticality Diagnostics")

		case ngapType.ProtocolIEIDAMFTNLAssociationFailedToSetupList:
			aMFTNLAssociationFailedToSetupList = ie.Value.AMFTNLAssociationFailedToSetupList
			lbConn.Log.Traceln("Decode IE AMFTNLAssociationFailedToSetupList")
			if aMFTNLAssociationFailedToSetupList == nil {
				lbConn.Log.Errorf("AMFTNLAssociationFailedToSetupList is nil")
				return
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

func HandleErrorIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("InitiatingMessage is nil")
		return
	}
	errorIndication := initiatingMessage.Value.ErrorIndication
	if errorIndication == nil {
		lbConn.Log.Errorf("ErrorIndication is nil")
		return
	}

	for _, ie := range errorIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func HandleCellTrafficTrace(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		logger.NgapLog.Errorf("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // ignore
	if initiatingMessage == nil {
		lbConn.Log.Errorf("InitiatingMessage is nil")
		return
	}
	cellTrafficTrace := initiatingMessage.Value.CellTrafficTrace
	if cellTrafficTrace == nil {
		lbConn.Log.Errorf("CellTrafficTrace is nil")
		return
	}

	lbConn.Log.Infoln("Handle Cell Traffic Trace")

	for _, ie := range cellTrafficTrace.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Traceln("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Errorf("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
			lbConn.Log.Traceln("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Errorf("RanUeNgapID is nil")
			} else {
				gnb := lbConn.RanPointer
				ue, ok := gnb.FindUeByUeRanID(rANUENGAPIDInt)
				if !ok {
					lbConn.Log.Errorf("UE not registered")
					return
				}
				ie.Value.RANUENGAPID.Value = ue.UeLbID
				context.ForwardToAmf(message, ue)
			}
		}
	}
}

func printAndGetCause(lbConn *context.LBConn, cause *ngapType.Cause) (present int, value aper.Enumerated) {
	present = cause.Present
	switch cause.Present {
	case ngapType.CausePresentRadioNetwork:
		lbConn.Log.Warnf("Cause RadioNetwork[%d]", cause.RadioNetwork.Value)
		value = cause.RadioNetwork.Value
	case ngapType.CausePresentTransport:
		lbConn.Log.Warnf("Cause Transport[%d]", cause.Transport.Value)
		value = cause.Transport.Value
	case ngapType.CausePresentProtocol:
		lbConn.Log.Warnf("Cause Protocol[%d]", cause.Protocol.Value)
		value = cause.Protocol.Value
	case ngapType.CausePresentNas:
		lbConn.Log.Warnf("Cause Nas[%d]", cause.Nas.Value)
		value = cause.Nas.Value
	case ngapType.CausePresentMisc:
		lbConn.Log.Warnf("Cause Misc[%d]", cause.Misc.Value)
		value = cause.Misc.Value
	default:
		lbConn.Log.Errorf("Invalid Cause group[%d]", cause.Present)
	}
	return
}

func printCriticalityDiagnostics(lbConn *context.LBConn, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {
	lbConn.Log.Trace("Criticality Diagnostics")

	if criticalityDiagnostics.ProcedureCriticality != nil {
		switch criticalityDiagnostics.ProcedureCriticality.Value {
		case ngapType.CriticalityPresentReject:
			lbConn.Log.Trace("Procedure Criticality: Reject")
		case ngapType.CriticalityPresentIgnore:
			lbConn.Log.Trace("Procedure Criticality: Ignore")
		case ngapType.CriticalityPresentNotify:
			lbConn.Log.Trace("Procedure Criticality: Notify")
		}
	}

	if criticalityDiagnostics.IEsCriticalityDiagnostics != nil {
		for _, ieCriticalityDiagnostics := range criticalityDiagnostics.IEsCriticalityDiagnostics.List {
			lbConn.Log.Tracef("IE ID: %d", ieCriticalityDiagnostics.IEID.Value)

			switch ieCriticalityDiagnostics.IECriticality.Value {
			case ngapType.CriticalityPresentReject:
				lbConn.Log.Trace("Criticality Reject")
			case ngapType.CriticalityPresentNotify:
				lbConn.Log.Trace("Criticality Notify")
			}

			switch ieCriticalityDiagnostics.TypeOfError.Value {
			case ngapType.TypeOfErrorPresentNotUnderstood:
				lbConn.Log.Trace("Type of error: Not understood")
			case ngapType.TypeOfErrorPresentMissing:
				lbConn.Log.Trace("Type of error: Missing")
			}
		}
	}
}

func buildCriticalityDiagnostics(
	procedureCode *int64,
	triggeringMessage *aper.Enumerated,
	procedureCriticality *aper.Enumerated,
	iesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList) (
	criticalityDiagnostics ngapType.CriticalityDiagnostics) {
	if procedureCode != nil {
		criticalityDiagnostics.ProcedureCode = new(ngapType.ProcedureCode)
		criticalityDiagnostics.ProcedureCode.Value = *procedureCode
	}

	if triggeringMessage != nil {
		criticalityDiagnostics.TriggeringMessage = new(ngapType.TriggeringMessage)
		criticalityDiagnostics.TriggeringMessage.Value = *triggeringMessage
	}

	if procedureCriticality != nil {
		criticalityDiagnostics.ProcedureCriticality = new(ngapType.Criticality)
		criticalityDiagnostics.ProcedureCriticality.Value = *procedureCriticality
	}

	if iesCriticalityDiagnostics != nil {
		criticalityDiagnostics.IEsCriticalityDiagnostics = iesCriticalityDiagnostics
	}

	return criticalityDiagnostics
}

func buildCriticalityDiagnosticsIEItem(ieCriticality aper.Enumerated, ieID int64, typeOfErr aper.Enumerated) (
	item ngapType.CriticalityDiagnosticsIEItem) {
	item = ngapType.CriticalityDiagnosticsIEItem{
		IECriticality: ngapType.Criticality{
			Value: ieCriticality,
		},
		IEID: ngapType.ProtocolIEID{
			Value: ieID,
		},
		TypeOfError: ngapType.TypeOfError{
			Value: typeOfErr,
		},
	}

	return item
}
