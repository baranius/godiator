package godiatr

type (
	IPipeline interface {
		Next() IPipeline
		SetNext(handler IPipeline)
		Handle(request interface{}, params ...interface{}) (interface{}, error)
	}

	Pipeline struct {
		IPipeline
		next IPipeline
	}
)

func (h *Pipeline) SetNext(handler IPipeline) {
	h.next = handler
}

func (h *Pipeline) Next() IPipeline {
	return h.next
}
