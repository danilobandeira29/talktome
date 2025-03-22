package attendant

import (
	"testing"
)

func TestChatAttendant(t *testing.T) {
	c := New("1234", "12", "11", "Cliente não consegue cancelar a assinatura pela plataforma")
	messages := []map[string]interface{}{
		{
			"type":        "client",
			"message":     "Estou tentando cancelar minha assinatura, mas não consigo",
			"finish_chat": false,
		},
		{
			"type":    "attendant",
			"message": "Certo. Essa conversa está sendo gravada...",
		},
		{
			"type":    "attendant",
			"message": "Qual o motivo do cancelamento?",
		},
		{
			"type":        "client",
			"message":     "Me senti enganado com a plataforma",
			"finish_chat": false,
		},
		{
			"type":    "attendant",
			"message": "Entendi. Acabei de cancelar sua assinatura e o estorno será feito em até 7 dias úteis. Consigo ajudar você em algo mais?",
		},
		{
			"type":        "client",
			"message":     "Já? Obrigado, isso já é suficiente",
			"finish_chat": true,
		},
	}
	for _, m := range messages {
		if err := c.SendMessage(m); err != nil {
			t.Errorf("error message: %v\n", err)
			return
		}
	}
	if finished := c.Finished(); !finished {
		t.Errorf("error chat attendant should be finished")
		return
	}
}
