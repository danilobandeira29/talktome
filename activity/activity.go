package activity

import (
	"fmt"
	"github.com/danilobandeira29/talktome/chat"
	"time"
)

type Activity struct {
	ID     string
	Type   string
	StepID string
	Order  int8
}

func New(id, t, stepID string, order int8) *Activity {
	return &Activity{
		ID:     id,
		Type:   t,
		StepID: stepID,
		Order:  order,
	}
}

type ChatAIActivity struct {
	chat       *chat.Chat
	activity   *Activity
	StudentID  string
	CreatedAt  time.Time
	finishedAt *time.Time
	PromptURL  string
	Step       uint8
	MaxStep    uint8
}

func (c *ChatAIActivity) SendMessage(m chat.Message) error {
	whoSendMessage, ok := m["type"]
	if !ok {
		return fmt.Errorf("chat ai activity: property type is mandatory to send a message")
	}
	isFirstMessage := len(c.chat.History()) == 0
	if isFirstMessage && whoSendMessage != "tutor" {
		return fmt.Errorf("chat ai activity: only the tutor can start a chat")
	}
	if c.finishedAt != nil {
		return fmt.Errorf("chat ai activity: you cannot send message because this chat has been finished")
	}
	lastStudentGrade, ok := m["lastStudentGrade"]
	if whoSendMessage == "tutor" && !ok {
		return fmt.Errorf("chat ai activity: to send a message the tutor needs to provide the last student grade")
	}
	if err := c.chat.SendMessage(m); err != nil {
		return err
	}
	if whoSendMessage == "tutor" && lastStudentGrade == "high" {
		c.Step++
	}
	chatAIActivityShouldFinished := whoSendMessage == "tutor" && c.Step >= c.MaxStep
	if chatAIActivityShouldFinished {
		t := time.Now()
		c.finishedAt = &t
	}
	return nil
}

func (c *ChatAIActivity) Finished() bool {
	if c.finishedAt != nil {
		return true
	}
	return false
}
