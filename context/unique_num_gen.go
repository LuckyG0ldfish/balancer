package context

import (
	"container/list"
	
	"github.com/LuckyG0ldfish/balancer/logger"
)

// Generator for UE-IDs
type UniqueNumberGen struct{
	ListEmpty 		bool
	Recycled 		*list.List
	RegularID		int64
}

// Creates, initializes and returns a new *UeIdGen
func NewUniqueNumberGen(StartNumber int64) *UniqueNumberGen{
	var UeIdGen UniqueNumberGen
	UeIdGen.ListEmpty = true 
	UeIdGen.Recycled = list.New()
	UeIdGen.RegularID = int64(StartNumber)
	return &UeIdGen
}

// Selects the next available LB-internal ID for a UE 
func (gen *UniqueNumberGen) NextNumber() int64 {
	// default case at the moment 
	if gen.ListEmpty {
		return gen.addOne()
	} 

	// unused
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

// not used at the moment 
// Takes unused IDs and makes them available for reuse 
func (gen *UniqueNumberGen) RecycleNumber(id int64) {
	gen.Recycled.InsertAfter(id, gen.Recycled.Back())
	logger.ContextLog.Tracef("Number reusable: %d", id)
}

// checks whether there are integers that can be reused 
func (gen *UniqueNumberGen) checkEmpty() {
	if gen.Recycled == nil || gen.Recycled.Len() == 0 {
		gen.ListEmpty = true 
		return 
	}
	gen.ListEmpty = false 
}

// adds a new integer by incrementing the current 
func (gen *UniqueNumberGen) addOne() int64{
	id := gen.RegularID
	gen.RegularID++
	return id
}