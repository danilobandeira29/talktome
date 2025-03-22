package activity

import (
	"github.com/danilobandeira29/talktome/chat"
	"testing"
	"time"
)

func TestChatAIActivity(t *testing.T) {
	c := &ChatAIActivity{
		chat:      chat.New("1234"),
		activity:  New("1", "CHAT_AI", "2", 1),
		StudentID: "Danilo Bandeira",
		CreatedAt: time.Now(),
		PromptURL: "Aja como um professor que dara uma aula sobre o assunto XYZ...",
		MaxStep:   2,
	}
	messages := []map[string]interface{}{
		{
			"type":             "tutor",
			"message":          "Ola, Danilo! Agora você irá aprender sobre o assunto XYZ! Está animado e preparado?",
			"lastStudentGrade": "high",
		},
		{
			"type":    "student",
			"message": "Estou um pouco nervoso, mas estou animado! Vamos!",
		},
		{
			"type":             "tutor",
			"message":          "É normal está nervoso, mas eu posso te ajudar. Iremos com calma, em casos de dúvida, basta me falar, ok? Vamos começar...",
			"lastStudentGrade": "high",
		},
	}
	for _, m := range messages {
		if err := c.SendMessage(m); err != nil {
			t.Errorf("error sending tutor message: %v\n", err)
			return
		}
	}
	if finished := c.Finished(); !finished {
		t.Errorf("error chat ai activity should be finished")
		return
	}
}
