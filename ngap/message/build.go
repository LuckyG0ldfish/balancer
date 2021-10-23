package message

import(
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/ngap"
	"github.com/LuckyG0ldfish/balancer/context"
)

func BuildNGSetupResponse() ([]byte, error) {
	lbSelf := context.LB_Self()
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeNGSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject
	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentNGSetupResponse
	successfulOutcome.Value.NGSetupResponse = new(ngapType.NGSetupResponse)

	nGSetupResponse := successfulOutcome.Value.NGSetupResponse
	nGSetupResponseIEs := &nGSetupResponse.ProtocolIEs

	// AMFName
	ie := ngapType.NGSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFName
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupResponseIEsPresentAMFName
	ie.Value.AMFName = new(ngapType.AMFName)

	aMFName := ie.Value.AMFName
	aMFName.Value = lbSelf.Name

	nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	// ServedGUAMIList
	ie = ngapType.NGSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDServedGUAMIList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupResponseIEsPresentServedGUAMIList
	ie.Value.ServedGUAMIList = new(ngapType.ServedGUAMIList)

	// servedGUAMIList := ie.Value.ServedGUAMIList
	// for _, guami := range lbSelf.ServedGuamiList {
	// 	servedGUAMIItem := ngapType.ServedGUAMIItem{}
	// 	servedGUAMIItem.GUAMI.PLMNIdentity = ngapConvert.PlmnIdToNgap(*guami.PlmnId)
	// 	regionId, setId, prtId := ngapConvert.AmfIdToNgap(guami.AmfId)
	// 	servedGUAMIItem.GUAMI.AMFRegionID.Value = regionId
	// 	servedGUAMIItem.GUAMI.AMFSetID.Value = setId
	// 	servedGUAMIItem.GUAMI.AMFPointer.Value = prtId
	// 	servedGUAMIList.List = append(servedGUAMIList.List, servedGUAMIItem)
	// }

	nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	// relativeAMFCapacity
	ie = ngapType.NGSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRelativeAMFCapacity
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGSetupResponseIEsPresentRelativeAMFCapacity
	ie.Value.RelativeAMFCapacity = new(ngapType.RelativeAMFCapacity)
	// relativeAMFCapacity := ie.Value.RelativeAMFCapacity
	// relativeAMFCapacity.Value = lbSelf.RelativeCapacity

	nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	// ServedGUAMIList
	ie = ngapType.NGSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPLMNSupportList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupResponseIEsPresentPLMNSupportList
	ie.Value.PLMNSupportList = new(ngapType.PLMNSupportList)

	// pLMNSupportList := ie.Value.PLMNSupportList
	// for _, plmnItem := range lbSelf.PlmnSupportList {
	// 	pLMNSupportItem := ngapType.PLMNSupportItem{}
	// 	pLMNSupportItem.PLMNIdentity = ngapConvert.PlmnIdToNgap(plmnItem.PlmnId)
	// 	for _, snssai := range plmnItem.SNssaiList {
	// 		sliceSupportItem := ngapType.SliceSupportItem{}
	// 		sliceSupportItem.SNSSAI = ngapConvert.SNssaiToNgap(snssai)
	// 		pLMNSupportItem.SliceSupportList.List = append(pLMNSupportItem.SliceSupportList.List, sliceSupportItem)
	// 	}
	// 	pLMNSupportList.List = append(pLMNSupportList.List, pLMNSupportItem)
	// }

	nGSetupResponseIEs.List = append(nGSetupResponseIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildNGSetupFailure(cause ngapType.Cause) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeNGSetup
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject
	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentNGSetupFailure
	unsuccessfulOutcome.Value.NGSetupFailure = new(ngapType.NGSetupFailure)

	nGSetupFailure := unsuccessfulOutcome.Value.NGSetupFailure
	nGSetupFailureIEs := &nGSetupFailure.ProtocolIEs

	// Cause
	ie := ngapType.NGSetupFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGSetupFailureIEsPresentCause
	ie.Value.Cause = &cause

	nGSetupFailureIEs.List = append(nGSetupFailureIEs.List, ie)

	return ngap.Encoder(pdu)
}