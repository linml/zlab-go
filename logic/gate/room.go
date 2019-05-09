package gate

import "github.com/lonng/nano"

type Room struct {
	RoomID string
	group  *nano.Group // 组播通道
}
