package gnb_ngap

// This handles messages incoming from AMF with the functions of the GNBs handler 

import (
	"github.com/LuckyG0ldfish/balancer/context"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/ngap/ngapType"
)

func HandleNGSetupResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	logger.GNBHandlerLog.Debugln("[gNB] Handle NG Setup Response")

	LB := context.LB_Self()

	var servedGUAMIList *ngapType.ServedGUAMIList
	var plmnSupportList *ngapType.PLMNSupportList

	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}

	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Errorf("Successful Outcome is nil")
		return
	}

	ngSetupResponse := successfulOutcome.Value.NGSetupResponse
	if ngSetupResponse == nil {
		lbConn.Log.Errorf("ngSetupResponse is nil")
		return
	}

	for _, ie := range ngSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFName:
			lbConn.Log.Traceln("[NGAP] Decode IE AMFName")
		case ngapType.ProtocolIEIDServedGUAMIList:
			lbConn.Log.Traceln("[NGAP] Decode IE ServedGUAMIList")
			servedGUAMIList = ie.Value.ServedGUAMIList
			if servedGUAMIList == nil {
				lbConn.Log.Errorf("ServedGUAMIList is nil")
			}
			LB.ServedGuamiList = servedGUAMIList
		case ngapType.ProtocolIEIDRelativeAMFCapacity:
			lbConn.Log.Traceln("[NGAP] Decode IE RelativeAMFCapacity")
			relativeAMFCapacity := ie.Value.RelativeAMFCapacity
			amf, ok := LB.LbAmfFindByConn(lbConn.Conn)
			if !ok {
				lbConn.Log.Errorf("AMF not found -> Capacity not set")
				
			} else {
				amf.RelativeCapacity = relativeAMFCapacity.Value
				lbConn.Log.Traceln("[NGAP] AMFs RelativeAMFCapacity set to %d", relativeAMFCapacity.Value)
			}
		case ngapType.ProtocolIEIDPLMNSupportList:
			lbConn.Log.Traceln("[NGAP] Decode IE PLMNSupportList")
			plmnSupportList = ie.Value.PLMNSupportList
			if plmnSupportList == nil {
				lbConn.Log.Errorf("PLMNSupportList is nil")
			}
			LB.PlmnSupportList = plmnSupportList
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			lbConn.Log.Traceln("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}
}

func HandleInitialContextSetupRequest(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	logger.GNBHandlerLog.Debugln("[gNB] Handle Initial Context Setup Request")
	
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}

	initialContextSetupRequest := initiatingMessage.Value.InitialContextSetupRequest
	if initialContextSetupRequest == nil {
		lbConn.Log.Errorf("InitialContextSetupRequest is nil")
		return
	}

	var aMFUENGAPIDInt int64
	var amfIDPresent bool = false 

	for _, ie := range initialContextSetupRequest.ProtocolIEs.List {
		switch ie.Id.Value {
			case ngapType.ProtocolIEIDAMFUENGAPID: // reject
				aMFUENGAPID = ie.Value.AMFUENGAPID
				lbConn.Log.Traceln("Decode IE AmfUeNgapID")
				if aMFUENGAPID == nil {
					lbConn.Log.Errorf("AmfUeNgapID is nil")
				} else {
				aMFUENGAPIDInt = aMFUENGAPID.Value
				amfIDPresent = true
				}
			case ngapType.ProtocolIEIDRANUENGAPID: // reject
				rANUENGAPID = ie.Value.RANUENGAPID
				rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
				lbConn.Log.Traceln("Decode IE RanUeNgapID")
				if rANUENGAPID == nil {
					lbConn.Log.Errorf("RanUeNgapID is nil")
				} else {
					amf := lbConn.AmfPointer
					ue, ok := amf.FindUeByUeID(rANUENGAPIDInt)
					if !ok {
						lbConn.Log.Errorf("UE not registered")
						return 
					}
					ie.Value.RANUENGAPID.Value = ue.UeRanID
					if amfIDPresent {
						ue.UeAmfId = aMFUENGAPIDInt
					}
					context.ForwardToGnb(message, ue)
				}
		}
	}
}

// TODO 
func HandleUEContextReleaseCommand(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	logger.GNBHandlerLog.Debugln("[gNB] Handle UE Context Release Command TODO")

	var ueNgapIDs *ngapType.UENGAPIDs
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	// var lbCtx = context.LB_Self()

	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}

	ueContextReleaseCommand := initiatingMessage.Value.UEContextReleaseCommand
	if ueContextReleaseCommand == nil {
		logger.NgapLog.Error("UEContextReleaseCommand is nil")
		return
	}

	for _, ie := range ueContextReleaseCommand.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDUENGAPIDs:
			logger.NgapLog.Traceln("[NGAP] Decode IE UENGAPIDs")
			ueNgapIDs = ie.Value.UENGAPIDs
			if ueNgapIDs == nil {
				logger.NgapLog.Errorf("UENGAPIDs is nil")
				return 
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		// TODO: send error indication
		return
	}

	switch ueNgapIDs.Present {
	case ngapType.UENGAPIDsPresentUENGAPIDPair:
		// TODO
		// ueCtx, _ = lbCtx.UePoolLoad(ueNgapIDs.UENGAPIDPair.RANUENGAPID.Value)
	case ngapType.UENGAPIDsPresentAMFUENGAPID:
		// TODO: find UE according to specific AMF
		// The implementation here may have error when N3IWF need to
		// connect multiple AMFs.
		// Use UEpool in AMF context can solve this problem
		// ueCtx = amf.FindUeByAmfUeNgapID(ueNgapIDs.AMFUENGAPID.Value)
	}
}

func HandleDownlinkNASTransport(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	logger.GNBHandlerLog.Debugln("[gNB] Handle Downlink NAS Transport")

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	

	if message == nil {
		logger.NgapLog.Errorf("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Errorf("Initiating Message is nil")
		return
	}

	downlinkNASTransport := initiatingMessage.Value.DownlinkNASTransport
	if downlinkNASTransport == nil {
		logger.NgapLog.Errorf("DownlinkNASTransport is nil")
		return
	}

	var aMFUENGAPIDInt int64
	var amfIDPresent bool = false

	for _, ie := range downlinkNASTransport.ProtocolIEs.List {
		switch ie.Id.Value {
			case ngapType.ProtocolIEIDAMFUENGAPID: // reject
				aMFUENGAPID = ie.Value.AMFUENGAPID
				lbConn.Log.Trace("Decode IE AmfUeNgapID")
				if aMFUENGAPID == nil {
					lbConn.Log.Errorf("AmfUeNgapID is nil")
				} else {
				aMFUENGAPIDInt = aMFUENGAPID.Value
				amfIDPresent = true
				}
			case ngapType.ProtocolIEIDRANUENGAPID: // reject
				rANUENGAPID = ie.Value.RANUENGAPID
				rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
				lbConn.Log.Trace("Decode IE RanUeNgapID")
				if rANUENGAPID == nil {
					lbConn.Log.Errorf("RanUeNgapID is nil")
				} else {
					amf := lbConn.AmfPointer
					ue, ok := amf.FindUeByUeID(rANUENGAPIDInt)
					if !ok {
						lbConn.Log.Errorf("UE not registered")
						return 
					}
					ie.Value.RANUENGAPID.Value = ue.UeRanID
					if amfIDPresent {
						ue.UeAmfId = aMFUENGAPIDInt
					}
					context.ForwardToGnb(message, ue)
				}
		}	
	}
}

func HandlePDUSessionResourceSetupRequest(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	logger.GNBHandlerLog.Debugln("[gNB] Handle PDU Session Resource Setup Request")

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	
	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}

	pduSessionResourceSetupRequest := initiatingMessage.Value.PDUSessionResourceSetupRequest
	if pduSessionResourceSetupRequest == nil {
		lbConn.Log.Errorf("PDUSessionResourceSetupRequest is nil")
		return
	}

	var aMFUENGAPIDInt int64
	var amfIDPresent bool = false

	for _, ie := range pduSessionResourceSetupRequest.ProtocolIEs.List {
		switch ie.Id.Value {
			case ngapType.ProtocolIEIDAMFUENGAPID: // reject
				aMFUENGAPID = ie.Value.AMFUENGAPID
				lbConn.Log.Traceln("Decode IE AmfUeNgapID")
				if aMFUENGAPID == nil {
					lbConn.Log.Errorf("AmfUeNgapID is nil")
				} else {
				aMFUENGAPIDInt = aMFUENGAPID.Value
				amfIDPresent = true
				}
			case ngapType.ProtocolIEIDRANUENGAPID: // reject
				rANUENGAPID = ie.Value.RANUENGAPID
				rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
				lbConn.Log.Traceln("Decode IE RanUeNgapID")
				if rANUENGAPID == nil {
					lbConn.Log.Errorf("RanUeNgapID is nil")
				} else {
					amf := lbConn.AmfPointer
					ue, ok := amf.FindUeByUeID(rANUENGAPIDInt)
					if !ok {
						lbConn.Log.Errorf("UE not registered")
						return 
					}
					ie.Value.RANUENGAPID.Value = ue.UeRanID
					if amfIDPresent {
						ue.UeAmfId = aMFUENGAPIDInt
					}
					context.ForwardToGnb(message, ue)
				}
		}
	}
}

// TODO
func HandlePDUSessionResourceReleaseCommand(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	logger.GNBHandlerLog.Debugln("[gNB] Handle PDU Session Resource Release Command")
	
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	if message == nil {
		lbConn.Log.Errorf("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Errorf("Initiating Message is nil")
		return
	}

	pDUSessionResourceReleaseCommand := initiatingMessage.Value.PDUSessionResourceReleaseCommand
	if pDUSessionResourceReleaseCommand == nil {
		lbConn.Log.Errorf("pDUSessionResourceReleaseCommand is nil")
		return
	}

	var aMFUENGAPIDInt int64
	var amfIDPresent bool = false

	for _, ie := range pDUSessionResourceReleaseCommand.ProtocolIEs.List {
		switch ie.Id.Value {
			case ngapType.ProtocolIEIDAMFUENGAPID: // reject
				aMFUENGAPID = ie.Value.AMFUENGAPID
				lbConn.Log.Traceln("Decode IE AmfUeNgapID")
				if aMFUENGAPID == nil {
					lbConn.Log.Errorf("AmfUeNgapID is nil")
				} else {
				aMFUENGAPIDInt = aMFUENGAPID.Value
				amfIDPresent = true
				}
			case ngapType.ProtocolIEIDRANUENGAPID: // reject
				rANUENGAPID = ie.Value.RANUENGAPID
				rANUENGAPIDInt := ie.Value.RANUENGAPID.Value
				lbConn.Log.Traceln("Decode IE RanUeNgapID")
				if rANUENGAPID == nil {
					lbConn.Log.Errorf("RanUeNgapID is nil")
				} else {
					amf := lbConn.AmfPointer
					ue, ok := amf.FindUeByUeID(rANUENGAPIDInt)
					if !ok {
						logger.NgapLog.Errorf("UE not registered")
						return 
					}
					ie.Value.RANUENGAPID.Value = ue.UeRanID
					if amfIDPresent {
						ue.UeAmfId = aMFUENGAPIDInt
					}
					context.ForwardToGnb(message, ue)
				}
		}
	}
}
