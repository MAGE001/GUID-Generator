package generator

// Generator generate global unique ids
type Generator interface {
	// batch get n guid
	NextIds(n int) []int64
}
