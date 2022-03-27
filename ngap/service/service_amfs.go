package service

import (
	"strconv"
	"time"

	// "git.cs.nctu.edu.tw/calee/sctp"
	"github.com/ishidawataru/sctp"
	"github.com/LuckyG0ldfish/balancer/context"
	"github.com/LuckyG0ldfish/balancer/logger"
	ngap_message "github.com/LuckyG0ldfish/balancer/ngap/message"
	"github.com/free5gc/ngap"
)

// Continuesly checks whether new AMFs have to be added
func InitAmfs(ngapHandler NGAPHandler) {
	self := context.LB_Self()
	for {
		if !self.Running { return }
		if self.NewAmf {
			// starting connections to each individual AMF 
			for i := 0; i < len(self.NewRegistAmfIpList); i++  {
				ip := self.NewRegistAmfIpList[i]
				logger.NgapLog.Debugf("connecting to: " + ip + ":" + strconv.Itoa(amfPort))
				go CreateAndStartAmf(ip, amfPort, ngapHandler, context.TypeIdRegist)
			}
			if self.DifferentAmfTypes == 3 {
				for i := 0; i < len(self.NewRegularAmfIpList); i++  {
					ip := self.NewRegistAmfIpList[i]
					logger.NgapLog.Debugf("connecting to: " + ip + ":" + strconv.Itoa(amfPort))
					go CreateAndStartAmf(ip, amfPort, ngapHandler, context.TypeIdRegular)
				}
			}
			if self.DifferentAmfTypes >= 2 {
				for i := 0; i < len(self.NewDeregistAmfIpList); i++  {
					ip := self.NewRegistAmfIpList[i]
					logger.NgapLog.Debugf("connecting to: " + ip + ":" + strconv.Itoa(amfPort))
					go CreateAndStartAmf(ip, amfPort, ngapHandler, context.TypeIdDeregist)
				}
			}
			if !self.ContinuesAmfRegistration {
				logger.NgapLog.Tracef("No further AMFs to accept")
				return 
			}
			// Resets to accept more 
			self.NewRegistAmfIpList = []string{}
			self.NewAmf = false
			logger.NgapLog.Tracef("Waiting for new Amfs")
		}
		time.Sleep(2 * time.Second)
	}
}

// Creates AMF and initializes the starting process
func CreateAndStartAmf(amfIP string, amfPort int, ngapHandler NGAPHandler, amfType int) {
	self := context.LB_Self()
	amf := context.CreateAndAddAmfToLB(amfType)
	addr, err := context.GenSCTPAddr(self.LbIP, int(numberGen.NextNumber()))
	if err != nil {
		logger.NgapLog.Errorln("LB-SCTP-Address couldn't be build")
		return 
	}
	StartAmf(amf, addr, amfIP, amfPort, ngapHandler)
}

// Initializes LB to AMF communication and starts handling the connection 
func StartAmf(amf *context.LbAmf, lbaddr *sctp.SCTPAddr, amfIP string, amfPort int, handler NGAPHandler) {
	self := context.LB_Self()
	logger.NgapLog.Debugf("Connecting to amf")
	for {
		conn, err := ConnectToAmf(lbaddr, amfIP, amfPort)
		if err == nil {
			amf.LbConn.Conn = conn
			logger.NgapLog.Debugf("Connected to amf")
			if amf.AmfTypeIdent == context.TypeIdRegist {
				self.Next_Regist_Amf = amf
				self.LbRegistAmfPool.Store(amf.LbConn.Conn, amf)
			} else if amf.AmfTypeIdent == context.TypeIdDeregist {
				self.Next_Deregist_Amf = amf
				self.LbDeregistAmfPool.Store(amf.LbConn.Conn, amf)
			} else {
				self.Next_Regular_Amf = amf
				self.LbRegularAmfPool.Store(amf.LbConn.Conn, amf)
			}
			connections.Store(amf.LbConn, *amf.LbConn)
			ngap_message.SendNGSetupRequest(amf.LbConn)
			handleConnection(amf.LbConn, readBufSize, handler)
			return 
		}
		time.Sleep(2 * time.Second)
	}
}

// Establishes a SCTP connection to an AMF 
func ConnectToAmf(lbaddr *sctp.SCTPAddr, amfIP string, amfPort int) (*sctp.SCTPConn, error) {
	amfAddr, _ := context.GenSCTPAddr(amfIP, amfPort)
	conn, err := sctp.DialSCTP("sctp", lbaddr, amfAddr)
	if err != nil {
		logger.NgapLog.Warnf("Connection Failed: Dial failed")
		return nil, err
	}
	info, err := conn.GetDefaultSentParam()
	if err != nil {
		logger.NgapLog.Warnf("Connection Failed: failed to get DefaultSentParam")
		return nil, err
	}
	info.PPID = ngap.PPID
	err = conn.SetDefaultSentParam(info)
	if err != nil {
		logger.NgapLog.Warnf("Connection Failed: failed to set DefaultSentParam")
		return nil, err
	}
	err = conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
	if err != nil {
		logger.NgapLog.Warnf("Connection Failed: failed to SubscribeEvents")
		return nil, err
	}

	// Change the value of the SCTP_NODELAY flag to disable the Nagle Algorithm and send packets as soon as they're available.
	value, err := conn.GetNoDelay()
	if err != nil {
		logger.AppLog.Errorf("Error getting SCTP_NODELAY: %v", err)
	}
	logger.AppLog.Infof("[BEFORE] Boolean value of the SCTP_NODELAY flag: %v", value)


	err = conn.SetNoDelay(1)
	if err != nil {
		logger.AppLog.Errorf("Error setting SCTP_NODELAY: %v", err)
	}

	value, err = conn.GetNoDelay()
	if err != nil {
		logger.AppLog.Errorf("Error getting SCTP_NODELAY: %v", err)
	}
	logger.AppLog.Infof("[AFTER] Boolean value of the SCTP_NODELAY flag: %v", value)

	logger.NgapLog.Infoln("connected to amf: IP " + amfIP + " Port: " + strconv.Itoa(amfPort))
	return conn, nil
}