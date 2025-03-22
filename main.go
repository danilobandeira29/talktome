package main

import (
	"bufio"
	"fmt"
	"github.com/danilobandeira29/talktome/ai"
	"github.com/danilobandeira29/talktome/attendant"
	"github.com/danilobandeira29/talktome/chat"
	"log"
	"os"
)

func main() {
	c := attendant.New("1", "23", "12", "Cliente não consegue cancelar assinatura pela plataforma")
	scanner := bufio.NewScanner(os.Stdin)
	var i int
	log.Println("chat started")
	for {
		// TODO: works, it's just a poc. take it easy
		var t string
		isClient := i%2 == 0
		if isClient {
			t = "client"
		} else {
			t = "attendant"
		}
		fmt.Printf("%s: ", t)
		scanner.Scan()
		clientMessage := chat.Message{
			"type":    t,
			"message": scanner.Text(),
		}
		if isClient {
			clientMessage["finish_chat"] = false
			intention, err := ai.DetectIntent(c.History())
			if err != nil {
				log.Printf("error dectecting intent: %v\n", err)
				return
			}
			if intention.Intention == "finish_chat" {
				clientMessage["finish_chat"] = true
			}
		}
		if errSend := c.SendMessage(clientMessage); errSend != nil {
			if _, errF := fmt.Fprintf(os.Stdout, "%s\n", errSend.Error()); errF != nil {
				log.Printf("error sending message to writer: %v\n", errF)
				return
			}
			return
		}
		if c.Finished() {
			break
		}
		i++
	}
	log.Println("chat finished")
}
