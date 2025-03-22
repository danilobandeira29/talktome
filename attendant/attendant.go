package attendant

import (
	"fmt"
	"github.com/danilobandeira29/talktome/chat"
	"time"
)

type ChatAttendant struct {
	chat        *chat.Chat
	AttendantID string
	ClientID    string
	finishedAt  *time.Time
	Problem     string
}

func New(chatID, attendantID, clientID, problem string) *ChatAttendant {
	return &ChatAttendant{
		chat:        chat.New(chatID),
		AttendantID: attendantID,
		ClientID:    clientID,
		finishedAt:  nil,
		Problem:     problem,
	}
}

func (c *ChatAttendant) SendMessage(m chat.Message) error {
	whoSendMessage, ok := m["type"]
	if !ok {
		return fmt.Errorf("chat attendant: property type is mandatory to send a message")
	}
	isFirstMessage := len(c.chat.History()) == 0
	if isFirstMessage && whoSendMessage != "client" {
		return fmt.Errorf("chat attendant: only the client can start a chat")
	}
	if c.Finished() {
		return fmt.Errorf("chat attendant: you cannot send message because this chat has been finished")
	}
	if err := c.chat.SendMessage(m); err != nil {
		return fmt.Errorf("chat attendant: %v", err)
	}
	intent, ok := m["finish_chat"].(bool)
	if whoSendMessage == "client" && !ok {
		return fmt.Errorf("chat attendant finish: cannot finish chat because property finish_chat is mandatory as boolean")
	}
	if whoSendMessage == "client" && intent {
		t := time.Now()
		c.finishedAt = &t
	}
	return nil
}

func (c *ChatAttendant) History() []chat.Message {
	return c.chat.Messages
}

func (c *ChatAttendant) Finished() bool {
	if c.finishedAt != nil {
		return true
	}
	return false
}
