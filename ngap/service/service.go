package service

import (
	"encoding/hex"
	"io"
	"sync"
	"syscall"
	"time"

	"github.com/ishidawataru/sctp"

	"github.com/LuckyG0ldfish/balancer/logger"
	"github.com/free5gc/ngap"

	"github.com/LuckyG0ldfish/balancer/context"
)

type NGAPHandler struct {
	HandleMessage      func(lbConn *context.LBConn, msg []byte, startTime int64)
	HandleNotification func(conn *context.LBConn, notification sctp.Notification)
}

const readBufSize uint32 = 256 // 8192
const amfPort int = 38412

// set default read timeout to 2 seconds
var readTimeout syscall.Timeval = syscall.Timeval{Sec: 2, Usec: 0}

var (
	sctpListener 	*sctp.SCTPListener
	connections  	sync.Map	
	numberGen 		context.UniqueNumberGen
)

var sctpConfig sctp.SocketConfig = sctp.SocketConfig{
	InitMsg:   sctp.InitMsg{NumOstreams: 3, MaxInstreams: 5, MaxAttempts: 2, MaxInitTimeout: 2},
	RtoInfo:   &sctp.RtoInfo{SrtoAssocID: 0, SrtoInitial: 500, SrtoMax: 1500, StroMin: 100},
	AssocInfo: &sctp.AssocInfo{AsocMaxRxt: 4},
}

// Starts all NGAP related services
func Run(addr *sctp.SCTPAddr, handler NGAPHandler) {
	// creating number generator for ports to connect to AMFs 
	self := context.LB_Self()
	numberGen = *context.NewUniqueNumberGen(int64(self.LbToAmfPort))
	
	// All AMFs related services started 
	go InitAmfs(handler)
	
	// All GNBs related services started 
	go listenAndServeGNBs(addr, handler)
}

// Handling all the the LBs open connections (AMFs + GNBs)
func handleConnection(lbConn *context.LBConn, bufsize uint32, handler NGAPHandler) {// conn *sctp.SCTPConn, bufsize uint32, handler NGAPHandler) {
	logger.NgapLog.Tracef("Waiting for message")
	for !lbConn.Closed {
		buf := make([]byte, bufsize)

		n, info, notification, err := lbConn.Conn.SCTPRead(buf)
		
		if err != nil {
			switch err {
			case io.EOF, io.ErrUnexpectedEOF:
				logger.NgapLog.Debugln("Read EOF from client")
				return
			case syscall.EAGAIN:
				logger.NgapLog.Traceln("SCTP read timeout")
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
		if notification != nil {
			if handler.HandleNotification != nil {
				handler.HandleNotification(lbConn, notification)
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
			time := int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Microsecond)
			go handler.HandleMessage(lbConn, buf[:n], time)
		}
	}
}
// Removes a LB_Conn from the list of connections and the related AMF/GNB from their pool 
func RemoveLBConnection(conn *context.LBConn) {
	conn.Closed = true 
	if conn.TypeID == context.TypeAmf {
		conn.AmfPointer.RemoveAmfContext()
	} else {
		conn.RanPointer.RemoveGnbContext()
	}
	time.Sleep(1 * time.Second)
	connections.Delete(conn.ID)
}

// Closes all connections 
func Stop() {

	logger.NgapLog.Infof("Close SCTP server...")
	if err := sctpListener.Close(); err != nil {
		logger.NgapLog.Error(err)
		logger.NgapLog.Infof("SCTP server may not close normally.")
	}
	
	connections.Range(func(key, value interface{}) bool {
		lbConn, ok := value.(context.LBConn)
		if !ok {
			logger.NgapLog.Errorf("couldn't be converted")
		}
		if lbConn.Conn == nil {
			return true 
		}
		if err := lbConn.Conn.Close(); err != nil {
			logger.NgapLog.Error(err)
		}
		return true
	})

	logger.NgapLog.Infof("SCTP server closed")
}