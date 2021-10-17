package context 

import (
	"net"
	

	"git.cs.nctu.edu.tw/calee/sctp"
)

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