package context

import (
	"container/list"
	
	"github.com/LuckyG0ldfish/balancer/logger"
)

// Generator for UE-IDs
type UeIdGen struct{
	ListEmpty 		bool
	Recycled 		*list.List
	RegularID		int64
}

// Creates, initializes and returns a new *UeIdGen
func NewUeIdGen() *UeIdGen{
	var UeIdGen UeIdGen
	UeIdGen.ListEmpty = true 
	UeIdGen.Recycled = list.New()
	UeIdGen.RegularID = 4
	return &UeIdGen
}

// Selects the next available LB-internal ID for a UE 
func (gen *UeIdGen) NextID() int64 {
	if gen.ListEmpty {
		return gen.addOne()
	} 

	e := gen.Recycled.Remove(gen.Recycled.Front())
	if e == nil {
		return gen.addOne()
	}
	id, ok := e.(int64)
	if !ok {
		return gen.addOne()
	} else {
		go gen.checkEmpty()
		return id
	}
}

// Takes unused IDs and makes them available for reuse 
func (gen *UeIdGen) RecycleID(id int64) {
	gen.Recycled.InsertAfter(id, gen.Recycled.Back())
	logger.ContextLog.Traceln("ID recycled: %d", id)
}

func (gen *UeIdGen) checkEmpty() {
	if gen.Recycled == nil || gen.Recycled.Len() == 0 {
		gen.ListEmpty = true 
		return 
	}
	gen.ListEmpty = false 
}

func (gen *UeIdGen) addOne() int64{
	id := gen.RegularID
	gen.RegularID++
	logger.ContextLog.Traceln("New ID generated")
	return id
}