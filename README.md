# Godiator

[![CI](https://github.com/baranius/godiator/actions/workflows/main.yml/badge.svg)](https://github.com/baranius/godiator/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/baranius/godiator/branch/master/graph/badge.svg)](https://codecov.io/gh/baranius/godiator)
[![Go Report Card](https://goreportcard.com/badge/github.com/baranius/godiator)](https://goreportcard.com/report/github.com/baranius/godiator)
[![GoDoc](https://godoc.org/github.com/baranius/godiator?status.svg)](https://godoc.org/github.com/baranius/godiator)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/baranius/godiator)](go.mod)

A lightweight, generic-based mediator pattern implementation for Go. Godiator enables loosely coupled communication between components using **Request/Response** (unicast) and **Publish/Subscribe** (multicast) patterns.

## Features

- **Simple API**: Minimal boilerplate using Go 1.18+ Generics.
- **Request/Response**: Direct synchronous communication between components.
- **Publish/Subscribe**: Asynchronous event broadcasting.
- **Pipelines**: Middleware support for cross-cutting concerns (logging, validation, etc.).
- **Testing Friendly**: Built-in `mockiator` package for easy mocking in unit tests.
- **Thread-Safe**: Safe for concurrent use.
- **Zero Dependencies**: Lightweight and focused.

## Installation

```bash
go get github.com/baranius/godiator
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/baranius/godiator"
)

// 1. Define Request and Response
type GetUserRequest struct {
	ID int
}

type GetUserResponse struct {
	Name string
}

// 2. Define a Handler
type GetUserHandler struct{}

func (h *GetUserHandler) Handle(req GetUserRequest, params ...any) (GetUserResponse, error) {
	// Simulate logic
	if req.ID == 0 {
		return GetUserResponse{}, fmt.Errorf("invalid user ID")
	}
	return GetUserResponse{Name: "John Doe"}, nil
}

func main() {
	// 3. Register the Handler
	godiator.RegisterHandler[GetUserRequest, GetUserResponse](&GetUserHandler{})

	// 4. Send a Request
	req := GetUserRequest{ID: 1}
	resp, err := godiator.Send[GetUserRequest, GetUserResponse](req)
	
	if err != nil {
		panic(err)
	}

	fmt.Printf("User Name: %s\n", resp.Name)
}
```

## Usage

### Request / Response

Handlers process a specific request type and return a response. Only one handler can be registered per request type.

#### defining a handler
Implement the `Handle` method matching the `Handler[TRequest, TResponse]` interface:

```go
type MyHandler struct{}

func (h *MyHandler) Handle(request MyRequest, params ...any) (MyResponse, error) {
    // Logic here...
    return MyResponse{...}, nil
}
```

#### Registration and Sending
Use generic functions to register handlers and send requests:

```go
// Register (usually in main or init)
godiator.RegisterHandler[MyRequest, MyResponse](&MyHandler{})

// Send
response, err := godiator.Send[MyRequest, MyResponse](MyRequest{Id: 10})
```

### Publish / Subscribe

Subscribers listen for specific events. Multiple subscribers can be registered for the same request/event type. They are executed asynchronously (fire-and-forget).

#### Defining a Subscriber
Implement the `Handle` method matching the `Subscriber[TRequest]` interface:

```go
type EmailSubscriber struct{}

func (s *EmailSubscriber) Handle(event UserCreatedEvent, params ...any) {
    // Send email logic...
}
```

#### Registration and Publishing

```go
// Register generic subscribers
godiator.RegisterSubscriber[UserCreatedEvent](&EmailSubscriber{})
godiator.RegisterSubscriber[UserCreatedEvent](&AnalyticsSubscriber{})

// Publish event
godiator.Publish(UserCreatedEvent{UserID: 123})
```

### Pipelines (Middleware)

Pipelines intercept requests before they reach the handler. They are useful for logging, authentication, validation, etc. Pipelines are executed in **reverse order** of registration (LIFO) - the last registered pipeline runs first (wrapping others).

To create a pipeline, embed `pipeline.BasePipeline` and override `Handle`.

```go
import (
    "github.com/baranius/godiator"
    "github.com/baranius/godiator/pipeline"
)

type LoggingPipeline struct {
    pipeline.BasePipeline
}

func (p *LoggingPipeline) Handle(request any, params ...any) (any, error) {
    fmt.Println("Before Handler")
    
    // Call next in chain
    response, err := p.Next().Handle(request, params...)
    
    fmt.Println("After Handler")
    return response, err
}

// Register (applies to ALL requests)
godiator.RegisterPipeline(&LoggingPipeline{})
```

### Mocking for Tests

The `mockiator` package provides utilities to mock handlers and subscribers in unit tests without implementing the full interfaces manually.

```go
import (
    "testing"
    "github.com/baranius/godiator/mockiator"
    "github.com/stretchr/testify/assert"
)

func TestMyLogic(t *testing.T) {
    // Mock a response
    mock := mockiator.OnSend[MyRequest, MyResponse](func(req MyRequest, params ...any) (MyResponse, error) {
        return MyResponse{Result: "Mocked"}, nil
    })

    // Run your code that calls godiator.Send...
    // ...

    // Assertions
    assert.True(t, mock.IsCalled)
    assert.Equal(t, 1, mock.TimesCalled)
}
```

## Contributing

Contributions are welcome!
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Verify tests pass (`go test ./...`)
5. Push to the branch
6. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
