package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/codecrafters-io/claude-code-starter-go/internal/pkg/assistant"
	"github.com/sirupsen/logrus"
)

func main() {
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
		logrus.Fatal(err)
	}

	for {
		prompt := assistant.Prompt()
		_, err := assistant.Process(ctx, prompt)
		if err != nil {
			logrus.Error(err)
		}
	}
}
