package pipe

// Func is a series of pipe functions arguments contract
// when using pipe the function callback must comply with this contract
type Func[Args any] func(args Args, responses Responses) (response any, err error)

// P to start compose functions with Pipe
// use when initiate pipe at the beginning
func P[Args any](funcs ...Func[Args]) func(args Args) (response any, err error) {
	return func(args Args) (response any, err error) {
		exec := Pipe(funcs...)
		return exec(args, nil)
	}
}

// Pipe series of functions into one processing unit
// ordering our logic in form of pipe for readable and clean code
func Pipe[Args any](funcs ...Func[Args]) Func[Args] {
	return func(args Args, responses Responses) (response any, err error) {
		_, isFromPipe := responses.(pipeResponse)
		if responses == nil {
			responses = pipeResponse{}
		}

		var newResponses Responses = pipeResponse{}
		for _, f := range funcs {
			response, err = f(args, responses)
			if err != nil {
				return nil, err
			}
			responses = responses.Add(response)

			if isFromPipe {
				newResponses = newResponses.Add(response)
			}
		}

		if isFromPipe {
			return newResponses, nil
		}

		return responses, nil
	}
}

// PipeGo enhance the serial processing of Pipe with Go rountine concurrency
// saving most of the time by utilize this function instead of writing manual Go routine code
func PipeGo[Args any](funcs ...Func[Args]) Func[Args] {
	return func(args Args, responses Responses) (response any, err error) {
		_, isFromPipe := responses.(pipeResponse)
		if responses == nil {
			responses = pipeResponse{}
		}
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

		var newResponses Responses = pipeResponse{}
		for i := 0; i < len(mres); i++ {
			responses = responses.Add(mres[i])
			if isFromPipe {
				newResponses = newResponses.Add(mres[i])
			}
		}
		if isFromPipe {
			return newResponses, nil
		}

		return responses, nil
	}
}
