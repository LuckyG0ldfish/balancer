package context

import (
	"encoding/csv"
	"os"
	"strconv"
	
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
	dur 		int64
}

type AmfCounter struct {
	Amf 		*LbAmf
	IndivUE		int 
	Traffic		int
}

type metricsUE struct {
	id 			int64 
	time 		int64
}

func (r *Routing_Table) Print() {
	self := LB_Self()
	if self.MetricsLevel == 2 {
		printRouting(r)
		return 
	} else if self.MetricsLevel == 1 {
		printUETimings(r)
		return 
	}
	printRouting(r)
	printUETimings(r)
}

func printUETimings(r *Routing_Table) {
	var ueRegistrationList []*metricsUE
	var ueRegularList []*metricsUE
	var ueDeregistrationList []*metricsUE
	
	for i := 0; i < len(r.traces); i++ {
		temp := r.traces[i]
		
		if temp.ue_State == TypeIdRegist {
			ok, slot := idPresent(temp.ueID, ueRegistrationList)
			if !ok {
				ue := newMetricsUE(temp.ueID, temp.dur)
				ueRegistrationList = append(ueRegistrationList, ue)
			} else {
				new := ueRegistrationList[slot].time + temp.dur
				ueRegistrationList[slot].time = new
			}
		} else if temp.ue_State == TypeIdRegular {
			ok, slot := idPresent(temp.ueID, ueRegularList)
			if !ok {
				ue := newMetricsUE(temp.ueID, temp.dur)
				ueRegularList = append(ueRegularList, ue)
			} else {
				new := ueRegularList[slot].time + temp.dur
				ueRegularList[slot].time = new
			}
		} else {
			ok, slot := idPresent(temp.ueID, ueDeregistrationList)
			if !ok {
				ue := newMetricsUE(temp.ueID, temp.dur)
				ueDeregistrationList = append(ueDeregistrationList, ue)
			} else {
				new := ueDeregistrationList[slot].time + temp.dur
				ueDeregistrationList[slot].time = new
			}
		}
		
	}

	if len(ueRegistrationList) != 0 {
		sortedRegist := sortList(ueRegistrationList)
		output := createOutputList(sortedRegist)
		createAndWriteCSV(output, "./config/ueRegistTimings.csv")
	}
	if len(ueRegularList) != 0 {
		sortedRegular := sortList(ueRegularList)
		output := createOutputList(sortedRegular)
		createAndWriteCSV(output, "./config/ueRegularTimings.csv")
	}
	if len(ueRegistrationList) != 0 {
		sortedDeregist := sortList(ueDeregistrationList)
		output := createOutputList(sortedDeregist)
		createAndWriteCSV(output, "./config/ueDeregistTimings.csv")
	}
}

func createOutputList(sorted []*metricsUE) [][]string{
	var output [][]string 
	heads := []string{"GnbUeId", "Duration"}
	output = append(output, heads)
	for i := 0; i < len(sorted); i++ {
		ue := sorted[i]
		dur := strconv.Itoa(int(ue.time)/1000) // to millisecounds
		id := strconv.Itoa(int(ue.id))
		row := []string {id, dur}
		output = append(output, row) 
	}
	return output
}

func sortList(ueList []*metricsUE) []*metricsUE {
    for i := 1; i < len(ueList); i++ {
        var j = i
        for j >= 1 && ueList[j].id < ueList[j - 1].id {
            ueList[j], ueList[j - 1] = ueList[j - 1], ueList[j]
            j--
        }
    }
	return ueList
}

func idPresent(id int64, slice []*metricsUE) (bool, int) {
	for i := 0; i < len(slice); i++ {
		if slice[i].id == id {
			return true, i
		}
	}
	return false, 0
}

func newMetricsUE(id int64, dur int64) (*metricsUE){
	var ue metricsUE
	ue.id = id
	ue.time = dur
	return &ue
}

func printRouting(r *Routing_Table) {
	var output [][]string 
	heads := []string{"LbUeId", "GNB-ID", "AMF-ID", "Delay", "State"}
	output = append(output, heads)
	
	for i := 0; i < len(r.traces); i++ {
		trace := r.traces[i]
		state := strconv.FormatInt(int64(trace.ue_State), 10)
		id := strconv.FormatInt(trace.ueID, 10)
		origin := strconv.FormatInt(trace.origin, 10)
		destination := strconv.FormatInt(trace.destination, 10)
		time := strconv.FormatInt(trace.dur, 10)

		if trace.d_type == TypeAmf {
			row := []string{id, origin, destination, time, state}
			output = append(output, row)
		} else if trace.d_type == TypeGnb {
			row := []string{id, destination, origin, time, state}
			output = append(output, row)
		}	
	}

	// for i := 0; i < len(r.amfs); i++ {
	// 	amfC := r.amfs[i]
	// 	id := amfC.Amf.AmfID
	// 	logger.ContextLog.Infof("AMF %d : | individuel UEs %d | total traffic %d \n", uint64(id), amfC.IndivUE, amfC.Traffic)
	// }

	createAndWriteCSV(output, "./config/routing.csv")
}

func (r *Routing_Table) AddRouting_Element(origin int64, ueID int64, destination int64, d_type int, ue_State int, time int64) {
	trace := newTrace(origin, ueID, destination, d_type, ue_State, time)
	r.traces = append(r.traces, trace)
}

func NewTable() (*Routing_Table){
	var t Routing_Table
	t.traces = []*trace{}
	return &t
}

func newTrace(origin int64, ueID int64, destination int64, d_type int, ue_State	int, time int64) (*trace){
	var t trace
	t.id = trafficNum
	trafficNum ++
	t.origin = origin 
	t.ueID = ueID
	t.destination = destination
	t.d_type = d_type
	t.ue_State = ue_State
	t.dur = time
	return &t
}

func (r *Routing_Table) AddAmfCounter(amf *LbAmf) {
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

func createAndWriteCSV(input [][]string, location string) {
	file, err := os.Create(location)
	if err != nil {
		logger.ContextLog.Fatalf("failed creating file: %s", err)
	}
	writer := csv.NewWriter(file)
	for _, row := range input {
		_ = writer.Write(row)
	}
	writer.Flush()
	file.Close()
}