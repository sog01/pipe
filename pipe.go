package pipe

type Func[Args any] func(args Args, responses []any) (response any, err error)

// Pipe series of functions into one processing unit
// ordering our logic in form of pipe for readable and clean code
func Pipe[Args any](funcs ...Func[Args]) Func[Args] {
	return func(args Args, responses []any) (response any, err error) {
		for _, f := range funcs {
			response, err = f(args, responses)
			if err != nil {
				return nil, err
			}
			responses = append(responses, response)
		}
		return responses, nil
	}
}

// PipeGo enhance the serial processing of Pipe with Go rountine concurrency
// saving most of the time by utilize this function instead of writing manual Go routine code
func PipeGo[Args any](funcs ...Func[Args]) Func[Args] {
	return func(args Args, responses []any) (response any, err error) {
		c := make(chan struct {
			index    int
			response any
			err      error
		})
		for i, f := range funcs {
			go func(f Func[Args], i int) {
				response, err = f(args, responses)
				c <- struct {
					index    int
					response any
					err      error
				}{
					index:    i,
					response: response,
					err:      err,
				}
			}(f, i)
		}

		mres := make(map[int]any)
		for i := 0; i < len(funcs); i++ {
			resp := <-c
			if resp.err != nil {
				return responses, resp.err
			}
			mres[resp.index] = resp.response
		}

		for i := 0; i < len(mres); i++ {
			responses = append(responses, mres[i])
		}

		return responses, nil
	}
}
