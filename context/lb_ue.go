package context

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"regexp"
	// "time"

	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/UeauCommon"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/openapi/models"
	"github.com/sirupsen/logrus"
)

// Type, that stores all relevant information of UEs
type LB_UE struct{
	UeStateIdent 	int			// Identifies the state of the UE 

	GNB_UE_ID 		int64		// ID given to the UE by GNB/RAN
	LB_UE_ID 		int64		// ID given to the UE by LB
	AMF_UE_ID 		int64		// ID given to the UE by AMF
	
	GnbID			int64		// LB-internal ID of GNB that issued the UE 
	GnbPointer 		*LB_GNB		// ponter to the connected GNB

	AmfID		 	int64		// LB-internal ID of AMF that processes the UE  
	AmfPointer		*LB_AMF		// ponter to the connected AMF

	UplinkFlag		bool 		// set true when Uplink-NAS-RegistrationComplete is done
	ResponseFlag	bool		// set true when InitialContextSetupResponse is done
	DeregFlag		bool		// set true when Uplink-NAS-Deregistration-Accept is done 

	/* nas decrypt */ // TODO
	RRCECause 		string
	ULCount			security.Count 	//TODO amf_ue L728 | gmm HandleRegist HandleServiceRequest (only get())
	DLCount			security.Count	// TODO set in CopyDataFromUeContextModel | .AddOne() in nas Encode()
	Kamf            string	
	Kgnb            []uint8			// 32 byte
	Kn3iwf          []uint8   		// 32 byte
	CipheringAlg    uint8			 
	IntegrityAlg    uint8			
	KnasInt         [16]uint8 		// 16 byte 
	KnasEnc         [16]uint8 		// 16 byte 
	Supi            string			
	ABBA            []uint8			
	Kseaf           string		
	UESecurityCapability     nasType.UESecurityCapability // for security command
	// TODO set in NGAP HandlePathSwitchRequest
	MacFailed       bool      // set to true if the integrity check of current NAS message is failed


	/* logger */
	Log 			*logrus.Entry
}

// Creates, initializes and returns a new *LbUe
func NewUE() (*LB_UE){
	var ue LB_UE
	ue.UeStateIdent = TypeIdRegist
	ue.Log = logger.UELog
	// set ABBA value as described at TS 33.501 Annex A.7.1
	ue.ABBA = []uint8{0x00, 0x00} // set in GMM AuthenticationProcedure + AuthenticationFailure
	ue.UplinkFlag = false 
	ue.ResponseFlag = false 
	ue.DeregFlag = false 
	return &ue
}

// Removes LbUe from AMF and RAN Context withing LB  
func (ue *LB_UE) RemoveUeEntirely() {
	ue.RemoveUeFromAMF()
	ue.RemoveUeFromGNB()
	ue = nil 
}

// Removes LbUe from AMF Context withing LB 
func (ue *LB_UE) RemoveUeFromAMF() {
	if ue.AmfPointer != nil {
		ue.AmfPointer.NumberOfConnectedUEs -= 1
		ue.AmfPointer.Ues.Delete(ue.LB_UE_ID) // sync.Map key here is the LB internal UE-ID 
		ue.AmfPointer.Log.Debugf("LB_UE_ID %d context removed from AMF", ue.LB_UE_ID)
	}
}

// Removes LbUe from RAN Context withing LB 
func (ue *LB_UE) RemoveUeFromGNB() {
	if ue.GnbPointer != nil {
		ue.GnbPointer.Ues.Delete(ue.GNB_UE_ID) // sync.Map key here is the RAN UE-ID
		ue.GnbPointer.Log.Debugf("LB_UE_ID %d context removed from GNB", ue.LB_UE_ID)
	}
}

// Sets UEs values and adds it to the Amfs UE-Map
func (ue *LB_UE) AddUeToAmf(next *LB_AMF) {
	ue.AmfID = next.AmfID
	ue.AmfPointer = next
	next.Ues.Store(ue.LB_UE_ID, ue)
	next.NumberOfConnectedUEs += 1
	logger.ContextLog.Tracef("GNB_UE_ID: %d added to AMF %d", ue.LB_UE_ID, next.AmfID)
}


func (ue *LB_UE) RegistrationComplete() {
	if ue.UeStateIdent == TypeIdRegist && ue.ResponseFlag && ue.UplinkFlag {
		self := LB_Self()
		if self.DifferentAmfTypes == 3 {
			ue.UeStateIdent = TypeIdRegular
			next := self.Next_Regular_Amf
			ue.RemoveUeFromAMF()
			ue.AddUeToAmf(next)
			self.SelectNextRegularAmf()
			return 
		} else if self.DifferentAmfTypes == 2 {
			next := self.Next_Deregist_Amf
			ue.RemoveUeFromAMF()
			ue.AddUeToAmf(next)
			self.SelectNextDeregistAmf()
		}			
	}
}	

// TODO
// Kamf Derivation function defined in TS 33.501 Annex A.7
// gmm handler HandleAuthenticationResponse L1943 + 1978
func (ue *LB_UE) DerivateKamf() {
	supiRegexp, err := regexp.Compile("(?:imsi|supi)-([0-9]{5,15})")
	if err != nil {
		logger.ContextLog.Error(err)
		return
	}
	// could probably be solved differently 
	groups := supiRegexp.FindStringSubmatch(ue.Supi)
	if groups == nil {
		logger.NASLog.Errorln("supi is not correct")
		return
	}

	P0 := []byte(groups[1])
	L0 := UeauCommon.KDFLen(P0)
	P1 := ue.ABBA
	L1 := UeauCommon.KDFLen(P1)

	KseafDecode, err := hex.DecodeString(ue.Kseaf)
	if err != nil {
		logger.ContextLog.Error(err)
		return
	}
	KamfBytes := UeauCommon.GetKDFValue(KseafDecode, UeauCommon.FC_FOR_KAMF_DERIVATION, P0, L0, P1, L1)
	ue.Kamf = hex.EncodeToString(KamfBytes)
}

// TODO
// Access Network key Derivation function defined in TS 33.501 Annex A.9
func (ue *LB_UE) DerivateAnKey(anType models.AccessType) {
	accessType := security.AccessType3GPP // Defalut 3gpp
	P0 := make([]byte, 4)
	binary.BigEndian.PutUint32(P0, ue.ULCount.Get())
	L0 := UeauCommon.KDFLen(P0)
	if anType == models.AccessType_NON_3_GPP_ACCESS {
		accessType = security.AccessTypeNon3GPP
	}
	P1 := []byte{accessType}
	L1 := UeauCommon.KDFLen(P1)

	KamfBytes, err := hex.DecodeString(ue.Kamf) //TODO Kamf
	if err != nil {
		logger.ContextLog.Error(err)
		return
	}
	key := UeauCommon.GetKDFValue(KamfBytes, UeauCommon.FC_FOR_KGNB_KN3IWF_DERIVATION, P0, L0, P1, L1)
	switch accessType {
	case security.AccessType3GPP:
		ue.Kgnb = key
	case security.AccessTypeNon3GPP:
		ue.Kn3iwf = key
	}
}

// TODO
// Algorithm key Derivation function defined in TS 33.501 Annex A.9
func (ue *LB_UE) DerivateAlgKey() {
	// Security Key
	P0 := []byte{security.NNASEncAlg}
	L0 := UeauCommon.KDFLen(P0)
	P1 := []byte{ue.CipheringAlg}
	L1 := UeauCommon.KDFLen(P1)

	KamfBytes, err := hex.DecodeString(ue.Kamf)
	if err != nil {
		logger.ContextLog.Error(err)
		return
	}
	kenc := UeauCommon.GetKDFValue(KamfBytes, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	copy(ue.KnasEnc[:], kenc[16:32])

	// Integrity Key
	P0 = []byte{security.NNASIntAlg}
	L0 = UeauCommon.KDFLen(P0)
	P1 = []byte{ue.IntegrityAlg}
	L1 = UeauCommon.KDFLen(P1)

	kint := UeauCommon.GetKDFValue(KamfBytes, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	copy(ue.KnasInt[:], kint[16:32])
}

// TODO
func (ue *LB_UE) SelectSecurityAlg(intOrder, encOrder []uint8) {
	ue.CipheringAlg = security.AlgCiphering128NEA0
	ue.IntegrityAlg = security.AlgIntegrity128NIA0

	ueSupported := uint8(0)
	for _, intAlg := range intOrder {
		switch intAlg {
		case security.AlgIntegrity128NIA0:
			ueSupported = ue.UESecurityCapability.GetIA0_5G()
		case security.AlgIntegrity128NIA1:
			ueSupported = ue.UESecurityCapability.GetIA1_128_5G()
		case security.AlgIntegrity128NIA2:
			ueSupported = ue.UESecurityCapability.GetIA2_128_5G()
		case security.AlgIntegrity128NIA3:
			ueSupported = ue.UESecurityCapability.GetIA3_128_5G()
		}
		if ueSupported == 1 {
			ue.IntegrityAlg = intAlg
			break
		}
	}

	ueSupported = uint8(0)
	for _, encAlg := range encOrder {
		switch encAlg {
		case security.AlgCiphering128NEA0:
			ueSupported = ue.UESecurityCapability.GetEA0_5G()
		case security.AlgCiphering128NEA1:
			ueSupported = ue.UESecurityCapability.GetEA1_128_5G()
		case security.AlgCiphering128NEA2:
			ueSupported = ue.UESecurityCapability.GetEA2_128_5G()
		case security.AlgCiphering128NEA3:
			ueSupported = ue.UESecurityCapability.GetEA3_128_5G()
		}
		if ueSupported == 1 {
			ue.CipheringAlg = encAlg
			break
		}
	}
}

// TODO
func (ue *LB_UE) CopyDataFromUeContextModel(ueContext models.UeContext) {
	if ueContext.Supi != "" {
		ue.Supi = ueContext.Supi
	}

	if len(ueContext.MmContextList) > 0 {
		for _, mmContext := range ueContext.MmContextList {
			if mmContext.AccessType == models.AccessType__3_GPP_ACCESS {
				if nasSecurityMode := mmContext.NasSecurityMode; nasSecurityMode != nil {
					switch nasSecurityMode.IntegrityAlgorithm {
					case models.IntegrityAlgorithm_NIA0:
						ue.IntegrityAlg = security.AlgIntegrity128NIA0
					case models.IntegrityAlgorithm_NIA1:
						ue.IntegrityAlg = security.AlgIntegrity128NIA1
					case models.IntegrityAlgorithm_NIA2:
						ue.IntegrityAlg = security.AlgIntegrity128NIA2
					case models.IntegrityAlgorithm_NIA3:
						ue.IntegrityAlg = security.AlgIntegrity128NIA3
					}

					switch nasSecurityMode.CipheringAlgorithm {
					case models.CipheringAlgorithm_NEA0:
						ue.CipheringAlg = security.AlgCiphering128NEA0
					case models.CipheringAlgorithm_NEA1:
						ue.CipheringAlg = security.AlgCiphering128NEA1
					case models.CipheringAlgorithm_NEA2:
						ue.CipheringAlg = security.AlgCiphering128NEA2
					case models.CipheringAlgorithm_NEA3:
						ue.CipheringAlg = security.AlgCiphering128NEA3
					}

					if mmContext.NasDownlinkCount != 0 {
						overflow := uint16((uint32(mmContext.NasDownlinkCount) & 0x00ffff00) >> 8)
						sqn := uint8(uint32(mmContext.NasDownlinkCount & 0x000000ff))
						ue.DLCount.Set(overflow, sqn)
					}

					if mmContext.NasUplinkCount != 0 {
						overflow := uint16((uint32(mmContext.NasUplinkCount) & 0x00ffff00) >> 8)
						sqn := uint8(uint32(mmContext.NasUplinkCount & 0x000000ff))
						ue.ULCount.Set(overflow, sqn)
					}

					// TS 29.518 Table 6.1.6.3.2.1
					if mmContext.UeSecurityCapability != "" {
						// ue.SecurityCapabilities
						buf, err := base64.StdEncoding.DecodeString(mmContext.UeSecurityCapability)
						if err != nil {
							logger.ContextLog.Error(err)
							return
						}
						ue.UESecurityCapability.Buffer = buf
						ue.UESecurityCapability.SetLen(uint8(len(buf)))
					}
				}
			}
		}
	}
}