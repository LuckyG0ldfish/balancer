package ngap 

import (
	"reflect"
	"github.com/free5gc/ngap/ngapType"
)

func getAmfID(pdu *ngapType.NGAPPDU) (amfID int, amfpoint *int){
	var value interface{}
	var present int
	switch pdu.Present { 
	case ngapType.NGAPPDUPresentInitiatingMessage:
		value = pdu.InitiatingMessage
		
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		value = pdu.SuccessfulOutcome
		
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		value = pdu.UnsuccessfulOutcome
		
	}
	// temp := reflect.ValueOf(value)
	// fields := make([]interface{}, temp.NumField())

	// needed := pduFields

	return 
}

func getInitiatingPointer(message *ngapType.InitiatingMessage) (amf, ran *int) {
	switch message.ProcedureCode.Value {
	case ngapType.ProcedureCodeNGSetup:
		message.Value.ProcedureCodeNGSetup
	case ngapType.ProcedureCodeInitialUEMessage:
		HandleInitialUEMessage(ran, pdu)
	case ngapType.ProcedureCodeUplinkNASTransport:
		HandleUplinkNasTransport(ran, pdu)
	case ngapType.ProcedureCodeNGReset:
		HandleNGReset(ran, pdu)
	case ngapType.ProcedureCodeHandoverCancel:
		HandleHandoverCancel(ran, pdu)
	case ngapType.ProcedureCodeUEContextReleaseRequest:
		HandleUEContextReleaseRequest(ran, pdu)
	case ngapType.ProcedureCodeNASNonDeliveryIndication:
		HandleNasNonDeliveryIndication(ran, pdu)
	case ngapType.ProcedureCodeLocationReportingFailureIndication:
		HandleLocationReportingFailureIndication(ran, pdu)
	case ngapType.ProcedureCodeErrorIndication:
		HandleErrorIndication(ran, pdu)
	case ngapType.ProcedureCodeUERadioCapabilityInfoIndication:
		HandleUERadioCapabilityInfoIndication(ran, pdu)
	case ngapType.ProcedureCodeHandoverNotification:
		HandleHandoverNotify(ran, pdu)
	case ngapType.ProcedureCodeHandoverPreparation:
		HandleHandoverRequired(ran, pdu)
	case ngapType.ProcedureCodeRANConfigurationUpdate:
		HandleRanConfigurationUpdate(ran, pdu)
	case ngapType.ProcedureCodeRRCInactiveTransitionReport:
		HandleRRCInactiveTransitionReport(ran, pdu)
	case ngapType.ProcedureCodePDUSessionResourceNotify:
		HandlePDUSessionResourceNotify(ran, pdu)
	case ngapType.ProcedureCodePathSwitchRequest:
		HandlePathSwitchRequest(ran, pdu)
	case ngapType.ProcedureCodeLocationReport:
		HandleLocationReport(ran, pdu)
	case ngapType.ProcedureCodeUplinkUEAssociatedNRPPaTransport:
		HandleUplinkUEAssociatedNRPPATransport(ran, pdu)
	case ngapType.ProcedureCodeUplinkRANConfigurationTransfer:
		HandleUplinkRanConfigurationTransfer(ran, pdu)
	case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
		HandlePDUSessionResourceModifyIndication(ran, pdu)
	case ngapType.ProcedureCodeCellTrafficTrace:
		HandleCellTrafficTrace(ran, pdu)
	case ngapType.ProcedureCodeUplinkRANStatusTransfer:
		HandleUplinkRanStatusTransfer(ran, pdu)
	case ngapType.ProcedureCodeUplinkNonUEAssociatedNRPPaTransport:
		HandleUplinkNonUEAssociatedNRPPATransport(ran, pdu)
	default:
		ran.Log.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, initiatingMessage.ProcedureCode.Value)
	}

	return
}

func main() {
    x := struct{Foo string; Bar int }{"foo", 2}

    v := reflect.ValueOf(x)

    values := make([]interface{}, v.NumField())

    for i := 0; i < v.NumField(); i++ {
        values[i] = v.Field(i).Interface()
    }

    fmt.Println(values)
}