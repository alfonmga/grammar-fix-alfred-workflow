package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/sashabaranov/go-openai"
)

var wf *aw.Workflow
var openaiClient *openai.Client

func init() {
	if os.Getenv("LOG_LEVEL") == "debug" {
		slog.SetDefault(
			slog.New(
				slog.NewTextHandler(
					os.Stderr,
					&slog.HandlerOptions{Level: slog.LevelDebug},
				),
			),
		)
	}
	openaiClient = openai.NewClient(os.Getenv("OPENAI_KEY"))
	slog.Debug("OpenAI client initialized", "KEY", os.Getenv("OPENAI_KEY"))
}

func main() {
	wf = aw.New()
	wf.Run(run)
}
func run() {
	userSelectedText := wf.Args()[0]
	if userSelectedText == "" {
		return // skip if no text is selected by user
	}
	slog.Debug("User selected text", "text", userSelectedText)

	prompt := fmt.Sprintf(
		`Correct the grammar of the following text:\n\n%s`,
		userSelectedText,
	)
	slog.Debug("Generated GPT prompt", "prompt", prompt)

	ctx, cancel := context.WithTimeout(context.Background(), 11*time.Second)
	defer cancel()
	promptResp, err := openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		panic(fmt.Sprintf("Error while calling OpenAI API: %s", err))
	}
	promptRespContent := promptResp.Choices[0].Message.Content
	slog.Debug("OpenAI prompt response", "response", promptRespContent)

	fmt.Println(promptRespContent)
}
