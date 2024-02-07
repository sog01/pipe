# Pipe

Pipe is a [Go](http://golang.org/) library for organize our code to be more readable and clean in form of functions series.

## Installation

Install Pipe using "go get" command:

    go get github.com/sog01/pipe

## Example Usage

We can use pipe to serialize our logic into several functions with single responsibilty principle.

```go
package main

import (
	"errors"
	"strings"

	"github.com/sog01/pipe"
)

func main() {
	e := pipe.Pipe(
		isUserEmailExists,
		validateUserEmail,
		insertUser,
	)

	_, err := e(UserInput{
		Email: "john.doe@gmail.com",
	}, nil)
	if err != nil {
		panic(err)
	}
}

var DB = make(map[string]any)

type UserInput struct {
	Email    string
	Password string
}

func isUserEmailExists(args UserInput, responses pipe.Responses) (response any, err error) {
	_, exists := DB[args.Email]
	if exists {
		return nil, errors.New("email already exists")
	}
	return nil, nil
}

func validateUserEmail(args UserInput, responses pipe.Responses) (response any, err error) {
	if !strings.Contains(args.Email, "@") {
		return nil, errors.New("incorrect email address")
	}
	return nil, nil
}

func insertUser(args UserInput, responses pipe.Responses) (response any, err error) {
	DB[args.Email] = args
	return nil, nil
}
```

As a other developer / reviewer we can understand the flow in the first place, then jump to each function to get deeper context.

## Responses Use Cases

As a pipe function we can use a previous response as an input for the next function. Like an example below :

```go
package main

import (
	"errors"
	"strings"

	"github.com/sog01/pipe"
)

func main() {
	e := pipe.Pipe(
		getBlacklistUsers,
		isBlacklistUser,
		...
	)

	_, err := e(UserInput{
		Email: "john.blacklist@bmail.com",
	}, nil)
	if err != nil {
		panic(err)
	}
}

var (
	DB = make(map[string]any)
)

type UserInput struct {
	Email    string
	Password string
}

func getBlacklistUsers(args UserInput, responses pipe.Responses) (response any, err error) {
	return map[string]any{
		"john.blacklist@bmail.com": struct{}{},
		"doe.blacklist@bmail.com":  struct{}{},
	}, nil
}

func isBlacklistUser(args UserInput, responses pipe.Responses) (response any, err error) {
	blacklistUsers := pipe.Get[map[string]any](responses)
	_, isBlacklist := blacklistUsers[args.Email]
	if isBlacklist {
		return nil, errors.New("this email is from blacklist")
	}
	return nil, nil
}
```

The scenario is to validate the incoming users whether is blackisted or not, so at the beginning we can get the blacklist users. Then, we utilize the `responses` on the next function to validate the users.

## Concurrency in Pipe with PipeGo

The pipe also support concurrency by simply use `PipeGo` instead. This abstract concurrency under the hood, so we don't need to write Go routine manually:

```go
package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sog01/pipe"
)

func main() {
	e := pipe.Pipe(
		getBlacklistUsers,
		isBlacklistUser,
		isUserEmailExists,
		validateUserEmail,
		pipe.PipeGo(
			insertUser,
			sendNotification,
			trackUserBehaviour,
		),
	)

	_, err := e(UserInput{
		Email: "john.doe@gmail.com",
	}, nil)
	if err != nil {
		panic(err)
	}
}

func insertUser(args UserInput, responses pipe.Responses) (response any, err error) {
	DB[args.Email] = args
	return nil, nil
}

func sendNotification(args UserInput, responses pipe.Responses) (response any, err error) {
	fmt.Println("send notification")
	return nil, nil
}

func trackUserBehaviour(args UserInput, responses pipe.Responses) (response any, err error) {
	fmt.Println("track user behaviour")
	return nil, nil
}
```
