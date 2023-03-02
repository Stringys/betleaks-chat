package models

import "fmt"

type Chat struct {
	Users    map[string]*User
	Messages chan *Message
	Join     chan *User
	Leave    chan *User
	ID       int
}

func NewChat(chatID int) *Chat {
	return &Chat{
		Users:    make(map[string]*User),
		Messages: make(chan *Message),
		Join:     make(chan *User),
		Leave:    make(chan *User),
		ID:       chatID,
	}
}

func (c *Chat) Run() {
	for {
		select {
		case user := <-c.Join:
			c.add(user)
		case message := <-c.Messages:
			c.broadcast(message)
		case user := <-c.Leave:
			c.disconnect(user)
		}
	}
}

func (c *Chat) broadcast(msg *Message) {
	for _, user := range c.Users {
		user.WriteMessage(msg)
	}
}

func (c *Chat) add(user *User) {
	if _, found := c.Users[user.Username]; !found {
		c.Users[user.Username] = user

		body := fmt.Sprintf("%s join the chat", user.Username)
		c.broadcast(NewMessage(body, "Server"))
	}
}

func (c *Chat) disconnect(user *User) {
	if _, found := c.Users[user.Username]; found {
		defer user.Conn.Close()
		delete(c.Users, user.Username)

		body := fmt.Sprintf("%s left the chat", user.Username)
		c.broadcast(NewMessage(body, "Server"))
	}
}
