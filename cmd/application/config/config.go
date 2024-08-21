package config

import (
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
	"github.com/Yaroher2442/FamilySyncHub/internal/repositories/pg"
)

type Config struct {
	Postgres *pg.Config       `mapstructure:"postgres"`
	Telegram *telegram.Config `mapstructure:"telegram"`
}
