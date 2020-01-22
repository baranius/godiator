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
___

## Handlers

Handlers are the main objects that runs your business logic.

#### Creating Handlers
Creating a handler in godiatr requires a method named as **Handle**. Take a look to [API](#handler) reference for details.

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

#### [More >>>](https://github.com/baranx/godiatr/tree/master/examples/handler)

___

## Notifications

Notifications are observers. The way they run is quite similar to Pub-Sub. 

#### Creating Notifications
Creating a notification handler in godiatr requires a method named as **Handle**. Take a look to [API](#notification) for details

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

#### Registering Notifications

You should call **RegisterNotificationHandler** method to register notification handlers. It takes request model reference
and handler's initializer method as arguments.

```go
func RegisterHandlers() {
    g := godiatr.GetInstance()
    
    g.RegisterNotificationHandler(&NotificationRequest{}, NewNotification)
}
```

#### Calling Notifications 

You should call **Notify** method to run registered notification handlers. 

```go
type GetItemController struct {
    g godiatr.IGodiatr
}

func NewGetItemController() *GetItemController {
    return &GetItemController{g: godiatr.GetInstance()}
}

func (c *GetItemController) GetItem() {
    payloadValue := "sample_value"
    request := &NotificationRequest{PayloadString: &payloadValue}
    
    c.g.Notify(request)
}
```

#### [More >>>](https://github.com/baranx/godiatr/tree/master/examples/notification)

___

## Pipelines

A pipeline is a method that runs before each godiatr handler. It's quite similar to interceptor design pattern.

#### Creating Pipelines

Pipelines should be derived from **godiatr.Pipeline** struct and include a **Handle** method just like handlers. . Take a look to [API](#pipeline) for details

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

#### [More >>>](https://github.com/baranx/godiatr/tree/master/examples/pipelines)

___

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

## Tests

You can download the git repository an run the test in the **example** folder for testing or debugging.

## Contribution

You're very welcome to contribute the project. Please feel free to contribute or asking questions. 
