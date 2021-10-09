package context

import (
	"sync"
)

var (
	lbContext 					= LBContext{}
)

func init() {
	
}

type LBContext struct{
	LbRanPool		sync.Map 		// gNBs connected to the LB 
	LbAmfClients	sync.Map		// clients (each connected to AMF 1:1) connected to LB
}

func LB_Self() *LBContext {
	return &lbContext
}