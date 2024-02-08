package pipe

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	type createUserRequest struct {
		Username          string
		Email             string
		ProfilePictureUrl string
	}

	type args struct {
		args  createUserRequest
		funcs []Func[createUserRequest]
	}
	tests := []struct {
		name string
		args args
		want Responses
	}{
		{
			name: "create user",
			args: args{
				args: createUserRequest{
					Username:          "username",
					Email:             "email",
					ProfilePictureUrl: "profilePictureUrl",
				},
				funcs: []Func[createUserRequest]{
					func(args createUserRequest, responses Responses) (response any, err error) {
						return args.Username, nil
					},
					func(args createUserRequest, responses Responses) (response any, err error) {
						username := Get[string](responses)
						fmt.Println(username)
						return args.Email, nil
					},
					func(args createUserRequest, responses Responses) (response any, err error) {
						email := Get[string](responses)
						fmt.Println(email)
						return args.ProfilePictureUrl, nil
					},
					func(args createUserRequest, responses Responses) (response any, err error) {
						return nil, nil
					},
				},
			},
			want: pipeResponse{
				resp: []any{
					"username",
					"email",
					"profilePictureUrl",
					nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := Pipe(tt.args.funcs...)
			got, _ := exec(tt.args.args, nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPipeGo(t *testing.T) {
	type createUserRequest struct {
		Username          string
		Email             string
		ProfilePictureUrl string
	}

	type args struct {
		args  createUserRequest
		funcs []Func[createUserRequest]
	}
	tests := []struct {
		name string
		args args
		want Responses
	}{
		{
			name: "create user",
			args: args{
				args: createUserRequest{
					Username:          "username",
					Email:             "email",
					ProfilePictureUrl: "profilePictureUrl",
				},
				funcs: []Func[createUserRequest]{
					func(args createUserRequest, responses Responses) (response any, err error) {
						time.Sleep(1 * time.Second)
						return args.Username, nil
					},
					func(args createUserRequest, responses Responses) (response any, err error) {
						time.Sleep(1 * time.Second)
						return args.Email, nil
					},
					func(args createUserRequest, responses Responses) (response any, err error) {
						time.Sleep(1 * time.Second)
						return args.ProfilePictureUrl, nil
					},
				},
			},
			want: pipeResponse{
				resp: []any{
					"username",
					"email",
					"profilePictureUrl",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := PipeGo(tt.args.funcs...)
			got, _ := exec(tt.args.args, nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPipeAndPipeGo(t *testing.T) {
	type createUserRequest struct {
		Username          string
		Email             string
		ProfilePictureUrl string
	}

	type args struct {
		args  createUserRequest
		funcs []Func[createUserRequest]
	}
	tests := []struct {
		name string
		args args
		want Responses
	}{
		{
			name: "create user",
			args: args{
				args: createUserRequest{
					Username:          "username",
					Email:             "email",
					ProfilePictureUrl: "profilePictureUrl",
				},
				funcs: []Func[createUserRequest]{
					func(args createUserRequest, responses Responses) (response any, err error) {
						return args.Username, nil
					},
					func(args createUserRequest, responses Responses) (response any, err error) {
						username := Get[string](responses)
						fmt.Println(username)
						return args.Email, nil
					},
					func(args createUserRequest, responses Responses) (response any, err error) {
						email := Get[string](responses)
						fmt.Println(email)
						return args.ProfilePictureUrl, nil
					},
					PipeGo(
						func(args createUserRequest, responses Responses) (response any, err error) {
							time.Sleep(1 * time.Second)
							username, _ := Index[string](responses, 0)
							return "go" + username, nil
						},
						func(args createUserRequest, responses Responses) (response any, err error) {
							time.Sleep(1 * time.Second)
							email, _ := Index[string](responses, 1)
							return "go" + email, nil
						},
						func(args createUserRequest, responses Responses) (response any, err error) {
							time.Sleep(1 * time.Second)
							return "go" + args.ProfilePictureUrl, nil
						},
					),
				},
			},
			want: pipeResponse{
				[]any{
					"username",
					"email",
					"profilePictureUrl",
					"gousername",
					"goemail",
					"goprofilePictureUrl",
				},
			},
		},
		{
			name: "create with several pipe",
			args: args{
				args: createUserRequest{
					Username:          "username",
					Email:             "email",
					ProfilePictureUrl: "profilePictureUrl",
				},
				funcs: []Func[createUserRequest]{
					func(args createUserRequest, responses Responses) (response any, err error) {
						return args.Username, nil
					},
					func(args createUserRequest, responses Responses) (response any, err error) {
						username := Get[string](responses)
						fmt.Println(username)
						return args.Email, nil
					},
					func(args createUserRequest, responses Responses) (response any, err error) {
						email := Get[string](responses)
						fmt.Println(email)
						return args.ProfilePictureUrl, nil
					},
					PipeGo(
						func(args createUserRequest, responses Responses) (response any, err error) {
							time.Sleep(1 * time.Second)
							username, _ := Index[string](responses, 0)
							return "go" + username, nil
						},
						func(args createUserRequest, responses Responses) (response any, err error) {
							time.Sleep(1 * time.Second)
							email, _ := Index[string](responses, 1)
							return "go" + email, nil
						},
						func(args createUserRequest, responses Responses) (response any, err error) {
							time.Sleep(1 * time.Second)
							return "go" + args.ProfilePictureUrl, nil
						},
					),
					Pipe(
						func(args createUserRequest, responses Responses) (response any, err error) {
							username, _ := Index[string](responses, 0)
							return "p" + username, nil
						},
						func(args createUserRequest, responses Responses) (response any, err error) {
							email, _ := Index[string](responses, 1)
							return "p" + email, nil
						},
						func(args createUserRequest, responses Responses) (response any, err error) {
							return "p" + args.ProfilePictureUrl, nil
						},
						Pipe(
							func(args createUserRequest, responses Responses) (response any, err error) {
								username, _ := Index[string](responses, 0)
								return "p1" + username, nil
							},
							func(args createUserRequest, responses Responses) (response any, err error) {
								email, _ := Index[string](responses, 1)
								return "p1" + email, nil
							},
							func(args createUserRequest, responses Responses) (response any, err error) {
								return "p1" + args.ProfilePictureUrl, nil
							},
						),
					),
				},
			},
			want: pipeResponse{
				[]any{
					"username",
					"email",
					"profilePictureUrl",
					"gousername",
					"goemail",
					"goprofilePictureUrl",
					"pusername",
					"pemail",
					"pprofilePictureUrl",
					"p1username",
					"p1email",
					"p1profilePictureUrl",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := Pipe(tt.args.funcs...)
			got, _ := exec(tt.args.args, nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
