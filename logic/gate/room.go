package gate

import "zlab/library/nano"

type Room struct {
	RoomID string
	group  *nano.Group // 组播通道
}
