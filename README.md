# godiator

Easy to use mediator implementation in Golang. Provides in-process messaging ability to your apps.

Implementation examples:
- [Echo] (https://github.com/baranius/godiator-echo)

## Installation

You should install godiator via Go package manager:

```
$ go get -v https://github.com/baranius/godiator
```

## Usage

Use **GetInstance()** method to get **godiator** object. Take a look to complete [API](#api) reference for details.

```go
type GetItemController struct {
	g godiator.IGodiator
}

func NewGetItemController() *GetItemController {
	return &GetItemController{g: godiator.GetInstance()}
}
```

## Handlers

[Handlers](#handler) are the main objects that runs your business logic.

A handler is a struct that contains **Handle** method. 

**Handle** method should get a request object and return a tuple of response object and error.

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

//Registering handlers in godiator requires an init function.
func NewSampleHandler() interface{} {
	return &SampleHandler{}
}

func (h *SampleHandler) Handle(request *SampleRequest) (*SampleResponse, error){
	return &SampleResponse{ResultString: request.PayloadString}, nil
}
```

**Handle** method can take additional parameters 

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

func (h *SampleHandler) Handle(request *SampleRequest, ctx context.Context) (*SampleResponse, error){
	return &SampleResponse{ResultString: request.PayloadString}, nil
}
``` 

#### Registering Handlers

**Register** method takes request model
and handler's initializer method as arguments.

```go
func RegisterHandlers() {
    g := godiator.GetInstance()
    
    g.Register(&SampleRequest{}, NewSampleHandler)
}
```

#### Calling Handlers 

Calling **Send** method with request object executes the related handler. 

```go
type GetItemController struct {
    g godiator.IGodiator
}

func NewGetItemController() *GetItemController {
    return &GetItemController{g: godiator.GetInstance()}
}

func (c *GetItemController) GetItem() {
    payloadValue := "sample_value"
    request := &SampleRequest{PayloadString: &payloadValue}
    
    response, err := c.g.Send(request)

    // If your handle method takes additional parameters, don't forget to pass them to 'Send'
    response, err := c.g.Send(request, context.TODO())
}
```
___

## Subscriptions

[Subscriptions](#subscription) are observers. 
They operate quite similar to message broadcasting.

```go
type (
	SubscriptionRequest struct {
        PayloadString *string
    }
    
    SubscriptionHandler struct {}
) 

func NewSubscriptionHandler() interface{} {
	return &SubscriptionHandler{}
}

func (n *SubscriptionHandler) Handle(request interface{}) {
	r := request.(*SubscriptionRequest)
	fmt.Printf("Subscription called with payload : '%v'", *r.PayloadString)
}
```

#### Registering Subscriptions

**RegisterSubscription** method takes request model 
and related handler(s).

```go
func RegisterSubscriptions() {
    g := godiator.GetInstance()
    
    g.RegisterSubscription(&SubscriptionRequest{}, NewSubscriptionHandler, NewSubscriptionHandlerA, NewSubscriptionHandlerB)
}
```

#### Calling Subscriptions 

Calling **Publish** method with request object notifies the related handler(s). 

```go
type GetItemController struct {
    g godiator.IGodiator
}

func NewGetItemController() *GetItemController {
    return &GetItemController{g: godiator.GetInstance()}
}

func (c *GetItemController) GetItem() {
    payloadValue := "sample_value"
    request := &SubscriptionRequest{PayloadString: &payloadValue}
    
    c.g.Publish(request)
}
```
___

## Pipelines

[Pipelines](#pipeline) are interceptors that runs before each handler call.

They should be derived from **godiator.Pipeline** struct and include a **Handle** method just like handlers. 
     
*Important:  **params ...interface{}** is optional. But you should define them in your Pipeline's Handle method*


```go 
type ValidationPipeline struct {
	godiator.Pipeline
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

Call **RegisterPipeline** method to register pipelines.

*Important: Pipelines run in the ***definition order***.*

```go
func RegisterPipelines() {
    g := godiator.GetInstance()
    
    g.RegisterPipeline(&ValidationPipeline{}) // Will run first
    g.RegisterPipeline(&SomeOtherPipeline{})
}
```

## Mocking

You can mock Send and Notify methods of **godiator** via **MockGodiator**.

The only thing you need to do is delegating **OnSend** or **OnNotify** methods.

#### OnSend
```go
mockGodiator := MockGodiator{}

mockGodiator.OnSend = func(request interface{}, params ...interface{}) (i interface{}, err error){
    return &Response{}, nil
}
```

#### OnNotify

```go
mockGodiator := MockGodiator{}

mockGodiator.OnNotify = func(request interface{}, params ...interface{}) {
    fmt.Print("Called")
}
```

## API

### Godiator

##### RegisterPipeline(h IPipeline)

- **h:** `A struct derived from godiator.Pipeline`

##### Register(request interface{}, handler func()interface{})

- **request:** `Interface`
- **handler:** `Initialize function which returns Handler instance`

##### RegisterSubscription(request interface{}, handler func()interface{})
 
 - **request:** `Interface`
 - **handler:** `Initialize function which returns Subscription handler instance`
 
##### Send(request interface{}, params ...interface{}) (interface{}, error)

- **request:** `Handler's request model`
- **params (optional):** `Optional list of objects`

##### Publish(request interface{}, params ...interface{})

- **request:** `Handler's request model`
- **params (optional):** `Optional list of objects`

### Handler

##### Handle(request interface{}, params ...interface{}) (interface{}, error)

- **request:** `Handler's request model`
- **params (optional):** `Optional list of objects`

### Subscription

##### Publish(request interface{}, params ...interface{}) 

- **request:** `Handler's request model`
- **params (optional):** `Optional list of objects`

### Pipeline

##### Handle(request interface{}, params ...interface{}) (interface{}, error) 

- **request:** `Handler's request model`
- **params (optional):** `Optional list of objects`

## Contribution

You're very welcome to contribute the project. Please feel free to contribute or asking questions. 
