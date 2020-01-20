package godiatr

import (
	"reflect"
	"sync"
)

var (
	godiatr *Godiatr
	once sync.Once
)

type IGodiatr interface {
	RegisterPipeline(h IPipeline)
	Register(request interface{}, handler func()interface{})
	Send(request interface{}, params ...interface{}) (interface{}, error)
}

type Godiatr struct {
	registries map[reflect.Type]func()interface{}
	pipelines  []IPipeline
}

func GetInstance() *Godiatr {
	once.Do(func() {
		godiatr = &Godiatr{registries: make(map[reflect.Type]func()interface{})}
	})
	return godiatr
}

func (m *Godiatr) RegisterPipeline(h IPipeline) {
	m.pipelines = append(m.pipelines, h)
}

func (m *Godiatr) Register(request interface{}, handler func()interface{}) {
	m.registries[reflect.TypeOf(request)] = handler
}

func (m *Godiatr) Send(request interface{}, params ...interface{}) (interface{}, error) {
	// Initialize an anonymous handler
	businessPipeline := new(businessPipeline)
	businessPipeline.mediator = m

	if len(m.pipelines) > 0 {
		// Loop through pipelines by reverse if exists and bind them to each other
		var mainPipeline IPipeline
		for i := len(m.pipelines) - 1; i >= 0; i-- {
			if i == len(m.pipelines)-1 {
				pipeline := m.pipelines[i]
				pipeline.SetNext(businessPipeline)
				mainPipeline = m.pipelines[i]
			} else {
				m.pipelines[i].SetNext(mainPipeline)
				mainPipeline = m.pipelines[i]
			}
		}
		if mainPipeline != nil {
			// Call nested handlers w/+ given parameters
			return mainPipeline.Handle(request, params...)
		} else {
			panic("Pipeline error")
		}
	} else {
		// Call handler w/- pipeline if there is no pipeline
		return businessPipeline.Handle(request, params...)
	}
}
