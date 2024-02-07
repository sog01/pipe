package pipe

import (
	"errors"
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
		want []any
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
					func(args createUserRequest, responses []any) (response any, err error) {
						return args.Username, nil
					},
					func(args createUserRequest, responses []any) (response any, err error) {
						username := responses[0].(string)
						fmt.Println(username)
						return args.Email, nil
					},
					func(args createUserRequest, responses []any) (response any, err error) {
						email := responses[1].(string)
						fmt.Println(email)
						return args.ProfilePictureUrl, nil
					},
					func(args createUserRequest, responses []any) (response any, err error) {
						return nil, nil
					},
				},
			},
			want: []any{
				"username",
				"email",
				"profilePictureUrl",
				nil,
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
		want []any
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
					func(args createUserRequest, responses []any) (response any, err error) {
						time.Sleep(1 * time.Second)
						return args.Username, nil
					},
					func(args createUserRequest, responses []any) (response any, err error) {
						time.Sleep(1 * time.Second)
						return args.Email, nil
					},
					func(args createUserRequest, responses []any) (response any, err error) {
						time.Sleep(1 * time.Second)
						return args.ProfilePictureUrl, nil
					},
				},
			},
			want: []any{
				"username",
				"email",
				"profilePictureUrl",
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
		name    string
		args    args
		wantErr error
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
					func(args createUserRequest, responses []any) (response any, err error) {
						return args.Username, nil
					},
					func(args createUserRequest, responses []any) (response any, err error) {
						username := responses[0].(string)
						fmt.Println(username)
						return args.Email, nil
					},
					func(args createUserRequest, responses []any) (response any, err error) {
						email := responses[1].(string)
						fmt.Println(email)
						return args.ProfilePictureUrl, nil
					},
					PipeGo(
						func(args createUserRequest, responses []any) (response any, err error) {
							time.Sleep(1 * time.Second)
							username := responses[0].(string)
							return "go" + username, nil
						},
						func(args createUserRequest, responses []any) (response any, err error) {
							time.Sleep(1 * time.Second)
							email := responses[1].(string)
							return "go" + email, nil
						},
						func(args createUserRequest, responses []any) (response any, err error) {
							time.Sleep(1 * time.Second)
							return "go" + args.ProfilePictureUrl, nil
						},
					),
				},
			},
			wantErr: nil,
		},
		{
			name: "create user error",
			args: args{
				args: createUserRequest{
					Username:          "username",
					Email:             "email",
					ProfilePictureUrl: "profilePictureUrl",
				},
				funcs: []Func[createUserRequest]{
					func(args createUserRequest, responses []any) (response any, err error) {
						return args.Username, nil
					},
					func(args createUserRequest, responses []any) (response any, err error) {
						username := responses[0].(string)
						fmt.Println(username)
						return args.Email, nil
					},
					func(args createUserRequest, responses []any) (response any, err error) {
						email := responses[1].(string)
						fmt.Println(email)
						return args.ProfilePictureUrl, nil
					},
					PipeGo(
						func(args createUserRequest, responses []any) (response any, err error) {
							time.Sleep(1 * time.Second)
							return "go" + args.Username, errors.New("errors go username")
						},
						func(args createUserRequest, responses []any) (response any, err error) {
							time.Sleep(1 * time.Second)
							return "go" + args.Email, nil
						},
						func(args createUserRequest, responses []any) (response any, err error) {
							time.Sleep(1 * time.Second)
							return "go" + args.ProfilePictureUrl, nil
						},
					),
				},
			},
			wantErr: errors.New("errors go username"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := Pipe(tt.args.funcs...)
			_, err := exec(tt.args.args, nil)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
