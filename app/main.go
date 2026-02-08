package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/codecrafters-io/claude-code-starter-go/internal/assistant"
	"github.com/codecrafters-io/claude-code-starter-go/internal/logger"
)

func main() {
	logger.Setup()
	logger.PrintBanner()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		cancel()
	}()

	assistant, err := assistant.New()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	for {
		logger.Prompt()
		prompt, err := assistant.Prompt(ctx)
		if err != nil {
			if err == context.Canceled {
				fmt.Println("\nGoodbye!")
				return
			}
			logger.Error(err)
			continue
		}
		err = assistant.Process(ctx, prompt)
		if err != nil {
			logger.Error(err)
		}
	}
}
