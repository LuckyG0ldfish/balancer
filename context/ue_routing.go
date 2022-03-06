package context

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/LuckyG0ldfish/balancer/logger"
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
	time 		time.Duration
}

type AmfCounter struct {
	Amf 		*LbAmf
	IndivUE		int 
	Traffic		int
}

type metricsUE struct {
	id 			int64 
	time 		time.Time
}

func (r *Routing_Table) Print() {
	self := LB_Self()
	if self.Metrics > 1{
		printRouting(r)
		return 
	} else if self.Metrics == 1 {
		printUETimings(r)
		return 
	}
	printRouting(r)
	printUETimings(r)
}

func printUETimings(r *Routing_Table) {
	start := time.Now()
	var ueList sync.Map
	for i := 0; i < len(r.traces); i++ {
		temp := r.traces[i]
		ex, ok := ueList.Load(temp.ueID)
		if !ok {
			ue := newMetricsUE(temp.ueID, start.Add(temp.time))
			ueList.Store(ue.id, ue)
		} else {
			mUE, ok := ex.(metricsUE)
			if !ok {
				logger.ContextLog.Error("ue timings print failed")
				return 
			}
			mUE.time.Add(temp.time)
		}
	}

	var output string

	ueList.Range(func(key, value interface{}) bool {
		ue, ok := value.(metricsUE)
		if !ok {
			logger.NgapLog.Errorf("ue timings print failed")
			return false 
		}
		duration := ue.time.Sub(start)
		dur := duration.String()
		id := strconv.Itoa(int(ue.id))
		output = output + " (" + id + ": " + dur + ")"
		return true
	})
}

func newMetricsUE(id int64, time time.Time) (*metricsUE){
	var ue metricsUE
	ue.id = id
	ue.time = time 
	return &ue
}

func printRouting(r *Routing_Table) {
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

	for i := 0; i < len(r.amfs); i++ {
		amfC := r.amfs[i]
		id := amfC.Amf.AmfID
		fmt.Printf("AMF %d : | individuel UEs %d | total traffic %d \n", uint64(id), amfC.IndivUE, amfC.Traffic)
	}
}

func (r *Routing_Table) AddRouting_Element(origin int64, ueID int64, destination int64, d_type int, ue_State int, time time.Duration) {
	trace := newTrace(origin, ueID, destination, d_type, ue_State, time)
	r.traces = append(r.traces, trace)
}

func NewTable() (*Routing_Table){
	var t Routing_Table
	t.traces = []*trace{}
	return &t
}

func newTrace(origin int64, ueID int64, destination int64, d_type int, ue_State	int, time time.Duration) (*trace){
	var t trace
	t.id = trafficNum
	trafficNum ++
	t.origin = origin 
	t.ueID = ueID
	t.destination = destination
	t.d_type = d_type
	t.ue_State = ue_State
	t.time = time
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

func (r *Routing_Table) incrementAmfTraffic(amf *LbAmf) {
	for i := 0; i < len(r.amfs); i++ {
		if r.amfs[i].Amf.AmfID == amf.AmfID {
			r.amfs[i].Traffic++
			return
		}
	}
}

func (r *Routing_Table) incrementAmfIndividualUEs(amf *LbAmf) {
	for i := 0; i < len(r.amfs); i++ {
		if r.amfs[i].Amf.AmfID == amf.AmfID {
			r.amfs[i].IndivUE++
			return
		}
	}
}

