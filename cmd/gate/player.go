package main

import (
	"zlab/library/nano/session"
)

//Player ..
type Player struct {
	uid     string
	session *session.Session
	block   *Block
}

//NewPlayer ..
func NewPlayer(sess *session.Session, uid string) *Player {
	p := &Player{
		uid: uid,
	}
	p.bindSession(sess)
	return p
}

func (p *Player) bindSession(sess *session.Session) {
	p.session = sess
	p.session.Set(kcurPlayer, p)
}

func (p *Player) removeSession() {
	p.session.Remove(kcurPlayer)
	p.session = nil
}
