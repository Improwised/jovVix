package structs

// All response structs
// Response struct have Res prefix

type SocketResponseFormat struct {
	EventName string      `json:"event"`
	Data      interface{} `json:"data"`
}
