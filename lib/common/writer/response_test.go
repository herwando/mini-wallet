package writer

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herwando/mini-wallet/lib/common/commonerr"
)

var errFoo = errors.New("errFoo")

func TestWriteOK(t *testing.T) {
	w := httptest.NewRecorder()
	type args struct {
		ctx  context.Context
		w    http.ResponseWriter
		data interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				w:    w,
				data: "ok",
			},
		},
		{
			args: args{
				w:    w,
				data: make(chan int),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteOK(tt.args.ctx, tt.args.w, tt.args.data)
		})
	}
}

func TestWriteStrOK(t *testing.T) {
	w := httptest.NewRecorder()
	type args struct {
		ctx context.Context
		w   http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				w: w,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteStrOK(tt.args.ctx, tt.args.w)
		})
	}
}

func TestWriteJSONAPIError(t *testing.T) {
	w := httptest.NewRecorder()
	type args struct {
		ctx context.Context
		w   http.ResponseWriter
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				ctx: context.WithValue(context.Background(), ErrorCtxKey, nil),
				w:   w,
				err: errFoo,
			},
		},
		{
			args: args{
				ctx: context.WithValue(context.Background(), ErrorCtxKey, nil),
				w:   w,
				err: commonerr.SetNewInternalError(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteJSONAPIError(tt.args.ctx, tt.args.w, tt.args.err)
		})
	}
}
