package config

import "net"

type commandID int

const (
	CMD_NICK = iota
	CMD_JOIN
	CMD_ROOM
	CMD_MSG
	CMD_QUIT
)

type command struct {
	commandID commandID
	listener  *client
	args      []string
}

type rooms struct {
	name    string
	members map[net.Addr]*client
}

func (r *rooms) BroadcastMsg(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.SendMsg(msg)
		}
	}
}

func HaltIfEmpty(value int) bool {
	return value < 1
}
