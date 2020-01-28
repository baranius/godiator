# godiatr

Easy to use mediator implementation in Golang. Provides in-process messaging ability to your apps.

## Installation

You should install godiatr via Go package manager:

```
$ go get -v https://github.com/baranx/godiatr
```

## Usage

Use **GetInstance()** method to get **godiatr** object. Take a look to complete [API](#api) reference for details.

```go
type GetItemController struct {
	g godiatr.IGodiatr
}

func NewGetItemController() *GetItemController {
	return &GetItemController{g: godiatr.GetInstance()}
}
```

## Handlers

[Handlers](#handler) are the main objects that runs your business logic.

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

##### [Detailed Example](https://github.com/baranx/godiatr/tree/master/examples/handler)
___

## Notifications

[Notifications](#notification) are observers. The way they run is quite similar to Pub-Sub. 

```go
type (
    Notification struct {}

    NotificationRequest struct {
        PayloadString *string
    }
) 

func NewNotification() interface{} {
	return &Notification{}
}

func (n *Notification) Handle(request interface{}, params ...interface{}) {
	r := request.(*NotificationRequest)
	fmt.Printf("Notification called with payload : '%v'", *r.PayloadString)
}
```

##### [Detailed Example](https://github.com/baranx/godiatr/tree/master/examples/notification)

___

## Pipelines

[Pipelines](#pipeline) are routines that runs before each godiatr handler. It's quite similar to interceptor design pattern.

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

**Important!!! :** *Pipelines will run in the definition order.*

##### [Detailed Example](https://github.com/baranx/godiatr/tree/master/examples/pipelines)

## Mocking

You can mock Send and Notify methods of **godiatr** via **MockGodiatr** in mock package.

The only thing you need to do is delegating **OnSend** or **OnNotify** methods.


#### OnSend
```go
mockGodiatr.OnSend = func(request interface{}, params ...interface{}) (i interface{}, err error){
    return &Response{}, nil
}
```

[Complete Example](https://github.com/baranx/godiatr/tree/master/examples/mocking/send)

#### OnNotify

```go
mockGodiatr.OnNotify = func(request interface{}, params ...interface{}) {
    fmt.Print("Called")
}
```
[Complete Example](https://github.com/baranx/godiatr/tree/master/examples/mocking/notify)

## API

### Godiatr

##### RegisterPipeline(h IPipeline)

- **h:** `A struct derived from godiatr.Pipeline`

##### RegisterHandler(request interface{}, handler func()interface{})

- **request:** `Interface`
- **handler:** `Initialize function which returns Handler instance`

##### RegisterNotificationHandler(request interface{}, handler func()interface{})
 
 - **request:** `Interface`
 - **handler:** `Initialize function which returns Notification handler instance`
 
##### Send(request interface{}, params ...interface{}) (interface{}, error)

- **request:** `Handler's request model`
- **params (optional):** `Optional list of objects`

##### Notify(request interface{}, params ...interface{})

- **request:** `Handler's request model`
- **params (optional):** `Optional list of objects`

### Handler

##### Handle(request interface{}, params ...interface{}) (interface{}, error)

- **request:** `Handler's request model`
- **params (optional):** `Optional list of objects`

### Notification

##### Notify(request interface{}, params ...interface{}) 

- **request:** `Handler's request model`
- **params (optional):** `Optional list of objects`

### Pipeline

##### Handle(request interface{}, params ...interface{}) (interface{}, error) 

- **request:** `Handler's request model`
- **params (optional):** `Optional list of objects`

## Contribution

You're very welcome to contribute the project. Please feel free to contribute or asking questions. 
