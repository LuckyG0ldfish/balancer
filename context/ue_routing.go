package context

import (
	"fmt"
) 

const TypeAmf int = 0 
const TypeGnb int = 1 

var trafficNum int = 1 

type Routing_Table struct {
	traces 		[]*trace
	amfs 		[]*AmfCounter
}

type trace struct {
	id 			int
	origin 		int64
	ueID 		int64
	destination int64
	d_type 		int 
	ue_State	int
}

type AmfCounter struct {
	Amf 		*LbAmf
	IndivUE		int 
	Traffic		int
}


func (r *Routing_Table) Print() {
	for i := 0; i < len(r.traces); i++ {
		p := r.traces[i]
		var s string
		if p.ue_State == TypeIdRegist {
			s = "Registration"
		} else if p.ue_State == TypeIdRegular {
			s = "Regular Traffic"
		} else {
			s = "Deregistration"
		}
		if p.d_type == TypeAmf {
			fmt.Printf("LbUeID: %d | GNB: %d -> AMF: %d || %s \n", uint64(p.ueID), uint64(p.origin), uint64(p.destination), s)
		} else if p.d_type == TypeGnb {
			fmt.Printf("LbUeID: %d | GNB: %d <- AMF: %d || %s \n", uint64(p.ueID), uint64(p.destination), uint64(p.origin), s)
		}	
	}
}

func (r *Routing_Table) AddRouting_Element(origin int64, ueID int64, destination int64, d_type int, ue_State int) {
	trace := newTrace(origin, ueID, destination, d_type, ue_State)
	r.traces = append(r.traces, trace)
}

func NewTable() (*Routing_Table){
	var t Routing_Table
	t.traces = []*trace{}
	return &t
}

func newTrace(origin int64, ueID int64, destination int64, d_type int, ue_State	int) (*trace){
	var t trace
	t.id = trafficNum
	trafficNum ++
	t.origin = origin 
	t.ueID = ueID
	t.destination = destination
	t.d_type = d_type
	t.ue_State = ue_State
	return &t
}

func (r *Routing_Table) addAmfCounter(amf *LbAmf) {
	r.amfs = append(r.amfs, newAmfCounter(amf))
}

func newAmfCounter(amf *LbAmf) *AmfCounter{
	var amfC AmfCounter
	amfC.Amf = amf
	amfC.IndivUE = 0 
	amfC.Traffic = 0 
	return &amfC
}

