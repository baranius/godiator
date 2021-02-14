package godiatr

type (
	IGodiatr interface {
		GetHandler(request interface{}) interface{}
		GetHandlerResponse(request interface{}) interface{}
		RegisterPipeline(h IPipeline)
		Register(request interface{}, handler func() interface{})
		RegisterNotification(request interface{}, handler func() interface{})
		Send(request interface{}, params ...interface{}) (interface{}, error)
		Publish(request interface{}, params ...interface{})
	}

	IPipeline interface {
		Next() IPipeline
		SetNext(handler IPipeline)
		Handle(request interface{}, params ...interface{}) (interface{}, error)
	}
)
