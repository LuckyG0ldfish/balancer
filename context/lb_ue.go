package context

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"regexp"
	"time"

	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/UeauCommon"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/openapi/models"
	"github.com/sirupsen/logrus"
)

// Type, that stores all relevant information of UEs
type LbUe struct{
	UeStateIdent 	int			// Identifies the state of the UE 

	UeRanID 		int64		// ID given to the UE by GNB/RAN
	UeLbID 			int64		// ID given to the UE by LB
	UeAmfID 		int64		// ID given to the UE by AMF
	
	RanID			int64		// LB-internal ID of GNB that issued the UE 
	RanPointer 		*LbGnb

	AmfID		 	int64		// LB-internal ID of AMF that processes the UE  
	AmfPointer		*LbAmf

	UplinkFlag		bool 		// set true when Uplink-NAS-RegistrationComplete is done
	ResponseFlag	bool		// set true when InitialContextSetupResponse is done
	// DeregFlag		bool		// set true when Uplink-NAS-RegistrationComplete 

	/* nas decrypt */
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
func NewUE() (*LbUe){
	var ue LbUe
	ue.UeStateIdent = TypeIdRegist
	ue.Log = logger.UELog
	// set ABBA value as described at TS 33.501 Annex A.7.1
	ue.ABBA = []uint8{0x00, 0x00} // set in GMM AuthenticationProcedure + AuthenticationFailure
	ue.UplinkFlag = false 
	ue.ResponseFlag = false 
	return &ue
}

// Removes LbUe from AMF and RAN Context withing LB  
func (ue *LbUe) RemoveUeEntirely() {
	time.Sleep(10 * time.Millisecond)
	ue.RemoveUeFromAMF()
	ue.RemoveUeFromGNB()
}

// Removes LbUe from AMF Context withing LB 
func (ue *LbUe) RemoveUeFromAMF() {
	if ue.AmfPointer != nil {
		ue.AmfPointer.Ues.Delete(ue.UeLbID) // sync.Map key here is the LB internal UE-ID 
		ue.AmfPointer.Log.Traceln("UE context removed from AMF")
		ue.AmfPointer = nil 
		ue.AmfID = 0 
		ue.UeAmfID = 0 
	}
}

// Removes LbUe from RAN Context withing LB 
func (ue *LbUe) RemoveUeFromGNB() {
	if ue.RanPointer != nil {
		ue.RanPointer.Ues.Delete(ue.UeRanID) // sync.Map key here is the RAN UE-ID
		ue.RanPointer.Log.Traceln("UE context removed from GNB")
		ue.RanPointer = nil 
		ue.RanID = 0 
	}
}

// Sets UEs values and adds it to the Amfs UE-Map
func (ue *LbUe) AddUeToAmf(next *LbAmf) {
	ue.AmfID = next.AmfID
	ue.AmfPointer = next
	next.Ues.Store(ue.UeLbID, ue)
}

func (ue *LbUe) RegistrationComplete() {
	if ue.UeStateIdent == TypeIdRegist && ue.ResponseFlag && ue.UplinkFlag {
		ue.UeStateIdent = TypeIdRegular
		// next := self.Next_Regular_Amf
		// ue.RemoveUeFromAMF()
		// ue.AddUeToAmf(next)
		// go context.SelectNextRegularAmf()			
	}
}	

// TODO
// Kamf Derivation function defined in TS 33.501 Annex A.7
// gmm handler HandleAuthenticationResponse L1943 + 1978
func (ue *LbUe) DerivateKamf() {
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
func (ue *LbUe) DerivateAnKey(anType models.AccessType) {
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
func (ue *LbUe) DerivateAlgKey() {
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
func (ue *LbUe) SelectSecurityAlg(intOrder, encOrder []uint8) {
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
func (ue *LbUe) CopyDataFromUeContextModel(ueContext models.UeContext) {
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