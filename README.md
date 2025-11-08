# godiator

Easy to use mediator implementation in Golang. Provides in-process messaging ability to your apps with full type safety using Go generics.

## Features

- **Type-safe**: Uses Go generics for compile-time type checking
- **Simple API**: Clean and intuitive interface
- **Handlers**: Request/Response pattern for synchronous operations
- **Subscriptions**: Fire-and-forget pattern for asynchronous notifications
- **Pipelines**: Interceptor pattern for cross-cutting concerns
- **Mocking**: Built-in mocking support for testing
- **Zero dependencies**: No external dependencies beyond standard library

## Installation

Install godiator via Go package manager:

```bash
go get -v github.com/baranius/godiator
```

**Requirements**: Go 1.23 or higher

## Quick Start

```go
package main

import (
	"fmt"
	
	"github.com/baranius/godiator"
	"github.com/baranius/godiator/register"
)

// Define request and response types
type GetUserRequest struct {
	UserID int
}

type GetUserResponse struct {
	UserID   int
	Username string
	Email    string
}

// Define handler
type GetUserHandler struct{}

func (h *GetUserHandler) Handle(request GetUserRequest, params ...any) (GetUserResponse, error) {
	// Your business logic here
	return GetUserResponse{
		UserID:   request.UserID,
		Username: "john_doe",
		Email:    "john@example.com",
	}, nil
}

func main() {
	// Register handler
	register.Handler(GetUserRequest{}, &GetUserHandler{})
	
	// Send request
	request := GetUserRequest{UserID: 123}
	response, err := godiator.Send[GetUserRequest, GetUserResponse](request)
	
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("User: %s (%s)\n", response.Username, response.Email)
}
```

## Usage Examples

Complete usage examples:
- [Echo](https://github.com/baranius/godiator-echo)
- [gin](https://github.com/baranius/godiator-gin)

## Handlers

Handlers implement the request/response pattern. They are the main objects that execute your business logic.

A handler is a struct that implements the `Handler[TRequest, TResponse]` interface by providing a `Handle` method.

### Handler Interface

```go
type Handler[TRequest any, TResponse any] interface {
	Handle(request TRequest, params ...any) (TResponse, error)
}
```

### Basic Handler Example

```go
type (
	SampleRequest struct {
		PayloadString *string
	}

	SampleResponse struct {
		ResultString *string
	}

	SampleHandler struct{}
)

func (h *SampleHandler) Handle(request SampleRequest, params ...any) (SampleResponse, error) {
	return SampleResponse{
		ResultString: request.PayloadString,
	}, nil
}
```

### Handler with Additional Parameters

The `Handle` method can accept additional parameters (like `context.Context`) through the variadic `params` parameter:

```go
import "context"

func (h *SampleHandler) Handle(request SampleRequest, params ...any) (SampleResponse, error) {
	// Extract context if provided
	var ctx context.Context
	if len(params) > 0 {
		if c, ok := params[0].(context.Context); ok {
			ctx = c
		}
	}
	
	// Use context in your business logic
	_ = ctx
	
	return SampleResponse{
		ResultString: request.PayloadString,
	}, nil
}
```

### Registering Handlers

Use `register.Handler()` to register your handler. The first parameter is a zero-value instance of your request type used for type registration.

```go
func init() {
	register.Handler(SampleRequest{}, &SampleHandler{})
}
```

You can also register handlers dynamically:

```go
func RegisterHandlers() {
	register.Handler(SampleRequest{}, &SampleHandler{})
	register.Handler(OtherRequest{}, &OtherHandler{})
}
```

### Calling Handlers

Use `godiator.Send[TRequest, TResponse]()` to execute the handler for a given request:

```go
payloadValue := "sample_value"
request := SampleRequest{PayloadString: &payloadValue}

response, err := godiator.Send[SampleRequest, SampleResponse](request)

if err != nil {
	// Handle error
}

// Use response
fmt.Println(*response.ResultString)
```

### Calling Handlers with Additional Parameters

Pass additional parameters as variadic arguments:

```go
import "context"

ctx := context.Background()
response, err := godiator.Send[SampleRequest, SampleResponse](request, ctx)
```

---

## Subscriptions

Subscriptions implement the observer pattern (fire-and-forget). They operate similar to message broadcasting and are executed asynchronously in goroutines.

### Subscriber Interface

```go
type Subscriber[TRequest any] interface {
	Handle(request TRequest, params ...any)
}
```

**Note**: Subscribers do not return values or errors. They are fire-and-forget operations.

### Basic Subscriber Example

```go
import (
	"fmt"
	"time"
)

type (
	UserCreatedEvent struct {
		UserID   int
		Username string
		Email    string
		CreatedAt  time.Time
	}

	EmailNotificationSubscriber struct{}
	LoggingSubscriber struct{}
)

func (s *EmailNotificationSubscriber) Handle(request UserCreatedEvent, params ...any) {
	// Send welcome email
	fmt.Printf("Sending welcome email to %s\n", request.Email)
}

func (s *LoggingSubscriber) Handle(request UserCreatedEvent, params ...any) {
	// Log the event
	fmt.Printf("User created: %d at %s\n", request.UserID, request.CreatedAt)
}
```

### Registering Subscriptions

Use `register.Subscriber()` to register one or more subscribers for a request type:

```go
func RegisterSubscriptions() {
	register.Subscriber(
		UserCreatedEvent{},
		&EmailNotificationSubscriber{},
		&LoggingSubscriber{},
	)
}
```

### Calling Subscriptions

Use `godiator.Publish[TRequest]()` to notify all registered subscribers:

```go
import (
	"github.com/baranius/godiator"
	"time"
)

event := UserCreatedEvent{
	UserID:   123,
	Username: "john_doe",
	Email:    "john@example.com",
	CreatedAt: time.Now(),
}

godiator.Publish[UserCreatedEvent](event)
```

**Note**: Subscribers are executed asynchronously in separate goroutines. If you need to wait for completion, you must implement your own synchronization mechanism.

### Calling Subscriptions with Additional Parameters

```go
import "context"

ctx := context.Background()
godiator.Publish[UserCreatedEvent](event, ctx)
```

---

## Pipelines

Pipelines are interceptors that run before each handler call. They allow you to implement cross-cutting concerns like validation, logging, authentication, etc.

### Pipeline Structure

Pipelines must embed `core.Pipeline` and implement the `Handle` method:

```go
import (
	"errors"
	
	"github.com/baranius/godiator/core"
)

type ValidationPipeline struct {
	core.Pipeline
}

func (p *ValidationPipeline) Handle(request any, params ...any) (any, error) {
	// Type assert to your request type
	r := request.(*SampleRequest)
	
	// Perform validation
	if r.PayloadString == nil {
		return nil, errors.New("PayloadString should not be null")
	}
	
	// Call next pipeline or handler
	return p.Next().Handle(request, params...)
}
```

**Important**: 
- Always call `p.Next().Handle()` to continue the pipeline chain
- If validation fails, return an error without calling `Next()`
- The `params ...any` parameter is optional but should be defined in your Pipeline's Handle method signature

### Registering Pipelines

Use `register.Pipeline()` to register pipelines. **Pipelines run in the order they are registered**.

```go
import "github.com/baranius/godiator/register"

func RegisterPipelines() {
	register.Pipeline(&ValidationPipeline{})      // Runs first
	register.Pipeline(&LoggingPipeline{})         // Runs second
	register.Pipeline(&AuthenticationPipeline{})  // Runs third
	// Handler executes after all pipelines
}
```

### Pipeline Execution Order

Pipelines are executed in the order they are registered:
1. First registered pipeline
2. Second registered pipeline
3. ...
4. Last registered pipeline
5. Handler

Each pipeline can:
- Modify the request before passing it to the next pipeline
- Validate the request and return an error
- Log or monitor the execution
- Modify the response after the handler executes

---

## Mocking

godiator provides a built-in mocking package (`mockiator`) for testing your handlers and subscribers.

### Mocking Handlers

Use `mockiator.OnSend()` to mock handler execution:

```go
import (
	"github.com/baranius/godiator"
	"github.com/baranius/godiator/mockiator"
)

func TestHandler(t *testing.T) {
	// Setup mock
	mockHandler := mockiator.OnSend[SampleRequest, SampleResponse](
		SampleRequest{},
		func(request SampleRequest, params ...any) (SampleResponse, error) {
			return SampleResponse{
				ResultString: request.PayloadString,
			}, nil
		},
	)
	
	// Execute
	request := SampleRequest{PayloadString: stringPtr("test")}
	response, err := godiator.Send[SampleRequest, SampleResponse](request)
	
	// Verify
	assert.NoError(t, err)
	assert.True(t, mockHandler.IsCalled)
	assert.Equal(t, 1, mockHandler.TimesCalled)
	assert.Equal(t, "test", *response.ResultString)
}

func stringPtr(s string) *string {
	return &s
}
```

### Mocking Subscribers

Use `mockiator.OnPublish()` to mock subscriber execution:

```go
import (
	"github.com/baranius/godiator"
	"github.com/baranius/godiator/mockiator"
)

func TestSubscriber(t *testing.T) {
	// Setup mock
	mockSubscriber := mockiator.OnPublish[UserCreatedEvent](
		UserCreatedEvent{},
		func(request UserCreatedEvent, params ...any) {
			// Mock behavior
			fmt.Printf("Mock: User %d created\n", request.UserID)
		},
	)
	
	// Execute
	event := UserCreatedEvent{UserID: 123}
	godiator.Publish[UserCreatedEvent](event)
	
	// Wait for async execution
	time.Sleep(100 * time.Millisecond)
	
	// Verify
	assert.True(t, mockSubscriber.IsCalled)
	assert.Equal(t, 1, mockSubscriber.TimesCalled)
}
```

### Mock Tracking

Both `OnSend()` and `OnPublish()` return mock objects with tracking capabilities:

- `IsCalled`: Boolean indicating if the mock was called
- `TimesCalled`: Number of times the mock was called

---

## API Reference

### Package: `godiator`

#### `Send[TRequest, TResponse](request TRequest, params ...any) (TResponse, error)`

Executes the handler registered for the given request type.

- **Type Parameters**:
  - `TRequest`: The request type (must match registered request type)
  - `TResponse`: The expected response type
- **Parameters**:
  - `request`: The request object to send
  - `params`: Optional variadic parameters (e.g., `context.Context`)
- **Returns**:
  - `TResponse`: The response from the handler
  - `error`: Error if handler execution fails or handler not found (panics if handler not found)

**Example**:
```go
response, err := godiator.Send[GetUserRequest, GetUserResponse](request)
```

#### `Publish[TRequest](request TRequest, params ...any)`

Notifies all subscribers registered for the given request type. Subscribers are executed asynchronously in goroutines.

- **Type Parameters**:
  - `TRequest`: The request type (must match registered request type)
- **Parameters**:
  - `request`: The request object to publish
  - `params`: Optional variadic parameters
- **Panics**: If no subscribers are registered for the request type

**Example**:
```go
godiator.Publish[UserCreatedEvent](event)
```

### Package: `register`

#### `Handler[TRequest, TResponse](request TRequest, handler core.Handler[TRequest, TResponse])`

Registers a handler for a specific request type.

- **Type Parameters**:
  - `TRequest`: The request type
  - `TResponse`: The response type
- **Parameters**:
  - `request`: Zero-value instance of the request type (used for type registration)
  - `handler`: Handler instance implementing `Handler[TRequest, TResponse]`

**Example**:
```go
register.Handler(GetUserRequest{}, &GetUserHandler{})
```

#### `Subscriber[TRequest](request TRequest, subscribers ...core.Subscriber[TRequest])`

Registers one or more subscribers for a specific request type.

- **Type Parameters**:
  - `TRequest`: The request type
- **Parameters**:
  - `request`: Zero-value instance of the request type
  - `subscribers`: One or more subscriber instances implementing `Subscriber[TRequest]`

**Example**:
```go
register.Subscriber(UserCreatedEvent{}, &EmailSubscriber{}, &LoggingSubscriber{})
```

#### `Pipeline(p *core.Pipeline)`

Registers a pipeline interceptor. Pipelines are executed in registration order before handlers.

- **Parameters**:
  - `p`: Pipeline instance embedding `core.Pipeline`

**Example**:
```go
register.Pipeline(&ValidationPipeline{})
```

### Package: `core`

#### `Handler[TRequest, TResponse]` Interface

Interface that handlers must implement.

```go
type Handler[TRequest any, TResponse any] interface {
	Handle(request TRequest, params ...any) (TResponse, error)
}
```

#### `Subscriber[TRequest]` Interface

Interface that subscribers must implement.

```go
type Subscriber[TRequest any] interface {
	Handle(request TRequest, params ...any)
}
```

#### `Pipeline` Struct

Base struct for pipeline implementations. Must be embedded in your pipeline types.

```go
type Pipeline struct {
	nextPipeline pipeline
}

func (p *Pipeline) Next() pipeline
func (p *Pipeline) SetNext(nextPipeline pipeline)
func (p *Pipeline) Handle(request any, params ...any) (any, error)
```

### Package: `mockiator`

#### `OnSend[TRequest, TResponse](request TRequest, handler func(TRequest, ...any) (TResponse, error)) *mockHandler[TRequest, TResponse]`

Creates a mock handler for testing.

- **Returns**: Mock handler with `IsCalled` and `TimesCalled` properties

**Example**:
```go
mock := mockiator.OnSend[SampleRequest, SampleResponse](SampleRequest{}, handlerFunc)
```

#### `OnPublish[TRequest](request TRequest, handler func(TRequest, ...any)) *mockSubscriber[TRequest]`

Creates a mock subscriber for testing.

- **Returns**: Mock subscriber with `IsCalled` and `TimesCalled` properties

**Example**:
```go
mock := mockiator.OnPublish[UserCreatedEvent](UserCreatedEvent{}, handlerFunc)
```

---

## Best Practices

1. **Registration**: Register all handlers, subscribers, and pipelines during application initialization (e.g., in `init()` functions or setup functions called from `main()`)

2. **Type Safety**: Always use concrete types for requests and responses. Avoid using `interface{}` for better type safety

3. **Error Handling**: Always check errors returned from `Send()` operations

4. **Context**: Use `context.Context` as the first parameter when you need cancellation or timeout support

5. **Subscribers**: Remember that subscribers run asynchronously. Don't rely on execution order between subscribers

6. **Pipelines**: Keep pipeline logic focused and single-purpose. Chain multiple pipelines for complex cross-cutting concerns

7. **Testing**: Use `mockiator` package for unit testing to avoid dependencies on actual handlers

---

## Examples

See the `samples/` directory for complete working examples:
- Handler examples (`samples/handler.go`)
- Subscriber examples (`samples/subscription.go`)
- Test examples (`samples/*_test.go`)

---

## License

See [LICENSE](LICENSE) file for details.

## Contribution

Contributions are welcome! Please feel free to:
- Open issues for bugs or feature requests
- Submit pull requests
- Ask questions or provide feedback

---

## Changelog

### Current Version

- **Type-safe API**: Full support for Go generics
- **Simplified API**: No singleton pattern, direct function calls
- **Pipeline support**: Interceptor pattern for cross-cutting concerns
- **Mocking support**: Built-in testing utilities
- **Async subscriptions**: Fire-and-forget pattern with goroutines
