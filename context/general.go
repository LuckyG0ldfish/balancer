package context 

import (
	"net"
	

	"git.cs.nctu.edu.tw/calee/sctp"
)

func SendByteToConn(conn *sctp.SCTPConn, message []byte) {
	conn.Write(message)
}

func FindUeInSlice(slice []*LbUe, UeAmfID int64) (*LbUe, int){
	if len(slice) == 0 {
		return nil, 0
	} 
	
	for _, ue := range slice {
		if ue.UeAmfId == UeAmfID {
			return ue, 1
		}
	}
	for _, ue := range slice {
		if ue.UeAmfId == 0 {
			return ue, 2
		}
	}
	return nil, 3
}

func GenSCTPAddr(ip string, port int) (lbAddr *sctp.SCTPAddr, err error){
	ips := []net.IPAddr{}
	if ip, err1 := net.ResolveIPAddr("ip", ip); err1 != nil {
		//err := fmt.Errorf("Error resolving address '%s': %v", ip, err1)
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

