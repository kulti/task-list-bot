package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"go.uber.org/zap"

	"github.com/kulti/task-list-bot/internal/processor"
)

type botFlags struct {
	Token   string `env:"BOT_TOKEN,required"`
	OwnerID int    `env:"BOT_OWNER_ID,required"`
	Debug   bool   `env:"BOT_DEBUG"`
}

func main() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	zapLogger, err := config.Build()
	if err != nil {
		fmt.Println("failed to init logger: ", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(zapLogger)

	var botFlags botFlags
	if err := env.Parse(&botFlags); err != nil {
		fmt.Println("failed to parse bot flags: ", err)
		os.Exit(1)
	}

	bot, err := createBot(botFlags, processor.New())
	if err != nil {
		zap.L().Fatal("failed to create bot", zap.Error(err))
	}

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	bot.Start()
	<-interruptCh
	bot.Stop()
}
