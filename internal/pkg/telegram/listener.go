package telegram

import (
	"context"
	"fmt"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

type Update struct {
	ChatId int64
	Update *tgbotapi.Update
	Bot    *tgbotapi.BotAPI
}

type RouterInterface interface {
	Update(ctx context.Context, update *Update)
}

type Asyncer interface {
	Go(fn func())
}

type DebugSyncAsyncer struct{}

func (d *DebugSyncAsyncer) Go(fn func()) { fn() }

type Listener struct {
	config *Config
	logger *zap.Logger
	router RouterInterface
	async  Asyncer
}

type Config struct {
	Token string `mapstructure:"token"`
	Debug bool   `mapstructure:"debug"`
}

func NewListener(cfg *Config, router RouterInterface) *Listener {
	var async Asyncer = pool.New()
	if cfg.Debug {
		async = &DebugSyncAsyncer{}
	}
	return &Listener{
		router: router,
		config: cfg,
		logger: logger.NewStructLogger("telegram"),
		async:  async,
	}
}

func (l *Listener) Start() (context.CancelFunc, error) {
	bot, err := tgbotapi.NewBotAPI(l.config.Token)
	if err != nil {
		return nil, fmt.Errorf("fail create bot: %w", err)
	}

	bot.Debug = true

	l.logger.Info("Authorized on account", zap.String("username", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case update := <-updates:
				upd := update
				l.async.Go(func() {
					l.logger.Debug("Update received", zap.Int64("chat", update.FromChat().ID))
					l.router.Update(context.Background(), &Update{Update: &upd, Bot: bot, ChatId: update.FromChat().ID})
					l.logger.Debug("Update processed", zap.Int64("chat", update.FromChat().ID))
				})
			case <-ctx.Done():
				return
			}
		}
	}()
	return cancel, err

}
