package ai

import (
	"bytes"
	"encoding/json"
	"github.com/danilobandeira29/talktome/chat"
	"io"
	"log"
	"net/http"
)

const openAIKey = ""

type DetectIntentResponse struct {
	Intention string `json:"intention"`
}

type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type requestBody struct {
	Model    string       `json:"model"`
	Messages []GPTMessage `json:"messages"`
}

type responseBody struct {
	Choices []struct {
		Message map[string]interface{} `json:"message"`
	} `json:"choices"`
}

func DetectIntent(history []chat.Message) (*DetectIntentResponse, error) {
	url := "https://api.openai.com/v1/chat/completions"
	var mgs []GPTMessage
	for _, m := range history {
		t := m["type"].(string)
		if t == "attendant" {
			t = "assistant"
		} else {
			t = "user"
		}
		mgs = append(mgs, GPTMessage{
			Role:    t,
			Content: m["message"].(string),
		})
	}
	body := requestBody{
		Model: "gpt-4",
		// TODO: this prompt doesn't works sometimes
		Messages: append([]GPTMessage{{Role: "system", Content: `Você é um classificador de intenção. Analise a última mensagem do cliente e determine se ele deseja encerrar a conversa.  
Se a última mensagem do usuário indicar um encerramento (por exemplo: "não, só isso", "obrigado", "até mais", "tchau"), retorne exatamente "finish_chat".  
Caso contrário, retorne exatamente "continue".`}}, mgs...),
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+openAIKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		log.Println(string(res))
		return nil, err
	}
	var openaiResp responseBody
	err = json.Unmarshal(res, &openaiResp)
	if err != nil {
		return nil, err
	}
	answer := openaiResp.Choices[0].Message["content"]
	if answer == "finish_chat" {
		return &DetectIntentResponse{Intention: "finish_chat"}, nil
	}
	return &DetectIntentResponse{Intention: "continue"}, nil
}
