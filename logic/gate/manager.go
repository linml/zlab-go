package gate

import (
	"time"

	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
)

const kickResetBacklog = 8

type Manager struct {
	component.Base
	group      *nano.Group        // 广播channel
	players    map[string]*Player // 所有的玩家
	chKick     chan string        // 退出队列
	chReset    chan string        // 重置队列
	chRecharge chan RechargeInfo  // 充值信息
}
type RechargeInfo struct {
	Uid  string // 用户ID
	Coin int64  // 房卡数量
}

var defaultManager = NewManager()

func NewManager() *Manager {
	return &Manager{
		group:      nano.NewGroup("_SYSTEM_MESSAGE_BROADCAST"),
		players:    map[string]*Player{},
		chKick:     make(chan string, kickResetBacklog),
		chReset:    make(chan string, kickResetBacklog),
		chRecharge: make(chan RechargeInfo, 32),
	}
}

func (m *Manager) AfterInit() {
	session.Lifetime.OnClosed(func(s *session.Session) {
		m.group.Leave(s)
	})

	nano.NewTimer(time.Second, func() {
	ctrl:
		for {
			select {
			case uid := <-m.chKick:
				p, ok := m.player(uid)
				if !ok || p.session == nil {
					logger.Errorf("玩家%s不在线", uid)
				}
				p.session.Close()
				logger.Infof("踢出玩家, UID=%d", uid)
			case uid := <-m.chReset:
				p, ok := defaultManager.player(uid)
				if !ok {
					return
				}
				if p.session != nil {
					logger.Errorf("玩家正在游戏中，不能重置: %d", uid)
					return
				}
				//p.desk = nil
				logger.Infof("重置玩家, UID=%d", uid)
			case rc := <-m.chRecharge:
				player, ok := m.player(rc.Uid)
				// 如果玩家在线
				if s := player.session; ok && s != nil {
					//s.Push("onCoinChange", &protocol.CoinChangeInformation{Coin: ri.Coin})
				}
			default:
				break ctrl
			}

		}
	})
}

func (m *Manager) player(uid string) (*Player, bool) {
	p, ok := m.players[uid]
	return p, ok
}
