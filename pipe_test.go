package pipe

import (
	"fmt"
	"testing"
	"time"

	util "github.com/sog01/high-order-funcs"
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
		want map[string]struct{}
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
			want: map[string]struct{}{
				"username":          {},
				"email":             {},
				"profilePictureUrl": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := PipeGo(tt.args.funcs...)
			got, _ := exec(tt.args.args, nil)
			gotMap := util.Reduce(got.([]interface{}), func(res map[string]struct{}, response any) map[string]struct{} {
				res[response.(string)] = struct{}{}
				return res
			}, map[string]struct{}{})
			assert.Equal(t, tt.want, gotMap)
		})
	}
}
