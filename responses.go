package pipe

// Responses define a functions responses inside a pipe
type Responses []any

// Get a last response from the previous pipe
func Get[T any](r Responses) (t T) {
	res, ok := r[len(r)-1].(T)
	if ok {
		return res
	}
	return t
}

// Index used to get a responses from a given index
func Index[T any](r Responses, i int) (t T, valid bool) {
	if i > len(r)-1 {
		return t, false
	}
	res, ok := r[i].(T)
	if ok {
		return res, true
	}
	return t, false
}
