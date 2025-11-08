package interfaces

type (
	// Main Handler interface
	Handler[TRequest any, TResponse any] interface {
		Handle(request TRequest, params ...any) (TResponse, error)
	}

	// Fire & Forget handler interface
	Subscriber[TRequest any] interface {
		Handle(request TRequest, params ...any)
	}

	// Pipeline handler interface
	Pipeline interface {
		Next() Pipeline
		SetNext(p Pipeline)
		Handle(request any, params ...any) (any, error)
	}
)
