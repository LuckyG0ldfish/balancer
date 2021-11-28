package service

import (
	"encoding/hex"
	"io"
	"strconv"
	"sync"
	"syscall"
	"time"

	"git.cs.nctu.edu.tw/calee/sctp"

	//"github.com/LuckyG0ldfish/balancer/ngap"
	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/ngap"

	"github.com/LuckyG0ldfish/balancer/context"
	ngap_message "github.com/LuckyG0ldfish/balancer/ngap/message"
)

type NGAPHandler struct {
	HandleMessage      func(lbConn *context.LBConn, msg []byte)
	HandleNotification func(conn *sctp.SCTPConn, notification sctp.Notification)
}

const readBufSize uint32 = 8192

// set default read timeout to 2 seconds
var readTimeout syscall.Timeval = syscall.Timeval{Sec: 2, Usec: 0}

var (
	sctpListener 	*sctp.SCTPListener
	connections  	sync.Map
	ID 				int64
)

var sctpConfig sctp.SocketConfig = sctp.SocketConfig{
	InitMsg:   sctp.InitMsg{NumOstreams: 3, MaxInstreams: 5, MaxAttempts: 2, MaxInitTimeout: 2},
	RtoInfo:   &sctp.RtoInfo{SrtoAssocID: 0, SrtoInitial: 500, SrtoMax: 1500, StroMin: 100},
	AssocInfo: &sctp.AssocInfo{AsocMaxRxt: 4},
}

// Starting listen and serving of GNBs 
func Run(addr *sctp.SCTPAddr, handler NGAPHandler) {
	go listenAndServeGNBs(addr, handler)
}

// Listens to incoming connections and servers them 
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

		if err := newConn.SetReadBuffer(int(readBufSize)); err != nil {
			logger.NgapLog.Errorf("Set read buffer error: %+v, accept failed", err)
			if err = newConn.Close(); err != nil {
				logger.NgapLog.Errorf("Close error: %+v", err)
			}
			continue
		} else {
			logger.NgapLog.Debugf("Set read buffer to %d bytes", readBufSize)
		}

		if err := newConn.SetReadTimeout(readTimeout); err != nil {
			logger.NgapLog.Errorf("Set read timeout error: %+v, accept failed", err)
			if err = newConn.Close(); err != nil {
				logger.NgapLog.Errorf("Close error: %+v", err)
			}
			continue
		} else {
			logger.NgapLog.Debugf("Set read timeout: %+v", readTimeout)
		}

		logger.NgapLog.Infof("[LB] SCTP Accept from: %s", newConn.RemoteAddr().String())
		connections.Store(newConn, newConn)
		
		// add connection as new GNBConn 
		lbSelf := context.LB_Self()
		ran := lbSelf.AddNewGnbToLB(newConn)
		logger.ContextLog.Tracef("LB_GNB created")
		go handleConnection(ran.LbConn, readBufSize, handler)
		ran.LbConn.ID = ID 
		connections.Store(ID, ran.LbConn)	
		ID ++ 
	}
}

// Closes all connections 
func Stop() {
	logger.NgapLog.Infof("Close SCTP server...")
	if err := sctpListener.Close(); err != nil {
		logger.NgapLog.Error(err)
		logger.NgapLog.Infof("SCTP server may not close normally.")
	}

	connections.Range(func(key, value interface{}) bool {
		lbConn := value.(*context.LBConn)
		if err := lbConn.Conn.Close(); err != nil {
			logger.NgapLog.Error(err)
		}
		return true
	})

	logger.NgapLog.Infof("SCTP server closed")
}

// Handling all the the LBs open connections (AMFs + GNBs)
func handleConnection(lbConn *context.LBConn, bufsize uint32, handler NGAPHandler) {// conn *sctp.SCTPConn, bufsize uint32, handler NGAPHandler) {
	logger.NgapLog.Tracef("Waiting for message")
	for {
		buf := make([]byte, bufsize)

		n, info, notification, err := lbConn.Conn.SCTPRead(buf)
		
		if err != nil {
			switch err {
			case io.EOF, io.ErrUnexpectedEOF:
				logger.NgapLog.Debugln("Read EOF from client")
				return
			case syscall.EAGAIN:
				logger.NgapLog.Debugln("SCTP read timeout")
				continue
			case syscall.EINTR:
				logger.NgapLog.Debugf("SCTPRead: %+v", err)
				continue
			default:
				logger.NgapLog.Errorf("Handle connection[addr: %+v] error: %+v", lbConn.Conn.RemoteAddr(), err)
				return
			}
		}
		if lbConn.TypeID == context.TypeIdAMFConn {
			logger.NgapLog.Debugf("AMF message recieved")
		} else if lbConn.TypeID == context.TypeIdGNBConn {
			logger.NgapLog.Debugf("RAN message recieved")
		} else {
			logger.NgapLog.Errorf("unidientified message recieved") // TODO 
			break 
		}
		logger.NgapLog.Tracef("length: " + strconv.Itoa(n))
		if notification != nil {
			if handler.HandleNotification != nil {
				handler.HandleNotification(lbConn.Conn, notification)
			} else {
				logger.NgapLog.Warnf("Received sctp notification[type 0x%x] but not handled", notification.Type())
			}
		} else {
			// TODO no info recieved 
			if info == nil {
				
				logger.NgapLog.Warnln("info == nil")
				// continue
			} else if info.PPID != ngap.PPID{
				logger.NgapLog.Warnln("Received SCTP PPID != 60, discard this packet") 
			}
			logger.NgapLog.Tracef("Read %d bytes", n)
			logger.NgapLog.Tracef("Packet content:\n%+v", hex.Dump(buf[:n]))
			go handler.HandleMessage(lbConn, buf[:n])
		}
	}
}

// Initializes LB to AMF communication and starts handling the connection 
func StartAmf(amf *context.LbAmf, lbaddr *sctp.SCTPAddr, amfIP string, amfPort int, handler NGAPHandler) {
	self := context.LB_Self()
	logger.NgapLog.Debugf("Connecting to amf")
	for {
		conn, err := ConnectToAmf(lbaddr, amfIP, amfPort)
		if err == nil {
			amf.LbConn.Conn = conn
			ngap_message.SendNGSetupRequest(amf.LbConn)
			go handleConnection(amf.LbConn, readBufSize, handler)
			logger.NgapLog.Debugf("Connected to amf")
			self.Next_Amf = amf
			amf.LbConn.ID = ID
			connections.Store(ID, amf.LbConn)
			ID++
			break
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
	return conn, nil
}

