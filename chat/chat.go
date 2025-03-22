package chat

import "fmt"

type Message map[string]interface{}

type Chat struct {
	ID       string
	Messages []Message
}

func New(id string) *Chat {
	return &Chat{
		ID:       id,
		Messages: []Message{},
	}
}

func (c *Chat) SendMessage(m Message) error {
	message, ok := m["message"]
	if !ok {
		return fmt.Errorf("chat send message: property message is mandatory to send a message")
	}
	if message == "" {
		return fmt.Errorf("chat send message: cannot send an empty message")
	}
	c.Messages = append(c.Messages, m)
	return nil
}

func (c *Chat) History() []Message {
	return c.Messages
}
