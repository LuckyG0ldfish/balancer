package context

import (
	"sync"
)

var (
	lbContext 					= LBContext{}
)

func init() {
	LB_Self().Name = "lb"
	//LB_Self().NetworkName.Full = "free5GC"
}

type LBContext struct{
	Name 			string 
	//NetworkName   factory.NetworkName

	LbRanPool		sync.Map 		// gNBs connected to the LB 
	LbAmfPool		sync.Map		// amfs (each connected to AMF 1:1) connected to LB
}

func LB_Self() *LBContext {
	return &lbContext
}