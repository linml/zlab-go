package gate

import "zlab/library/nano/component"

type RoomManager struct {
	component.Base
	rooms map[string]*Room
}

var defaultRoomManager = NewRoomManager()

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: map[string]*Room{},
	}
}
