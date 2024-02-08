package pipe

import "context"

// FuncCtx is a series of pipe functions arguments contract
// when using pipe the function callback must comply with this contract
type FuncCtx[Args any] func(ctx context.Context, args Args, responses Responses) (response any, err error)

// Func is a series of pipe functions arguments contract
// when using pipe the function callback must comply with this contract
type Func[Args any] func(args Args, responses Responses) (response any, err error)

// P to start compose functions with Pipe
// use when initiate pipe at the beginning
func P[Args any](funcs ...Func[Args]) func(args Args) (responses Responses, err error) {
	return func(args Args) (responses Responses, err error) {
		exec := Pipe(funcs...)
		resp, err := exec(args, nil)
		r, _ := resp.(Responses)
		return r, err
	}
}

// PCxt to start compose functions with Pipe and Context
// use when initiate pipe at the beginning
func PCtx[Args any](funcs ...FuncCtx[Args]) func(ctx context.Context, args Args) (responses Responses, err error) {
	return func(ctx context.Context, args Args) (responses Responses, err error) {
		exec := PipeCtx(funcs...)
		resp, err := exec(ctx, args, nil)
		r, _ := resp.(Responses)
		return r, err
	}
}

// PipeCtx series of functions into one processing unit
// ordering our logic in form of pipe for readable and clean code
func PipeCtx[Args any](funcs ...FuncCtx[Args]) FuncCtx[Args] {
	return func(ctx context.Context, args Args, responses Responses) (response any, err error) {
		_, isFromPipe := responses.(pipeResponse)
		if responses == nil {
			responses = pipeResponse{}
		}

		var newResponses Responses = pipeResponse{}
		for _, f := range funcs {
			response, err = f(ctx, args, responses)
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

// PipeGoCtx enhance the serial processing of Pipe with Go rountine concurrency
// saving most of the time by utilize this function instead of writing manual Go routine code
func PipeGoCtx[Args any](funcs ...FuncCtx[Args]) FuncCtx[Args] {
	return func(ctx context.Context, args Args, responses Responses) (response any, err error) {
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
			go func(f FuncCtx[Args], i int) {
				response, err = f(ctx, args, responses)
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

// Pipe series of functions into one processing unit
// ordering our logic in form of pipe for readable and clean code
func Pipe[Args any](funcs ...Func[Args]) Func[Args] {
	var funcCtx []FuncCtx[Args]
	for _, f := range funcs {
		ff := f
		funcCtx = append(funcCtx, func(ctx context.Context, args Args, responses Responses) (response any, err error) {
			return ff(args, responses)
		})
	}

	fctx := PipeCtx(funcCtx...)
	return func(args Args, responses Responses) (response any, err error) {
		return fctx(context.TODO(), args, responses)
	}
}

// PipeGo enhance the serial processing of Pipe with Go rountine concurrency
// saving most of the time by utilize this function instead of writing manual Go routine code
func PipeGo[Args any](funcs ...Func[Args]) Func[Args] {
	var funcCtx []FuncCtx[Args]
	for _, f := range funcs {
		ff := f
		funcCtx = append(funcCtx, func(ctx context.Context, args Args, responses Responses) (response any, err error) {
			return ff(args, responses)
		})
	}

	fctx := PipeGoCtx(funcCtx...)
	return func(args Args, responses Responses) (response any, err error) {
		return fctx(context.TODO(), args, responses)
	}
}
