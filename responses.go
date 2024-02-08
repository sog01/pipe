package pipe

// Responses define a functions responses inside a pipe
type Responses interface {
	Add(any) Responses
	Get() []any
}

type pipeResponse struct {
	resp []any
}

func (p pipeResponse) Add(a any) Responses {
	if r, ok := a.(pipeResponse); ok {
		p.resp = append(p.resp, r.resp...)
	} else {
		p.resp = append(p.resp, a)
	}
	return p
}

func (p pipeResponse) Get() []any {
	return p.resp
}

// Get a last response from the previous pipe
func Get[T any](r Responses) (t T) {
	rr := r.Get()
	res, ok := rr[len(rr)-1].(T)
	if ok {
		return res
	}
	return t
}

// Index used to get a responses from a given index
func Index[T any](r Responses, i int) (t T, valid bool) {
	rr := r.Get()
	if i > len(rr)-1 {
		return t, false
	}
	res, ok := rr[i].(T)
	if ok {
		return res, true
	}
	return t, false
}

// Find used to find responses from a given generic type
func Find[T any](r Responses) (t T, found bool) {
	rr := r.Get()
	for _, rr := range rr {
		if t, ok := rr.(T); ok {
			return t, true
		}
	}
	return t, false
}
