package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	// get host string for url
	fmt.Print("\n> enter host to connect to (ex: \"localhost:8080\") \n")
	host, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("\n> an error occured while reading input. Please try again", err)
		return
	}

	// get path string for url
	fmt.Print("\n> [optional] enter path (ex: \"/ws-test\") \n")
	path, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("> an error occured while reading input. Please try again", err)
		return
	}

	// trim line break
	host = strings.TrimSuffix(host, "\n")
	path = strings.TrimSuffix(path, "\n")

	// build url string
	u := url.URL{Scheme: "ws", Host: host, Path: path}
	log.Printf("\n> connecting to %s\n", u.String())

	// init ws conn
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("\n> dial: %v\n", err)
	}
	defer c.Close()

	// allow user to send messages
	sendCh := make(chan string)
	fmt.Print("\n> send: ")
	go func() {
		for {
			msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				fmt.Println("\n> an error occured while reading input. Please try again", err)
				return
			}
			msg = strings.TrimSuffix(msg, "\n")
			sendCh <- msg
		}
	}()

	// listen for incoming msg
	done := make(chan struct{})
	receiveCh := make(chan string)
	go func() {
		defer close(done)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("\n> read:", err)
				return
			}
			receiveCh <- string(msg)
		}
	}()

	// listen to send channel to write new messages, listen to done channel
	for {
		select {
		case <-done:
			return
		case msg := <-receiveCh:
			fmt.Printf("\n> received:\n%s\n\n> reply: ", msg)
		case msg := <-sendCh:
			err := c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("\n> write:", err)
				return
			}
		}
	}
}
