package context

import (
	"container/list"
	"fmt" 
) 

type UeIdGen struct{
	Recycled 	*list.List
	RegularID		int64
}

func NewUeIdGen() *UeIdGen{
	var UeIdGen UeIdGen
	UeIdGen.Recycled = list.New()
	UeIdGen.RegularID = 4
	return &UeIdGen
}

func (gen *UeIdGen) NextID() int64 {
	if gen.Recycled == nil || gen.Recycled.Len() == 0 {
		return gen.AddOne()
	} 
	e := gen.Recycled.Remove(gen.Recycled.Front())
	if e == nil {
		return gen.AddOne()
	} else {
		id, ok := e.(int64)
		if !ok {
			return gen.AddOne()
		} else {
			return id
		}
	}
}

func (gen *UeIdGen) RecycleID(id int64) {
	gen.Recycled.InsertAfter(id, gen.Recycled.Back())
	fmt.Println("ID recycled: %d", id)
}

func (gen *UeIdGen) AddOne() int64{
	id := gen.RegularID
	gen.RegularID++
	fmt.Println("New ID generated")
	return id
}