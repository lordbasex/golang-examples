package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const timeout = 100 * time.Second
const maxTokens = 4000

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(new(NullWriter))
	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		log.Fatal("Missing API_KEY")
	}

	client := gpt3.NewClient(apiKey)

	ctx := context.Background()

	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)

			for {
				fmt.Print(">> ")

				if !scanner.Scan() {
					break
				}

				question := scanner.Text()

				if question == "quit" || question == "exit" {
					return
				}

				response, err := getResponse(ctx, client, question)
				if err != nil {
					log.Println(err)
					continue
				}

				color.Green(response)
				if runtime.GOOS == "darwin" {
					if err := executeSayCommand(response); err != nil {
						log.Fatal(err)
					}
				}

				fmt.Println()
			}
		},
	}

	rootCmd.Execute()
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

// executeSayCommand ejecuta el comando 'say' de macOS con el texto dado.
func executeSayCommand(text string) error {
	cmd := exec.Command("say", text)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start say command: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait for say command: %w", err)
	}

	return nil
}
