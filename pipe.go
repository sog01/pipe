package pipe

type Func[Args any] func(args Args, responses []any) (response any, err error)

func Pipe[Args any](funcs ...Func[Args]) Func[Args] {
	return func(args Args, responses []any) (response any, err error) {
		for _, f := range funcs {
			response, err = f(args, responses)
			if err != nil {
				return nil, err
			}
			if response == nil {
				continue
			}
			responses = append(responses, response)
		}
		return responses, nil
	}
}

func PipeGo[Args any](funcs ...Func[Args]) Func[Args] {
	return func(args Args, responses []any) (response any, err error) {
		c := make(chan struct {
			response any
			err      error
		})
		for _, f := range funcs {
			go func(f Func[Args]) {
				response, err = f(args, responses)
				c <- struct {
					response any
					err      error
				}{
					response: response,
					err:      err,
				}
			}(f)
		}

		for i := 0; i < len(funcs); i++ {
			resp := <-c
			if resp.err != nil {
				return responses, resp.err
			}
			responses = append(responses, resp.response)
		}

		return responses, nil
	}
}
