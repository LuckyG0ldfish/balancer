package ngap

import (
	"encoding/hex"
	"strconv"

	"github.com/LuckyG0ldfish/balancer/context" 
	"github.com/LuckyG0ldfish/balancer/logger" 
	"github.com/free5gc/ngap/ngapType"

	// "github.com/free5gc/amf/consumer"
	// gmm_common "github.com/free5gc/amf/gmm/common"
	// gmm_message "github.com/free5gc/amf/gmm/message"
	// "github.com/free5gc/amf/logger"
	// "github.com/free5gc/amf/nas"
	// ngap_message "github.com/free5gc/amf/ngap/message"
	// "github.com/free5gc/aper"
	// "github.com/free5gc/nas/nasMessage"
	// libngap "github.com/free5gc/ngap"
	// "github.com/free5gc/ngap/ngapConvert"
	// "github.com/free5gc/ngap/ngapType"
	// "github.com/free5gc/openapi/models"
)


//TODO
func HandleNGSetupRequest(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var globalRANNodeID *ngapType.GlobalRANNodeID
	var rANNodeName *ngapType.RANNodeName
	var supportedTAList *ngapType.SupportedTAList
	var pagingDRX *ngapType.PagingDRX

	var cause ngapType.Cause

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		ran.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ran.Log.Error("Initiating Message is nil")
		return
	}
	nGSetupRequest := initiatingMessage.Value.NGSetupRequest
	if nGSetupRequest == nil {
		ran.Log.Error("NGSetupRequest is nil")
		return
	}
	ran.Log.Info("Handle NG Setup request")
	for i := 0; i < len(nGSetupRequest.ProtocolIEs.List); i++ {
		ie := nGSetupRequest.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDGlobalRANNodeID:
			globalRANNodeID = ie.Value.GlobalRANNodeID
			ran.Log.Trace("Decode IE GlobalRANNodeID")
			if globalRANNodeID == nil {
				ran.Log.Error("GlobalRANNodeID is nil")
				return
			}
		case ngapType.ProtocolIEIDSupportedTAList:
			supportedTAList = ie.Value.SupportedTAList
			ran.Log.Trace("Decode IE SupportedTAList")
			if supportedTAList == nil {
				ran.Log.Error("SupportedTAList is nil")
				return
			}
		case ngapType.ProtocolIEIDRANNodeName:
			rANNodeName = ie.Value.RANNodeName
			ran.Log.Trace("Decode IE RANNodeName")
			if rANNodeName == nil {
				ran.Log.Error("RANNodeName is nil")
				return
			}
		case ngapType.ProtocolIEIDDefaultPagingDRX:
			pagingDRX = ie.Value.DefaultPagingDRX
			ran.Log.Trace("Decode IE DefaultPagingDRX")
			if pagingDRX == nil {
				ran.Log.Error("DefaultPagingDRX is nil")
				return
			}
		}
	}

	ran.SetRanId(globalRANNodeID)
	if rANNodeName != nil {
		ran.Name = rANNodeName.Value
	}
	if pagingDRX != nil {
		ran.Log.Tracef("PagingDRX[%d]", pagingDRX.Value)
	}

	for i := 0; i < len(supportedTAList.List); i++ {
		supportedTAItem := supportedTAList.List[i]
		tac := hex.EncodeToString(supportedTAItem.TAC.Value)
		capOfSupportTai := cap(ran.SupportedTAList)
		for j := 0; j < len(supportedTAItem.BroadcastPLMNList.List); j++ {
			supportedTAI := context.NewSupportedTAI()
			supportedTAI.Tai.Tac = tac
			broadcastPLMNItem := supportedTAItem.BroadcastPLMNList.List[j]
			plmnId := ngapConvert.PlmnIdToModels(broadcastPLMNItem.PLMNIdentity)
			supportedTAI.Tai.PlmnId = &plmnId
			capOfSNssaiList := cap(supportedTAI.SNssaiList)
			for k := 0; k < len(broadcastPLMNItem.TAISliceSupportList.List); k++ {
				tAISliceSupportItem := broadcastPLMNItem.TAISliceSupportList.List[k]
				if len(supportedTAI.SNssaiList) < capOfSNssaiList {
					supportedTAI.SNssaiList = append(supportedTAI.SNssaiList, ngapConvert.SNssaiToModels(tAISliceSupportItem.SNSSAI))
				} else {
					break
				}
			}
			ran.Log.Tracef("PLMN_ID[MCC:%s MNC:%s] TAC[%s]", plmnId.Mcc, plmnId.Mnc, tac)
			if len(ran.SupportedTAList) < capOfSupportTai {
				ran.SupportedTAList = append(ran.SupportedTAList, supportedTAI)
			} else {
				break
			}
		}
	}

	if len(ran.SupportedTAList) == 0 {
		ran.Log.Warn("NG-Setup failure: No supported TA exist in NG-Setup request")
		cause.Present = ngapType.CausePresentMisc
		cause.Misc = &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		}
	} else {
		var found bool
		for i, tai := range ran.SupportedTAList {
			if context.InTaiList(tai.Tai, context.AMF_Self().SupportTaiLists) {
				ran.Log.Tracef("SERVED_TAI_INDEX[%d]", i)
				found = true
				break
			}
		}
		if !found {
			ran.Log.Warn("NG-Setup failure: Cannot find Served TAI in AMF")
			cause.Present = ngapType.CausePresentMisc
			cause.Misc = &ngapType.CauseMisc{
				Value: ngapType.CauseMiscPresentUnknownPLMN,
			}
		}
	}

	if cause.Present == ngapType.CausePresentNothing {
		ngap_message.SendNGSetupResponse(ran)
	} else {
		ngap_message.SendNGSetupFailure(ran, cause)
	}
}


func HandleUplinkNasTransport(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var nASPDU *ngapType.NASPDU
	var userLocationInformation *ngapType.UserLocationInformation

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
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			// lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				// lbConn.Log.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			// lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				// lbConn.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDNASPDU:
			nASPDU = ie.Value.NASPDU
			// lbConn.Log.Trace("Decode IE NasPdu")
			if nASPDU == nil {
				// lbConn.Log.Error("nASPDU is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation:
			userLocationInformation = ie.Value.UserLocationInformation
			// lbConn.Log.Trace("Decode IE UserLocationInformation")
			if userLocationInformation == nil {
				// lbConn.Log.Error("UserLocationInformation is nil")
				return
			}
		}
	}

	ranUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		// lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
	amfUe := ranUe.AmfUe
	if amfUe == nil {
		err := ranUe.Remove()
		if err != nil {
			// lbConn.Log.Errorf(err.Error())
		}
		// lbConn.Log.Errorf("No UE Context of RanUe with RANUENGAPID[%d] AMFUENGAPID[%d] ",
			// rANUENGAPID.Value, aMFUENGAPID.Value)
		return
	}

	ranUe.Log.Infof("Uplink NAS Transport (RAN UE NGAP ID: %d)", ranUe.RanUeNgapId)

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	// nas.HandleNAS(ranUe, ngapType.ProcedureCodeUplinkNASTransport, nASPDU.Value)
}

func HandleNGReset(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var cause *ngapType.Cause
	var resetType *ngapType.ResetType

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("Initiating Message is nil")
		return
	}
	nGReset := initiatingMessage.Value.NGReset
	if nGReset == nil {
		lbConn.Log.Error("NGReset is nil")
		return
	}

	lbConn.Log.Info("Handle NG Reset")

	for _, ie := range nGReset.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			lbConn.Log.Trace("Decode IE Cause")
			if cause == nil {
				lbConn.Log.Error("Cause is nil")
				return
			}
		case ngapType.ProtocolIEIDResetType:
			resetType = ie.Value.ResetType
			lbConn.Log.Trace("Decode IE ResetType")
			if resetType == nil {
				lbConn.Log.Error("ResetType is nil")
				return
			}
		}
	}

	printAndGetCause(lbConn, cause)

	switch resetType.Present {
	case ngapType.ResetTypePresentNGInterface:
		lbConn.Log.Trace("ResetType Present: NG Interface")
		lbConn.RemoveAllUeInRan()
		ngap_message.SendNGResetAcknowledge(lbConn, nil, nil)
	case ngapType.ResetTypePresentPartOfNGInterface:
		lbConn.Log.Trace("ResetType Present: Part of NG Interface")

		partOfNGInterface := resetType.PartOfNGInterface
		if partOfNGInterface == nil {
			lbConn.Log.Error("PartOfNGInterface is nil")
			return
		}

		var ranUe *context.RanUe

		for _, ueAssociatedLogicalNGConnectionItem := range partOfNGInterface.List {
			if ueAssociatedLogicalNGConnectionItem.AMFUENGAPID != nil {
				lbConn.Log.Tracef("AmfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
				for _, ue := range lbConn.RanUeList {
					if ue.AmfUeNgapId == ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value {
						ranUe = ue
						break
					}
				}
			} else if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
				lbConn.Log.Tracef("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
				ranUe = lbConn.RanUeFindByRanUeNgapID(ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
			}

			if ranUe == nil {
				lbConn.Log.Warn("Cannot not find UE Context")
				if ueAssociatedLogicalNGConnectionItem.AMFUENGAPID != nil {
					lbConn.Log.Warnf("AmfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
				}
				if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
					lbConn.Log.Warnf("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
				}
			}

			err := ranUe.Remove()
			if err != nil {
				lbConn.Log.Error(err.Error())
			}
		}
		ngap_message.SendNGResetAcknowledge(lbConn, partOfNGInterface, nil)
	default:
		lbConn.Log.Warnf("Invalid ResetType[%d]", resetType.Present)
	}
}

func HandleNGResetAcknowledge(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var uEAssociatedLogicalNGConnectionList *ngapType.UEAssociatedLogicalNGConnectionList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	nGResetAcknowledge := successfulOutcome.Value.NGResetAcknowledge
	if nGResetAcknowledge == nil {
		lbConn.Log.Error("NGResetAcknowledge is nil")
		return
	}

	lbConn.Log.Info("Handle NG Reset Acknowledge")

	for _, ie := range nGResetAcknowledge.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDUEAssociatedLogicalNGConnectionList:
			uEAssociatedLogicalNGConnectionList = ie.Value.UEAssociatedLogicalNGConnectionList
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
		}
	}

	if uEAssociatedLogicalNGConnectionList != nil {
		lbConn.Log.Tracef("%d UE association(s) has been reset", len(uEAssociatedLogicalNGConnectionList.List))
		for i, item := range uEAssociatedLogicalNGConnectionList.List {
			if item.AMFUENGAPID != nil && item.RANUENGAPID != nil {
				lbConn.Log.Tracef("%d: AmfUeNgapID[%d] RanUeNgapID[%d]", i+1, item.AMFUENGAPID.Value, item.RANUENGAPID.Value)
			} else if item.AMFUENGAPID != nil {
				lbConn.Log.Tracef("%d: AmfUeNgapID[%d] RanUeNgapID[-1]", i+1, item.AMFUENGAPID.Value)
			} else if item.RANUENGAPID != nil {
				lbConn.Log.Tracef("%d: AmfUeNgapID[-1] RanUeNgapID[%d]", i+1, item.RANUENGAPID.Value)
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
	var userLocationInformation *ngapType.UserLocationInformation
	var infoOnRecommendedCellsAndRANNodesForPaging *ngapType.InfoOnRecommendedCellsAndRANNodesForPaging
	var pDUSessionResourceList *ngapType.PDUSessionResourceListCxtRelCpl
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	uEContextReleaseComplete := successfulOutcome.Value.UEContextReleaseComplete
	if uEContextReleaseComplete == nil {
		lbConn.Log.Error("NGResetAcknowledge is nil")
		return
	}

	lbConn.Log.Info("Handle UE Context Release Complete")

	for _, ie := range uEContextReleaseComplete.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation:
			userLocationInformation = ie.Value.UserLocationInformation
			lbConn.Log.Trace("Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDInfoOnRecommendedCellsAndRANNodesForPaging:
			infoOnRecommendedCellsAndRANNodesForPaging = ie.Value.InfoOnRecommendedCellsAndRANNodesForPaging
			lbConn.Log.Trace("Decode IE InfoOnRecommendedCellsAndRANNodesForPaging")
			if infoOnRecommendedCellsAndRANNodesForPaging != nil {
				lbConn.Log.Warn("IE infoOnRecommendedCellsAndRANNodesForPaging is not support")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceListCxtRelCpl:
			pDUSessionResourceList = ie.Value.PDUSessionResourceListCxtRelCpl
			lbConn.Log.Trace("Decode IE PDUSessionResourceList")
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	ranUe := context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No RanUe Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, &cause, nil)
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}
	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		lbConn.Log.Infof("Release UE Context : RanUe[AmfUeNgapId: %d]", ranUe.AmfUeNgapId)
		err := ranUe.Remove()
		if err != nil {
			lbConn.Log.Errorln(err.Error())
		}
		return
	}
	// TODO: AMF shall, if supported, store it and may use it for subsequent paging
	if infoOnRecommendedCellsAndRANNodesForPaging != nil {
		amfUe.InfoOnRecommendedCellsAndRanNodesForPaging = new(context.InfoOnRecommendedCellsAndRanNodesForPaging)

		recommendedCells := &amfUe.InfoOnRecommendedCellsAndRanNodesForPaging.RecommendedCells
		for _, item := range infoOnRecommendedCellsAndRANNodesForPaging.RecommendedCellsForPaging.RecommendedCellList.List {
			recommendedCell := context.RecommendedCell{}

			switch item.NGRANCGI.Present {
			case ngapType.NGRANCGIPresentNRCGI:
				recommendedCell.NgRanCGI.Present = context.NgRanCgiPresentNRCGI
				recommendedCell.NgRanCGI.NRCGI = new(models.Ncgi)
				plmnID := ngapConvert.PlmnIdToModels(item.NGRANCGI.NRCGI.PLMNIdentity)
				recommendedCell.NgRanCGI.NRCGI.PlmnId = &plmnID
				recommendedCell.NgRanCGI.NRCGI.NrCellId = ngapConvert.BitStringToHex(&item.NGRANCGI.NRCGI.NRCellIdentity.Value)
			case ngapType.NGRANCGIPresentEUTRACGI:
				recommendedCell.NgRanCGI.Present = context.NgRanCgiPresentEUTRACGI
				recommendedCell.NgRanCGI.EUTRACGI = new(models.Ecgi)
				plmnID := ngapConvert.PlmnIdToModels(item.NGRANCGI.EUTRACGI.PLMNIdentity)
				recommendedCell.NgRanCGI.EUTRACGI.PlmnId = &plmnID
				recommendedCell.NgRanCGI.EUTRACGI.EutraCellId = ngapConvert.BitStringToHex(
					&item.NGRANCGI.EUTRACGI.EUTRACellIdentity.Value)
			}

			if item.TimeStayedInCell != nil {
				recommendedCell.TimeStayedInCell = new(int64)
				*recommendedCell.TimeStayedInCell = *item.TimeStayedInCell
			}

			*recommendedCells = append(*recommendedCells, recommendedCell)
		}

		recommendedRanNodes := &amfUe.InfoOnRecommendedCellsAndRanNodesForPaging.RecommendedRanNodes
		ranNodeList := infoOnRecommendedCellsAndRANNodesForPaging.RecommendRANNodesForPaging.RecommendedRANNodeList.List
		for _, item := range ranNodeList {
			recommendedRanNode := context.RecommendRanNode{}

			switch item.AMFPagingTarget.Present {
			case ngapType.AMFPagingTargetPresentGlobalRANNodeID:
				recommendedRanNode.Present = context.RecommendRanNodePresentRanNode
				recommendedRanNode.GlobalRanNodeId = new(models.GlobalRanNodeId)
				// TODO: recommendedRanNode.GlobalRanNodeId = ngapConvert.RanIdToModels(item.AMFPagingTarget.GlobalRANNodeID)
			case ngapType.AMFPagingTargetPresentTAI:
				recommendedRanNode.Present = context.RecommendRanNodePresentTAI
				tai := ngapConvert.TaiToModels(*item.AMFPagingTarget.TAI)
				recommendedRanNode.Tai = &tai
			}
			*recommendedRanNodes = append(*recommendedRanNodes, recommendedRanNode)
		}
	}

	// for each pduSessionID invoke Nsmf_PDUSession_UpdateSMContext Request
	var cause context.CauseAll
	if tmp, exist := amfUe.ReleaseCause[lbConn.AnType]; exist {
		cause = *tmp
	}
	if amfUe.State[lbConn.AnType].Is(context.Registered) {
		ranUe.Log.Info("Rel Ue Context in GMM-Registered")
		if cause.NgapCause != nil && pDUSessionResourceList != nil {
			for _, pduSessionReourceItem := range pDUSessionResourceList.List {
				pduSessionID := int32(pduSessionReourceItem.PDUSessionID.Value)
				smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
				if !ok {
					ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
					// TODO: Check if doing error handling here
					continue
				}
				response, _, _, err := consumer.SendUpdateSmContextDeactivateUpCnxState(amfUe, smContext, cause)
				if err != nil {
					lbConn.Log.Errorf("Send Update SmContextDeactivate UpCnxState Error[%s]", err.Error())
				} else if response == nil {
					lbConn.Log.Errorln("Send Update SmContextDeactivate UpCnxState Error")
				}
			}
		}
	}

	// Remove UE N2 Connection
	delete(amfUe.ReleaseCause, ran.AnType)
	switch ranUe.ReleaseAction {
	case context.UeContextN2NormalRelease:
		lbConn.Log.Infof("Release UE[%s] Context : N2 Connection Release", amfUe.Supi)
		// amfUe.DetachRanUe(ran.AnType)
		err := ranUe.Remove()
		if err != nil {
			lbConn.Log.Errorln(err.Error())
		}
	case context.UeContextReleaseUeContext:
		lbConn.Log.Infof("Release UE[%s] Context : Release Ue Context", amfUe.Supi)
		gmm_common.RemoveAmfUe(amfUe)
	case context.UeContextReleaseHandover:
		lbConn.Log.Infof("Release UE[%s] Context : Release for Handover", amfUe.Supi)
		// TODO: it's a workaround, need to fix it.
		targetRanUe := context.AMF_Self().RanUeFindByAmfUeNgapID(ranUe.TargetUe.AmfUeNgapId)

		context.DetachSourceUeTargetUe(ranUe)
		err := ranUe.Remove()
		if err != nil {
			lbConn.Log.Errorln(err.Error())
		}
		amfUe.AttachRanUe(targetRanUe)
		// Todo: remove indirect tunnel
	default:
		lbConn.Log.Errorf("Invalid Release Action[%d]", ranUe.ReleaseAction)
	}
}

func HandlePDUSessionResourceReleaseResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceReleasedList *ngapType.PDUSessionResourceReleasedListRelRes
	var userLocationInformation *ngapType.UserLocationInformation
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceReleaseResponse := successfulOutcome.Value.PDUSessionResourceReleaseResponse
	if pDUSessionResourceReleaseResponse == nil {
		lbConn.Log.Error("PDUSessionResourceReleaseResponse is nil")
		return
	}

	lbConn.Log.Info("Handle PDU Session Resource Release Response")

	for _, ie := range pDUSessionResourceReleaseResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUENgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDPDUSessionResourceReleasedListRelRes:
			pDUSessionResourceReleasedList = ie.Value.PDUSessionResourceReleasedListRelRes
			lbConn.Log.Trace("Decode IE PDUSessionResourceReleasedList")
			if pDUSessionResourceReleasedList == nil {
				lbConn.Log.Error("PDUSessionResourceReleasedList is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation:
			userLocationInformation = ie.Value.UserLocationInformation
			lbConn.Log.Trace("Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	ranUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ranUe.Log.Error("amfUe is nil")
		return
	}
	if pDUSessionResourceReleasedList != nil {
		ranUe.Log.Trace("Send PDUSessionResourceReleaseResponseTransfer to SMF")

		for _, item := range pDUSessionResourceReleasedList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceReleaseResponseTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			_, responseErr, problemDetail, err := consumer.SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_REL_RSP, transfer)
			// TODO: error handling
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceReleaseResponse] Error: %+v", err)
			} else if responseErr != nil && responseErr.JsonData.Error != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceReleaseResponse] Error: %+v",
					responseErr.JsonData.Error.Cause)
			} else if problemDetail != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceReleaseResponse] Failed: %+v", problemDetail)
			}
		}
	}
}

func HandleUERadioCapabilityCheckResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var iMSVoiceSupportIndicator *ngapType.IMSVoiceSupportIndicator
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	var ranUe *context.RanUe

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}

	uERadioCapabilityCheckResponse := successfulOutcome.Value.UERadioCapabilityCheckResponse
	if uERadioCapabilityCheckResponse == nil {
		lbConn.Log.Error("UERadioCapabilityCheckResponse is nil")
		return
	}
	lbConn.Log.Info("Handle UE Radio Capability Check Response")

	for i := 0; i < len(uERadioCapabilityCheckResponse.ProtocolIEs.List); i++ {
		ie := uERadioCapabilityCheckResponse.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDIMSVoiceSupportIndicator:
			iMSVoiceSupportIndicator = ie.Value.IMSVoiceSupportIndicator
			lbConn.Log.Trace("Decode IE IMSVoiceSupportIndicator")
			if iMSVoiceSupportIndicator == nil {
				lbConn.Log.Error("iMSVoiceSupportIndicator is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	ranUe = lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	// TODO: handle iMSVoiceSupportIndicator

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

func HandleLocationReportingFailureIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var ranUe *context.RanUe

	var cause *ngapType.Cause

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("Initiating Message is nil")
		return
	}
	locationReportingFailureIndication := initiatingMessage.Value.LocationReportingFailureIndication
	if locationReportingFailureIndication == nil {
		lbConn.Log.Error("LocationReportingFailureIndication is nil")
		return
	}

	lbConn.Log.Info("Handle Location Reporting Failure Indication")

	for i := 0; i < len(locationReportingFailureIndication.ProtocolIEs.List); i++ {
		ie := locationReportingFailureIndication.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			lbConn.Log.Trace("Decode IE Cause")
			if cause == nil {
				lbConn.Log.Error("Cause is nil")
				return
			}
		}
	}

	printAndGetCause(lbConn, cause)

	ranUe = lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
}

func HandleInitialUEMessage(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	amfSelf := context.AMF_Self()

	var rANUENGAPID *ngapType.RANUENGAPID
	var nASPDU *ngapType.NASPDU
	var userLocationInformation *ngapType.UserLocationInformation
	var rRCEstablishmentCause *ngapType.RRCEstablishmentCause
	var fiveGSTMSI *ngapType.FiveGSTMSI
	// var aMFSetID *ngapType.AMFSetID
	var uEContextRequest *ngapType.UEContextRequest
	// var allowedNSSAI *ngapType.AllowedNSSAI

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("Initiating Message is nil")
		return
	}
	initialUEMessage := initiatingMessage.Value.InitialUEMessage
	if initialUEMessage == nil {
		lbConn.Log.Error("InitialUEMessage is nil")
		return
	}

	lbConn.Log.Info("Handle Initial UE Message")

	for _, ie := range initialUEMessage.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
					ngapType.ProtocolIEIDRANUENGAPID, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDNASPDU: // reject
			nASPDU = ie.Value.NASPDU
			lbConn.Log.Trace("Decode IE NasPdu")
			if nASPDU == nil {
				lbConn.Log.Error("NasPdu is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDNASPDU,
					ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // reject
			userLocationInformation = ie.Value.UserLocationInformation
			lbConn.Log.Trace("Decode IE UserLocationInformation")
			if userLocationInformation == nil {
				lbConn.Log.Error("UserLocationInformation is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
					ngapType.ProtocolIEIDUserLocationInformation, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRRCEstablishmentCause: // ignore
			rRCEstablishmentCause = ie.Value.RRCEstablishmentCause
			lbConn.Log.Trace("Decode IE RRCEstablishmentCause")
		case ngapType.ProtocolIEIDFiveGSTMSI: // optional, reject
			fiveGSTMSI = ie.Value.FiveGSTMSI
			lbConn.Log.Trace("Decode IE 5G-S-TMSI")
		case ngapType.ProtocolIEIDAMFSetID: // optional, ignore
			// aMFSetID = ie.Value.AMFSetID
			lbConn.Log.Trace("Decode IE AmfSetID")
		case ngapType.ProtocolIEIDUEContextRequest: // optional, ignore
			uEContextRequest = ie.Value.UEContextRequest
			lbConn.Log.Trace("Decode IE UEContextRequest")
		case ngapType.ProtocolIEIDAllowedNSSAI: // optional, reject
			// allowedNSSAI = ie.Value.AllowedNSSAI
			lbConn.Log.Trace("Decode IE Allowed NSSAI")
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		lbConn.Log.Trace("Has missing reject IE(s)")

		procedureCode := ngapType.ProcedureCodeInitialUEMessage
		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
		procedureCriticality := ngapType.CriticalityPresentIgnore
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality,
			&iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(lbConn, nil, rANUENGAPID, nil, &criticalityDiagnostics)
	}

	ranUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe != nil && ranUe.AmfUe == nil {
		err := ranUe.Remove()
		if err != nil {
			lbConn.Log.Errorln(err.Error())
		}
		ranUe = nil
	}
	if ranUe == nil {
		var err error
		ranUe, err = lbConn.NewRanUe(rANUENGAPID.Value)
		if err != nil {
			lbConn.Log.Errorf("NewRanUe Error: %+v", err)
		}
		lbConn.Log.Debugf("New RanUe [RanUeNgapID: %d]", ranUe.RanUeNgapId)

		if fiveGSTMSI != nil {
			ranUe.Log.Debug("Receive 5G-S-TMSI")

			servedGuami := amfSelf.ServedGuamiList[0]

			// <5G-S-TMSI> := <AMF Set ID><AMF Pointer><5G-TMSI>
			// GUAMI := <MCC><MNC><AMF Region ID><AMF Set ID><AMF Pointer>
			// 5G-GUTI := <GUAMI><5G-TMSI>
			tmpReginID, _, _ := ngapConvert.AmfIdToNgap(servedGuami.AmfId)
			amfID := ngapConvert.AmfIdToModels(tmpReginID, fiveGSTMSI.AMFSetID.Value, fiveGSTMSI.AMFPointer.Value)
			tmsi := hex.EncodeToString(fiveGSTMSI.FiveGTMSI.Value)
			guti := servedGuami.PlmnId.Mcc + servedGuami.PlmnId.Mnc + amfID + tmsi

			// TODO: invoke Namf_Communication_UEContextTransfer if serving AMF has changed since
			// last Registration Request procedure
			// Described in TS 23.502 4.2.2.2.2 step 4 (without UDSF deployment)

			if amfUe, ok := amfSelf.AmfUeFindByGuti(guti); !ok {
				ranUe.Log.Warnf("Unknown UE [GUTI: %s]", guti)
			} else {
				ranUe.Log.Tracef("find AmfUe [GUTI: %s]", guti)
				ranUe.Log.Debugf("AmfUe Attach RanUe [RanUeNgapID: %d]", ranUe.RanUeNgapId)
				amfUe.AttachRanUe(ranUe)
			}
		}
	} else {
		ranUe.AmfUe.AttachRanUe(ranUe)
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	if rRCEstablishmentCause != nil {
		ranUe.Log.Tracef("[Initial UE Message] RRC Establishment Cause[%d]", rRCEstablishmentCause.Value)
		ranUe.RRCEstablishmentCause = strconv.Itoa(int(rRCEstablishmentCause.Value))
	}

	if uEContextRequest != nil {
		lbConn.Log.Debug("Trigger initial Context Setup procedure")
		ranUe.UeContextRequest = true
		// TODO: Trigger Initial Context Setup procedure
	} else {
		ranUe.UeContextRequest = false
	}

	// TS 23.502 4.2.2.2.3 step 6a Nnrf_NFDiscovery_Request (NF type, AMF Set)
	// if aMFSetID != nil {
	// TODO: This is a rerouted message
	// TS 38.413: AMF shall, if supported, use the IE as described in TS 23.502
	// }

	// ng-ran propagate allowedNssai in the rerouted initial ue message (TS 38.413 8.6.5)
	// TS 23.502 4.2.2.2.3 step 4a Nnssf_NSSelection_Get
	// if allowedNSSAI != nil {
	// TODO: AMF should use it as defined in TS 23.502
	// }

	pdu, err := libngap.Encoder(*message)
	if err != nil {
		lbConn.Log.Errorf("libngap Encoder Error: %+v", err)
	}
	ranUe.InitialUEMessage = pdu
	nas.HandleNAS(ranUe, ngapType.ProcedureCodeInitialUEMessage, nASPDU.Value)
}

func HandlePDUSessionResourceSetupResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceSetupResponseList *ngapType.PDUSessionResourceSetupListSURes
	var pDUSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListSURes
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var ranUe *context.RanUe

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceSetupResponse
	if pDUSessionResourceSetupResponse == nil {
		lbConn.Log.Error("PDUSessionResourceSetupResponse is nil")
		return
	}

	lbConn.Log.Info("Handle PDU Session Resource Setup Response")

	for _, ie := range pDUSessionResourceSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDPDUSessionResourceSetupListSURes: // ignore
			pDUSessionResourceSetupResponseList = ie.Value.PDUSessionResourceSetupListSURes
			lbConn.Log.Trace("Decode IE PDUSessionResourceSetupListSURes")
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListSURes: // ignore
			pDUSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListSURes
			lbConn.Log.Trace("Decode IE PDUSessionResourceFailedToSetupListSURes")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // optional, ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	if rANUENGAPID != nil {
		ranUe = lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		if ranUe == nil {
			lbConn.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		}
	}

	if aMFUENGAPID != nil {
		ranUe = context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			lbConn.Log.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
			return
		}
	}

	if ranUe != nil {
		ranUe.Log.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)
		amfUe := ranUe.AmfUe
		if amfUe == nil {
			ranUe.Log.Error("amfUe is nil")
			return
		}

		if pDUSessionResourceSetupResponseList != nil {
			ranUe.Log.Trace("Send PDUSessionResourceSetupResponseTransfer to SMF")

			for _, item := range pDUSessionResourceSetupResponseList.List {
				pduSessionID := int32(item.PDUSessionID.Value)
				transfer := item.PDUSessionResourceSetupResponseTransfer
				smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
				if !ok {
					ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				}
				_, _, _, err := consumer.SendUpdateSmContextN2Info(amfUe, smContext,
					models.N2SmInfoType_PDU_RES_SETUP_RSP, transfer)
				if err != nil {
					ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupResponseTransfer] Error: %+v", err)
				}
				// RAN initiated QoS Flow Mobility in subclause 5.2.2.3.7
				// if response != nil && response.BinaryDataN2SmInformation != nil {
				// TODO: n2SmInfo send to RAN
				// } else if response == nil {
				// TODO: error handling
				// }
			}
		}

		if pDUSessionResourceFailedToSetupList != nil {
			ranUe.Log.Trace("Send PDUSessionResourceSetupUnsuccessfulTransfer to SMF")

			for _, item := range pDUSessionResourceFailedToSetupList.List {
				pduSessionID := int32(item.PDUSessionID.Value)
				transfer := item.PDUSessionResourceSetupUnsuccessfulTransfer
				smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
				if !ok {
					ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				}
				_, _, _, err := consumer.SendUpdateSmContextN2Info(amfUe, smContext,
					models.N2SmInfoType_PDU_RES_SETUP_FAIL, transfer)
				if err != nil {
					ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupUnsuccessfulTransfer] Error: %+v", err)
				}

				// if response != nil && response.BinaryDataN2SmInformation != nil {
				// TODO: n2SmInfo send to RAN
				// } else if response == nil {
				// TODO: error handling
				// }
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

func HandlePDUSessionResourceModifyResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pduSessionResourceModifyResponseList *ngapType.PDUSessionResourceModifyListModRes
	var pduSessionResourceFailedToModifyList *ngapType.PDUSessionResourceFailedToModifyListModRes
	var userLocationInformation *ngapType.UserLocationInformation
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var ranUe *context.RanUe

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceModifyResponse := successfulOutcome.Value.PDUSessionResourceModifyResponse
	if pDUSessionResourceModifyResponse == nil {
		lbConn.Log.Error("PDUSessionResourceModifyResponse is nil")
		return
	}

	lbConn.Log.Info("Handle PDU Session Resource Modify Response")

	for _, ie := range pDUSessionResourceModifyResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDPDUSessionResourceModifyListModRes: // ignore
			pduSessionResourceModifyResponseList = ie.Value.PDUSessionResourceModifyListModRes
			lbConn.Log.Trace("Decode IE PDUSessionResourceModifyListModRes")
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToModifyListModRes: // ignore
			pduSessionResourceFailedToModifyList = ie.Value.PDUSessionResourceFailedToModifyListModRes
			lbConn.Log.Trace("Decode IE PDUSessionResourceFailedToModifyListModRes")
		case ngapType.ProtocolIEIDUserLocationInformation: // optional, ignore
			userLocationInformation = ie.Value.UserLocationInformation
			lbConn.Log.Trace("Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // optional, ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	if rANUENGAPID != nil {
		ranUe = lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		if ranUe == nil {
			lbConn.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		}
	}

	if aMFUENGAPID != nil {
		ranUe = context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			lbConn.Log.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
			return
		}
	}

	if ranUe != nil {
		ranUe.Log.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)
		amfUe := ranUe.AmfUe
		if amfUe == nil {
			ranUe.Log.Error("amfUe is nil")
			return
		}

		if pduSessionResourceModifyResponseList != nil {
			ranUe.Log.Trace("Send PDUSessionResourceModifyResponseTransfer to SMF")

			for _, item := range pduSessionResourceModifyResponseList.List {
				pduSessionID := int32(item.PDUSessionID.Value)
				transfer := item.PDUSessionResourceModifyResponseTransfer
				smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
				if !ok {
					ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				}
				_, _, _, err := consumer.SendUpdateSmContextN2Info(amfUe, smContext,
					models.N2SmInfoType_PDU_RES_MOD_RSP, transfer)
				if err != nil {
					ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceModifyResponseTransfer] Error: %+v", err)
				}
				// if response != nil && response.BinaryDataN2SmInformation != nil {
				// TODO: n2SmInfo send to RAN
				// } else if response == nil {
				// TODO: error handling
				// }
			}
		}

		if pduSessionResourceFailedToModifyList != nil {
			ranUe.Log.Trace("Send PDUSessionResourceModifyUnsuccessfulTransfer to SMF")

			for _, item := range pduSessionResourceFailedToModifyList.List {
				pduSessionID := int32(item.PDUSessionID.Value)
				transfer := item.PDUSessionResourceModifyUnsuccessfulTransfer
				smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
				if !ok {
					ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				}
				// response, _, _, err := consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID,
				_, _, _, err := consumer.SendUpdateSmContextN2Info(amfUe, smContext,
					models.N2SmInfoType_PDU_RES_MOD_FAIL, transfer)
				if err != nil {
					ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceModifyUnsuccessfulTransfer] Error: %+v", err)
				}
				// if response != nil && response.BinaryDataN2SmInformation != nil {
				// TODO: n2SmInfo send to RAN
				// } else if response == nil {
				// TODO: error handling
				// }
			}
		}

		if userLocationInformation != nil {
			ranUe.UpdateLocation(userLocationInformation)
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

func HandlePDUSessionResourceNotify(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceNotifyList *ngapType.PDUSessionResourceNotifyList
	var pDUSessionResourceReleasedListNot *ngapType.PDUSessionResourceReleasedListNot
	var userLocationInformation *ngapType.UserLocationInformation

	var ranUe *context.RanUe

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	PDUSessionResourceNotify := initiatingMessage.Value.PDUSessionResourceNotify
	if PDUSessionResourceNotify == nil {
		lbConn.Log.Error("PDUSessionResourceNotify is nil")
		return
	}

	for _, ie := range PDUSessionResourceNotify.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID // reject
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID // reject
			lbConn.Log.Trace("Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDPDUSessionResourceNotifyList: // reject
			pDUSessionResourceNotifyList = ie.Value.PDUSessionResourceNotifyList
			lbConn.Log.Trace("Decode IE pDUSessionResourceNotifyList")
			if pDUSessionResourceNotifyList == nil {
				lbConn.Log.Error("pDUSessionResourceNotifyList is nil")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceReleasedListNot: // ignore
			pDUSessionResourceReleasedListNot = ie.Value.PDUSessionResourceReleasedListNot
			lbConn.Log.Trace("Decode IE PDUSessionResourceReleasedListNot")
			if pDUSessionResourceReleasedListNot == nil {
				lbConn.Log.Error("PDUSessionResourceReleasedListNot is nil")
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // optional, ignore
			userLocationInformation = ie.Value.UserLocationInformation
			lbConn.Log.Trace("Decode IE userLocationInformation")
			if userLocationInformation == nil {
				lbConn.Log.Warn("userLocationInformation is nil [optional]")
			}
		}
	}

	ranUe = lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
	}

	ranUe = context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
		return
	}

	ranUe.Log.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)
	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ranUe.Log.Error("amfUe is nil")
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	ranUe.Log.Trace("Send PDUSessionResourceNotifyTransfer to SMF")

	for _, item := range pDUSessionResourceNotifyList.List {
		pduSessionID := int32(item.PDUSessionID.Value)
		transfer := item.PDUSessionResourceNotifyTransfer
		smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
		if !ok {
			ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
		}
		response, errResponse, problemDetail, err := consumer.SendUpdateSmContextN2Info(amfUe, smContext,
			models.N2SmInfoType_PDU_RES_NTY, transfer)
		if err != nil {
			ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceNotifyTransfer] Error: %+v", err)
		}

		if response != nil {
			responseData := response.JsonData
			n2Info := response.BinaryDataN1SmMessage
			n1Msg := response.BinaryDataN2SmInformation
			if n2Info != nil {
				switch responseData.N2SmInfoType {
				case models.N2SmInfoType_PDU_RES_MOD_REQ:
					ranUe.Log.Debugln("AMF Transfer NGAP PDU Resource Modify Req from SMF")
					var nasPdu []byte
					if n1Msg != nil {
						pduSessionId := uint8(pduSessionID)
						nasPdu, err =
							gmm_message.BuildDLNASTransport(amfUe, lbConn.AnType, nasMessage.PayloadContainerTypeN1SMInfo,
								n1Msg, pduSessionId, nil, nil, 0)
						if err != nil {
							ranUe.Log.Warnf("GMM Message build DL NAS Transport filaed: %v", err)
						}
					}
					list := ngapType.PDUSessionResourceModifyListModReq{}
					ngap_message.AppendPDUSessionResourceModifyListModReq(&list, pduSessionID, nasPdu, n2Info)
					ngap_message.SendPDUSessionResourceModifyRequest(ranUe, list)
				}
			}
		} else if errResponse != nil {
			errJSON := errResponse.JsonData
			n1Msg := errResponse.BinaryDataN2SmInformation
			ranUe.Log.Warnf("PDU Session Modification is rejected by SMF[pduSessionId:%d], Error[%s]\n",
				pduSessionID, errJSON.Error.Cause)
			if n1Msg != nil {
				gmm_message.SendDLNASTransport(
					ranUe, nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, pduSessionID, 0, nil, 0)
			}
			// TODO: handle n2 info transfer
		} else if err != nil {
			return
		} else {
			// TODO: error handling
			ranUe.Log.Errorf("Failed to Update smContext[pduSessionID: %d], Error[%v]", pduSessionID, problemDetail)
			return
		}
	}

	if pDUSessionResourceReleasedListNot != nil {
		ranUe.Log.Trace("Send PDUSessionResourceNotifyReleasedTransfer to SMF")
		for _, item := range pDUSessionResourceReleasedListNot.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceNotifyReleasedTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			response, errResponse, problemDetail, err := consumer.SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_NTY_REL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceNotifyReleasedTransfer] Error: %+v", err)
			}
			if response != nil {
				responseData := response.JsonData
				n2Info := response.BinaryDataN1SmMessage
				n1Msg := response.BinaryDataN2SmInformation
				if n2Info != nil {
					switch responseData.N2SmInfoType {
					case models.N2SmInfoType_PDU_RES_REL_CMD:
						ranUe.Log.Debugln("AMF Transfer NGAP PDU Session Resource Rel Co from SMF")
						var nasPdu []byte
						if n1Msg != nil {
							pduSessionId := uint8(pduSessionID)
							nasPdu, err = gmm_message.BuildDLNASTransport(amfUe, lbConn.AnType,
								nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, pduSessionId, nil, nil, 0)
							if err != nil {
								ranUe.Log.Warnf("GMM Message build DL NAS Transport filaed: %v", err)
							}
						}
						list := ngapType.PDUSessionResourceToReleaseListRelCmd{}
						ngap_message.AppendPDUSessionResourceToReleaseListRelCmd(&list, pduSessionID, n2Info)
						ngap_message.SendPDUSessionResourceReleaseCommand(ranUe, nasPdu, list)
					}
				}
			} else if errResponse != nil {
				errJSON := errResponse.JsonData
				n1Msg := errResponse.BinaryDataN2SmInformation
				ranUe.Log.Warnf("PDU Session Release is rejected by SMF[pduSessionId:%d], Error[%s]\n",
					pduSessionID, errJSON.Error.Cause)
				if n1Msg != nil {
					gmm_message.SendDLNASTransport(
						ranUe, nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, pduSessionID, 0, nil, 0)
				}
			} else if err != nil {
				return
			} else {
				// TODO: error handling
				ranUe.Log.Errorf("Failed to Update smContext[pduSessionID: %d], Error[%v]", pduSessionID, problemDetail)
				return
			}
		}
	}
}

func HandlePDUSessionResourceModifyIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pduSessionResourceModifyIndicationList *ngapType.PDUSessionResourceModifyListModInd

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	var ranUe *context.RanUe

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // reject
	if initiatingMessage == nil {
		lbConn.Log.Error("InitiatingMessage is nil")
		cause := ngapType.Cause{
			Present: ngapType.CausePresentProtocol,
			Protocol: &ngapType.CauseProtocol{
				Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorReject,
			},
		}
		ngap_message.SendErrorIndication(lbConn, nil, nil, &cause, nil)
		return
	}
	pDUSessionResourceModifyIndication := initiatingMessage.Value.PDUSessionResourceModifyIndication
	if pDUSessionResourceModifyIndication == nil {
		lbConn.Log.Error("PDUSessionResourceModifyIndication is nil")
		cause := ngapType.Cause{
			Present: ngapType.CausePresentProtocol,
			Protocol: &ngapType.CauseProtocol{
				Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorReject,
			},
		}
		ngap_message.SendErrorIndication(lbConn, nil, nil, &cause, nil)
		return
	}

	lbConn.Log.Info("Handle PDU Session Resource Modify Indication")

	for _, ie := range pDUSessionResourceModifyIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
					ngapType.ProtocolIEIDAMFUENGAPID, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
					ngapType.ProtocolIEIDRANUENGAPID, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDPDUSessionResourceModifyListModInd: // reject
			pduSessionResourceModifyIndicationList = ie.Value.PDUSessionResourceModifyListModInd
			lbConn.Log.Trace("Decode IE PDUSessionResourceModifyListModInd")
			if pduSessionResourceModifyIndicationList == nil {
				lbConn.Log.Error("PDUSessionResourceModifyListModInd is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
					ngapType.ProtocolIEIDPDUSessionResourceModifyListModInd, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		lbConn.Log.Error("Has missing reject IE(s)")

		procedureCode := ngapType.ProcedureCodePDUSessionResourceModifyIndication
		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
		procedureCriticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality,
			&iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, nil, &criticalityDiagnostics)
		return
	}

	ranUe = lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, &cause, nil)
		return
	}

	lbConn.Log.Tracef("UE Context AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		lbConn.Log.Error("AmfUe is nil")
		return
	}

	pduSessionResourceModifyListModCfm := ngapType.PDUSessionResourceModifyListModCfm{}
	pduSessionResourceFailedToModifyListModCfm := ngapType.PDUSessionResourceFailedToModifyListModCfm{}

	lbConn.Log.Trace("Send PDUSessionResourceModifyIndicationTransfer to SMF")
	for _, item := range pduSessionResourceModifyIndicationList.List {
		pduSessionID := int32(item.PDUSessionID.Value)
		transfer := item.PDUSessionResourceModifyIndicationTransfer
		smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
		if !ok {
			ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
		}
		response, errResponse, _, err := consumer.SendUpdateSmContextN2Info(amfUe, smContext,
			models.N2SmInfoType_PDU_RES_MOD_IND, transfer)
		if err != nil {
			lbConn.Log.Errorf("SendUpdateSmContextN2Info Error:\n%s", err.Error())
		}

		if response != nil && response.BinaryDataN2SmInformation != nil {
			ngap_message.AppendPDUSessionResourceModifyListModCfm(&pduSessionResourceModifyListModCfm, int64(pduSessionID),
				response.BinaryDataN2SmInformation)
		}
		if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
			ngap_message.AppendPDUSessionResourceFailedToModifyListModCfm(&pduSessionResourceFailedToModifyListModCfm,
				int64(pduSessionID), errResponse.BinaryDataN2SmInformation)
		}
	}

	ngap_message.SendPDUSessionResourceModifyConfirm(ranUe, pduSessionResourceModifyListModCfm,
		pduSessionResourceFailedToModifyListModCfm, nil)
}

func HandleInitialContextSetupResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceSetupResponseList *ngapType.PDUSessionResourceSetupListCxtRes
	var pDUSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListCxtRes
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	if initialContextSetupResponse == nil {
		lbConn.Log.Error("InitialContextSetupResponse is nil")
		return
	}

	lbConn.Log.Info("Handle Initial Context Setup Response")

	for _, ie := range initialContextSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Warn("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Warn("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceSetupListCxtRes:
			pDUSessionResourceSetupResponseList = ie.Value.PDUSessionResourceSetupListCxtRes
			lbConn.Log.Trace("Decode IE PDUSessionResourceSetupResponseList")
			if pDUSessionResourceSetupResponseList == nil {
				lbConn.Log.Warn("PDUSessionResourceSetupResponseList is nil")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListCxtRes:
			pDUSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListCxtRes
			lbConn.Log.Trace("Decode IE PDUSessionResourceFailedToSetupList")
			if pDUSessionResourceFailedToSetupList == nil {
				lbConn.Log.Warn("PDUSessionResourceFailedToSetupList is nil")
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE Criticality Diagnostics")
			if criticalityDiagnostics == nil {
				lbConn.Log.Warn("Criticality Diagnostics is nil")
			}
		}
	}

	ranUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
	amfUe := ranUe.AmfUe
	if amfUe == nil {
		lbConn.Log.Error("amfUe is nil")
		return
	}

	lbConn.Log.Tracef("RanUeNgapID[%d] AmfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)

	if pDUSessionResourceSetupResponseList != nil {
		ranUe.Log.Trace("Send PDUSessionResourceSetupResponseTransfer to SMF")

		for _, item := range pDUSessionResourceSetupResponseList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupResponseTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			// response, _, _, err := consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID,
			_, _, _, err := consumer.SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_SETUP_RSP, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupResponseTransfer] Error: %+v", err)
			}
			// RAN initiated QoS Flow Mobility in subclause 5.2.2.3.7
			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}

	if pDUSessionResourceFailedToSetupList != nil {
		ranUe.Log.Trace("Send PDUSessionResourceSetupUnsuccessfulTransfer to SMF")

		for _, item := range pDUSessionResourceFailedToSetupList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupUnsuccessfulTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			// response, _, _, err := consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID,
			_, _, _, err := consumer.SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_SETUP_FAIL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupUnsuccessfulTransfer] Error: %+v", err)
			}

			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}

	if ranUe.Ran.AnType == models.AccessType_NON_3_GPP_ACCESS {
		ngap_message.SendDownlinkNasTransport(ranUe, amfUe.RegistrationAcceptForNon3GPPAccess, nil)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

func HandleInitialContextSetupFailure(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListCxtFail
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		lbConn.Log.Error("UnsuccessfulOutcome is nil")
		return
	}
	initialContextSetupFailure := unsuccessfulOutcome.Value.InitialContextSetupFailure
	if initialContextSetupFailure == nil {
		lbConn.Log.Error("InitialContextSetupFailure is nil")
		return
	}

	lbConn.Log.Info("Handle Initial Context Setup Failure")

	for _, ie := range initialContextSetupFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Warn("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Warn("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListCxtFail:
			pDUSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListCxtFail
			lbConn.Log.Trace("Decode IE PDUSessionResourceFailedToSetupList")
			if pDUSessionResourceFailedToSetupList == nil {
				lbConn.Log.Warn("PDUSessionResourceFailedToSetupList is nil")
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			lbConn.Log.Trace("Decode IE Cause")
			if cause == nil {
				lbConn.Log.Warn("Cause is nil")
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE Criticality Diagnostics")
			if criticalityDiagnostics == nil {
				lbConn.Log.Warn("CriticalityDiagnostics is nil")
			}
		}
	}

	printAndGetCause(lbConn, cause)

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
	ranUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
	amfUe := ranUe.AmfUe
	if amfUe == nil {
		lbConn.Log.Error("amfUe is nil")
		return
	}

	if pDUSessionResourceFailedToSetupList != nil {
		ranUe.Log.Trace("Send PDUSessionResourceSetupUnsuccessfulTransfer to SMF")

		for _, item := range pDUSessionResourceFailedToSetupList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupUnsuccessfulTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			_, _, _, err := consumer.SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_SETUP_FAIL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupUnsuccessfulTransfer] Error: %+v", err)
			}

			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}
}

func HandleUEContextReleaseRequest(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceList *ngapType.PDUSessionResourceListCxtRelReq
	var cause *ngapType.Cause

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	uEContextReleaseRequest := initiatingMessage.Value.UEContextReleaseRequest
	if uEContextReleaseRequest == nil {
		lbConn.Log.Error("UEContextReleaseRequest is nil")
		return
	}

	lbConn.Log.Info("UE Context Release Request")

	for _, ie := range uEContextReleaseRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDPDUSessionResourceListCxtRelReq:
			pDUSessionResourceList = ie.Value.PDUSessionResourceListCxtRelReq
			lbConn.Log.Trace("Decode IE Pdu Session Resource List")
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			lbConn.Log.Trace("Decode IE Cause")
			if cause == nil {
				lbConn.Log.Warn("Cause is nil")
			}
		}
	}

	ranUe := context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No RanUe Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
		cause = &ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, cause, nil)
		return
	}

	lbConn.Log.Tracef("RanUeNgapID[%d] AmfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)

	causeGroup := ngapType.CausePresentRadioNetwork
	causeValue := ngapType.CauseRadioNetworkPresentUnspecified
	if cause != nil {
		causeGroup, causeValue = printAndGetCause(lbConn, cause)
	}

	amfUe := ranUe.AmfUe
	if amfUe != nil {
		causeAll := context.CauseAll{
			NgapCause: &models.NgApCause{
				Group: int32(causeGroup),
				Value: int32(causeValue),
			},
		}
		if amfUe.State[lbConn.AnType].Is(context.Registered) {
			ranUe.Log.Info("Ue Context in GMM-Registered")
			if pDUSessionResourceList != nil {
				for _, pduSessionReourceItem := range pDUSessionResourceList.List {
					pduSessionID := int32(pduSessionReourceItem.PDUSessionID.Value)
					smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
					if !ok {
						ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
						// TODO: Check if doing error handling here
						continue
					}
					response, _, _, err := consumer.SendUpdateSmContextDeactivateUpCnxState(amfUe, smContext, causeAll)
					if err != nil {
						ranUe.Log.Errorf("Send Update SmContextDeactivate UpCnxState Error[%s]", err.Error())
					} else if response == nil {
						ranUe.Log.Errorln("Send Update SmContextDeactivate UpCnxState Error")
					}
				}
			}
		} else {
			ranUe.Log.Info("Ue Context in Non GMM-Registered")
			amfUe.SmContextList.Range(func(key, value interface{}) bool {
				smContext := value.(*context.SmContext)
				detail, err := consumer.SendReleaseSmContextRequest(amfUe, smContext, &causeAll, "", nil)
				if err != nil {
					ranUe.Log.Errorf("Send ReleaseSmContextRequest Error[%s]", err.Error())
				} else if detail != nil {
					ranUe.Log.Errorf("Send ReleaseSmContextRequeste Error[%s]", detail.Cause)
				}
				return true
			})
			ngap_message.SendUEContextReleaseCommand(ranUe, context.UeContextReleaseUeContext, causeGroup, causeValue)
			return
		}
	}
	ngap_message.SendUEContextReleaseCommand(ranUe, context.UeContextN2NormalRelease, causeGroup, causeValue)
}

func HandleUEContextModificationResponse(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var rRCState *ngapType.RRCState
	var userLocationInformation *ngapType.UserLocationInformation
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var ranUe *context.RanUe

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	uEContextModificationResponse := successfulOutcome.Value.UEContextModificationResponse
	if uEContextModificationResponse == nil {
		lbConn.Log.Error("UEContextModificationResponse is nil")
		return
	}

	lbConn.Log.Info("Handle UE Context Modification Response")

	for _, ie := range uEContextModificationResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Warn("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Warn("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRRCState: // optional, ignore
			rRCState = ie.Value.RRCState
			lbConn.Log.Trace("Decode IE RRCState")
		case ngapType.ProtocolIEIDUserLocationInformation: // optional, ignore
			userLocationInformation = ie.Value.UserLocationInformation
			lbConn.Log.Trace("Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // optional, ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	if rANUENGAPID != nil {
		ranUe = lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		if ranUe == nil {
			lbConn.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		}
	}

	if aMFUENGAPID != nil {
		ranUe = context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			lbConn.Log.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
			return
		}
	}

	if ranUe != nil {
		ranUe.Log.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)

		if rRCState != nil {
			switch rRCState.Value {
			case ngapType.RRCStatePresentInactive:
				ranUe.Log.Trace("UE RRC State: Inactive")
			case ngapType.RRCStatePresentConnected:
				ranUe.Log.Trace("UE RRC State: Connected")
			}
		}

		if userLocationInformation != nil {
			ranUe.UpdateLocation(userLocationInformation)
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

func HandleUEContextModificationFailure(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var ranUe *context.RanUe

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		lbConn.Log.Error("UnsuccessfulOutcome is nil")
		return
	}
	uEContextModificationFailure := unsuccessfulOutcome.Value.UEContextModificationFailure
	if uEContextModificationFailure == nil {
		lbConn.Log.Error("UEContextModificationFailure is nil")
		return
	}

	lbConn.Log.Info("Handle UE Context Modification Failure")

	for _, ie := range uEContextModificationFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Warn("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Warn("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDCause: // ignore
			cause = ie.Value.Cause
			lbConn.Log.Trace("Decode IE Cause")
			if cause == nil {
				lbConn.Log.Warn("Cause is nil")
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // optional, ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	if rANUENGAPID != nil {
		ranUe = lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		if ranUe == nil {
			lbConn.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		}
	}

	if aMFUENGAPID != nil {
		ranUe = context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			lbConn.Log.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
		}
	}

	if ranUe != nil {
		lbConn.Log.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)
	}

	if cause != nil {
		printAndGetCause(lbConn, cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

func HandleRRCInactiveTransitionReport(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var rRCState *ngapType.RRCState
	var userLocationInformation *ngapType.UserLocationInformation

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("Initiating Message is nil")
		return
	}

	rRCInactiveTransitionReport := initiatingMessage.Value.RRCInactiveTransitionReport
	if rRCInactiveTransitionReport == nil {
		lbConn.Log.Error("RRCInactiveTransitionReport is nil")
		return
	}
	lbConn.Log.Info("Handle RRC Inactive Transition Report")

	for i := 0; i < len(rRCInactiveTransitionReport.ProtocolIEs.List); i++ {
		ie := rRCInactiveTransitionReport.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRRCState: // ignore
			rRCState = ie.Value.RRCState
			lbConn.Log.Trace("Decode IE RRCState")
			if rRCState == nil {
				lbConn.Log.Error("RRCState is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // ignore
			userLocationInformation = ie.Value.UserLocationInformation
			lbConn.Log.Trace("Decode IE UserLocationInformation")
			if userLocationInformation == nil {
				lbConn.Log.Error("UserLocationInformation is nil")
				return
			}
		}
	}

	ranUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
	} else {
		lbConn.Log.Tracef("RANUENGAPID[%d] AMFUENGAPID[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)

		if rRCState != nil {
			switch rRCState.Value {
			case ngapType.RRCStatePresentInactive:
				lbConn.Log.Trace("UE RRC State: Inactive")
			case ngapType.RRCStatePresentConnected:
				lbConn.Log.Trace("UE RRC State: Connected")
			}
		}
		ranUe.UpdateLocation(userLocationInformation)
	}
}

func HandleHandoverNotify(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var userLocationInformation *ngapType.UserLocationInformation

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("Initiating Message is nil")
		return
	}
	HandoverNotify := initiatingMessage.Value.HandoverNotify
	if HandoverNotify == nil {
		lbConn.Log.Error("HandoverNotify is nil")
		return
	}

	lbConn.Log.Info("Handle Handover notification")

	for i := 0; i < len(HandoverNotify.ProtocolIEs.List); i++ {
		ie := HandoverNotify.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AMFUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RANUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation:
			userLocationInformation = ie.Value.UserLocationInformation
			lbConn.Log.Trace("Decode IE userLocationInformation")
			if userLocationInformation == nil {
				lbConn.Log.Error("userLocationInformation is nil")
				return
			}
		}
	}

	targetUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)

	if targetUe == nil {
		lbConn.Log.Errorf("No RanUe Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, &cause, nil)
		return
	}

	if userLocationInformation != nil {
		targetUe.UpdateLocation(userLocationInformation)
	}
	amfUe := targetUe.AmfUe
	if amfUe == nil {
		lbConn.Log.Error("AmfUe is nil")
		return
	}
	sourceUe := targetUe.SourceUe
	if sourceUe == nil {
		// TODO: Send to S-AMF
		// Desciibed in (23.502 4.9.1.3.3) [conditional] 6a.Namf_Communication_N2InfoNotify.
		lbConn.Log.Error("N2 Handover between AMF has not been implemented yet")
	} else {
		lbConn.Log.Info("Handle Handover notification Finshed ")
		for _, pduSessionid := range targetUe.SuccessPduSessionId {
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionid)
			if !ok {
				sourceUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionid)
			}
			_, _, _, err := consumer.SendUpdateSmContextN2HandoverComplete(amfUe, smContext, "", nil)
			if err != nil {
				lbConn.Log.Errorf("Send UpdateSmContextN2HandoverComplete Error[%s]", err.Error())
			}
		}
		amfUe.AttachRanUe(targetUe)

		ngap_message.SendUEContextReleaseCommand(sourceUe, context.UeContextReleaseHandover, ngapType.CausePresentNas,
			ngapType.CauseNasPresentNormalRelease)
	}

	// TODO: The UE initiates Mobility Registration Update procedure as described in clause 4.2.2.2.2.
}

// TS 23.502 4.9.1
func HandlePathSwitchRequest(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var rANUENGAPID *ngapType.RANUENGAPID
	var sourceAMFUENGAPID *ngapType.AMFUENGAPID
	var userLocationInformation *ngapType.UserLocationInformation
	var uESecurityCapabilities *ngapType.UESecurityCapabilities
	var pduSessionResourceToBeSwitchedInDLList *ngapType.PDUSessionResourceToBeSwitchedDLList
	var pduSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListPSReq

	var ranUe *context.RanUe

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	pathSwitchRequest := initiatingMessage.Value.PathSwitchRequest
	if pathSwitchRequest == nil {
		lbConn.Log.Error("PathSwitchRequest is nil")
		return
	}

	lbConn.Log.Info("Handle Path Switch Request")

	for _, ie := range pathSwitchRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDSourceAMFUENGAPID: // reject
			sourceAMFUENGAPID = ie.Value.SourceAMFUENGAPID
			lbConn.Log.Trace("Decode IE SourceAmfUeNgapID")
			if sourceAMFUENGAPID == nil {
				lbConn.Log.Error("SourceAmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // ignore
			userLocationInformation = ie.Value.UserLocationInformation
			lbConn.Log.Trace("Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDUESecurityCapabilities: // ignore
			uESecurityCapabilities = ie.Value.UESecurityCapabilities
			lbConn.Log.Trace("Decode IE UESecurityCapabilities")
		case ngapType.ProtocolIEIDPDUSessionResourceToBeSwitchedDLList: // reject
			pduSessionResourceToBeSwitchedInDLList = ie.Value.PDUSessionResourceToBeSwitchedDLList
			lbConn.Log.Trace("Decode IE PDUSessionResourceToBeSwitchedDLList")
			if pduSessionResourceToBeSwitchedInDLList == nil {
				lbConn.Log.Error("PDUSessionResourceToBeSwitchedDLList is nil")
				return
			}
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListPSReq: // ignore
			pduSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListPSReq
			lbConn.Log.Trace("Decode IE PDUSessionResourceFailedToSetupListPSReq")
		}
	}

	if sourceAMFUENGAPID == nil {
		lbConn.Log.Error("SourceAmfUeNgapID is nil")
		return
	}
	ranUe = context.AMF_Self().RanUeFindByAmfUeNgapID(sourceAMFUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("Cannot find UE from sourceAMfUeNgapID[%d]", sourceAMFUENGAPID.Value)
		ngap_message.SendPathSwitchRequestFailure(lbConn, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
		return
	}

	lbConn.Log.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ranUe.Log.Error("AmfUe is nil")
		ngap_message.SendPathSwitchRequestFailure(lbConn, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
		return
	}

	if amfUe.SecurityContextIsValid() {
		// Update NH
		amfUe.UpdateNH()
	} else {
		ranUe.Log.Errorf("No Security Context : SUPI[%s]", amfUe.Supi)
		ngap_message.SendPathSwitchRequestFailure(lbConn, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
		return
	}

	if uESecurityCapabilities != nil {
		amfUe.UESecurityCapability.SetEA1_128_5G(uESecurityCapabilities.NRencryptionAlgorithms.Value.Bytes[0] & 0x80)
		amfUe.UESecurityCapability.SetEA2_128_5G(uESecurityCapabilities.NRencryptionAlgorithms.Value.Bytes[0] & 0x40)
		amfUe.UESecurityCapability.SetEA3_128_5G(uESecurityCapabilities.NRencryptionAlgorithms.Value.Bytes[0] & 0x20)
		amfUe.UESecurityCapability.SetIA1_128_5G(uESecurityCapabilities.NRintegrityProtectionAlgorithms.Value.Bytes[0] & 0x80)
		amfUe.UESecurityCapability.SetIA2_128_5G(uESecurityCapabilities.NRintegrityProtectionAlgorithms.Value.Bytes[0] & 0x40)
		amfUe.UESecurityCapability.SetIA3_128_5G(uESecurityCapabilities.NRintegrityProtectionAlgorithms.Value.Bytes[0] & 0x20)
		// not support any E-UTRA algorithms
	}

	if rANUENGAPID != nil {
		ranUe.RanUeNgapId = rANUENGAPID.Value
	}

	ranUe.UpdateLocation(userLocationInformation)

	var pduSessionResourceSwitchedList ngapType.PDUSessionResourceSwitchedList
	var pduSessionResourceReleasedListPSAck ngapType.PDUSessionResourceReleasedListPSAck
	var pduSessionResourceReleasedListPSFail ngapType.PDUSessionResourceReleasedListPSFail

	if pduSessionResourceToBeSwitchedInDLList != nil {
		for _, item := range pduSessionResourceToBeSwitchedInDLList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PathSwitchRequestTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			response, errResponse, _, err := consumer.SendUpdateSmContextXnHandover(amfUe, smContext,
				models.N2SmInfoType_PATH_SWITCH_REQ, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextXnHandover[PathSwitchRequestTransfer] Error:\n%s", err.Error())
			}
			if response != nil && response.BinaryDataN2SmInformation != nil {
				pduSessionResourceSwitchedItem := ngapType.PDUSessionResourceSwitchedItem{}
				pduSessionResourceSwitchedItem.PDUSessionID.Value = int64(pduSessionID)
				pduSessionResourceSwitchedItem.PathSwitchRequestAcknowledgeTransfer = response.BinaryDataN2SmInformation
				pduSessionResourceSwitchedList.List = append(pduSessionResourceSwitchedList.List, pduSessionResourceSwitchedItem)
			}
			if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
				pduSessionResourceReleasedItem := ngapType.PDUSessionResourceReleasedItemPSFail{}
				pduSessionResourceReleasedItem.PDUSessionID.Value = int64(pduSessionID)
				pduSessionResourceReleasedItem.PathSwitchRequestUnsuccessfulTransfer = errResponse.BinaryDataN2SmInformation
				pduSessionResourceReleasedListPSFail.List = append(pduSessionResourceReleasedListPSFail.List,
					pduSessionResourceReleasedItem)
			}
		}
	}

	if pduSessionResourceFailedToSetupList != nil {
		for _, item := range pduSessionResourceFailedToSetupList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PathSwitchRequestSetupFailedTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
			}
			response, errResponse, _, err := consumer.SendUpdateSmContextXnHandoverFailed(amfUe, smContext,
				models.N2SmInfoType_PATH_SWITCH_SETUP_FAIL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextXnHandoverFailed[PathSwitchRequestSetupFailedTransfer] Error: %+v", err)
			}
			if response != nil && response.BinaryDataN2SmInformation != nil {
				pduSessionResourceReleasedItem := ngapType.PDUSessionResourceReleasedItemPSAck{}
				pduSessionResourceReleasedItem.PDUSessionID.Value = int64(pduSessionID)
				pduSessionResourceReleasedItem.PathSwitchRequestUnsuccessfulTransfer = response.BinaryDataN2SmInformation
				pduSessionResourceReleasedListPSAck.List = append(pduSessionResourceReleasedListPSAck.List,
					pduSessionResourceReleasedItem)
			}
			if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
				pduSessionResourceReleasedItem := ngapType.PDUSessionResourceReleasedItemPSFail{}
				pduSessionResourceReleasedItem.PDUSessionID.Value = int64(pduSessionID)
				pduSessionResourceReleasedItem.PathSwitchRequestUnsuccessfulTransfer = errResponse.BinaryDataN2SmInformation
				pduSessionResourceReleasedListPSFail.List = append(pduSessionResourceReleasedListPSFail.List,
					pduSessionResourceReleasedItem)
			}
		}
	}

	// TS 23.502 4.9.1.2.2 step 7: send ack to Target NG-RAN. If none of the requested PDU Sessions have been switched
	// successfully, the AMF shall send an N2 Path Switch Request Failure message to the Target NG-RAN
	if len(pduSessionResourceSwitchedList.List) > 0 {
		// TODO: set newSecurityContextIndicator to true if there is a new security context
		err := ranUe.SwitchToRan(lbConn, rANUENGAPID.Value)
		if err != nil {
			ranUe.Log.Error(err.Error())
			return
		}
		ngap_message.SendPathSwitchRequestAcknowledge(ranUe, pduSessionResourceSwitchedList,
			pduSessionResourceReleasedListPSAck, false, nil, nil, nil)
	} else if len(pduSessionResourceReleasedListPSFail.List) > 0 {
		ngap_message.SendPathSwitchRequestFailure(lbConn, sourceAMFUENGAPID.Value, rANUENGAPID.Value,
			&pduSessionResourceReleasedListPSFail, nil)
	} else {
		ngap_message.SendPathSwitchRequestFailure(lbConn, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
	}
}

func HandleHandoverRequestAcknowledge(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceAdmittedList *ngapType.PDUSessionResourceAdmittedList
	var pDUSessionResourceFailedToSetupListHOAck *ngapType.PDUSessionResourceFailedToSetupListHOAck
	var targetToSourceTransparentContainer *ngapType.TargetToSourceTransparentContainer
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	handoverRequestAcknowledge := successfulOutcome.Value.HandoverRequestAcknowledge // reject
	if handoverRequestAcknowledge == nil {
		lbConn.Log.Error("HandoverRequestAcknowledge is nil")
		return
	}

	lbConn.Log.Info("Handle Handover Request Acknowledge")

	for _, ie := range handoverRequestAcknowledge.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDPDUSessionResourceAdmittedList: // ignore
			pDUSessionResourceAdmittedList = ie.Value.PDUSessionResourceAdmittedList
			lbConn.Log.Trace("Decode IE PduSessionResourceAdmittedList")
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListHOAck: // ignore
			pDUSessionResourceFailedToSetupListHOAck = ie.Value.PDUSessionResourceFailedToSetupListHOAck
			lbConn.Log.Trace("Decode IE PduSessionResourceFailedToSetupListHOAck")
		case ngapType.ProtocolIEIDTargetToSourceTransparentContainer: // reject
			targetToSourceTransparentContainer = ie.Value.TargetToSourceTransparentContainer
			lbConn.Log.Trace("Decode IE TargetToSourceTransparentContainer")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}
	if targetToSourceTransparentContainer == nil {
		lbConn.Log.Error("TargetToSourceTransparentContainer is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
			ngapType.ProtocolIEIDTargetToSourceTransparentContainer, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if len(iesCriticalityDiagnostics.List) > 0 {
		lbConn.Log.Error("Has missing reject IE(s)")

		procedureCode := ngapType.ProcedureCodeHandoverResourceAllocation
		triggeringMessage := ngapType.TriggeringMessagePresentSuccessfulOutcome
		procedureCriticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage,
			&procedureCriticality, &iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, nil, &criticalityDiagnostics)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}

	targetUe := context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
	if targetUe == nil {
		lbConn.Log.Errorf("No UE Context[AMFUENGAPID: %d]", aMFUENGAPID.Value)
		return
	}

	if rANUENGAPID != nil {
		targetUe.RanUeNgapId = rANUENGAPID.Value
	}
	lbConn.Log.Debugf("Target Ue RanUeNgapID[%d] AmfUeNgapID[%d]", targetUe.RanUeNgapId, targetUe.AmfUeNgapId)

	amfUe := targetUe.AmfUe
	if amfUe == nil {
		targetUe.Log.Error("amfUe is nil")
		return
	}

	var pduSessionResourceHandoverList ngapType.PDUSessionResourceHandoverList
	var pduSessionResourceToReleaseList ngapType.PDUSessionResourceToReleaseListHOCmd

	// describe in 23.502 4.9.1.3.2 step11
	if pDUSessionResourceAdmittedList != nil {
		for _, item := range pDUSessionResourceAdmittedList.List {
			pduSessionID := item.PDUSessionID.Value
			transfer := item.HandoverRequestAcknowledgeTransfer
			pduSessionId := int32(pduSessionID)
			if smContext, exist := amfUe.SmContextFindByPDUSessionID(pduSessionId); exist {
				response, errResponse, problemDetails, err := consumer.SendUpdateSmContextN2HandoverPrepared(amfUe,
					smContext, models.N2SmInfoType_HANDOVER_REQ_ACK, transfer)
				if err != nil {
					targetUe.Log.Errorf("Send HandoverRequestAcknowledgeTransfer error: %v", err)
				}
				if problemDetails != nil {
					targetUe.Log.Warnf("ProblemDetails[status: %d, Cause: %s]", problemDetails.Status, problemDetails.Cause)
				}
				if response != nil && response.BinaryDataN2SmInformation != nil {
					handoverItem := ngapType.PDUSessionResourceHandoverItem{}
					handoverItem.PDUSessionID = item.PDUSessionID
					handoverItem.HandoverCommandTransfer = response.BinaryDataN2SmInformation
					pduSessionResourceHandoverList.List = append(pduSessionResourceHandoverList.List, handoverItem)
					targetUe.SuccessPduSessionId = append(targetUe.SuccessPduSessionId, pduSessionId)
				}
				if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
					releaseItem := ngapType.PDUSessionResourceToReleaseItemHOCmd{}
					releaseItem.PDUSessionID = item.PDUSessionID
					releaseItem.HandoverPreparationUnsuccessfulTransfer = errResponse.BinaryDataN2SmInformation
					pduSessionResourceToReleaseList.List = append(pduSessionResourceToReleaseList.List, releaseItem)
				}
			}
		}
	}

	if pDUSessionResourceFailedToSetupListHOAck != nil {
		for _, item := range pDUSessionResourceFailedToSetupListHOAck.List {
			pduSessionID := item.PDUSessionID.Value
			transfer := item.HandoverResourceAllocationUnsuccessfulTransfer
			pduSessionId := int32(pduSessionID)
			if smContext, exist := amfUe.SmContextFindByPDUSessionID(pduSessionId); exist {
				_, _, problemDetails, err := consumer.SendUpdateSmContextN2HandoverPrepared(amfUe, smContext,
					models.N2SmInfoType_HANDOVER_RES_ALLOC_FAIL, transfer)
				if err != nil {
					targetUe.Log.Errorf("Send HandoverResourceAllocationUnsuccessfulTransfer error: %v", err)
				}
				if problemDetails != nil {
					targetUe.Log.Warnf("ProblemDetails[status: %d, Cause: %s]", problemDetails.Status, problemDetails.Cause)
				}
			}
		}
	}

	sourceUe := targetUe.SourceUe
	if sourceUe == nil {
		// TODO: Send Namf_Communication_CreateUEContext Response to S-AMF
		lbConn.Log.Error("handover between different Ue has not been implement yet")
	} else {
		lbConn.Log.Tracef("Source: RanUeNgapID[%d] AmfUeNgapID[%d]", sourceUe.RanUeNgapId, sourceUe.AmfUeNgapId)
		lbConn.Log.Tracef("Target: RanUeNgapID[%d] AmfUeNgapID[%d]", targetUe.RanUeNgapId, targetUe.AmfUeNgapId)
		if len(pduSessionResourceHandoverList.List) == 0 {
			targetUe.Log.Info("Handle Handover Preparation Failure [HoFailure In Target5GC NgranNode Or TargetSystem]")
			cause := &ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem,
				},
			}
			ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
			return
		}
		ngap_message.SendHandoverCommand(sourceUe, pduSessionResourceHandoverList, pduSessionResourceToReleaseList,
			*targetToSourceTransparentContainer, nil)
	}
}

func HandleHandoverFailure(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var cause *ngapType.Cause
	var targetUe *context.RanUe
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}

	unsuccessfulOutcome := message.UnsuccessfulOutcome // reject
	if unsuccessfulOutcome == nil {
		lbConn.Log.Error("Unsuccessful Message is nil")
		return
	}

	handoverFailure := unsuccessfulOutcome.Value.HandoverFailure
	if handoverFailure == nil {
		lbConn.Log.Error("HandoverFailure is nil")
		return
	}

	for _, ie := range handoverFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDCause: // ignore
			cause = ie.Value.Cause
			lbConn.Log.Trace("Decode IE Cause")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	causePresent := ngapType.CausePresentRadioNetwork
	causeValue := ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem
	if cause != nil {
		causePresent, causeValue = printAndGetCause(lbConn, cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}

	targetUe = context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)

	if targetUe == nil {
		lbConn.Log.Errorf("No UE Context[AmfUENGAPID: %d]", aMFUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, nil, &cause, nil)
		return
	}

	sourceUe := targetUe.SourceUe
	if sourceUe == nil {
		// TODO: handle N2 Handover between AMF
		lbConn.Log.Error("N2 Handover between AMF has not been implemented yet")
	} else {
		amfUe := targetUe.AmfUe
		if amfUe != nil {
			amfUe.SmContextList.Range(func(key, value interface{}) bool {
				pduSessionID := key.(int32)
				smContext := value.(*context.SmContext)
				causeAll := context.CauseAll{
					NgapCause: &models.NgApCause{
						Group: int32(causePresent),
						Value: int32(causeValue),
					},
				}
				_, _, _, err := consumer.SendUpdateSmContextN2HandoverCanceled(amfUe, smContext, causeAll)
				if err != nil {
					lbConn.Log.Errorf("Send UpdateSmContextN2HandoverCanceled Error for PduSessionId[%d]", pduSessionID)
				}
				return true
			})
		}
		ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, criticalityDiagnostics)
	}

	ngap_message.SendUEContextReleaseCommand(targetUe, context.UeContextReleaseHandover, causePresent, causeValue)
}

func HandleHandoverRequired(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var handoverType *ngapType.HandoverType
	var cause *ngapType.Cause
	var targetID *ngapType.TargetID
	var pDUSessionResourceListHORqd *ngapType.PDUSessionResourceListHORqd
	var sourceToTargetTransparentContainer *ngapType.SourceToTargetTransparentContainer
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("Initiating Message is nil")
		return
	}
	HandoverRequired := initiatingMessage.Value.HandoverRequired
	if HandoverRequired == nil {
		lbConn.Log.Error("HandoverRequired is nil")
		return
	}

	lbConn.Log.Info("Handle HandoverRequired\n")
	for i := 0; i < len(HandoverRequired.ProtocolIEs.List); i++ {
		ie := HandoverRequired.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID // reject
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDHandoverType: // reject
			handoverType = ie.Value.HandoverType
			lbConn.Log.Trace("Decode IE HandoverType")
		case ngapType.ProtocolIEIDCause: // ignore
			cause = ie.Value.Cause
			lbConn.Log.Trace("Decode IE Cause")
		case ngapType.ProtocolIEIDTargetID: // reject
			targetID = ie.Value.TargetID
			lbConn.Log.Trace("Decode IE TargetID")
		case ngapType.ProtocolIEIDPDUSessionResourceListHORqd: // reject
			pDUSessionResourceListHORqd = ie.Value.PDUSessionResourceListHORqd
			lbConn.Log.Trace("Decode IE PDUSessionResourceListHORqd")
		case ngapType.ProtocolIEIDSourceToTargetTransparentContainer: // reject
			sourceToTargetTransparentContainer = ie.Value.SourceToTargetTransparentContainer
			lbConn.Log.Trace("Decode IE SourceToTargetTransparentContainer")
		}
	}

	if aMFUENGAPID == nil {
		lbConn.Log.Error("AmfUeNgapID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDAMFUENGAPID,
			ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if rANUENGAPID == nil {
		lbConn.Log.Error("RanUeNgapID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDRANUENGAPID,
			ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}

	if handoverType == nil {
		lbConn.Log.Error("handoverType is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDHandoverType,
			ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if targetID == nil {
		lbConn.Log.Error("targetID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDTargetID,
			ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if pDUSessionResourceListHORqd == nil {
		lbConn.Log.Error("pDUSessionResourceListHORqd is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
			ngapType.ProtocolIEIDPDUSessionResourceListHORqd, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if sourceToTargetTransparentContainer == nil {
		lbConn.Log.Error("sourceToTargetTransparentContainer is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
			ngapType.ProtocolIEIDSourceToTargetTransparentContainer, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procedureCode := ngapType.ProcedureCodeHandoverPreparation
		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
		procedureCriticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage,
			&procedureCriticality, &iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, nil, &criticalityDiagnostics)
		return
	}

	sourceUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if sourceUe == nil {
		lbConn.Log.Errorf("Cannot find UE for RAN_UE_NGAP_ID[%d] ", rANUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, &cause, nil)
		return
	}
	amfUe := sourceUe.AmfUe
	if amfUe == nil {
		lbConn.Log.Error("Cannot find amfUE from sourceUE")
		return
	}

	if targetID.Present != ngapType.TargetIDPresentTargetRANNodeID {
		lbConn.Log.Errorf("targetID type[%d] is not supported", targetID.Present)
		return
	}
	amfUe.SetOnGoing(sourceUe.Ran.AnType, &context.OnGoing{
		Procedure: context.OnGoingProcedureN2Handover,
	})
	if !amfUe.SecurityContextIsValid() {
		sourceUe.Log.Info("Handle Handover Preparation Failure [Authentication Failure]")
		cause = &ngapType.Cause{
			Present: ngapType.CausePresentNas,
			Nas: &ngapType.CauseNas{
				Value: ngapType.CauseNasPresentAuthenticationFailure,
			},
		}
		ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
		return
	}
	aMFSelf := context.AMF_Self()
	targetRanNodeId := ngapConvert.RanIdToModels(targetID.TargetRANNodeID.GlobalRANNodeID)
	targetRan, ok := aMFSelf.AmfRanFindByRanID(targetRanNodeId)
	if !ok {
		// handover between different AMF
		sourceUe.Log.Warnf("Handover required : cannot find target Ran Node Id[%+v] in this AMF", targetRanNodeId)
		sourceUe.Log.Error("Handover between different AMF has not been implemented yet")
		return
		// TODO: Send to T-AMF
		// Described in (23.502 4.9.1.3.2) step 3.Namf_Communication_CreateUEContext Request
	} else {
		// Handover in same AMF
		sourceUe.HandOverType.Value = handoverType.Value
		tai := ngapConvert.TaiToModels(targetID.TargetRANNodeID.SelectedTAI)
		targetId := models.NgRanTargetId{
			RanNodeId: &targetRanNodeId,
			Tai:       &tai,
		}
		var pduSessionReqList ngapType.PDUSessionResourceSetupListHOReq
		for _, pDUSessionResourceHoItem := range pDUSessionResourceListHORqd.List {
			pduSessionId := int32(pDUSessionResourceHoItem.PDUSessionID.Value)
			if smContext, exist := amfUe.SmContextFindByPDUSessionID(pduSessionId); exist {
				response, _, _, err := consumer.SendUpdateSmContextN2HandoverPreparing(amfUe, smContext,
					models.N2SmInfoType_HANDOVER_REQUIRED, pDUSessionResourceHoItem.HandoverRequiredTransfer, "", &targetId)
				if err != nil {
					sourceUe.Log.Errorf("consumer.SendUpdateSmContextN2HandoverPreparing Error: %+v", err)
				}
				if response == nil {
					sourceUe.Log.Errorf("SendUpdateSmContextN2HandoverPreparing Error for PduSessionId[%d]", pduSessionId)
					continue
				} else if response.BinaryDataN2SmInformation != nil {
					ngap_message.AppendPDUSessionResourceSetupListHOReq(&pduSessionReqList, pduSessionId,
						smContext.Snssai(), response.BinaryDataN2SmInformation)
				}
			}
		}
		if len(pduSessionReqList.List) == 0 {
			sourceUe.Log.Info("Handle Handover Preparation Failure [HoFailure In Target5GC NgranNode Or TargetSystem]")
			cause = &ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem,
				},
			}
			ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
			return
		}
		// Update NH
		amfUe.UpdateNH()
		ngap_message.SendHandoverRequest(sourceUe, targetRan, *cause, pduSessionReqList,
			*sourceToTargetTransparentContainer, false)
	}
}

func HandleHandoverCancel(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var cause *ngapType.Cause

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("Initiating Message is nil")
		return
	}
	HandoverCancel := initiatingMessage.Value.HandoverCancel
	if HandoverCancel == nil {
		lbConn.Log.Error("Handover Cancel is nil")
		return
	}

	lbConn.Log.Info("Handle Handover Cancel")
	for i := 0; i < len(HandoverCancel.ProtocolIEs.List); i++ {
		ie := HandoverCancel.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AMFUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RANUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			lbConn.Log.Trace("Decode IE Cause")
			if cause == nil {
				lbConn.Log.Error(cause, "cause is nil")
				return
			}
		}
	}

	sourceUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if sourceUe == nil {
		lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, &cause, nil)
		return
	}

	if sourceUe.AmfUeNgapId != aMFUENGAPID.Value {
		lbConn.Log.Warnf("Conflict AMF_UE_NGAP_ID : %d != %d", sourceUe.AmfUeNgapId, aMFUENGAPID.Value)
	}
	lbConn.Log.Tracef("Source: RAN_UE_NGAP_ID[%d] AMF_UE_NGAP_ID[%d]", sourceUe.RanUeNgapId, sourceUe.AmfUeNgapId)

	causePresent := ngapType.CausePresentRadioNetwork
	causeValue := ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem
	if cause != nil {
		causePresent, causeValue = printAndGetCause(lbConn, cause)
	}
	targetUe := sourceUe.TargetUe
	if targetUe == nil {
		// Described in (23.502 4.11.1.2.3) step 2
		// Todo : send to T-AMF invoke Namf_UeContextReleaseRequest(targetUe)
		lbConn.Log.Error("N2 Handover between AMF has not been implemented yet")
	} else {
		lbConn.Log.Tracef("Target : RAN_UE_NGAP_ID[%d] AMF_UE_NGAP_ID[%d]", targetUe.RanUeNgapId, targetUe.AmfUeNgapId)
		amfUe := sourceUe.AmfUe
		if amfUe != nil {
			amfUe.SmContextList.Range(func(key, value interface{}) bool {
				pduSessionID := key.(int32)
				smContext := value.(*context.SmContext)
				causeAll := context.CauseAll{
					NgapCause: &models.NgApCause{
						Group: int32(causePresent),
						Value: int32(causeValue),
					},
				}
				_, _, _, err := consumer.SendUpdateSmContextN2HandoverCanceled(amfUe, smContext, causeAll)
				if err != nil {
					sourceUe.Log.Errorf("Send UpdateSmContextN2HandoverCanceled Error for PduSessionId[%d]", pduSessionID)
				}
				return true
			})
		}
		ngap_message.SendUEContextReleaseCommand(targetUe, context.UeContextReleaseHandover, causePresent, causeValue)
		ngap_message.SendHandoverCancelAcknowledge(sourceUe, nil)
	}
}

func HandleUplinkRanStatusTransfer(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var rANStatusTransferTransparentContainer *ngapType.RANStatusTransferTransparentContainer
	var ranUe *context.RanUe

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // ignore
	if initiatingMessage == nil {
		lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	uplinkRanStatusTransfer := initiatingMessage.Value.UplinkRANStatusTransfer
	if uplinkRanStatusTransfer == nil {
		lbConn.Log.Error("UplinkRanStatusTransfer is nil")
		return
	}

	lbConn.Log.Info("Handle Uplink Ran Status Transfer")

	for _, ie := range uplinkRanStatusTransfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANStatusTransferTransparentContainer: // reject
			rANStatusTransferTransparentContainer = ie.Value.RANStatusTransferTransparentContainer
			lbConn.Log.Trace("Decode IE RANStatusTransferTransparentContainer")
			if rANStatusTransferTransparentContainer == nil {
				lbConn.Log.Error("RANStatusTransferTransparentContainer is nil")
			}
		}
	}

	ranUe = lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	ranUe.Log.Tracef("UE Context AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ranUe.Log.Error("AmfUe is nil")
		return
	}
	// send to T-AMF using N1N2MessageTransfer (R16)
}

func HandleNasNonDeliveryIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var nASPDU *ngapType.NASPDU
	var cause *ngapType.Cause

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	nASNonDeliveryIndication := initiatingMessage.Value.NASNonDeliveryIndication
	if nASNonDeliveryIndication == nil {
		lbConn.Log.Error("NASNonDeliveryIndication is nil")
		return
	}

	lbConn.Log.Info("Handle Nas Non Delivery Indication")

	for _, ie := range nASNonDeliveryIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDNASPDU:
			nASPDU = ie.Value.NASPDU
			if nASPDU == nil {
				lbConn.Log.Error("NasPdu is nil")
				return
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			if cause == nil {
				lbConn.Log.Error("Cause is nil")
				return
			}
		}
	}

	ranUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	lbConn.Log.Tracef("RanUeNgapID[%d] AmfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)

	printAndGetCause(lbConn, cause)

	nas.HandleNAS(ranUe, ngapType.ProcedureCodeNASNonDeliveryIndication, nASPDU.Value)
}

func HandleRanConfigurationUpdate(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var rANNodeName *ngapType.RANNodeName
	var supportedTAList *ngapType.SupportedTAList
	var pagingDRX *ngapType.PagingDRX

	var cause ngapType.Cause

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}

	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("Initiating Message is nil")
		return
	}
	rANConfigurationUpdate := initiatingMessage.Value.RANConfigurationUpdate
	if rANConfigurationUpdate == nil {
		lbConn.Log.Error("RAN Configuration is nil")
		return
	}
	lbConn.Log.Info("Handle Ran Configuration Update")
	for i := 0; i < len(rANConfigurationUpdate.ProtocolIEs.List); i++ {
		ie := rANConfigurationUpdate.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANNodeName:
			rANNodeName = ie.Value.RANNodeName
			if rANNodeName == nil {
				lbConn.Log.Error("RAN Node Name is nil")
				return
			}
			lbConn.Log.Tracef("Decode IE RANNodeName = [%s]", rANNodeName.Value)
		case ngapType.ProtocolIEIDSupportedTAList:
			supportedTAList = ie.Value.SupportedTAList
			lbConn.Log.Trace("Decode IE SupportedTAList")
			if supportedTAList == nil {
				lbConn.Log.Error("Supported TA List is nil")
				return
			}
		case ngapType.ProtocolIEIDDefaultPagingDRX:
			pagingDRX = ie.Value.DefaultPagingDRX
			if pagingDRX == nil {
				lbConn.Log.Error("PagingDRX is nil")
				return
			}
			lbConn.Log.Tracef("Decode IE PagingDRX = [%d]", pagingDRX.Value)
		}
	}

	for i := 0; i < len(supportedTAList.List); i++ {
		supportedTAItem := supportedTAList.List[i]
		tac := hex.EncodeToString(supportedTAItem.TAC.Value)
		capOfSupportTai := cap(lbConn.SupportedTAList)
		for j := 0; j < len(supportedTAItem.BroadcastPLMNList.List); j++ {
			supportedTAI := context.NewSupportedTAI()
			supportedTAI.Tai.Tac = tac
			broadcastPLMNItem := supportedTAItem.BroadcastPLMNList.List[j]
			plmnId := ngapConvert.PlmnIdToModels(broadcastPLMNItem.PLMNIdentity)
			supportedTAI.Tai.PlmnId = &plmnId
			capOfSNssaiList := cap(supportedTAI.SNssaiList)
			for k := 0; k < len(broadcastPLMNItem.TAISliceSupportList.List); k++ {
				tAISliceSupportItem := broadcastPLMNItem.TAISliceSupportList.List[k]
				if len(supportedTAI.SNssaiList) < capOfSNssaiList {
					supportedTAI.SNssaiList = append(supportedTAI.SNssaiList, ngapConvert.SNssaiToModels(tAISliceSupportItem.SNSSAI))
				} else {
					break
				}
			}
			lbConn.Log.Tracef("PLMN_ID[MCC:%s MNC:%s] TAC[%s]", plmnId.Mcc, plmnId.Mnc, tac)
			if len(lbConn.SupportedTAList) < capOfSupportTai {
				lbConn.SupportedTAList = append(lbConn.SupportedTAList, supportedTAI)
			} else {
				break
			}
		}
	}

	if len(lbConn.SupportedTAList) == 0 {
		lbConn.Log.Warn("RanConfigurationUpdate failure: No supported TA exist in RanConfigurationUpdate")
		cause.Present = ngapType.CausePresentMisc
		cause.Misc = &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		}
	} else {
		var found bool
		for i, tai := range lbConn.SupportedTAList {
			if context.InTaiList(tai.Tai, context.AMF_Self().SupportTaiLists) {
				lbConn.Log.Tracef("SERVED_TAI_INDEX[%d]", i)
				found = true
				break
			}
		}
		if !found {
			lbConn.Log.Warn("RanConfigurationUpdate failure: Cannot find Served TAI in AMF")
			cause.Present = ngapType.CausePresentMisc
			cause.Misc = &ngapType.CauseMisc{
				Value: ngapType.CauseMiscPresentUnknownPLMN,
			}
		}
	}

	if cause.Present == ngapType.CausePresentNothing {
		lbConn.Log.Info("Handle RanConfigurationUpdateAcknowledge")
		ngap_message.SendRanConfigurationUpdateAcknowledge(lbConn, nil)
	} else {
		lbConn.Log.Info("Handle RanConfigurationUpdateAcknowledgeFailure")
		ngap_message.SendRanConfigurationUpdateFailure(lbConn, cause, nil)
	}
}

func HandleUplinkRanConfigurationTransfer(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var sONConfigurationTransferUL *ngapType.SONConfigurationTransfer

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	uplinkRANConfigurationTransfer := initiatingMessage.Value.UplinkRANConfigurationTransfer
	if uplinkRANConfigurationTransfer == nil {
		lbConn.Log.Error("ErrorIndication is nil")
		return
	}

	for _, ie := range uplinkRANConfigurationTransfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDSONConfigurationTransferUL: // optional, ignore
			sONConfigurationTransferUL = ie.Value.SONConfigurationTransferUL
			lbConn.Log.Trace("Decode IE SONConfigurationTransferUL")
			if sONConfigurationTransferUL == nil {
				lbConn.Log.Warn("sONConfigurationTransferUL is nil")
			}
		}
	}

	if sONConfigurationTransferUL != nil {
		targetRanNodeID := ngapConvert.RanIdToModels(sONConfigurationTransferUL.TargetRANNodeID.GlobalRANNodeID)

		if targetRanNodeID.GNbId.GNBValue != "" {
			lbConn.Log.Tracef("targerRanID [%s]", targetRanNodeID.GNbId.GNBValue)
		}

		aMFSelf := context.AMF_Self()

		targetRan, ok := aMFSelf.AmfRanFindByRanID(targetRanNodeID)
		if !ok {
			lbConn.Log.Warn("targetRan is nil")
		}

		ngap_message.SendDownlinkRanConfigurationTransfer(targetRan, sONConfigurationTransferUL)
	}
}

func HandleUplinkUEAssociatedNRPPATransport(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var routingID *ngapType.RoutingID
	var nRPPaPDU *ngapType.NRPPaPDU

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	uplinkUEAssociatedNRPPaTransport := initiatingMessage.Value.UplinkUEAssociatedNRPPaTransport
	if uplinkUEAssociatedNRPPaTransport == nil {
		lbConn.Log.Error("uplinkUEAssociatedNRPPaTransport is nil")
		return
	}

	lbConn.Log.Info("Handle Uplink UE Associated NRPPA Transpor")

	for _, ie := range uplinkUEAssociatedNRPPaTransport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE aMFUENGAPID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE rANUENGAPID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRoutingID: // reject
			routingID = ie.Value.RoutingID
			lbConn.Log.Trace("Decode IE routingID")
			if routingID == nil {
				lbConn.Log.Error("routingID is nil")
				return
			}
		case ngapType.ProtocolIEIDNRPPaPDU: // reject
			nRPPaPDU = ie.Value.NRPPaPDU
			lbConn.Log.Trace("Decode IE nRPPaPDU")
			if nRPPaPDU == nil {
				lbConn.Log.Error("nRPPaPDU is nil")
				return
			}
		}
	}

	ranUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	lbConn.Log.Tracef("RanUeNgapId[%d] AmfUeNgapId[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)

	ranUe.RoutingID = hex.EncodeToString(routingID.Value)

	// TODO: Forward NRPPaPDU to LMF
}

func HandleUplinkNonUEAssociatedNRPPATransport(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var routingID *ngapType.RoutingID
	var nRPPaPDU *ngapType.NRPPaPDU

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("Initiating Message is nil")
		return
	}
	uplinkNonUEAssociatedNRPPATransport := initiatingMessage.Value.UplinkNonUEAssociatedNRPPaTransport
	if uplinkNonUEAssociatedNRPPATransport == nil {
		lbConn.Log.Error("Uplink Non UE Associated NRPPA Transport is nil")
		return
	}

	lbConn.Log.Info("Handle Uplink Non UE Associated NRPPA Transport")

	for i := 0; i < len(uplinkNonUEAssociatedNRPPATransport.ProtocolIEs.List); i++ {
		ie := uplinkNonUEAssociatedNRPPATransport.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRoutingID:
			routingID = ie.Value.RoutingID
			lbConn.Log.Trace("Decode IE RoutingID")

		case ngapType.ProtocolIEIDNRPPaPDU:
			nRPPaPDU = ie.Value.NRPPaPDU
			lbConn.Log.Trace("Decode IE NRPPaPDU")
		}
	}

	if routingID == nil {
		lbConn.Log.Error("RoutingID is nil")
		return
	}
	// Forward routingID to LMF
	// Described in (23.502 4.13.5.6)

	if nRPPaPDU == nil {
		lbConn.Log.Error("NRPPaPDU is nil")
		return
	}
	// TODO: Forward NRPPaPDU to LMF
}

func HandleLocationReport(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var userLocationInformation *ngapType.UserLocationInformation
	var uEPresenceInAreaOfInterestList *ngapType.UEPresenceInAreaOfInterestList
	var locationReportingRequestType *ngapType.LocationReportingRequestType

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	locationReport := initiatingMessage.Value.LocationReport
	if locationReport == nil {
		lbConn.Log.Error("LocationReport is nil")
		return
	}

	lbConn.Log.Info("Handle Location Report")
	for _, ie := range locationReport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // ignore
			userLocationInformation = ie.Value.UserLocationInformation
			lbConn.Log.Trace("Decode IE userLocationInformation")
			if userLocationInformation == nil {
				lbConn.Log.Warn("userLocationInformation is nil")
			}
		case ngapType.ProtocolIEIDUEPresenceInAreaOfInterestList: // optional, ignore
			uEPresenceInAreaOfInterestList = ie.Value.UEPresenceInAreaOfInterestList
			lbConn.Log.Trace("Decode IE uEPresenceInAreaOfInterestList")
			if uEPresenceInAreaOfInterestList == nil {
				lbConn.Log.Warn("uEPresenceInAreaOfInterestList is nil [optional]")
			}
		case ngapType.ProtocolIEIDLocationReportingRequestType: // ignore
			locationReportingRequestType = ie.Value.LocationReportingRequestType
			lbConn.Log.Trace("Decode IE LocationReportingRequestType")
			if locationReportingRequestType == nil {
				lbConn.Log.Warn("LocationReportingRequestType is nil")
			}
		}
	}

	ranUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	ranUe.UpdateLocation(userLocationInformation)

	ranUe.Log.Tracef("Report Area[%d]", locationReportingRequestType.ReportArea.Value)

	switch locationReportingRequestType.EventType.Value {
	case ngapType.EventTypePresentDirect:
		ranUe.Log.Trace("To report directly")

	case ngapType.EventTypePresentChangeOfServeCell:
		ranUe.Log.Trace("To report upon change of serving cell")

	case ngapType.EventTypePresentUePresenceInAreaOfInterest:
		ranUe.Log.Trace("To report UE presence in the area of interest")
		for _, uEPresenceInAreaOfInterestItem := range uEPresenceInAreaOfInterestList.List {
			uEPresence := uEPresenceInAreaOfInterestItem.UEPresence.Value
			referenceID := uEPresenceInAreaOfInterestItem.LocationReportingReferenceID.Value

			for _, AOIitem := range locationReportingRequestType.AreaOfInterestList.List {
				if referenceID == AOIitem.LocationReportingReferenceID.Value {
					lbConn.Log.Tracef("uEPresence[%d], presence AOI ReferenceID[%d]", uEPresence, referenceID)
				}
			}
		}

	case ngapType.EventTypePresentStopChangeOfServeCell:
		ranUe.Log.Trace("To stop reporting at change of serving cell")
		ngap_message.SendLocationReportingControl(ranUe, nil, 0, locationReportingRequestType.EventType)
		// TODO: Clear location report

	case ngapType.EventTypePresentStopUePresenceInAreaOfInterest:
		ranUe.Log.Trace("To stop reporting UE presence in the area of interest")
		ranUe.Log.Tracef("ReferenceID To Be Cancelled[%d]",
			locationReportingRequestType.LocationReportingReferenceIDToBeCancelled.Value)
		// TODO: Clear location report

	case ngapType.EventTypePresentCancelLocationReportingForTheUe:
		ranUe.Log.Trace("To cancel location reporting for the UE")
		// TODO: Clear location report
	}
}

func HandleUERadioCapabilityInfoIndication(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	var uERadioCapability *ngapType.UERadioCapability
	var uERadioCapabilityForPaging *ngapType.UERadioCapabilityForPaging

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("Initiating Message is nil")
		return
	}
	uERadioCapabilityInfoIndication := initiatingMessage.Value.UERadioCapabilityInfoIndication
	if uERadioCapabilityInfoIndication == nil {
		lbConn.Log.Error("UERadioCapabilityInfoIndication is nil")
		return
	}

	lbConn.Log.Info("Handle UE Radio Capability Info Indication")

	for i := 0; i < len(uERadioCapabilityInfoIndication.ProtocolIEs.List); i++ {
		ie := uERadioCapabilityInfoIndication.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDUERadioCapability:
			uERadioCapability = ie.Value.UERadioCapability
			lbConn.Log.Trace("Decode IE UERadioCapability")
			if uERadioCapability == nil {
				lbConn.Log.Error("UERadioCapability is nil")
				return
			}
		case ngapType.ProtocolIEIDUERadioCapabilityForPaging:
			uERadioCapabilityForPaging = ie.Value.UERadioCapabilityForPaging
			lbConn.Log.Trace("Decode IE UERadioCapabilityForPaging")
			if uERadioCapabilityForPaging == nil {
				lbConn.Log.Error("UERadioCapabilityForPaging is nil")
				return
			}
		}
	}

	ranUe := lbConn.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		lbConn.Log.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
	lbConn.Log.Tracef("RanUeNgapID[%d] AmfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)
	amfUe := ranUe.AmfUe

	if amfUe == nil {
		ranUe.Log.Errorln("amfUe is nil")
		return
	}
	if uERadioCapability != nil {
		amfUe.UeRadioCapability = hex.EncodeToString(uERadioCapability.Value)
	}
	if uERadioCapabilityForPaging != nil {
		amfUe.UeRadioCapabilityForPaging = &context.UERadioCapabilityForPaging{}
		if uERadioCapabilityForPaging.UERadioCapabilityForPagingOfNR != nil {
			amfUe.UeRadioCapabilityForPaging.NR = hex.EncodeToString(
				uERadioCapabilityForPaging.UERadioCapabilityForPagingOfNR.Value)
		}
		if uERadioCapabilityForPaging.UERadioCapabilityForPagingOfEUTRA != nil {
			amfUe.UeRadioCapabilityForPaging.EUTRA = hex.EncodeToString(
				uERadioCapabilityForPaging.UERadioCapabilityForPagingOfEUTRA.Value)
		}
	}

	// TS 38.413 8.14.1.2/TS 23.502 4.2.8a step5/TS 23.501, clause 5.4.4.1.
	// send its most up to date UE Radio Capability information to the RAN in the N2 REQUEST message.
}

func HandleAMFconfigurationUpdateFailure(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		lbConn.Log.Error("Unsuccessful Message is nil")
		return
	}

	AMFconfigurationUpdateFailure := unsuccessfulOutcome.Value.AMFConfigurationUpdateFailure
	if AMFconfigurationUpdateFailure == nil {
		lbConn.Log.Error("AMFConfigurationUpdateFailure is nil")
		return
	}

	lbConn.Log.Info("Handle AMF Confioguration Update Failure")

	for _, ie := range AMFconfigurationUpdateFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			lbConn.Log.Trace("Decode IE Cause")
			if cause == nil {
				lbConn.Log.Error("Cause is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	//	TODO: Time To Wait

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}
}

func HandleAMFconfigurationUpdateAcknowledge(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFTNLAssociationSetupList *ngapType.AMFTNLAssociationSetupList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	var aMFTNLAssociationFailedToSetupList *ngapType.TNLAssociationList
	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		lbConn.Log.Error("SuccessfulOutcome is nil")
		return
	}
	aMFConfigurationUpdateAcknowledge := successfulOutcome.Value.AMFConfigurationUpdateAcknowledge
	if aMFConfigurationUpdateAcknowledge == nil {
		lbConn.Log.Error("AMFConfigurationUpdateAcknowledge is nil")
		return
	}

	lbConn.Log.Info("Handle AMF Configuration Update Acknowledge")

	for i := 0; i < len(aMFConfigurationUpdateAcknowledge.ProtocolIEs.List); i++ {
		ie := aMFConfigurationUpdateAcknowledge.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFTNLAssociationSetupList:
			aMFTNLAssociationSetupList = ie.Value.AMFTNLAssociationSetupList
			lbConn.Log.Trace("Decode IE AMFTNLAssociationSetupList")
			if aMFTNLAssociationSetupList == nil {
				lbConn.Log.Error("AMFTNLAssociationSetupList is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE Criticality Diagnostics")

		case ngapType.ProtocolIEIDAMFTNLAssociationFailedToSetupList:
			aMFTNLAssociationFailedToSetupList = ie.Value.AMFTNLAssociationFailedToSetupList
			lbConn.Log.Trace("Decode IE AMFTNLAssociationFailedToSetupList")
			if aMFTNLAssociationFailedToSetupList == nil {
				lbConn.Log.Error("AMFTNLAssociationFailedToSetupList is nil")
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
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	errorIndication := initiatingMessage.Value.ErrorIndication
	if errorIndication == nil {
		lbConn.Log.Error("ErrorIndication is nil")
		return
	}

	for _, ie := range errorIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				lbConn.Log.Error("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				lbConn.Log.Error("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			lbConn.Log.Trace("Decode IE Cause")
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			lbConn.Log.Trace("Decode IE CriticalityDiagnostics")
		}
	}

	if cause == nil && criticalityDiagnostics == nil {
		lbConn.Log.Error("[ErrorIndication] both Cause IE and CriticalityDiagnostics IE are nil, should have at least one")
		return
	}

	if cause != nil {
		printAndGetCause(lbConn, cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(lbConn, criticalityDiagnostics)
	}

	// TODO: handle error based on cause/criticalityDiagnostics
}

func HandleCellTrafficTrace(lbConn *context.LBConn, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var nGRANTraceID *ngapType.NGRANTraceID
	var nGRANCGI *ngapType.NGRANCGI
	var traceCollectionEntityIPAddress *ngapType.TransportLayerAddress

	var ranUe *context.RanUe

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if lbConn == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		lbConn.Log.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // ignore
	if initiatingMessage == nil {
		lbConn.Log.Error("InitiatingMessage is nil")
		return
	}
	cellTrafficTrace := initiatingMessage.Value.CellTrafficTrace
	if cellTrafficTrace == nil {
		lbConn.Log.Error("CellTrafficTrace is nil")
		return
	}

	lbConn.Log.Info("Handle Cell Traffic Trace")

	for _, ie := range cellTrafficTrace.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			lbConn.Log.Trace("Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			lbConn.Log.Trace("Decode IE RanUeNgapID")

		case ngapType.ProtocolIEIDNGRANTraceID: // ignore
			nGRANTraceID = ie.Value.NGRANTraceID
			lbConn.Log.Trace("Decode IE NGRANTraceID")
		case ngapType.ProtocolIEIDNGRANCGI: // ignore
			nGRANCGI = ie.Value.NGRANCGI
			lbConn.Log.Trace("Decode IE NGRANCGI")
		case ngapType.ProtocolIEIDTraceCollectionEntityIPAddress: // ignore
			traceCollectionEntityIPAddress = ie.Value.TraceCollectionEntityIPAddress
			lbConn.Log.Trace("Decode IE TraceCollectionEntityIPAddress")
		}
	}
	if aMFUENGAPID == nil {
		lbConn.Log.Error("AmfUeNgapID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDAMFUENGAPID,
			ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if rANUENGAPID == nil {
		lbConn.Log.Error("RanUeNgapID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDRANUENGAPID,
			ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procedureCode := ngapType.ProcedureCodeCellTrafficTrace
		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
		procedureCriticality := ngapType.CriticalityPresentIgnore
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality,
			&iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, nil, &criticalityDiagnostics)
		return
	}

	if aMFUENGAPID != nil {
		ranUe = context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			lbConn.Log.Errorf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
			cause := ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
				},
			}
			ngap_message.SendErrorIndication(lbConn, aMFUENGAPID, rANUENGAPID, &cause, nil)
			return
		}
	}

	lbConn.Log.Debugf("UE: AmfUeNgapID[%d], RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)

	ranUe.Trsr = hex.EncodeToString(nGRANTraceID.Value[6:])

	ranUe.Log.Tracef("TRSR[%s]", ranUe.Trsr)

	switch nGRANCGI.Present {
	case ngapType.NGRANCGIPresentNRCGI:
		plmnID := ngapConvert.PlmnIdToModels(nGRANCGI.NRCGI.PLMNIdentity)
		cellID := ngapConvert.BitStringToHex(&nGRANCGI.NRCGI.NRCellIdentity.Value)
		ranUe.Log.Debugf("NRCGI[plmn: %s, cellID: %s]", plmnID, cellID)
	case ngapType.NGRANCGIPresentEUTRACGI:
		plmnID := ngapConvert.PlmnIdToModels(nGRANCGI.EUTRACGI.PLMNIdentity)
		cellID := ngapConvert.BitStringToHex(&nGRANCGI.EUTRACGI.EUTRACellIdentity.Value)
		ranUe.Log.Debugf("EUTRACGI[plmn: %s, cellID: %s]", plmnID, cellID)
	}

	tceIpv4, tceIpv6 := ngapConvert.IPAddressToString(*traceCollectionEntityIPAddress)
	if tceIpv4 != "" {
		ranUe.Log.Debugf("TCE IP Address[v4: %s]", tceIpv4)
	}
	if tceIpv6 != "" {
		ranUe.Log.Debugf("TCE IP Address[v6: %s]", tceIpv6)
	}

	// TODO: TS 32.422 4.2.2.10
	// When AMF receives this new NG signalling message containing the Trace Recording Session Reference (TRSR)
	// and Trace Reference (TR), the AMF shall look up the SUPI/IMEI(SV) of the given call from its database and
	// shall send the SUPI/IMEI(SV) numbers together with the Trace Recording Session Reference and Trace Reference
	// to the Trace Collection Entity.
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
