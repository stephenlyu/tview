package graphs


type Graph interface {
	GetValueRange(startIndex int, endIndex int) (float64, float64)
	Update(startIndex int, endIndex int)
	Clear()
}
