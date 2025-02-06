package listners

import "context"

type ListenerServer interface {
	Start(ctx context.Context)
	Stop()
}
