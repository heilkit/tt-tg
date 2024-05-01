package internal

import (
	"fmt"
	"github.com/heilkit/tg"
	"github.com/heilkit/tg/scheduler"
	"log/slog"
	"os"
	"time"
)

type Manager struct {
	tg     *tg.Bot
	chat   *tg.Chat
	Config *Config
}

func NewManagerFromFile(filename string, api string, level slog.Level) (*Manager, error) {
	config, err := ConfigFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	bot, err := tg.NewBot(tg.Settings{
		URL:       api,
		Token:     config.Token,
		Scheduler: scheduler.ExtraConservative(),
		Retries:   3,
		Logger:    tg.LoggerSlog(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level, AddSource: level <= slog.LevelDebug}))),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &Manager{
		tg:     bot,
		chat:   &tg.Chat{ID: config.Chat},
		Config: config,
	}, nil
}

func (manager *Manager) Start(pollRate time.Duration) {
	slog.Info("Starting poll", "pollRate", pollRate)
	ticker := time.NewTicker(pollRate)
	for ; true; _ = <-ticker.C {
		for _, profile := range manager.Config.Profiles {
			if err := manager.Profile(profile); err != nil {
				slog.Error("failed to update profile", "profile", profile.Tag, "err", err)
			}
		}
	}
}
