package gate

import (
	"github.com/lonng/nano/session"
)

type Player struct {
	uid     string
	session *session.Session //用户数据
}
