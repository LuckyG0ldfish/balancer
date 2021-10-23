package ngap

import (

	// "strconv"
	//"aper"

	"fmt"

	"github.com/LuckyG0ldfish/balancer/context"
	"github.com/free5gc/aper"
	"github.com/free5gc/ngap/ngapType"
	// "github.com/free5gc/ngap"
	// "github.com/free5gc/amf/consumer"
	// gmm_common "github.com/free5gc/amf/gmm/common"
	// gmm_message "github.com/free5gc/amf/gmm/message"
	// "github.com/free5gc/amf/logger"
	// "github.com/free5gc/amf/nas"
	ngap_message "github.com/LuckyG0ldfish/balancer/ngap/message"
	// "github.com/free5gc/aper"
	// "github.com/free5gc/nas/nasMessage"
	// libngap "github.com/free5gc/ngap"
	// "github.com/free5gc/ngap/ngapConvert"
	// "github.com/free5gc/ngap/ngapType"
	// "github.com/free5gc/openapi/models"
)

var LB context.LBContext 

//TODO
func HandleNGSetupRequest(LbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var globalRANNodeID *ngapType.GlobalRANNodeID
	var rANNodeName *ngapType.RANNodeName
	var supportedTAList *ngapType.SupportedTAList
	var pagingDRX *ngapType.PagingDRX

	LB = *context.LB_Self()
	var cause ngapType.Cause

	if LbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// LbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// LbConn.Log.Error("Initiating Message is nil")
		return
	}
	nGSetupRequest := initiatingMessage.Value.NGSetupRequest
	if nGSetupRequest == nil {
		// LbConn.Log.Error("NGSetupRequest is nil")
		return
	}
	// LbConn.Log.Info("Handle NG Setup request")
	for i := 0; i < len(nGSetupRequest.ProtocolIEs.List); i++ {
		ie := nGSetupRequest.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDGlobalRANNodeID:
			globalRANNodeID = ie.Value.GlobalRANNodeID
			// LbConn.Log.Trace("Decode IE GlobalRANNodeID")
			if globalRANNodeID == nil {
				// LbConn.Log.Error("GlobalRANNodeID is nil")
				return
			}
		case ngapType.ProtocolIEIDSupportedTAList:
			supportedTAList = ie.Value.SupportedTAList
			// LbConn.Log.Trace("Decode IE SupportedTAList")
			if supportedTAList == nil {
				// LbConn.Log.Error("SupportedTAList is nil")
				return
			}
		case ngapType.ProtocolIEIDRANNodeName:
			rANNodeName = ie.Value.RANNodeName
			// LbConn.Log.Trace("Decode IE RANNodeName")
			if rANNodeName == nil {
				// LbConn.Log.Error("RANNodeName is nil")
				return
			}
		case ngapType.ProtocolIEIDDefaultPagingDRX:
			pagingDRX = ie.Value.DefaultPagingDRX
			// LbConn.Log.Trace("Decode IE DefaultPagingDRX")
			if pagingDRX == nil {
				// LbConn.Log.Error("DefaultPagingDRX is nil")
				return
			}
		}
	}

	// LbConn.SetRanId(globalRANNodeID)
	// if rANNodeName != nil {
	// 	// LbConn.Name = rANNodeName.Value
	// }
	// if pagingDRX != nil {
	// 	// LbConn.Log.Tracef("PagingDRX[%d]", pagingDRX.Value)
	// }

	// for i := 0; i < len(supportedTAList.List); i++ {
	// 	supportedTAItem := supportedTAList.List[i]
	// 	tac := hex.EncodeToString(supportedTAItem.TAC.Value)
	// 	capOfSupportTai := cap(LbConn.SupportedTAList)
	// 	for j := 0; j < len(supportedTAItem.BroadcastPLMNList.List); j++ {
	// 		supportedTAI := context.NewSupportedTAI()
	// 		supportedTAI.Tai.Tac = tac
	// 		broadcastPLMNItem := supportedTAItem.BroadcastPLMNList.List[j]
	// 		plmnId := ngapConvert.PlmnIdToModels(broadcastPLMNItem.PLMNIdentity)
	// 		supportedTAI.Tai.PlmnId = &plmnId
	// 		capOfSNssaiList := cap(supportedTAI.SNssaiList)
	// 		for k := 0; k < len(broadcastPLMNItem.TAISliceSupportList.List); k++ {
	// 			tAISliceSupportItem := broadcastPLMNItem.TAISliceSupportList.List[k]
	// 			if len(supportedTAI.SNssaiList) < capOfSNssaiList {
	// 				supportedTAI.SNssaiList = append(supportedTAI.SNssaiList, ngapConvert.SNssaiToModels(tAISliceSupportItem.SNSSAI))
	// 			} else {
	// 				break
	// 			}
	// 		}
	// 		// LbConn.Log.Tracef("PLMN_ID[MCC:%s MNC:%s] TAC[%s]", plmnId.Mcc, plmnId.Mnc, tac)
	// 		if len(LbConn.SupportedTAList) < capOfSupportTai {
	// 			LbConn.SupportedTAList = append(LbConn.SupportedTAList, supportedTAI)
	// 		} else {
	// 			break
	// 		}
	// 	}
	// }

	// if len(LbConn.SupportedTAList) == 0 {
	// 	LbConn.Log.Warn("NG-Setup failure: No supported TA exist in NG-Setup request")
	// 	cause.Present = ngapType.CausePresentMisc
	// 	cause.Misc = &ngapType.CauseMisc{
	// 		Value: ngapType.CauseMiscPresentUnspecified,
	// 	}
	// } else {
	// 	var found bool
	// 	for i, tai := range LbConn.SupportedTAList {
	// 		if context.InTaiList(tai.Tai, context.AMF_Self().SupportTaiLists) {
	// 			LbConn.Log.Tracef("SERVED_TAI_INDEX[%d]", i)
	// 			found = true
	// 			break
	// 		}
	// 	}
	// 	if !found {
	// 		LbConn.Log.Warn("NG-Setup failure: Cannot find Served TAI in AMF")
	// 		cause.Present = ngapType.CausePresentMisc
	// 		cause.Misc = &ngapType.CauseMisc{
	// 			Value: ngapType.CauseMiscPresentUnknownPLMN,
	// 		}
	// 	}
	// }

	if cause.Present == ngapType.CausePresentNothing {
		ngap_message.SendNGSetupResponse(LbConn)
	} else {
		ngap_message.SendNGSetupFailure(LbConn, cause)
	}
} 

func HandleUplinkNasTransport(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("Initiating Message is nil")
		return
	}

	uplinkNasTransport := initiatingMessage.Value.UplinkNASTransport
	if uplinkNasTransport == nil {
		// lbConn.Log.Error("UplinkNasTransport is nil")
		return
	}
	// lbConn.Log.Info("Handle Uplink Nas Transport")

	for i := 0; i < len(uplinkNasTransport.ProtocolIEs.List); i++ {
		ie := uplinkNasTransport.ProtocolIEs.List[i]
		switch ie.Id.Value {
			case ngapType.ProtocolIEIDAMFUENGAPID: // reject
				aMFUENGAPID = ie.Value.AMFUENGAPID
				// lbConn.Log.Trace("Decode IE AmfUeNgapID")
				if aMFUENGAPID == nil {
					// lbConn.Log.Error("AmfUeNgapID is nil")
					fmt.Println("AmfUeNgapID is nil")
				}
			case ngapType.ProtocolIEIDRANUENGAPID: // reject
				rANUENGAPID = ie.Value.RANUENGAPID
				// lbConn.Log.Trace("Decode IE RanUeNgapID")
				if rANUENGAPID == nil {
					// lbConn.Log.Error("RanUeNgapID is nil")
					fmt.Println("RanUeNgapID is nil")
				}
			}
	}


	// TODO
	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

// TODO
func HandleNGReset(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var cause *ngapType.Cause
	var resetType *ngapType.ResetType

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("Initiating Message is nil")
		return
	}
	nGReset := initiatingMessage.Value.NGReset
	if nGReset == nil {
		// lbConn.Log.Error("NGReset is nil")
		return
	}

	// lbConn.Log.Info("Handle NG Reset")

	for _, ie := range nGReset.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			// lbConn.Log.Trace("Decode IE Cause")
			if cause == nil {
				// lbConn.Log.Error("Cause is nil")
				return
			}
		case ngapType.ProtocolIEIDResetType:
			resetType = ie.Value.ResetType
			// lbConn.Log.Trace("Decode IE ResetType")
			if resetType == nil {
				// lbConn.Log.Error("ResetType is nil")
				return
			}
		}
	}

	printAndGetCause(lbConn, cause)

	switch resetType.Present {
	case ngapType.ResetTypePresentNGInterface:
		// lbConn.Log.Trace("ResetType Present: NG Interface")
		// lbConn.RemoveAllUeInRan()
		// ngap_message.SendNGResetAcknowledge(lbConn, nil, nil)
	case ngapType.ResetTypePresentPartOfNGInterface:
		// lbConn.Log.Trace("ResetType Present: Part of NG Interface")

		partOfNGInterface := resetType.PartOfNGInterface
		if partOfNGInterface == nil {
			// lbConn.Log.Error("PartOfNGInterface is nil")
			return
		}

		// 	var ranUe *context.RanUe

		// 	for _, ueAssociatedLogicalNGConnectionItem := range partOfNGInterface.List {
		// 		if ueAssociatedLogicalNGConnectionItem.AMFUENGAPID != nil {
		// 			// lbConn.Log.Tracef("AmfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
		// 			for _, ue := range lbConn.RanUeList {
		// 				if ue.AmfUeNgapId == ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value {
		// 					ranUe = ue
		// 					break
		// 				}
		// 			}
		// 		} else if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
		// 			lbConn.Log.Tracef("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
		// 			ranUe = lbConn.RanUeFindByRanUeNgapID(ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
		// 		}

		// 		if ranUe == nil {
		// 			// lbConn.Log.Warn("Cannot not find UE Context")
		// 			if ueAssociatedLogicalNGConnectionItem.AMFUENGAPID != nil {
		// 				// lbConn.Log.Warnf("AmfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
		// 			}
		// 			if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
		// 				// lbConn.Log.Warnf("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
		// 			}
		// 		}

		// 		err := ranUe.Remove()
		// 		if err != nil {
		// 			// lbConn.Log.Error(err.Error())
		// 		}
		// 	}
		// 	ngap_message.SendNGResetAcknowledge(lbConn, partOfNGInterface, nil)
		// default:
		// 	// lbConn.Log.Warnf("Invalid ResetType[%d]", resetType.Present)
	}
}

// TODO
func HandleNGResetAcknowledge(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var uEAssociatedLogicalNGConnectionList *ngapType.UEAssociatedLogicalNGConnectionList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		// lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	nGResetAcknowledge := successfulOutcome.Value.NGResetAcknowledge
	if nGResetAcknowledge == nil {
		// lbConn.Log.Error("NGResetAcknowledge is nil")
		return
	}

	// lbConn.Log.Info("Handle NG Reset Acknowledge")

	for _, ie := range nGResetAcknowledge.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDUEAssociatedLogicalNGConnectionList:
			uEAssociatedLogicalNGConnectionList = ie.Value.UEAssociatedLogicalNGConnectionList
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
		}
	}

	if uEAssociatedLogicalNGConnectionList != nil {
		// lbConn.Log.Tracef("%d UE association(s) has been reset", len(uEAssociatedLogicalNGConnectionList.List))
		for _, item := range uEAssociatedLogicalNGConnectionList.List {
			if item.AMFUENGAPID != nil && item.RANUENGAPID != nil {
				// lbConn.Log.Tracef("%d: AmfUeNgapID[%d] RanUeNgapID[%d]", i+1, item.AMFUENGAPID.Value, item.RANUENGAPID.Value)
			} else if item.AMFUENGAPID != nil {
				// lbConn.Log.Tracef("%d: AmfUeNgapID[%d] RanUeNgapID[-1]", i+1, item.AMFUENGAPID.Value)
			} else if item.RANUENGAPID != nil {
				// lbConn.Log.Tracef("%d: AmfUeNgapID[-1] RanUeNgapID[%d]", i+1, item.RANUENGAPID.Value)
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

func HandleUEContextReleaseComplete(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		// lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	uEContextReleaseComplete := successfulOutcome.Value.UEContextReleaseComplete
	if uEContextReleaseComplete == nil {
		// lbConn.Log.Error("NGResetAcknowledge is nil")
		return
	}

	// lbConn.Log.Info("Handle UE Context Release Complete")

	for _, ie := range uEContextReleaseComplete.ProtocolIEs.List {
		switch ie.Id.Value {
			case ngapType.ProtocolIEIDAMFUENGAPID: // reject
				aMFUENGAPID = ie.Value.AMFUENGAPID
				// lbConn.Log.Trace("Decode IE AmfUeNgapID")
				if aMFUENGAPID == nil {
					// lbConn.Log.Error("AmfUeNgapID is nil")
					fmt.Println("AmfUeNgapID is nil")
				}
			case ngapType.ProtocolIEIDRANUENGAPID: // reject
				rANUENGAPID = ie.Value.RANUENGAPID
				// lbConn.Log.Trace("Decode IE RanUeNgapID")
				if rANUENGAPID == nil {
					// lbConn.Log.Error("RanUeNgapID is nil")
					fmt.Println("RanUeNgapID is nil")
				}
			}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
				fmt.Println("UE Found // this should not happen") // 
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToAmf(lbConn, m2, ue)
		return
	}
}

func HandlePDUSessionResourceReleaseResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		// lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceReleaseResponse := successfulOutcome.Value.PDUSessionResourceReleaseResponse
	if pDUSessionResourceReleaseResponse == nil {
		// lbConn.Log.Error("PDUSessionResourceReleaseResponse is nil")
		return
	}

	// lbConn.Log.Info("Handle PDU Session Resource Release Response")

	for _, ie := range pDUSessionResourceReleaseResponse.ProtocolIEs.List {
		switch ie.Id.Value {
			case ngapType.ProtocolIEIDAMFUENGAPID: // reject
				aMFUENGAPID = ie.Value.AMFUENGAPID
				// lbConn.Log.Trace("Decode IE AmfUeNgapID")
				if aMFUENGAPID == nil {
					// lbConn.Log.Error("AmfUeNgapID is nil")
					fmt.Println("AmfUeNgapID is nil")
				}
			case ngapType.ProtocolIEIDRANUENGAPID: // reject
				rANUENGAPID = ie.Value.RANUENGAPID
				// lbConn.Log.Trace("Decode IE RanUeNgapID")
				if rANUENGAPID == nil {
					// lbConn.Log.Error("RanUeNgapID is nil")
					fmt.Println("RanUeNgapID is nil")
				}
			}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
				fmt.Println("UE Found // this should not happen") // 
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToAmf(lbConn, m2, ue)
		return
	}
}

func HandleUERadioCapabilityCheckResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		// lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}

	uERadioCapabilityCheckResponse := successfulOutcome.Value.UERadioCapabilityCheckResponse
	if uERadioCapabilityCheckResponse == nil {
		// lbConn.Log.Error("UERadioCapabilityCheckResponse is nil")
		return
	}
	// lbConn.Log.Info("Handle UE Radio Capability Check Response")

	for i := 0; i < len(uERadioCapabilityCheckResponse.ProtocolIEs.List); i++ {
		ie := uERadioCapabilityCheckResponse.ProtocolIEs.List[i]
		switch ie.Id.Value {
			case ngapType.ProtocolIEIDAMFUENGAPID: // reject
				aMFUENGAPID = ie.Value.AMFUENGAPID
				// lbConn.Log.Trace("Decode IE AmfUeNgapID")
				if aMFUENGAPID == nil {
					// lbConn.Log.Error("AmfUeNgapID is nil")
					fmt.Println("AmfUeNgapID is nil")
				}
			case ngapType.ProtocolIEIDRANUENGAPID: // reject
				rANUENGAPID = ie.Value.RANUENGAPID
				// lbConn.Log.Trace("Decode IE RanUeNgapID")
				if rANUENGAPID == nil {
					// lbConn.Log.Error("RanUeNgapID is nil")
					fmt.Println("RanUeNgapID is nil")
				}
			}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
				fmt.Println("UE Found // this should not happen") // 
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToAmf(lbConn, m2, ue)
		return
	}
}

func HandleLocationReportingFailureIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	// var ranUe *context.RanUe

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("Initiating Message is nil")
		return
	}
	locationReportingFailureIndication := initiatingMessage.Value.LocationReportingFailureIndication
	if locationReportingFailureIndication == nil {
		// lbConn.Log.Error("LocationReportingFailureIndication is nil")
		return
	}

	// lbConn.Log.Info("Handle Location Reporting Failure Indication")

	for i := 0; i < len(locationReportingFailureIndication.ProtocolIEs.List); i++ {
		ie := locationReportingFailureIndication.ProtocolIEs.List[i]
		switch ie.Id.Value {
			case ngapType.ProtocolIEIDAMFUENGAPID: // reject
				aMFUENGAPID = ie.Value.AMFUENGAPID
				// lbConn.Log.Trace("Decode IE AmfUeNgapID")
				if aMFUENGAPID == nil {
					// lbConn.Log.Error("AmfUeNgapID is nil")
					fmt.Println("AmfUeNgapID is nil")
				}
			case ngapType.ProtocolIEIDRANUENGAPID: // reject
				rANUENGAPID = ie.Value.RANUENGAPID
				// lbConn.Log.Trace("Decode IE RanUeNgapID")
				if rANUENGAPID == nil {
					// lbConn.Log.Error("RanUeNgapID is nil")
					fmt.Println("RanUeNgapID is nil")
				}
			}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
				fmt.Println("UE Found // this should not happen") // 
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToAmf(lbConn, m2, ue)
		return
	}
}

func HandleInitialUEMessage(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var rANUENGAPID *ngapType.RANUENGAPID
	var nASPDU *ngapType.NASPDU

	LB = *context.LB_Self()

	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("Initiating Message is nil")
		return
	}
	initialUEMessage := initiatingMessage.Value.InitialUEMessage
	if initialUEMessage == nil {
		// lbConn.Log.Error("InitialUEMessage is nil")
		return
	}

	// lbConn.Log.Info("Handle Initial UE Message")

	for _, ie := range initialUEMessage.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				fmt.Println("InitialUEMessage: rANUENGAPID == nil")
				return 
				// lbConn.Log.Error("RanUeNgapID is nil")
				// item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
				// 	ngapType.ProtocolIEIDRANUENGAPID, ngapType.TypeOfErrorPresentMissing)
				// iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDNASPDU: // reject
			nASPDU = ie.Value.NASPDU
			// lbConn.Log.Trace("Decode IE NasPdu")
			if nASPDU == nil {
				fmt.Println("InitialUEMessage: nASPDU == nil")
				// TODO
				// lbConn.Log.Error("NasPdu is nil")
				// item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDNASPDU,
				// 	ngapType.TypeOfErrorPresentMissing)
				// iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentGNBConn {
		gnb, ok := LB.LbGnbFindByConn(lbConn.Conn)
		ue := context.NewUE()		
		if ok {
			ue.UeRanID = rANUENGAPID.Value
			ue.RanID = gnb.GnbID
			var empty []*context.LbUe
			ues := append(empty, ue)
			gnb.Ues.Store(rANUENGAPID.Value, ues)
			LB.ForwardToNextAmf(lbConn, m2, ue)
		} else {
			fmt.Println("No GNB")
		}

		return
	}
}

func HandlePDUSessionResourceSetupResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	// var ranUe *context.RanUe

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		// lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceSetupResponse
	if pDUSessionResourceSetupResponse == nil {
		// lbConn.Log.Error("PDUSessionResourceSetupResponse is nil")
		return
	}

	// lbConn.Log.Info("Handle PDU Session Resource Setup Response")

	for _, ie := range pDUSessionResourceSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
				fmt.Println("UE Found // this should not happen") // 
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToAmf(lbConn, m2, ue)
		return
	}
}

func HandlePDUSessionResourceModifyResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	// var ranUe *context.RanUe

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		// lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceModifyResponse := successfulOutcome.Value.PDUSessionResourceModifyResponse
	if pDUSessionResourceModifyResponse == nil {
		// lbConn.Log.Error("PDUSessionResourceModifyResponse is nil")
		return
	}

	// lbConn.Log.Info("Handle PDU Session Resource Modify Response")

	for _, ie := range pDUSessionResourceModifyResponse.ProtocolIEs.List {
		switch ie.Id.Value {
			case ngapType.ProtocolIEIDAMFUENGAPID: // reject
				aMFUENGAPID = ie.Value.AMFUENGAPID
				// lbConn.Log.Trace("Decode IE AmfUeNgapID")
				if aMFUENGAPID == nil {
					// lbConn.Log.Error("AmfUeNgapID is nil")
					fmt.Println("AmfUeNgapID is nil")
				}
			case ngapType.ProtocolIEIDRANUENGAPID: // reject
				rANUENGAPID = ie.Value.RANUENGAPID
				// lbConn.Log.Trace("Decode IE RanUeNgapID")
				if rANUENGAPID == nil {
					// lbConn.Log.Error("RanUeNgapID is nil")
					fmt.Println("RanUeNgapID is nil")
				}
			}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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

	// if rANUENGAPID != nil {
	// 	ranUe = lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	// 	if ranUe == nil {
	// 		lbConn.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
	// 	}
	// }

	// if aMFUENGAPID != nil {
	// 	ranUe = context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
	// 	if ranUe == nil {
	// 		lbConn.Log.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
	// 		return
	// 	}
	// }
}

func HandlePDUSessionResourceNotify(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	// var ranUe *context.RanUe

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	PDUSessionResourceNotify := initiatingMessage.Value.PDUSessionResourceNotify
	if PDUSessionResourceNotify == nil {
		// lbConn.Log.Error("PDUSessionResourceNotify is nil")
		return
	}

	for _, ie := range PDUSessionResourceNotify.ProtocolIEs.List {
		switch ie.Id.Value {
			case ngapType.ProtocolIEIDAMFUENGAPID: // reject
				aMFUENGAPID = ie.Value.AMFUENGAPID
				// lbConn.Log.Trace("Decode IE AmfUeNgapID")
				if aMFUENGAPID == nil {
					// lbConn.Log.Error("AmfUeNgapID is nil")
					fmt.Println("AmfUeNgapID is nil")
				}
			case ngapType.ProtocolIEIDRANUENGAPID: // reject
				rANUENGAPID = ie.Value.RANUENGAPID
				// lbConn.Log.Trace("Decode IE RanUeNgapID")
				if rANUENGAPID == nil {
					// lbConn.Log.Error("RanUeNgapID is nil")
					fmt.Println("RanUeNgapID is nil")
				}
			}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

func HandlePDUSessionResourceModifyIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	// var ranUe *context.RanUe

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // reject
	if initiatingMessage == nil {
		// lbConn.Log.Error("InitiatingMessage is nil")
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
		// lbConn.Log.Error("PDUSessionResourceModifyIndication is nil")
		// cause := ngapType.Cause{
		// 	Present: ngapType.CausePresentProtocol,
		// 	Protocol: &ngapType.CauseProtocol{
		// 		Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorReject,
		// 	},
		// }
		// ngap_message.SendErrorIndication(lbConn, nil, nil, &cause, nil)
		return
	}

	// lbConn.Log.Info("Handle PDU Session Resource Modify Indication")

	for _, ie := range pDUSessionResourceModifyIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

func HandleInitialContextSetupResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		// lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	if initialContextSetupResponse == nil {
		// lbConn.Log.Error("InitialContextSetupResponse is nil")
		return
	}

	// lbConn.Log.Info("Handle Initial Context Setup Response")

	for _, ie := range initialContextSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Warn("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Warn("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

func HandleInitialContextSetupFailure(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	
	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		// lbConn.Log.Error("UnsuccessfulOutcome is nil")
		return
	}
	initialContextSetupFailure := unsuccessfulOutcome.Value.InitialContextSetupFailure
	if initialContextSetupFailure == nil {
		// lbConn.Log.Error("InitialContextSetupFailure is nil")
		return
	}

	// lbConn.Log.Info("Handle Initial Context Setup Failure")

	for _, ie := range initialContextSetupFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Warn("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Warn("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

func HandleUEContextReleaseRequest(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	uEContextReleaseRequest := initiatingMessage.Value.UEContextReleaseRequest
	if uEContextReleaseRequest == nil {
		// lbConn.Log.Error("UEContextReleaseRequest is nil")
		return
	}

	// lbConn.Log.Info("UE Context Release Request")

	for _, ie := range uEContextReleaseRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

func HandleUEContextModificationResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	// var ranUe *context.RanUe

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		// lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	uEContextModificationResponse := successfulOutcome.Value.UEContextModificationResponse
	if uEContextModificationResponse == nil {
		// lbConn.Log.Error("UEContextModificationResponse is nil")
		return
	}

	// lbConn.Log.Info("Handle UE Context Modification Response")

	for _, ie := range uEContextModificationResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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

	// if rANUENGAPID != nil {
	// 	ranUe = lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	// 	if ranUe == nil {
	// 		lbConn.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
	// 	}
	// }

	// if aMFUENGAPID != nil {
	// 	ranUe = context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
	// 	if ranUe == nil {
	// 		lbConn.Log.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
	// 		return
	// 	}
	// }
}

func HandleUEContextModificationFailure(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	
	LB = *context.LB_Self()
		// var ranUe *context.RanUe

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		// lbConn.Log.Error("UnsuccessfulOutcome is nil")
		return
	}
	uEContextModificationFailure := unsuccessfulOutcome.Value.UEContextModificationFailure
	if uEContextModificationFailure == nil {
		// lbConn.Log.Error("UEContextModificationFailure is nil")
		return
	}

	// lbConn.Log.Info("Handle UE Context Modification Failure")

	for _, ie := range uEContextModificationFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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

	// if rANUENGAPID != nil {
	// 	ranUe = lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	// 	if ranUe == nil {
	// 		lbConn.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
	// 	}
	// }

	// if aMFUENGAPID != nil {
	// 	ranUe = context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
	// 	if ranUe == nil {
	// 		lbConn.Log.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
	// 	}
	// }
}

func HandleRRCInactiveTransitionReport(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	
	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("Initiating Message is nil")
		return
	}

	rRCInactiveTransitionReport := initiatingMessage.Value.RRCInactiveTransitionReport
	if rRCInactiveTransitionReport == nil {
		// lbConn.Log.Error("RRCInactiveTransitionReport is nil")
		return
	}
	// lbConn.Log.Info("Handle RRC Inactive Transition Report")

	for i := 0; i < len(rRCInactiveTransitionReport.ProtocolIEs.List); i++ {
		ie := rRCInactiveTransitionReport.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

func HandleHandoverNotify(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("Initiating Message is nil")
		return
	}
	HandoverNotify := initiatingMessage.Value.HandoverNotify
	if HandoverNotify == nil {
		// lbConn.Log.Error("HandoverNotify is nil")
		return
	}

	// lbConn.Log.Info("Handle Handover notification")

	for i := 0; i < len(HandoverNotify.ProtocolIEs.List); i++ {
		ie := HandoverNotify.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

// TODO
// TS 23.502 4.9.1
func HandlePathSwitchRequest(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var rANUENGAPID *ngapType.RANUENGAPID
	var sourceAMFUENGAPID *ngapType.AMFUENGAPID
	// var userLocationInformation *ngapType.UserLocationInformation
	// var uESecurityCapabilities *ngapType.UESecurityCapabilities
	var pduSessionResourceToBeSwitchedInDLList *ngapType.PDUSessionResourceToBeSwitchedDLList
	// var pduSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListPSReq

	LB = *context.LB_Self()

	// var ranUe *context.RanUe

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	pathSwitchRequest := initiatingMessage.Value.PathSwitchRequest
	if pathSwitchRequest == nil {
		// lbConn.Log.Error("PathSwitchRequest is nil")
		return
	}

	// lbConn.Log.Info("Handle Path Switch Request")

	for _, ie := range pathSwitchRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDSourceAMFUENGAPID: // reject
			sourceAMFUENGAPID = ie.Value.SourceAMFUENGAPID
			// lbConn.Log.Trace("Decode IE SourceAmfUeNgapID")
			if sourceAMFUENGAPID == nil {
				// lbConn.Log.Error("SourceAmfUeNgapID is nil")
				return
			}
		// case ngapType.ProtocolIEIDUserLocationInformation: // ignore
		// 	userLocationInformation = ie.Value.UserLocationInformation
		// 	lbConn.Log.Trace("Decode IE UserLocationInformation")
		// case ngapType.ProtocolIEIDUESecurityCapabilities: // ignore
		// 	uESecurityCapabilities = ie.Value.UESecurityCapabilities
		// 	lbConn.Log.Trace("Decode IE UESecurityCapabilities")
		case ngapType.ProtocolIEIDPDUSessionResourceToBeSwitchedDLList: // reject
			pduSessionResourceToBeSwitchedInDLList = ie.Value.PDUSessionResourceToBeSwitchedDLList
			// lbConn.Log.Trace("Decode IE PDUSessionResourceToBeSwitchedDLList")
			if pduSessionResourceToBeSwitchedInDLList == nil {
				// lbConn.Log.Error("PDUSessionResourceToBeSwitchedDLList is nil")
				return
			}
			// case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListPSReq: // ignore
			// 	pduSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListPSReq
			// 	lbConn.Log.Trace("Decode IE PDUSessionResourceFailedToSetupListPSReq")
		}
	}

	//TODO

	// if lbConn.TypeID == context.TypeIdentAMFConn {
	// 	amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
	// 	amf.Ues.LoadOrStore(aMFUENGAPID.Value, context.NewUE(aMFUENGAPID.Value))
	// 	LB.ForwardToGnb(lbConn, message, rANUENGAPID.Value)
	// 	return
	// }
	if lbConn.TypeID == context.TypeIdentGNBConn {
		gnb, ok := LB.LbGnbFindByConn(lbConn.Conn)
		ue := context.NewUE()		
		if ok {
			ue.UeRanID = rANUENGAPID.Value
			ue.RanID = gnb.GnbID
			var empty []*context.LbUe
			ues := append(empty, ue)
			gnb.Ues.Store(rANUENGAPID.Value, ues)
			LB.ForwardToNextAmf(lbConn, m2, ue)
		} else {
			fmt.Println("No GNB")
		}

		return
	}

	// if sourceAMFUENGAPID == nil {
	// 	lbConn.Log.Error("SourceAmfUeNgapID is nil")
	// 	return
	// }
	// ranUe = context.AMF_Self().RanUeFindByAmfUeNgapID(sourceAMFUENGAPID.Value)
	// if ranUe == nil {
	// 	lbConn.Log.Errorf("Cannot find UE from sourceAMfUeNgapID[%d]", sourceAMFUENGAPID.Value)
	// 	ngap_message.SendPathSwitchRequestFailure(lbConn, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
	// 	return
	// }
}

func HandleHandoverRequestAcknowledge(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	
	LB = *context.LB_Self()

	// var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		// lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	handoverRequestAcknowledge := successfulOutcome.Value.HandoverRequestAcknowledge // reject
	if handoverRequestAcknowledge == nil {
		// lbConn.Log.Error("HandoverRequestAcknowledge is nil")
		return
	}

	// lbConn.Log.Info("Handle Handover Request Acknowledge")

	for _, ie := range handoverRequestAcknowledge.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

//TODO
func HandleHandoverFailure(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
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
func HandleHandoverRequired(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("Initiating Message is nil")
		return
	}
	HandoverRequired := initiatingMessage.Value.HandoverRequired
	if HandoverRequired == nil {
		// lbConn.Log.Error("HandoverRequired is nil")
		return
	}

	// lbConn.Log.Info("Handle HandoverRequired\n")
	for i := 0; i < len(HandoverRequired.ProtocolIEs.List); i++ {
		ie := HandoverRequired.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID // reject
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			// case ngapType.ProtocolIEIDHandoverType: // reject
			// 	handoverType = ie.Value.HandoverType
			// 	lbConn.Log.Trace("Decode IE HandoverType")
			// case ngapType.ProtocolIEIDCause: // ignore
			// 	cause = ie.Value.Cause
			// 	lbConn.Log.Trace("Decode IE Cause")
			// case ngapType.ProtocolIEIDTargetID: // reject
			// 	targetID = ie.Value.TargetID
			// 	lbConn.Log.Trace("Decode IE TargetID")
			// case ngapType.ProtocolIEIDPDUSessionResourceListHORqd: // reject
			// 	pDUSessionResourceListHORqd = ie.Value.PDUSessionResourceListHORqd
			// 	lbConn.Log.Trace("Decode IE PDUSessionResourceListHORqd")
			// case ngapType.ProtocolIEIDSourceToTargetTransparentContainer: // reject
			// 	sourceToTargetTransparentContainer = ie.Value.SourceToTargetTransparentContainer
			// 	lbConn.Log.Trace("Decode IE SourceToTargetTransparentContainer")
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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

	// if aMFUENGAPID == nil {
	// 	lbConn.Log.Error("AmfUeNgapID is nil")
	// 	item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDAMFUENGAPID,
	// 		ngapType.TypeOfErrorPresentMissing)
	// 	iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	// }
	// if rANUENGAPID == nil {
	// 	lbConn.Log.Error("RanUeNgapID is nil")
	// 	item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDRANUENGAPID,
	// 		ngapType.TypeOfErrorPresentMissing)
	// 	iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	// }

	// if handoverType == nil {
	// 	lbConn.Log.Error("handoverType is nil")
	// 	item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDHandoverType,
	// 		ngapType.TypeOfErrorPresentMissing)
	// 	iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	// }
	// if targetID == nil {
	// 	lbConn.Log.Error("targetID is nil")
	// 	item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDTargetID,
	// 		ngapType.TypeOfErrorPresentMissing)
	// 	iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	// }
	// if pDUSessionResourceListHORqd == nil {
	// 	lbConn.Log.Error("pDUSessionResourceListHORqd is nil")
	// 	item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
	// 		ngapType.ProtocolIEIDPDUSessionResourceListHORqd, ngapType.TypeOfErrorPresentMissing)
	// 	iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	// }
	// if sourceToTargetTransparentContainer == nil {
	// 	lbConn.Log.Error("sourceToTargetTransparentContainer is nil")
	// 	item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
	// 		ngapType.ProtocolIEIDSourceToTargetTransparentContainer, ngapType.TypeOfErrorPresentMissing)
	// 	iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	// }

	// if len(iesCriticalityDiagnostics.List) > 0 {
	// 	procedureCode := ngapType.ProcedureCodeHandoverPreparation
	// 	triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
	// 	procedureCriticality := ngapType.CriticalityPresentReject
	// 	criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage,
	// 		&procedureCriticality, &iesCriticalityDiagnostics)
	// 	ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, nil, &criticalityDiagnostics)
	// 	return
	// }

	// sourceUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	// if sourceUe == nil {
	// 	lbConn.Log.Errorf("Cannot find UE for RAN_UE_NGAP_ID[%d] ", rANUENGAPID.Value)
	// 	cause := ngapType.Cause{
	// 		Present: ngapType.CausePresentRadioNetwork,
	// 		RadioNetwork: &ngapType.CauseRadioNetwork{
	// 			Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
	// 		},
	// 	}
	// 	ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, &cause, nil)
	// 	return
	// }
	// amfUe := sourceUe.AmfUe
	// if amfUe == nil {
	// 	lbConn.Log.Error("Cannot find amfUE from sourceUE")
	// 	return
	// }

	// if targetID.Present != ngapType.TargetIDPresentTargetRANNodeID {
	// 	lbConn.Log.Errorf("targetID type[%d] is not supported", targetID.Present)
	// 	return
	// }
	// amfUe.SetOnGoing(sourceUe.Ran.AnType, &context.OnGoing{
	// 	Procedure: context.OnGoingProcedureN2Handover,
	// })
	// if !amfUe.SecurityContextIsValid() {
	// 	sourceUe.Log.Info("Handle Handover Preparation Failure [Authentication Failure]")
	// 	cause = &ngapType.Cause{
	// 		Present: ngapType.CausePresentNas,
	// 		Nas: &ngapType.CauseNas{
	// 			Value: ngapType.CauseNasPresentAuthenticationFailure,
	// 		},
	// 	}
	// 	ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
	// 	return
	// }
	// aMFSelf := context.AMF_Self()
	// targetRanNodeId := ngapConvert.RanIdToModels(targetID.TargetRANNodeID.GlobalRANNodeID)
	// targetRan, ok := aMFSelf.AmfRanFindByRanID(targetRanNodeId)
	// if !ok {
	// 	// handover between different AMF
	// 	sourceUe.Log.Warnf("Handover required : cannot find target Ran Node Id[%+v] in this AMF", targetRanNodeId)
	// 	sourceUe.Log.Error("Handover between different AMF has not been implemented yet")
	// 	return
	// 	// TODO: Send to T-AMF
	// 	// Described in (23.502 4.9.1.3.2) step 3.Namf_Communication_CreateUEContext Request
	// } else {
	// 	// Handover in same AMF
	// 	sourceUe.HandOverType.Value = handoverType.Value
	// 	tai := ngapConvert.TaiToModels(targetID.TargetRANNodeID.SelectedTAI)
	// 	targetId := models.NgRanTargetId{
	// 		RanNodeId: &targetRanNodeId,
	// 		Tai:       &tai,
	// 	}
	// 	var pduSessionReqList ngapType.PDUSessionResourceSetupListHOReq
	// 	for _, pDUSessionResourceHoItem := range pDUSessionResourceListHORqd.List {
	// 		pduSessionId := int32(pDUSessionResourceHoItem.PDUSessionID.Value)
	// 		if smContext, exist := amfUe.SmContextFindByPDUSessionID(pduSessionId); exist {
	// 			response, _, _, err := consumer.SendUpdateSmContextN2HandoverPreparing(amfUe, smContext,
	// 				models.N2SmInfoType_HANDOVER_REQUIRED, pDUSessionResourceHoItem.HandoverRequiredTransfer, "", &targetId)
	// 			if err != nil {
	// 				sourceUe.Log.Errorf("consumer.SendUpdateSmContextN2HandoverPreparing Error: %+v", err)
	// 			}
	// 			if response == nil {
	// 				sourceUe.Log.Errorf("SendUpdateSmContextN2HandoverPreparing Error for PduSessionId[%d]", pduSessionId)
	// 				continue
	// 			} else if response.BinaryDataN2SmInformation != nil {
	// 				ngap_message.AppendPDUSessionResourceSetupListHOReq(&pduSessionReqList, pduSessionId,
	// 					smContext.Snssai(), response.BinaryDataN2SmInformation)
	// 			}
	// 		}
	// 	}
	// 	if len(pduSessionReqList.List) == 0 {
	// 		sourceUe.Log.Info("Handle Handover Preparation Failure [HoFailure In Target5GC NgranNode Or TargetSystem]")
	// 		cause = &ngapType.Cause{
	// 			Present: ngapType.CausePresentRadioNetwork,
	// 			RadioNetwork: &ngapType.CauseRadioNetwork{
	// 				Value: ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem,
	// 			},
	// 		}
	// 		ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
	// 		return
	// 	}
	// 	// Update NH
	// 	amfUe.UpdateNH()
	// 	ngap_message.SendHandoverRequest(sourceUe, targetRan, *cause, pduSessionReqList,
	// 		*sourceToTargetTransparentContainer, false)
	// }
}

//TODO
func HandleHandoverCancel(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var cause *ngapType.Cause

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("Initiating Message is nil")
		return
	}
	HandoverCancel := initiatingMessage.Value.HandoverCancel
	if HandoverCancel == nil {
		// lbConn.Log.Error("Handover Cancel is nil")
		return
	}

	// lbConn.Log.Info("Handle Handover Cancel")
	for i := 0; i < len(HandoverCancel.ProtocolIEs.List); i++ {
		ie := HandoverCancel.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AMFUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RANUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			// lbConn.Log.Trace("Decode IE Cause")
			if cause == nil {
				// lbConn.Log.Error(cause, "cause is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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

	// sourceUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	// if sourceUe == nil {
	// 	lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
	// 	cause := ngapType.Cause{
	// 		Present: ngapType.CausePresentRadioNetwork,
	// 		RadioNetwork: &ngapType.CauseRadioNetwork{
	// 			Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
	// 		},
	// 	}
	// 	ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, &cause, nil)
	// 	return
	// }

	// if sourceUe.AmfUeNgapId != aMFUENGAPID.Value {
	// 	lbConn.Log.Warnf("Conflict AMF_UE_NGAP_ID : %d != %d", sourceUe.AmfUeNgapId, aMFUENGAPID.Value)
	// }
	// lbConn.Log.Tracef("Source: RAN_UE_NGAP_ID[%d] AMF_UE_NGAP_ID[%d]", sourceUe.RanUeNgapId, sourceUe.AmfUeNgapId)

	// causePresent := ngapType.CausePresentRadioNetwork
	// causeValue := ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem
	// if cause != nil {
	// 	causePresent, causeValue = printAndGetCause(lbConn, cause)
	// }
	// targetUe := sourceUe.TargetUe
	// if targetUe == nil {
	// 	// Described in (23.502 4.11.1.2.3) step 2
	// 	// Todo : send to T-AMF invoke Namf_UeContextReleaseRequest(targetUe)
	// 	lbConn.Log.Error("N2 Handover between AMF has not been implemented yet")
	// } else {
	// 	lbConn.Log.Tracef("Target : RAN_UE_NGAP_ID[%d] AMF_UE_NGAP_ID[%d]", targetUe.RanUeNgapId, targetUe.AmfUeNgapId)
	// 	amfUe := sourceUe.AmfUe
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
	// 				sourceUe.Log.Errorf("Send UpdateSmContextN2HandoverCanceled Error for PduSessionId[%d]", pduSessionID)
	// 			}
	// 			return true
	// 		})
	// 	}
	// 	ngap_message.SendUEContextReleaseCommand(targetUe, context.UeContextReleaseHandover, causePresent, causeValue)
	// 	ngap_message.SendHandoverCancelAcknowledge(sourceUe, nil)
	// }
}

func HandleUplinkRanStatusTransfer(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	
	// var ranUe *context.RanUe

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // ignore
	if initiatingMessage == nil {
		// lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	uplinkRanStatusTransfer := initiatingMessage.Value.UplinkRANStatusTransfer
	if uplinkRanStatusTransfer == nil {
		// lbConn.Log.Error("UplinkRanStatusTransfer is nil")
		return
	}

	// lbConn.Log.Info("Handle Uplink Ran Status Transfer")

	for _, ie := range uplinkRanStatusTransfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

func HandleNasNonDeliveryIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	nASNonDeliveryIndication := initiatingMessage.Value.NASNonDeliveryIndication
	if nASNonDeliveryIndication == nil {
		// lbConn.Log.Error("NASNonDeliveryIndication is nil")
		return
	}

	// lbConn.Log.Info("Handle Nas Non Delivery Indication")

	for _, ie := range nASNonDeliveryIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

//Todo
func HandleRanConfigurationUpdate(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var rANNodeName *ngapType.RANNodeName
	var supportedTAList *ngapType.SupportedTAList
	var pagingDRX *ngapType.PagingDRX

	LB = *context.LB_Self()

	// var cause ngapType.Cause

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}

	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("Initiating Message is nil")
		return
	}
	rANConfigurationUpdate := initiatingMessage.Value.RANConfigurationUpdate
	if rANConfigurationUpdate == nil {
		// lbConn.Log.Error("RAN Configuration is nil")
		return
	}
	// lbConn.Log.Info("Handle Ran Configuration Update")
	for i := 0; i < len(rANConfigurationUpdate.ProtocolIEs.List); i++ {
		ie := rANConfigurationUpdate.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANNodeName:
			rANNodeName = ie.Value.RANNodeName
			if rANNodeName == nil {
				// lbConn.Log.Error("RAN Node Name is nil")
				return
			}
			// lbConn.Log.Tracef("Decode IE RANNodeName = [%s]", rANNodeName.Value)
		case ngapType.ProtocolIEIDSupportedTAList:
			supportedTAList = ie.Value.SupportedTAList
			// lbConn.Log.Trace("Decode IE SupportedTAList")
			if supportedTAList == nil {
				// lbConn.Log.Error("Supported TA List is nil")
				return
			}
		case ngapType.ProtocolIEIDDefaultPagingDRX:
			pagingDRX = ie.Value.DefaultPagingDRX
			if pagingDRX == nil {
				// lbConn.Log.Error("PagingDRX is nil")
				return
			}
			// lbConn.Log.Tracef("Decode IE PagingDRX = [%d]", pagingDRX.Value)
		}
	}

	// TODO

	// for i := 0; i < len(supportedTAList.List); i++ {
	// 	supportedTAItem := supportedTAList.List[i]
	// 	tac := hex.EncodeToString(supportedTAItem.TAC.Value)
	// 	capOfSupportTai := cap(lbConn.SupportedTAList)
	// 	for j := 0; j < len(supportedTAItem.BroadcastPLMNList.List); j++ {
	// 		supportedTAI := context.NewSupportedTAI()
	// 		supportedTAI.Tai.Tac = tac
	// 		broadcastPLMNItem := supportedTAItem.BroadcastPLMNList.List[j]
	// 		plmnId := ngapConvert.PlmnIdToModels(broadcastPLMNItem.PLMNIdentity)
	// 		supportedTAI.Tai.PlmnId = &plmnId
	// 		capOfSNssaiList := cap(supportedTAI.SNssaiList)
	// 		for k := 0; k < len(broadcastPLMNItem.TAISliceSupportList.List); k++ {
	// 			tAISliceSupportItem := broadcastPLMNItem.TAISliceSupportList.List[k]
	// 			if len(supportedTAI.SNssaiList) < capOfSNssaiList {
	// 				supportedTAI.SNssaiList = append(supportedTAI.SNssaiList, ngapConvert.SNssaiToModels(tAISliceSupportItem.SNSSAI))
	// 			} else {
	// 				break
	// 			}
	// 		}
	// 		lbConn.Log.Tracef("PLMN_ID[MCC:%s MNC:%s] TAC[%s]", plmnId.Mcc, plmnId.Mnc, tac)
	// 		if len(lbConn.SupportedTAList) < capOfSupportTai {
	// 			lbConn.SupportedTAList = append(lbConn.SupportedTAList, supportedTAI)
	// 		} else {
	// 			break
	// 		}
	// 	}
	// }

	// if len(lbConn.SupportedTAList) == 0 {
	// 	lbConn.Log.Warn("RanConfigurationUpdate failure: No supported TA exist in RanConfigurationUpdate")
	// 	cause.Present = ngapType.CausePresentMisc
	// 	cause.Misc = &ngapType.CauseMisc{
	// 		Value: ngapType.CauseMiscPresentUnspecified,
	// 	}
	// } else {
	// 	var found bool
	// 	for i, tai := range lbConn.SupportedTAList {
	// 		if context.InTaiList(tai.Tai, context.AMF_Self().SupportTaiLists) {
	// 			lbConn.Log.Tracef("SERVED_TAI_INDEX[%d]", i)
	// 			found = true
	// 			break
	// 		}
	// 	}
	// 	if !found {
	// 		lbConn.Log.Warn("RanConfigurationUpdate failure: Cannot find Served TAI in AMF")
	// 		cause.Present = ngapType.CausePresentMisc
	// 		cause.Misc = &ngapType.CauseMisc{
	// 			Value: ngapType.CauseMiscPresentUnknownPLMN,
	// 		}
	// 	}
	// }

	// if cause.Present == ngapType.CausePresentNothing {
	// 	lbConn.Log.Info("Handle RanConfigurationUpdateAcknowledge")
	// 	ngap_message.SendRanConfigurationUpdateAcknowledge(lbConn, nil)
	// } else {
	// 	lbConn.Log.Info("Handle RanConfigurationUpdateAcknowledgeFailure")
	// 	ngap_message.SendRanConfigurationUpdateFailure(lbConn, cause, nil)
	// }
}

//TODO
func HandleUplinkRanConfigurationTransfer(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var sONConfigurationTransferUL *ngapType.SONConfigurationTransfer

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	uplinkRANConfigurationTransfer := initiatingMessage.Value.UplinkRANConfigurationTransfer
	if uplinkRANConfigurationTransfer == nil {
		// lbConn.Log.Error("ErrorIndication is nil")
		return
	}

	for _, ie := range uplinkRANConfigurationTransfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDSONConfigurationTransferUL: // optional, ignore
			sONConfigurationTransferUL = ie.Value.SONConfigurationTransferUL
			// lbConn.Log.Trace("Decode IE SONConfigurationTransferUL")
			if sONConfigurationTransferUL == nil {
				// lbConn.Log.Warn("sONConfigurationTransferUL is nil")
			}
		}
	}

	// if sONConfigurationTransferUL != nil {
	// 	targetRanNodeID := ngapConvert.RanIdToModels(sONConfigurationTransferUL.TargetRANNodeID.GlobalRANNodeID)

	// 	if targetRanNodeID.GNbId.GNBValue != "" {
	// 		// lbConn.Log.Tracef("targerRanID [%s]", targetRanNodeID.GNbId.GNBValue)
	// 	}

	// 	aMFSelf := context.AMF_Self()

	// 	targetRan, ok := aMFSelf.AmfRanFindByRanID(targetRanNodeID)
	// 	if !ok {
	// 		// lbConn.Log.Warn("targetRan is nil")
	// 	}

	// 	ngap_message.SendDownlinkRanConfigurationTransfer(targetRan, sONConfigurationTransferUL)
	// }
}

func HandleUplinkUEAssociatedNRPPATransport(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	uplinkUEAssociatedNRPPaTransport := initiatingMessage.Value.UplinkUEAssociatedNRPPaTransport
	if uplinkUEAssociatedNRPPaTransport == nil {
		// lbConn.Log.Error("uplinkUEAssociatedNRPPaTransport is nil")
		return
	}

	// lbConn.Log.Info("Handle Uplink UE Associated NRPPA Transpor")

	for _, ie := range uplinkUEAssociatedNRPPaTransport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

// TODO
func HandleUplinkNonUEAssociatedNRPPATransport(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var routingID *ngapType.RoutingID
	var nRPPaPDU *ngapType.NRPPaPDU

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("Initiating Message is nil")
		return
	}
	uplinkNonUEAssociatedNRPPATransport := initiatingMessage.Value.UplinkNonUEAssociatedNRPPaTransport
	if uplinkNonUEAssociatedNRPPATransport == nil {
		// lbConn.Log.Error("Uplink Non UE Associated NRPPA Transport is nil")
		return
	}

	// lbConn.Log.Info("Handle Uplink Non UE Associated NRPPA Transport")

	for i := 0; i < len(uplinkNonUEAssociatedNRPPATransport.ProtocolIEs.List); i++ {
		ie := uplinkNonUEAssociatedNRPPATransport.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRoutingID:
			routingID = ie.Value.RoutingID
			// lbConn.Log.Trace("Decode IE RoutingID")

		case ngapType.ProtocolIEIDNRPPaPDU:
			nRPPaPDU = ie.Value.NRPPaPDU
			// lbConn.Log.Trace("Decode IE NRPPaPDU")
		}
	}

	if routingID == nil {
		// lbConn.Log.Error("RoutingID is nil")
		return
	}
	// Forward routingID to LMF
	// Described in (23.502 4.13.5.6)

	if nRPPaPDU == nil {
		// lbConn.Log.Error("NRPPaPDU is nil")
		return
	}
	// TODO: Forward NRPPaPDU to LMF
}

func HandleLocationReport(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	locationReport := initiatingMessage.Value.LocationReport
	if locationReport == nil {
		// lbConn.Log.Error("LocationReport is nil")
		return
	}

	// lbConn.Log.Info("Handle Location Report")
	for _, ie := range locationReport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

func HandleUERadioCapabilityInfoIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("Initiating Message is nil")
		return
	}
	uERadioCapabilityInfoIndication := initiatingMessage.Value.UERadioCapabilityInfoIndication
	if uERadioCapabilityInfoIndication == nil {
		// lbConn.Log.Error("UERadioCapabilityInfoIndication is nil")
		return
	}

	// lbConn.Log.Info("Handle UE Radio Capability Info Indication")

	for i := 0; i < len(uERadioCapabilityInfoIndication.ProtocolIEs.List); i++ {
		ie := uERadioCapabilityInfoIndication.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

//TODO
func HandleAMFconfigurationUpdateFailure(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	
	LB = *context.LB_Self()
	
	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		// lbConn.Log.Error("Unsuccessful Message is nil")
		return
	}

	AMFconfigurationUpdateFailure := unsuccessfulOutcome.Value.AMFConfigurationUpdateFailure
	if AMFconfigurationUpdateFailure == nil {
		// lbConn.Log.Error("AMFConfigurationUpdateFailure is nil")
		return
	}

	// lbConn.Log.Info("Handle AMF Confioguration Update Failure")

	for _, ie := range AMFconfigurationUpdateFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			// lbConn.Log.Trace("Decode IE Cause")
			if cause == nil {
				// lbConn.Log.Error("Cause is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			// lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	//	TODO: Time To Wait

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

//TODO
func HandleAMFconfigurationUpdateAcknowledge(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFTNLAssociationSetupList *ngapType.AMFTNLAssociationSetupList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	var aMFTNLAssociationFailedToSetupList *ngapType.TNLAssociationList
	
	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		// lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	aMFConfigurationUpdateAcknowledge := successfulOutcome.Value.AMFConfigurationUpdateAcknowledge
	if aMFConfigurationUpdateAcknowledge == nil {
		// lbConn.Log.Error("AMFConfigurationUpdateAcknowledge is nil")
		return
	}

	// lbConn.Log.Info("Handle AMF Configuration Update Acknowledge")

	for i := 0; i < len(aMFConfigurationUpdateAcknowledge.ProtocolIEs.List); i++ {
		ie := aMFConfigurationUpdateAcknowledge.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFTNLAssociationSetupList:
			aMFTNLAssociationSetupList = ie.Value.AMFTNLAssociationSetupList
			// lbConn.Log.Trace("Decode IE AMFTNLAssociationSetupList")
			if aMFTNLAssociationSetupList == nil {
				// lbConn.Log.Error("AMFTNLAssociationSetupList is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			// lbConn.Log.Trace("Decode IE Criticality Diagnostics")

		case ngapType.ProtocolIEIDAMFTNLAssociationFailedToSetupList:
			aMFTNLAssociationFailedToSetupList = ie.Value.AMFTNLAssociationFailedToSetupList
			// lbConn.Log.Trace("Decode IE AMFTNLAssociationFailedToSetupList")
			if aMFTNLAssociationFailedToSetupList == nil {
				// lbConn.Log.Error("AMFTNLAssociationFailedToSetupList is nil")
				return
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

func HandleErrorIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	
	LB = *context.LB_Self()

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		// lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	errorIndication := initiatingMessage.Value.ErrorIndication
	if errorIndication == nil {
		// lbConn.Log.Error("ErrorIndication is nil")
		return
	}

	for _, ie := range errorIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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
}

func HandleCellTrafficTrace(lbConn *context.LBConn, message *ngapType.NGAPPDU, m2 []byte) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	LB = *context.LB_Self()
	
	// var ranUe *context.RanUe

	// var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if lbConn == nil {
		// logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		// lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // ignore
	if initiatingMessage == nil {
		// lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	cellTrafficTrace := initiatingMessage.Value.CellTrafficTrace
	if cellTrafficTrace == nil {
		// lbConn.Log.Error("CellTrafficTrace is nil")
		return
	}

	// lbConn.Log.Info("Handle Cell Traffic Trace")

	for _, ie := range cellTrafficTrace.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				fmt.Println("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				fmt.Println("RanUeNgapID is nil")
				return
			}
		}
	}

	if lbConn.TypeID == context.TypeIdentAMFConn {
		amf, _ := LB.LbAmfFindByConn(lbConn.Conn)
		fmt.Println("AMF Found")
		UEs, ok := amf.FindUeByUeRanID(rANUENGAPID.Value)
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
				fmt.Println("UE Found + adding aMFUENGAPID")
				ue.UeAmfId = aMFUENGAPID.Value
				ue = ue2
			case 3: 
				fmt.Println("no matching UE Found")
				return 
			}
		}
		LB.ForwardToGnb(lbConn, m2, ue)
		return
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

	// if aMFUENGAPID == nil {
	// 	lbConn.Log.Error("AmfUeNgapID is nil")
	// 	item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDAMFUENGAPID,
	// 		ngapType.TypeOfErrorPresentMissing)
	// 	iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	// }
	// if rANUENGAPID == nil {
	// 	lbConn.Log.Error("RanUeNgapID is nil")
	// 	item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDRANUENGAPID,
	// 		ngapType.TypeOfErrorPresentMissing)
	// 	iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	// }
}

func printAndGetCause(lbConn *context.LBConn, cause *ngapType.Cause) (present int, value aper.Enumerated) {
	present = cause.Present
	switch cause.Present {
	case ngapType.CausePresentRadioNetwork:
		// lbConn.Log.Warnf("Cause RadioNetwork[%d]", cause.RadioNetwork.Value)
		value = cause.RadioNetwork.Value
	case ngapType.CausePresentTransport:
		// lbConn.Log.Warnf("Cause Transport[%d]", cause.Transport.Value)
		value = cause.Transport.Value
	case ngapType.CausePresentProtocol:
		// lbConn.Log.Warnf("Cause Protocol[%d]", cause.Protocol.Value)
		value = cause.Protocol.Value
	case ngapType.CausePresentNas:
		// lbConn.Log.Warnf("Cause Nas[%d]", cause.Nas.Value)
		value = cause.Nas.Value
	case ngapType.CausePresentMisc:
		// lbConn.Log.Warnf("Cause Misc[%d]", cause.Misc.Value)
		value = cause.Misc.Value
	default:
		// lbConn.Log.Errorf("Invalid Cause group[%d]", cause.Present)
	}
	return
}

func printCriticalityDiagnostics(lbConn *context.LBConn, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {
	// lbConn.Log.Trace("Criticality Diagnostics")

	if criticalityDiagnostics.ProcedureCriticality != nil {
		switch criticalityDiagnostics.ProcedureCriticality.Value {
		case ngapType.CriticalityPresentReject:
			// lbConn.Log.Trace("Procedure Criticality: Reject")
		case ngapType.CriticalityPresentIgnore:
			// lbConn.Log.Trace("Procedure Criticality: Ignore")
		case ngapType.CriticalityPresentNotify:
			// lbConn.Log.Trace("Procedure Criticality: Notify")
		}
	}

	if criticalityDiagnostics.IEsCriticalityDiagnostics != nil {
		for _, ieCriticalityDiagnostics := range criticalityDiagnostics.IEsCriticalityDiagnostics.List {
			// lbConn.Log.Tracef("IE ID: %d", ieCriticalityDiagnostics.IEID.Value)

			switch ieCriticalityDiagnostics.IECriticality.Value {
			case ngapType.CriticalityPresentReject:
				// lbConn.Log.Trace("Criticality Reject")
			case ngapType.CriticalityPresentNotify:
				// lbConn.Log.Trace("Criticality Notify")
			}

			switch ieCriticalityDiagnostics.TypeOfError.Value {
			case ngapType.TypeOfErrorPresentNotUnderstood:
				// lbConn.Log.Trace("Type of error: Not understood")
			case ngapType.TypeOfErrorPresentMissing:
				// lbConn.Log.Trace("Type of error: Missing")
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
