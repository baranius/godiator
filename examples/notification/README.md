## Notifications
Creating a notification handler in godiatr requires a method named as **Handle**. 

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
