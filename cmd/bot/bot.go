package main

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type msgProcessor interface {
	Process(msg string) string
}

type bot struct {
	bot          *tgbotapi.BotAPI
	msgProcessor msgProcessor
	ownerID      int
	stopCtx      context.Context
	stopFn       context.CancelFunc
	stopWg       sync.WaitGroup
}

func createBot(opts botFlags, msgProcessor msgProcessor) (*bot, error) {
	ctx, cancelFn := context.WithCancel(context.Background())
	httpClient := &http.Client{
		Transport: contextRoundTripper{http.DefaultTransport, ctx},
	}

	b, err := tgbotapi.NewBotAPIWithClient(opts.Token, httpClient)
	if err != nil {
		cancelFn()
		return nil, err
	}

	b.Debug = opts.Debug

	return &bot{
		bot:          b,
		msgProcessor: msgProcessor,
		ownerID:      opts.OwnerID,
		stopCtx:      ctx,
		stopFn:       cancelFn,
	}, nil
}

func (b *bot) Start() {
	b.stopWg.Add(1)
	go func() {
		defer b.stopWg.Done()
		b.processUpdates()
	}()
}

func (b *bot) Stop() {
	b.stopFn()
	b.stopWg.Wait()
}

func (b *bot) processUpdates() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	for {
		select {
		case <-b.stopCtx.Done():
			return
		default:
		}

		updates, err := b.bot.GetUpdates(updateConfig)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			zap.L().Warn("Failed to get updates, retrying in 3 seconds...", zap.Error(err))
			time.Sleep(time.Second * 3)
			continue
		}

		for _, update := range updates {
			if update.UpdateID >= updateConfig.Offset {
				updateConfig.Offset = update.UpdateID + 1
				b.processUpdate(update)
			}
		}
	}
}

func (b *bot) processUpdate(update tgbotapi.Update) {
	msg := update.Message
	if msg == nil {
		return
	}

	if msg.From.ID != b.ownerID {
		msg := tgbotapi.NewMessage(msg.Chat.ID, unknownOwnerMessage)
		msg.ParseMode = tgbotapi.ModeMarkdown

		if _, err := b.bot.Send(msg); err != nil {
			zap.L().Warn("failed to send unknown-owner-message", zap.Error(err))
		}
		return
	}

	respMsg := b.msgProcessor.Process(msg.Text)
	if respMsg != "" {
		msg := tgbotapi.NewMessage(msg.Chat.ID, respMsg)
		msg.ParseMode = tgbotapi.ModeMarkdown
		if _, err := b.bot.Send(msg); err != nil {
			zap.L().Warn("failed to send response", zap.Error(err))
		}
	}
}

const unknownOwnerMessage = `Hi!

I'm not talking with strangers. Only the creator can talk with me.

If you want to create my sister or brother for yourself look at [here to how](https://github.com/kulti/task-list-bot).
`
