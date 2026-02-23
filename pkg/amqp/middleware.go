package amqp

// Middleware wraps a HandlerFunc with additional behaviour.
type Middleware func(HandlerFunc) HandlerFunc

// Chain applies middlewares to a handler in order: first middleware is outermost.
func Chain(handler HandlerFunc, mws ...Middleware) HandlerFunc {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler
}
