package service

import (
	"syscall"

	"github.com/ishidawataru/sctp"
	"github.com/LuckyG0ldfish/balancer/context"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/ngap"
)

// Listens to incoming connections and servers them (GNBs)
func listenAndServeGNBs(addr *sctp.SCTPAddr, handler NGAPHandler) {
	if listener, err := sctpConfig.Listen("sctp", addr); err != nil {
		logger.NgapLog.Errorf("Failed to listen: %+v", err)
		return
	} else {
		sctpListener = listener
	}
	logger.NgapLog.Infof("Listen on %s", sctpListener.Addr())

	for {
		newConn, err := sctpListener.AcceptSCTP()
		if err != nil {
			switch err {
			case syscall.EINTR, syscall.EAGAIN:
				logger.NgapLog.Debugf("AcceptSCTP: %+v", err)
			default:
				logger.NgapLog.Errorf("Failed to accept: %+v", err)
			}
			continue
		}

		var info *sctp.SndRcvInfo
		if infoTmp, err := newConn.GetDefaultSentParam(); err != nil {
			logger.NgapLog.Errorf("Get default sent param error: %+v, accept failed", err)
			if err = newConn.Close(); err != nil {
				logger.NgapLog.Errorf("Close error: %+v", err)
			}
			continue
		} else {
			info = infoTmp
			logger.NgapLog.Debugf("Get default sent param[value: %+v]", info)
		}

		info.PPID = ngap.PPID
		if err := newConn.SetDefaultSentParam(info); err != nil {
			logger.NgapLog.Errorf("Set default sent param error: %+v, accept failed", err)
			if err = newConn.Close(); err != nil {
				logger.NgapLog.Errorf("Close error: %+v", err)
			}
			continue
		} else {
			logger.NgapLog.Debugf("Set default sent param[value: %+v]", info)
		}

		events := sctp.SCTP_EVENT_DATA_IO | sctp.SCTP_EVENT_SHUTDOWN | sctp.SCTP_EVENT_ASSOCIATION
		if err := newConn.SubscribeEvents(events); err != nil {
			logger.NgapLog.Errorf("Failed to accept: %+v", err)
			if err = newConn.Close(); err != nil {
				logger.NgapLog.Errorf("Close error: %+v", err)
			}
			continue
		} else {
			logger.NgapLog.Debugln("Subscribe SCTP event[DATA_IO, SHUTDOWN_EVENT, ASSOCIATION_CHANGE]")
		}

		// if err := newConn.SetReadBuffer(int(readBufSize)); err != nil {
		// 	logger.NgapLog.Errorf("Set read buffer error: %+v, accept failed", err)
		// 	if err = newConn.Close(); err != nil {
		// 		logger.NgapLog.Errorf("Close error: %+v", err)
		// 	}
		// 	continue
		// } else {
		// 	logger.NgapLog.Infof("Set read buffer to %d bytes", readBufSize)
		// }

		if err := newConn.SetReadTimeout(readTimeout); err != nil {
			logger.NgapLog.Errorf("Set read timeout error: %+v, accept failed", err)
			if err = newConn.Close(); err != nil {
				logger.NgapLog.Errorf("Close error: %+v", err)
			}
			continue
		} else {
			logger.NgapLog.Infof("Set read timeout: %+v", readTimeout)
		}
		
		value, err := newConn.GetNoDelay()
		if err != nil {
			logger.AppLog.Errorf("Error getting SCTP_NODELAY: %v", err)
		}
		logger.AppLog.Infof("[BEFORE] Boolean value of the SCTP_NODELAY flag: %v", value)
	
	
		err = newConn.SetNoDelay(1)
		if err != nil {
			logger.AppLog.Errorf("Error setting SCTP_NODELAY: %v", err)
		}
	
		value, err = newConn.GetNoDelay()
		if err != nil {
			logger.AppLog.Errorf("Error getting SCTP_NODELAY: %v", err)
		}
		logger.AppLog.Infof("[AFTER] Boolean value of the SCTP_NODELAY flag: %v", value)

		logger.NgapLog.Infof("[LB] SCTP Accept from: %s", newConn.RemoteAddr().String())
		
		

		// add connection as new GNBConn 
		lb := context.LB_Self()
		gnb := context.CreateAndAddNewGnbToLB(newConn)
		logger.ContextLog.Tracef("LB_GNB created")
		connections.Store(gnb.LbConn, *gnb.LbConn)
		lb.LbRanPool.Store(gnb.GnbID, gnb)
		
		// Metrics
		if lb.MetricsLevel > 0 {
			mGNB := context.NewMetricsGNB(gnb.GnbID)
			lb.MetricsGNBs.Store(gnb.GnbID, mGNB)
			test, ok := lb.MetricsGNBs.Load(gnb.GnbID)
			if !ok {
				logger.ContextLog.Errorln("failed mgnb add #1")
			} 
			_, ok = test.(*context.MetricsGNB)
			if !ok {
				logger.ContextLog.Errorln("failed mgnb type cast #1")
			}
			_, ok = test.(context.MetricsGNB)
			if !ok {
				logger.ContextLog.Errorln("failed mgnb type cast #1.5")
			}
		} 
		
		
		go handleConnection(gnb.LbConn, readBufSize, handler)
	}
}