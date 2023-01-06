package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const timeout = 100 * time.Second
const maxTokens = 4000

func main() {

	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		log.Fatal("Missing API_KEY")
	}

	telegramApiKey, ok := os.LookupEnv("TELEGRAM_API_KEY")
	if !ok {
		log.Fatal("Missing TELEGRAM_API_KEY")
	}

	client := gpt3.NewClient(apiKey)

	ctx := context.Background()

	bot, err := tgbotapi.NewBotAPI(telegramApiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			//log.Print(update.Message.Text)
			color.Yellow(update.Message.Text)

			response, err := getResponse(ctx, client, update.Message.Text)
			if err != nil {
				log.Println(err)
				continue
			}

			color.Green(response)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

// getResponse obtiene la respuesta del cliente GPT3 para una pregunta dada.
func getResponse(ctx context.Context, client gpt3.Client, question string) (response string, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	sb := strings.Builder{}

	err = client.CompletionStreamWithEngine(
		ctx,
		gpt3.TextDavinci003Engine,
		gpt3.CompletionRequest{
			Prompt: []string{
				question,
			},
			MaxTokens:   gpt3.IntPtr(maxTokens),
			Temperature: gpt3.Float32Ptr(0),
		},
		func(resp *gpt3.CompletionResponse) {
			text := resp.Choices[0].Text

			sb.WriteString(text)
		},
	)
	if err != nil {
		return "", err
	}

	response = sb.String()
	response = strings.TrimLeft(response, "\n")

	return response, nil
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) {
	return 0, nil
}
