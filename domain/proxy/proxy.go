package proxy

import "context"

type ProxyRepository interface {
	Forward(ctx context.Context) error
}