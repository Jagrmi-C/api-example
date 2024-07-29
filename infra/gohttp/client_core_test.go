package gohttp_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jc88/api-example/infra/gohttp"
)

func TestNewClientWhitIAMAuth(t *testing.T) {
	t.Run(`The builder was builded without options`,
		func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
		},
	)
}

func TestDefaultBackoff(t *testing.T) {
	type args struct {
		min        time.Duration
		max        time.Duration
		attemptNum int
		resp       *http.Response
	}
	tests := []struct {
		args args
		want time.Duration
	}{
		{
			args: args{
				min:        time.Second,
				max:        5 * time.Second,
				attemptNum: 0,
			},
			want: time.Second,
		},
		{
			args: args{
				min:        time.Second,
				max:        5 * time.Minute,
				attemptNum: 1,
			},
			want: 2 * time.Second,
		},
		{
			args: args{
				min:        time.Second,
				max:        5 * time.Minute,
				attemptNum: 2,
			},
			want: 4 * time.Second,
		},
		{
			args: args{
				min:        time.Second,
				max:        5 * time.Minute,
				attemptNum: 3,
			},
			want: 8 * time.Second,
		},
		{
			args: args{
				min:        time.Second,
				max:        5 * time.Minute,
				attemptNum: 63,
			},
			want: 5 * time.Minute,
		},
		{
			args: args{
				min:        time.Second,
				max:        5 * time.Minute,
				attemptNum: 128,
			},
			want: 5 * time.Minute,
		},
		{
			args: args{
				min:        time.Second,
				max:        5 * time.Minute,
				attemptNum: 1,
				resp: &http.Response{
					StatusCode: http.StatusTooManyRequests,
					Header: http.Header{
						"Retry-After": []string{"120"},
					},
				},
			},
			want: 120 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("min: %v, max: %v", tt.args.min, tt.args.max), func(t *testing.T) {
			if got := gohttp.DefaultBackoff(tt.args.min, tt.args.max, tt.args.attemptNum, tt.args.resp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultBackoff() = %v, want %v", got, tt.want)
			}

			if v := gohttp.DefaultBackoff(
				tt.args.min,
				tt.args.max,
				tt.args.attemptNum,
				tt.args.resp,
			); v != tt.want {
				t.Fatalf("bad: %#v -> %s", tt, v)
			}
		})
	}
}

func TestDefaultRetryPolicy(t *testing.T) {
	type args struct {
		ctx  context.Context
		resp *http.Response
		err  error

		exceededCtx bool
	}
	tests := []struct {
		name        string
		args        args
		shouldRetry bool
		wantErr     bool
	}{
		{
			name: "with response that has 429 status code",
			args: args{
				ctx: context.Background(),
				resp: &http.Response{
					StatusCode: http.StatusTooManyRequests,
				},
			},
			shouldRetry: true,
		},
		{
			name: "with response that has 501 status code",
			args: args{
				ctx: context.Background(),
				resp: &http.Response{
					StatusCode: http.StatusNotImplemented,
				},
			},
			shouldRetry: false,
		},
		{
			name: "with response that has strange 0 status code",
			args: args{
				ctx: context.Background(),
				resp: &http.Response{
					StatusCode: 0,
				},
			},
			shouldRetry: true,
		},
		{
			name: "with response that has 503 status code",
			args: args{
				ctx: context.Background(),
				resp: &http.Response{
					StatusCode: 503,
				},
			},
			shouldRetry: true,
		},
		{
			name: "with response that has 404 status code",
			args: args{
				ctx: context.Background(),
				resp: &http.Response{
					StatusCode: 404,
				},
			},
			shouldRetry: false,
		},
		{
			name: "with exceeded context",
			args: args{
				ctx: context.Background(),
				resp: &http.Response{
					StatusCode: 404,
				},
				exceededCtx: true,
			},
			shouldRetry: false,
			wantErr:     true,
		},
		{
			name: "with redirect error",
			args: args{
				ctx:  context.Background(),
				resp: &http.Response{},
				err: &url.Error{
					Err: errors.New("stopped after 097171162745 redirects"),
				},
			},
			shouldRetry: false,
		},
		{
			name: "with schema error",
			args: args{
				ctx:  context.Background(),
				resp: &http.Response{},
				err: &url.Error{
					Err: errors.New("unsupported protocol scheme"),
				},
			},
			shouldRetry: false,
		},
		{
			name: "with certificate error",
			args: args{
				ctx:  context.Background(),
				resp: &http.Response{},
				err: &url.Error{
					Err: errors.New("certificate is not trusted"),
				},
			},
			shouldRetry: false,
		},
		{
			name: "with unexpected error",
			args: args{
				ctx:  context.Background(),
				resp: &http.Response{},
				err:  errors.New("an unexpected error"),
			},
			shouldRetry: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(tt.args.ctx)
			defer cancel()

			if tt.args.exceededCtx {
				cancel()
			}

			if tt.name == "with redirect error" {
				fmt.Println("")
			}

			got, err := gohttp.DefaultRetryPolicy(ctx, tt.args.resp, tt.args.err)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			if got != tt.shouldRetry {
				assert.Equal(t, tt.shouldRetry, got)
			}
		})
	}
}
