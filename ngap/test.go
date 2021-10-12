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
		//case ngapType.ProcedureCodeNGReset:
		//	handler.HandleNGReset(amf, pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			// handler.HandleInitialContextSetupRequest(emulatorCtx.MainSCTPConnection, pdu)
			
		//case ngapType.ProcedureCodeUEContextModification:
		//	handler.HandleUEContextModificationRequest(amf, pdu)
		case ngapType.ProcedureCodeUEContextRelease:
			handler.HandleUEContextReleaseCommand(emulatorCtx.MainSCTPConnection, pdu)
		case ngapType.ProcedureCodeDownlinkNASTransport:
			handler.HandleDownlinkNASTransport(emulatorCtx.MainSCTPConnection, pdu)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			handler.HandlePDUSessionResourceSetupRequest(emulatorCtx.MainSCTPConnection, pdu)
		// TODO: This will be commented for the time being, after adding other procedures will be uncommented.
		//case ngapType.ProcedureCodePDUSessionResourceModify:
		//	handler.HandlePDUSessionResourceModifyRequest(amf, pdu)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			handler.HandlePDUSessionResourceReleaseCommand(emulatorCtx.MainSCTPConnection, pdu)
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
			NGAPLog.Warnf("Not implemented NGAP message(initiatingMessage), procedureCode:%d]\n",
				initiatingMessage.ProcedureCode.Value)
		}

	return
}

func getSucessfulPointer(message *ngapType.InitiatingMessage) (amf, ran *int) {
	switch message.ProcedureCode.Value {
	case ngapType.ProcedureCodeNGSetup:
		handler.HandleNGSetupResponse(conn, pdu)
	//case ngapType.ProcedureCodeNGReset:
	//	handler.HandleNGResetAcknowledge(amf, pdu)
	//case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
	//	handler.HandlePDUSessionResourceModifyConfirm(amf, pdu)
	//case ngapType.ProcedureCodeRANConfigurationUpdate:
	//	handler.HandleRANConfigurationUpdateAcknowledge(amf, pdu)
	default:
		NGAPLog.Warnf("Not implemented NGAP message(successfulOutcome), procedureCode:%d]\n",
			successfulOutcome.ProcedureCode.Value)
	}

	return
}

func getUnsucessfulPointer(message *ngapType.InitiatingMessage) (amf, ran *int){
	switch message.ProcedureCode.Value {
		//case ngapType.ProcedureCodeNGSetup:
		//	handler.HandleNGSetupFailure(sctpAddr, conn, pdu)
		//case ngapType.ProcedureCodeRANConfigurationUpdate:
		//	handler.HandleRANConfigurationUpdateFailure(amf, pdu)
	default:
		NGAPLog.Warnf("Not implemented NGAP message(unsuccessfulOutcome), procedureCode:%d]\n",
			unsuccessfulOutcome.ProcedureCode.Value)
		
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