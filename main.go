package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	gpt3 "github.com/sashabaranov/go-gpt3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func Chat(apiKey string, query string) (error, string) {
	ctx := context.Background()
	c := gpt3.NewClient(apiKey)

	req := gpt3.CompletionRequest{
		Model:            gpt3.GPT3TextDavinci003,
		MaxTokens:        3000,
		Temperature:      0.3,
		TopP:             1.0,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.0,
		Prompt:           query,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return err, ""
	}
	rr := resp.Choices[0].Text
	return nil, rr
}

func GetEventHandler(client *whatsmeow.Client) func(interface{}) {
	return func(evt interface{}) {
		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			log.Fatalln("Missing API KEY")
		}
		switch v := evt.(type) {
		case *events.Message:
			if !v.Info.IsFromMe && !v.Info.IsGroup {
				if strings.Contains(v.Message.GetConversation(), "-ask") {
					gagal, jwban := Chat(apiKey, strings.Replace(v.Message.GetConversation(), "-ask", "", -1))
					if gagal != nil {
						client.SendMessage(context.Background(), v.Info.Sender, &waProto.Message{
							Conversation: proto.String(gagal.Error()),
						})
					} else {
						client.SendMessage(context.Background(), v.Info.Sender, &waProto.Message{
							Conversation: proto.String(jwban),
						})
					}

				} else {
					client.SendMessage(context.Background(), v.Info.Sender, &waProto.Message{
						Conversation: proto.String("Jangan lupa pakai trigger -ask"),
					})
				}
			}
		}
	}
}

func main() {
	godotenv.Load()

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(GetEventHandler(client))

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				// fmt.Println("QR code:", evt.Code)
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()

}
