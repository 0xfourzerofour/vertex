package proxy

import (
	"context"
	"reflect"

	"github.com/valyala/fasthttp"
)

type GQLResp struct {
	Data  map[string]interface{} `json:"data,omitempty"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type ProxyRepository interface {
	Forward(ctx context.Context, url string, body []byte) (*GQLResp, error)
}

func HttpConnError(err error) (string, bool) {
	errName := ""
	known := false
	if err == fasthttp.ErrTimeout {
		errName = "timeout"
		known = true
	} else if err == fasthttp.ErrNoFreeConns {
		errName = "conn_limit"
		known = true
	} else if err == fasthttp.ErrConnectionClosed {
		errName = "conn_close"
		known = true
	} else {
		errName = reflect.TypeOf(err).String()
		if errName == "*net.OpError" {
			errName = "timeout"
			known = true
		}
	}
	return errName, known
}