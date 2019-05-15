package gate

import (
	"zlab/library/nano/session"
)

type Player struct {
	uid     string
	session *session.Session //用户数据
}
