package listener

import "context"

type Handler interface {
	Setup(ctx *context.Context)
	Handle(bytes []byte)
}
