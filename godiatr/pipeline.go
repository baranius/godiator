package godiatr

import "errors"

type Pipeline struct {
	next IPipeline
}

func (h *Pipeline) SetNext(handler IPipeline) {
	h.next = handler
}

func (h *Pipeline) Next() IPipeline {
	return h.next
}

func (h *Pipeline) Handle(request interface{}, params ...interface{}) (interface{}, error) {
	return nil, errors.New("unhandled pipeline")
}
