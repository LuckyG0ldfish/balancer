package context

import (
	"encoding/csv"
	"os"
	"strconv"
	"sync"

	"github.com/LuckyG0ldfish/balancer/logger"
) 

const TypeAmf int = 0 
const TypeGnb int = 1 

var trafficNum int = 1 

type MetricsGNB struct {
	ID 			int64
	MetricsUEs 	*sync.Map
}

type metricsUE struct {
	id 			int64 
	
	regTime 	int64 
	deregTime 	int64

	routings 	[]*trace
}

type trace struct {
	id 			int
	origin 		int64
	ueID 		int64
	destination int64
	destType 		int 
	ue_State	int
	startTime int64
	endTime int64
}

func AddRouting_Element(mGNBs *sync.Map, origin int64, ueID int64, destination int64, destType int, ue_State int, startTime int64, endTime int64) {	
	var id int64
	if destType == TypeAmf {
		id = origin
	} else {
		id = destination
	}
	gnb, ok := mGNBs.Load(id)
	if !ok {
		logger.ContextLog.Warning("metricsGNB does not exist (failed lookup)")
		return 
	}
	metricsGNB, ok := gnb.(*MetricsGNB)
	if !ok {
		logger.ContextLog.Warning("metricsGNB does not exist (failed type cast)")
		return 
	}
	
	trace := newTrace(origin, ueID, destination, destType, ue_State, startTime, endTime)
	
	ue, ok := metricsGNB.MetricsUEs.Load(ueID)
	if ok {
		mue, test := ue.(*metricsUE)
		if !test {
			logger.ContextLog.Error("error while parsing metricsUE")
		}
		mue.routings = append(mue.routings, trace)
	} else {
		mue := newMetricsUE(ueID)
		mue.routings = append(mue.routings, trace)
		mGNBs.Store(ueID, mue)
	}
}

func Print(gnbs *sync.Map) {
	gnbs.Range(func(key, value interface{}) bool {
		tempGNB, ok := value.(*MetricsGNB)
		if !ok {
			logger.NgapLog.Errorf("error while parsing metricsUE")
		}
		PrintPerGNB(tempGNB)
		return true
	})
}

func PrintPerGNB(gnb *MetricsGNB) {
	self := LB_Self()
	sortedUEs, routingTable := prepareMapForOutput(gnb.MetricsUEs)
	if self.MetricsLevel == 2 {
		printRouting(routingTable, gnb.ID)
		return 
	} else if self.MetricsLevel == 1 {
		printUETimings(sortedUEs, gnb.ID)
		return 
	}
	printRouting(routingTable, gnb.ID)
	printUETimings(sortedUEs, gnb.ID)
}

func prepareMapForOutput(m *sync.Map) (sorted []*metricsUE, routingTraces []*trace) {
	var unsorted []*metricsUE

	m.Range(func(key, value interface{}) bool {
		tempUE, ok := value.(*metricsUE)
		if !ok {
			logger.NgapLog.Errorf("error while parsing metricsUE")
		}
		unsorted = append(unsorted, tempUE)
		return true
	})

	sorted = sortUEsByUEID(unsorted) 

	for i := 0; i < len(sorted); i++ {
		var registTraces []*trace
		var deregTraces []*trace
		
		tempUE := sorted[i]
		tempUE.routings = sortTracesByStartTime(tempUE.routings)
		
		for j := 0; j < len(tempUE.routings); j++ {
			tempTrace := tempUE.routings[j] // creating the routing table 
			routingTraces = append(routingTraces, tempTrace)
			if tempTrace.ue_State == TypeIdRegist {
				registTraces = append(registTraces, tempTrace)
			} else if tempTrace.ue_State == TypeIdDeregist {
				deregTraces = append(deregTraces, tempTrace)
			} 
		}
		tempUE.regTime = calcuateDuration(registTraces)
		tempUE.deregTime = calcuateDuration(deregTraces)

	}
	return 
}

func calcuateDuration(traces []*trace) int64 {
	var dur int64
	var start int64 
	var end int64

	for i := 0; i < len(traces); i++ {
		if i == 0 {
			start = traces[i].startTime
			end = traces[i].endTime
			dur = dur + (end-start) 
			continue
		}
		if traces[i].startTime >= end && traces[i].endTime > end {
			start = traces[i].startTime
			end = traces[i].endTime
			dur = dur + (end - start)
		} else if traces[i].startTime <= end && traces[i].endTime > end {
			dur = dur + (traces[i].endTime - end)
			end = traces[i].endTime
		}
	}

	return dur
}

func printUETimings(m []*metricsUE, id int64) {
	var registOutput [][]string 
	var deregOutput [][]string 
	
	heads := []string{"GnbUeId", "Duration"}
	registOutput = append(registOutput, heads)
	deregOutput = append(deregOutput, heads)
	
	for i := 0; i < len(m); i++ {
		temp := m[i]
		dur := strconv.Itoa(int(temp.regTime)) //1000) // to millisecounds
		id := strconv.Itoa(int(temp.id))
		row := []string {id, dur}
		registOutput = append(registOutput, row)
		
		dur = strconv.Itoa(int(temp.deregTime)) //1000) // to millisecounds
		row = []string {id, dur}
		deregOutput = append(deregOutput, row)
	}

	s := strconv.FormatInt(id, 10)
	if len(registOutput) != 0 {
		createAndWriteCSV(registOutput, "./config/ueRegistTimingsGNB" + s + ".csv")
	}
	if len(deregOutput) != 0 {
		createAndWriteCSV(deregOutput, "./config/ueDeregistTimingsGNB" + s + ".csv")
	}
}

func sortUEsByUEID(ueList []*metricsUE) []*metricsUE {
    for i := 1; i < len(ueList); i++ {
        var j = i
        for j >= 1 && ueList[j].id < ueList[j - 1].id {
            ueList[j], ueList[j - 1] = ueList[j - 1], ueList[j]
            j--
        }
    }
	return ueList
}

func sortTracesByStartTime(traceList []*trace) []*trace {
    for i := 1; i < len(traceList); i++ {
        var j = i
        for j >= 1 && traceList[j].startTime < traceList[j - 1].startTime {
            traceList[j], traceList[j - 1] = traceList[j - 1], traceList[j]
            j--
        }
    }
	return traceList
}

// func idPresent(id int64, slice []*metricsUE) (bool, int) {
// 	for i := 0; i < len(slice); i++ {
// 		if slice[i].id == id {
// 			return true, i
// 		}
// 	}
// 	return false, 0
// }

func printRouting(traces []*trace, id int64) {
	var output [][]string 
	heads := []string{"GNBUeId", "GNB-ID", "AMF-ID", "Delay", "State"}
	output = append(output, heads)
	
	for i := 0; i < len(traces); i++ {
		trace := traces[i]
		state := strconv.FormatInt(int64(trace.ue_State), 10)
		id := strconv.FormatInt(trace.ueID, 10)
		origin := strconv.FormatInt(trace.origin, 10)
		destination := strconv.FormatInt(trace.destination, 10)
		time := strconv.FormatInt(trace.endTime - trace.startTime, 10)

		if trace.destType == TypeAmf {
			row := []string{id, origin, destination, time, state}
			output = append(output, row)
		} else if trace.destType == TypeGnb {
			row := []string{id, destination, origin, time, state}
			output = append(output, row)
		}	
	}
	s := strconv.FormatInt(id, 10)
	createAndWriteCSV(output, "./config/routingGNB" + s + ".csv")
}

func NewMetricsGNB(id int64) (*MetricsGNB) {
	var metricsGNB MetricsGNB
	metricsGNB.ID = id 
	metricsGNB.MetricsUEs = NewMetricsMap()
	return &metricsGNB
}

func newMetricsUE(id int64) (*metricsUE){
	var t metricsUE
	t.id = id 
	t.regTime = 0 
	t.deregTime = 0 
	return &t
}

func newTrace(origin int64, ueID int64, destination int64, destType int, ue_State	int, startTime int64, endTime int64) (*trace){
	var t trace
	t.id = trafficNum
	trafficNum ++
	t.origin = origin 
	t.ueID = ueID
	t.destination = destination
	t.destType = destType
	t.ue_State = ue_State
	t.startTime = startTime
	t.endTime = endTime
	return &t
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