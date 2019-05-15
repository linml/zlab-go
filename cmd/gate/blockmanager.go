package main

import (
	"fmt"

	"zlab/library/nano/component"
	"zlab/library/nano/session"
)

//BlockManager ..
type BlockManager struct {
	component.Base
	//房间数据
	blocks map[string]*Block // 房间数据
}

var defaultBlockManager = NewBlockManager()

//NewBlockManager ..
func NewBlockManager() *BlockManager {
	return &BlockManager{
		blocks: map[string]*Block{},
	}
}

//AfterInit ..
func (manager *BlockManager) AfterInit() {
	session.Lifetime.OnClosed(func(s *session.Session) {
		// Fixed: 玩家WIFI切换到4G网络不断开, 重连时，将UID设置为illegalSessionUid
		if s.UID() != "" {
			manager.onPlayerDisconnect(s)
		}
	})
}

func (manager *BlockManager) onPlayerDisconnect(s *session.Session) error {
	uid := s.UID()
	p, ok := defaultManager.players[uid]
	if !ok {
		return fmt.Errorf("player %s is not online", uid)
	}

	// 移除session
	p.removeSession()

	defaultManager.offline(uid)

	d := p.block
	d.onPlayerExit(s)
	return nil
}
