package context 

import (
	"net"
	
	"git.cs.nctu.edu.tw/calee/sctp"
)

// Writing a slice of bytes to a sctp.SCTPConn
func SendByteToConn(conn *sctp.SCTPConn, message []byte) {
	conn.Write(message)
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

