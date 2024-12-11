package config

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	rooms    map[string]*rooms
	commands chan command
}

func NewServer() *Server {
	return &Server{
		rooms:    make(map[string]*rooms),
		commands: make(chan command),
	}
}

func (s *Server) Run() {
	for cmd := range s.commands {
		switch cmd.commandID {
		case CMD_NICK:
			s.nick(cmd.listener, cmd.args)
		case CMD_JOIN:
			s.join(cmd.listener, cmd.args)
		case CMD_ROOM:
			s.listRoom(cmd.listener)
		case CMD_MSG:
			s.msg(cmd.listener, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.listener)
		}
	}
}

func (s *Server) NewClient(conn net.Conn) *client {
	log.Printf("New Client has joining: %s\n", conn.RemoteAddr().String())
	return &client{
		conn:     conn,
		nickName: "anonymous",
		commands: s.commands,
	}
}

func (s *Server) nick(client *client, args []string) {
	if len(args) < 2 {
		client.SendErr(fmt.Errorf("Nick name is required. Please enter: /nick NAME"))
		return
	}
	client.nickName = args[1]
	client.SendMsg("All right, I will call you " + client.nickName)
}

func (s *Server) join(c *client, args []string) {
	if len(args) < 2 {
		c.SendErr(fmt.Errorf("Room name is required. Please enter: /join ROOMNAME"))
		return
	}
	roomName := args[1]
	room, ok := s.rooms[roomName]
	if !ok {
		room = &rooms{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = room
	}
	room.members[c.conn.RemoteAddr()] = c
	s.quitCurrentRoom(c)
	c.chatRoom = room
	room.BroadcastMsg(c, fmt.Sprintf("%s successfully joined the room %s", c.nickName, c.chatRoom.name))
	c.SendMsg(fmt.Sprintf("welcome to the room %s", roomName))
}

func (s *Server) listRoom(client *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	client.SendMsg(fmt.Sprintf("Available Rooms: %s", strings.Join(rooms, ", ")))
}

func (s *Server) msg(client *client, args []string) {
	if len(args) < 2 {
		client.SendErr(fmt.Errorf("Message is required. Please enter: /msg MESSAGE"))
		return
	}
	msg := strings.Join(args[1:], " ")
	client.chatRoom.BroadcastMsg(client, fmt.Sprintf("%s: %s", client.nickName, msg))
}

func (s *Server) quit(client *client) {
	log.Printf("client has left the room: %s", client.conn.RemoteAddr().String())
	s.quitCurrentRoom(client)
	client.SendMsg("You left the room successfully")
	client.conn.Close()
}

func (s *Server) quitCurrentRoom(client *client) {
	if client.chatRoom != nil {
		oldRoom := s.rooms[client.chatRoom.name]
		delete(s.rooms[client.chatRoom.name].members, client.conn.RemoteAddr())
		oldRoom.BroadcastMsg(client, fmt.Sprintf("%s has left the room", client.nickName))
	}
}
