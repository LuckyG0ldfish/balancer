package context

import (
	"net"
	// "syscall"
	// "sync/atomic"

	"git.cs.nctu.edu.tw/calee/sctp"

	"github.com/LuckyG0ldfish/balancer/logger"
)

const TypeIdRegist 			int	= 0
const TypeIdRegular 		int	= 1
const TypeIdDeregist		int = 2

// func GetNoDelay(c *sctp.SCTPConn) (int, error) {
// 	return syscall.GetsockoptInt(c.fd(), SOL_SCTP, SCTP_NODELAY)
// }

// func SetNoDelay(bytes int, c *sctp.SCTPConn) error {
// 	return syscall.SetsockoptInt(c.fd(), SOL_SCTP, SCTP_NODELAY, bytes)
// }

// func fd(c *sctp.SCTPConn) int {
// 	return int(atomic.LoadInt32(&c._fd))
// }

// Writing a slice of bytes to a sctp.SCTPConn
func SendByteToConn(conn *sctp.SCTPConn, message []byte) {
	n, err := conn.Write(message)
	if err != nil {
		logger.NgapLog.Errorf("Write to SCTP socket failed: [%+v]", err)
	} else {
		logger.NgapLog.Tracef("Wrote %d bytes", n)
	}
}

// Use IP and port to generate *sctp.SCTPAddr
func GenSCTPAddr(ip string, port int) (lbAddr *sctp.SCTPAddr, err error){
	ips := []net.IPAddr{}
	if ip, err1 := net.ResolveIPAddr("ip", ip); err1 != nil {
		return nil, err1
	} else {
		ips = append(ips, *ip)
	}
	lbAddr = &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    port,
	}
	return lbAddr, nil
}

