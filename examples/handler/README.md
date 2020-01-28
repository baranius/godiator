## Handlers
Creating a handler in godiatr requires a method named as **Handle**. 

**Handle** method should get a request objects and return a tuple of response and error objects 

```go
type (
	SampleRequest struct{
		PayloadString *string
	}

	SampleResponse struct{
		ResultString *string
	}

	SampleHandler struct{
	}
)

func NewSampleHandler() interface{} {
	return &SampleHandler{}
}

func (h *SampleHandler) Handle(request *SampleRequest, params ...interface{}) (*SampleResponse, error){
	return &SampleResponse{ResultString: request.PayloadString}, nil
}
```

#### Registering Handlers

You should call **RegisterHandler** method to register handlers. It takes request model reference
and handler's initializer method as arguments.

```go
func RegisterHandlers() {
    g := godiatr.GetInstance()
    
    g.RegisterHandler(&SampleRequest{}, NewSampleHandler)
}
```

#### Calling Handlers 

You should call **Send** method to run registered handlers. 

```go
type GetItemController struct {
    g godiatr.IGodiatr
}

func NewGetItemController() *GetItemController {
    return &GetItemController{g: godiatr.GetInstance()}
}

func (c *GetItemController) GetItem() {
    payloadValue := "sample_value"
    request := &SampleRequest{PayloadString: &payloadValue}
    
    response, err := c.g.Send(request)
}
```