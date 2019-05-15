package main

import (
	"zlab/library/nano"
	"zlab/library/nano/session"
)

//Block 代表一个分组 房间 / 俱乐部等
type Block struct {
	blockid string
	group   *nano.Group // 组播通道
}

//NewBlock ..
func NewBlock(blockid string) *Block {
	return &Block{
		blockid: blockid,
		group:   nano.NewGroup(blockid),
	}
}

func (b *Block) join(p *Player) {
	b.group.Add(p.session)
}

func (b *Block) leave(p *Player) {
	b.group.Leave(p.session)
}

func (b *Block) onPlayerExit(s *session.Session) {
	b.group.Leave(s)
}
