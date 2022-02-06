package gmm

import (
	// "encoding/hex"
	// "fmt"
	// "crypto/sha256"

	"github.com/LuckyG0ldfish/balancer/context"
	// "github.com/free5gc/nas"
	// "github.com/free5gc/amf/logger"
	"github.com/free5gc/nas/nasMessage"
	// "github.com/free5gc/nas/security"
	"github.com/free5gc/openapi/models"
)

// TS 24.501 5.4.1
// TODO
func HandleAuthenticationResponse(ue *context.LbUe, accessType models.AccessType,
	authenticationResponse *nasMessage.AuthenticationResponse) error {
	// logger.GmmLog.Info("Handle Authentication Response")

	// // if ue.T3560 != nil {
	// // 	ue.T3560.Stop()
	// // 	ue.T3560 = nil // clear the timer
	// // }

	// // if ue.AuthenticationCtx == nil {
	// // 	return fmt.Errorf("Ue Authentication Context is nil")
	// // }

	
	// switch ue.AuthenticationCtx.AuthType {
	// case models.AuthType__5_G_AKA:
	// 	var av5gAka models.Av5gAka
	// 	if err := mapstructure.Decode(ue.AuthenticationCtx.Var5gAuthData, &av5gAka); err != nil {
	// 		return fmt.Errorf("Var5gAuthData Convert Type Error")
	// 	}
	// 	resStar := authenticationResponse.AuthenticationResponseParameter.GetRES()

	// 	// Calculate HRES* (TS 33.501 Annex A.5)
	// 	p0, err := hex.DecodeString(av5gAka.Rand)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	p1 := resStar[:]
	// 	concat := append(p0, p1...)
	// 	hResStarBytes := sha256.Sum256(concat)
	// 	hResStar := hex.EncodeToString(hResStarBytes[16:])

	// 	if hResStar != av5gAka.HxresStar {
	// 		logger.GmmLog.Errorf("HRES* Validation Failure (received: %s, expected: %s)", hResStar, av5gAka.HxresStar)

	// 		if ue.IdentityTypeUsedForRegistration == nasMessage.MobileIdentity5GSType5gGuti && ue.IdentityRequestSendTimes == 0 {
	// 			ue.IdentityRequestSendTimes++
	// 			gmm_message.SendIdentityRequest(ue.RanUe[accessType], accessType, nasMessage.MobileIdentity5GSTypeSuci)
	// 			return nil
	// 		} else {
	// 			// gmm_message.SendAuthenticationReject(ue.RanUe[accessType], "")
	// 			return // GmmFSM.SendEvent(ue.State[accessType], AuthFailEvent, fsm.ArgsType{
	// 			// 	ArgAmfUe:      ue,
	// 			// 	ArgAccessType: accessType,
	// 			// })
	// 		}
	// 	}

	// 	response, problemDetails, err := consumer.SendAuth5gAkaConfirmRequest(ue, hex.EncodeToString(resStar[:]))
	// 	if err != nil {
	// 		return err
	// 	} else if problemDetails != nil {
	// 		logger.GmmLog.Debugf("Auth5gAkaConfirm Error[Problem Detail: %+v]", problemDetails)
	// 		return nil
	// 	}
	// 	switch response.AuthResult {
	// 	case models.AuthResult_SUCCESS:
	// 		// ue.UnauthenticatedSupi = false
	// 		ue.Kseaf = response.Kseaf
	// 		ue.Supi = response.Supi
	// 		ue.DerivateKamf()
	// 		logger.GmmLog.Debugln("ue.DerivateKamf()", ue.Kamf)
	// 		// return GmmFSM.SendEvent(ue.State[accessType], AuthSuccessEvent, fsm.ArgsType{
	// 		// 	ArgAmfUe:      ue,
	// 		// 	ArgAccessType: accessType,
	// 		// 	ArgEAPSuccess: false,
	// 		// 	ArgEAPMessage: "",
	// 		// })
	// 	// case models.AuthResult_FAILURE:
	// 	// 	if ue.IdentityTypeUsedForRegistration == nasMessage.MobileIdentity5GSType5gGuti && ue.IdentityRequestSendTimes == 0 {
	// 	// 		ue.IdentityRequestSendTimes++
	// 	// 		gmm_message.SendIdentityRequest(ue.RanUe[accessType], accessType, nasMessage.MobileIdentity5GSTypeSuci)
	// 	// 		return nil
	// 	// 	} else {
	// 	// 		gmm_message.SendAuthenticationReject(ue.RanUe[accessType], "")
	// 	// 		return GmmFSM.SendEvent(ue.State[accessType], AuthFailEvent, fsm.ArgsType{
	// 	// 			ArgAmfUe:      ue,
	// 	// 			ArgAccessType: accessType,
	// 	// 		})
	// 	// 	}
	// 	}
	// case models.AuthType_EAP_AKA_PRIME:
	// 	response, problemDetails, err := consumer.SendEapAuthConfirmRequest(ue, *authenticationResponse.EAPMessage)
	// 	if err != nil {
	// 		return err
	// 	} else if problemDetails != nil {
	// 		logger.GmmLog.Debugf("EapAuthConfirm Error[Problem Detail: %+v]", problemDetails)
	// 		return nil
	// 	}

	// 	switch response.AuthResult {
	// 	case models.AuthResult_SUCCESS:
	// 		// ue.UnauthenticatedSupi = false
	// 		ue.Kseaf = response.KSeaf
	// 		ue.Supi = response.Supi
	// 		ue.DerivateKamf()
	// 		// TODO: select enc/int algorithm based on ue security capability & amf's policy,
	// 		// then generate KnasEnc, KnasInt
	// 	// 	return GmmFSM.SendEvent(ue.State[accessType], AuthSuccessEvent, fsm.ArgsType{
	// 	// 		ArgAmfUe:      ue,
	// 	// 		ArgAccessType: accessType,
	// 	// 		ArgEAPSuccess: true,
	// 	// 		ArgEAPMessage: response.EapPayload,
	// 	// 	})
	// 	// case models.AuthResult_FAILURE:
	// 	// 	if ue.IdentityTypeUsedForRegistration == nasMessage.MobileIdentity5GSType5gGuti && ue.IdentityRequestSendTimes == 0 {
	// 	// 		ue.IdentityRequestSendTimes++
	// 	// 		gmm_message.SendAuthenticationResult(ue.RanUe[accessType], false, response.EapPayload)
	// 	// 		gmm_message.SendIdentityRequest(ue.RanUe[accessType], accessType, nasMessage.MobileIdentity5GSTypeSuci)
	// 	// 		return nil
	// 	// 	} else {
	// 	// 		gmm_message.SendAuthenticationReject(ue.RanUe[accessType], response.EapPayload)
	// 	// 		return GmmFSM.SendEvent(ue.State[accessType], AuthFailEvent, fsm.ArgsType{
	// 	// 			ArgAmfUe:      ue,
	// 	// 			ArgAccessType: accessType,
	// 	// 		})
	// 	// 	}
	// 	// case models.AuthResult_ONGOING:
	// 	// 	ue.AuthenticationCtx.Var5gAuthData = response.EapPayload
	// 	// 	if _, exists := response.Links["eap-session"]; exists {
	// 	// 		ue.AuthenticationCtx.Links = response.Links
	// 	// 	}
	// 	// 	gmm_message.SendAuthenticationRequest(ue.RanUe[accessType])
	// 	}
	// }

	return nil
}

// Handle cleartext IEs of Registration Request, which cleattext IEs defined in TS 24.501 4.4.6
// TODO 
func HandleRegistrationRequest(ue *context.LbUe, anType models.AccessType, procedureCode int64,
	registrationRequest *nasMessage.RegistrationRequest) error {
	// // var guamiFromUeGuti models.Guami
	// // self := context.LB_Self()

	// if ue == nil {
	// 	return fmt.Errorf("AmfUe is nil")
	// }
	// if registrationRequest.UESecurityCapability != nil {
	// 	ue.UESecurityCapability = *registrationRequest.UESecurityCapability
	// } else {
	// 	// gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMProtocolErrorUnspecified, "")
	// 	return fmt.Errorf("UESecurityCapability is nil")
	// }
	// // ue.GmmLog.Info("Handle Registration Request")

	// // if ue.RanUe[anType] == nil {
	// // 	return fmt.Errorf("RanUe is nil")
	// // }

	// // ue.SetOnGoing(anType, &context.OnGoing{
	// // 	Procedure: context.OnGoingProcedureRegistration,
	// // })

	// // if ue.T3513 != nil {
	// // 	ue.T3513.Stop()
	// // 	ue.T3513 = nil // clear the timer
	// // }
	// // if ue.T3565 != nil {
	// // 	ue.T3565.Stop()
	// // 	ue.T3565 = nil // clear the timer
	// // }

	// // TS 24.501 8.2.6.21: if the UE is sending a REGISTRATION REQUEST message as an initial NAS message,
	// // the UE has a valid 5G NAS security context and the UE needs to send non-cleartext IEs
	// // TS 24.501 4.4.6: When the UE sends a REGISTRATION REQUEST or SERVICE REQUEST message that includes a NAS message
	// // container IE, the UE shall set the security header type of the initial NAS message to "integrity protected"
	// // if registrationRequest.NASMessageContainer != nil && !ue.MacFailed {
	// // 	contents := registrationRequest.NASMessageContainer.GetNASMessageContainerContents()

	// // 	// TS 24.501 4.4.6: When the UE sends a REGISTRATION REQUEST or SERVICE REQUEST message that includes a NAS
	// // 	// message container IE, the UE shall set the security header type of the initial NAS message to
	// // 	// "integrity protected"; then the AMF shall decipher the value part of the NAS message container IE
	// // 	err := security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.ULCount.Get(), security.Bearer3GPP,
	// // 		security.DirectionUplink, contents)
	// // 	if err != nil {
	// // 		ue.SecurityContextAvailable = false
	// // 	} else {
	// // 		m := nas.NewMessage()
	// // 		if err := m.GmmMessageDecode(&contents); err != nil {
	// // 			return err
	// // 		}

	// // 		messageType := m.GmmMessage.GmmHeader.GetMessageType()
	// // 		if messageType != nas.MsgTypeRegistrationRequest {
	// // 			return fmt.Errorf("The payload of NAS Message Container is not Registration Request")
	// // 		}
	// // 		// TS 24.501 4.4.6: The AMF shall consider the NAS message that is obtained from the NAS message container
	// // 		// IE as the initial NAS message that triggered the procedure
	// // 		registrationRequest = m.RegistrationRequest
	// // 	}
	// // }
	// // // TS 33.501 6.4.6 step 3: if the initial NAS message was protected but did not pass the integrity check
	// // ue.RetransmissionOfInitialNASMsg = ue.MacFailed

	// // ue.RegistrationRequest = registrationRequest
	// // ue.RegistrationType5GS = registrationRequest.NgksiAndRegistrationType5GS.GetRegistrationType5GS()
	// // switch ue.RegistrationType5GS {
	// // case nasMessage.RegistrationType5GSInitialRegistration:
	// // 	ue.GmmLog.Debugf("RegistrationType: Initial Registration")
	// // case nasMessage.RegistrationType5GSMobilityRegistrationUpdating:
	// // 	ue.GmmLog.Debugf("RegistrationType: Mobility Registration Updating")
	// // 	if ue.State[anType].Is(context.Deregistered) {
	// // 		gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMImplicitlyDeregistered, "")
	// // 		return fmt.Errorf("Mobility Registration Updating was sent when the UE state was Deregistered")
	// // 	}
	// // case nasMessage.RegistrationType5GSPeriodicRegistrationUpdating:
	// // 	ue.GmmLog.Debugf("RegistrationType: Periodic Registration Updating")
	// // 	if ue.State[anType].Is(context.Deregistered) {
	// // 		gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMImplicitlyDeregistered, "")
	// // 		return fmt.Errorf("Periodic Registration Updating was sent when the UE state was Deregistered")
	// // 	}
	// // case nasMessage.RegistrationType5GSEmergencyRegistration:
	// // 	return fmt.Errorf("Not Supportted RegistrationType: Emergency Registration")
	// // case nasMessage.RegistrationType5GSReserved:
	// // 	ue.RegistrationType5GS = nasMessage.RegistrationType5GSInitialRegistration
	// // 	ue.GmmLog.Debugf("RegistrationType: Reserved")
	// // default:
	// // 	ue.GmmLog.Debugf("RegistrationType: %v, chage state to InitialRegistration", ue.RegistrationType5GS)
	// // 	ue.RegistrationType5GS = nasMessage.RegistrationType5GSInitialRegistration
	// // }

	// // mobileIdentity5GSContents := registrationRequest.MobileIdentity5GS.GetMobileIdentity5GSContents()
	// // // ue.IdentityTypeUsedForRegistration = nasConvert.GetTypeOfIdentity(mobileIdentity5GSContents[0])
	// // switch ue.IdentityTypeUsedForRegistration { // get type of identity
	// // case nasMessage.MobileIdentity5GSTypeNoIdentity:
	// // 	ue.GmmLog.Debugf("No Identity")
	// // case nasMessage.MobileIdentity5GSTypeSuci:
	// // 	var plmnId string
	// // 	ue.Suci, plmnId = nasConvert.SuciToString(mobileIdentity5GSContents)
	// // 	ue.PlmnId = util.PlmnIdStringToModels(plmnId)
	// // 	ue.GmmLog.Debugf("SUCI: %s", ue.Suci)
	// // case nasMessage.MobileIdentity5GSType5gGuti:
	// // 	guamiFromUeGutiTmp, guti := nasConvert.GutiToString(mobileIdentity5GSContents)
	// // 	guamiFromUeGuti = guamiFromUeGutiTmp
	// // 	ue.Guti = guti
	// // 	ue.GmmLog.Debugf("GUTI: %s", guti)

	// // 	servedGuami := amfSelf.ServedGuamiList[0]
	// // 	if reflect.DeepEqual(guamiFromUeGuti, servedGuami) {
	// // 		ue.ServingAmfChanged = false
	// // 	} else {
	// // 		ue.GmmLog.Debugf("Serving AMF has changed")
	// // 		ue.ServingAmfChanged = true
	// // 	}
	// // case nasMessage.MobileIdentity5GSTypeImei:
	// // 	imei := nasConvert.PeiToString(mobileIdentity5GSContents)
	// // 	ue.Pei = imei
	// // 	ue.GmmLog.Debugf("PEI: %s", imei)
	// // case nasMessage.MobileIdentity5GSTypeImeisv:
	// // 	imeisv := nasConvert.PeiToString(mobileIdentity5GSContents)
	// // 	ue.Pei = imeisv
	// // 	ue.GmmLog.Debugf("PEI: %s", imeisv)
	// // }

	// // NgKsi: TS 24.501 9.11.3.32
	// // switch registrationRequest.NgksiAndRegistrationType5GS.GetTSC() {
	// // case nasMessage.TypeOfSecurityContextFlagNative:
	// // 	ue.NgKsi.Tsc = models.ScType_NATIVE
	// // case nasMessage.TypeOfSecurityContextFlagMapped:
	// // 	ue.NgKsi.Tsc = models.ScType_MAPPED
	// // }
	// // ue.NgKsi.Ksi = int32(registrationRequest.NgksiAndRegistrationType5GS.GetNasKeySetIdentifiler())
	// // if ue.NgKsi.Tsc == models.ScType_NATIVE && ue.NgKsi.Ksi != 7 {
	// // } else {
	// // 	ue.NgKsi.Tsc = models.ScType_NATIVE
	// // 	ue.NgKsi.Ksi = 0
	// // }

	// // // Copy UserLocation from ranUe
	// // ue.Location = ue.RanUe[anType].Location
	// // ue.Tai = ue.RanUe[anType].Tai

	// // // Check TAI
	// // if !context.InTaiList(ue.Tai, amfSelf.SupportTaiLists) {
	// // 	gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMTrackingAreaNotAllowed, "")
	// // 	return fmt.Errorf("Registration Reject[Tracking area not allowed]")
	// // }

	

	// // TODO (TS 23.502 4.2.2.2 step 4): if UE's 5g-GUTI is included & serving AMF has changed
	// // since last registration procedure, new AMF may invoke Namf_Communication_UEContextTransfer
	// // to old AMF, including the complete registration request nas msg, to request UE's SUPI & UE Context
	// // if ue.ServingAmfChanged {
	// // 	var transferReason models.TransferReason
	// // 	switch ue.RegistrationType5GS {
	// // 	case nasMessage.RegistrationType5GSInitialRegistration:
	// // 		transferReason = models.TransferReason_INIT_REG
	// // 	case nasMessage.RegistrationType5GSMobilityRegistrationUpdating:
	// // 		fallthrough
	// // 	case nasMessage.RegistrationType5GSPeriodicRegistrationUpdating:
	// // 		transferReason = models.TransferReason_MOBI_REG
	// // 	}

	// // 	searchOpt := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
	// // 		Guami: optional.NewInterface(openapi.MarshToJsonString(guamiFromUeGuti)),
	// // 	}
	// // 	err := consumer.SearchAmfCommunicationInstance(ue, amfSelf.NrfUri, models.NfType_AMF, models.NfType_AMF, &searchOpt)
	// // 	if err != nil {
	// // 		ue.GmmLog.Errorf("[GMM] %+v", err)
	// // 		gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMUEIdentityCannotBeDerivedByTheNetwork, "")
	// // 		return err
	// // 	}

	// // 	ueContextTransferRspData, problemDetails, err := consumer.UEContextTransferRequest(ue, anType, transferReason)
	// // 	if problemDetails != nil {
	// // 		if problemDetails.Cause == "INTEGRITY_CHECK_FAIL" || problemDetails.Cause == "CONTEXT_NOT_FOUND" {
	// // 			ue.GmmLog.Warnf("Can not retrieve UE Context from old AMF[Cause: %s]", problemDetails.Cause)
	// // 		} else {
	// // 			ue.GmmLog.Warnf("UE Context Transfer Request Failed Problem[%+v]", problemDetails)
	// // 		}
	// // 		ue.SecurityContextAvailable = false // need to start authentication procedure later
	// // 	} else if err != nil {
	// // 		ue.GmmLog.Errorf("UE Context Transfer Request Error[%+v]", err)
	// // 		gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMUEIdentityCannotBeDerivedByTheNetwork, "")
	// // 	} else {
	// // 		ue.CopyDataFromUeContextModel(*ueContextTransferRspData.UeContext)
	// // 	}
	// // }
	return nil
}