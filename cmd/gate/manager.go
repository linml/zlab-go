package main

import (
	"zlab/library/nano"
	"zlab/library/nano/component"
	"zlab/library/nano/session"
)

//Manager ..
type Manager struct {
	component.Base
	group   *nano.Group        // 广播channel
	players map[string]*Player // 所有的玩家
}

var defaultManager = NewManager()

//NewManager ..
func NewManager() *Manager {
	return &Manager{
		group:   nano.NewGroup("_SYSTEM_MESSAGE_BROADCAST"),
		players: map[string]*Player{},
	}
}

//AfterInit ..
func (m *Manager) AfterInit() {
	session.Lifetime.OnClosed(func(s *session.Session) {
		m.group.Leave(s)
	})
}

func (m *Manager) player(id string) (*Player, bool) {
	p, ok := m.players[id]
	return p, ok
}

func (m *Manager) setPlayer(id string, player *Player) {
	if _, ok := m.players[id]; !ok {
		m.players[id] = player
	}
}

func (m *Manager) sessionCount() int {
	return len(m.players)
}

func (m *Manager) offline(id string) {
	delete(m.players, id)
}
