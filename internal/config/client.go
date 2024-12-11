package config

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	stringTermination = '\n'
	unwantedChars     = "\r\n"
)

const (
	nick = "/nick"
	join = "/join"
	room = "/room"
	msg  = "/msg"
	quit = "/quit"
)

type client struct {
	conn     net.Conn
	nickName string
	chatRoom *rooms
	commands chan<- command
}

func (c *client) ReadInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			log.Printf("Failed to read message from connection.Err: %+v", err)
			return
		}
		// Remove unwanted characters from string
		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])
		log.Printf("Client has entered command: %s\n", cmd)
		switch cmd {
		case "/nick":
			c.commands <- command{
				commandID: CMD_NICK,
				listener:  c,
				args:      args,
			}
		case "/join":
			c.commands <- command{
				commandID: CMD_JOIN,
				listener:  c,
				args:      args,
			}
		case "/room":
			c.commands <- command{
				commandID: CMD_ROOM,
				listener:  c,
			}
		case "/msg":
			c.commands <- command{
				commandID: CMD_MSG,
				listener:  c,
				args:      args,
			}
		case "/quit":
			c.commands <- command{
				commandID: CMD_QUIT,
				listener:  c,
			}
		default:
			c.SendErr(fmt.Errorf("Unknown command name found: %s", args[0]))
		}
	}
}

func (c *client) SendErr(err error) {
	c.conn.Write([]byte(fmt.Sprintf("Error: %+v\n", err.Error())))
}

func (c *client) SendMsg(msg string) {
	c.conn.Write([]byte(fmt.Sprintf("> %+v\n", msg)))
}
