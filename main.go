package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type messengerInput struct {
	Entry []struct {
		Time      uint64 `json:"time,omitempty"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender,omitempty"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient,omitempty"`
			Timestamp uint64 `json:"timestamp,omitempty"`
			Message   *struct {
				Mid  string `json:"mid,omitempty"`
				Seq  uint64 `json:"seq,omitempty"`
				Text string `json:"text"`
			} `json:"message,omitempty"`
		} `json:"messaging"`
	}
}

func main() {
	isDevEnv := os.Getenv("GO_ENV") == "development"
	if isDevEnv {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/webhook", messengerVerify)

	whereToListen := ":" + os.Getenv("PORT")
	if isDevEnv {
		whereToListen = "localhost" + whereToListen
	}
	fmt.Println("Starting Server on " + whereToListen)
	log.Fatal(http.ListenAndServe(whereToListen, nil))
}

func messengerVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		challenge := r.URL.Query().Get("hub.challenge")
		verifyToken := r.URL.Query().Get("hub.verify_token")

		if len(verifyToken) > 0 && len(challenge) > 0 && verifyToken == os.Getenv("FB_VERIFY_TOKEN") {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, challenge)
			return
		}
	} else if r.Method == "POST" {
		defer r.Body.Close()

		input := new(messengerInput)
		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			log.Fatal(err)
		}

		log.Println("got message:", input.Entry[0].Messaging[0].Message.Text)

		reply := input.Entry[0].Messaging[0]
		reply.Sender, reply.Recipient = reply.Recipient, reply.Sender

		reply.Message.Text = "I got your message:" + input.Entry[0].Messaging[0].Message.Text
		reply.Message.Seq = 0
		reply.Message.Mid = ""

		b, _ := json.Marshal(reply)
		http.Post("https://graph.facebook.com/v2.6/me/messages?access_token="+os.Getenv("FB_PAGE_TOKEN"), "application/json", bytes.NewReader(b))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(400)
	fmt.Fprintf(w, "Bad Request")
}
