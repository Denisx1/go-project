package middleware

import (
	logger "go_project/infrastructure/logging"
	messaging "go_project/infrastructure/messaging"
	"go_project/infrastructure/messaging/component/producer"
)

type Middleware func(next messaging.Handler) messaging.Handler

func Chain(h messaging.Handler, mws ...Middleware) messaging.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// Wrap: RequestID -> Trace -> Log -> Retry -> Recover -> handler
func Wrap(
	h messaging.Handler,
	log logger.Logger,
	writer producer.Writer,
	system string,
) messaging.Handler {
	return Chain(h,
		withRequestID(),
		withLog(log),
		withTrace(),
		withDLQ(writer),
		withRetry(3),
		withRecover(),
	)
}
