package context

type UeKey struct{
	UeRanID 	int64
	RanID		int64
}

func (k1 *UeKey) Compare(k2 *UeKey) bool{
	return k1.UeRanID == k2.UeRanID && k2.RanID == k1.RanID
}

func NewKey() *UeKey{
	var key UeKey
	return &key
}