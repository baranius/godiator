#### Creating Pipelines

Pipelines should be derived from **godiatr.Pipeline** struct and include a **Handle** method just like handlers. 

**Important!!! :** *Even though **params ...interface{}** is optional, you should define them in your Pipeline's Handle method*

```go 
type ValidationPipeline struct {
	godiatr.Pipeline
}

func (p *ValidationPipeline) Handle(request interface{}, params ...interface{}) (interface{}, error) {
	r := request.(*handler.SampleRequest)

	if r.PayloadString == nil {
		return nil, errors.New("PayloadString_should_not_be_null")
	}

	return p.Next().Handle(request, params...)
}
```

#### Registering Pipelines 

You should call **RegisterPipeline** method to register pipelines.

**Important!!! :** *Pipelines will run in the definition order.*

```go
func RegisterPipelines() {
    g := godiatr.GetInstance()
    
    g.RegisterPipeline(&ValidationPipeline{}) // This will run first
    g.RegisterPipeline(&SomeOtherPipeline{}) // This will run second
}
```