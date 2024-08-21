package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const defaultConfigFileName = "config.yml"

var ErrEmptyPath = errors.New("path to config must not be empty")

type Loader[T any] struct {
	path          string
	dotenvLoadErr bool
}

func NewLoader[T any]() *Loader[T] {
	return &Loader[T]{
		path: defaultConfigFileName,
	}
}

func (l *Loader[T]) WithDotenvLoadErr() *Loader[T] {
	loader := NewLoader[T]()
	loader.path = l.path
	loader.dotenvLoadErr = true

	return loader
}

func (l *Loader[T]) WithFilename(filename string) *Loader[T] {
	loader := NewLoader[T]()
	loader.path = filename
	loader.dotenvLoadErr = l.dotenvLoadErr

	return loader
}

func (l *Loader[T]) Load() (*T, error) {
	if l.path == "" {
		return nil, ErrEmptyPath
	}

	viperInstance := viper.New()
	ext := strings.TrimLeft(filepath.Ext(l.path), ".")
	viperInstance.SetConfigFile(l.path)
	viperInstance.SetConfigType(ext)

	err := viperInstance.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("v.ReadInConfig: %w", err)
	}

	err = godotenv.Load()
	if err != nil && l.dotenvLoadErr {
		return nil, fmt.Errorf("can't load environment variables from file: %w", err)
	}

	for _, key := range viperInstance.AllKeys() {
		value := viperInstance.GetString(key)
		if value == "" {
			continue
		}

		viperInstance.Set(key, os.ExpandEnv(value))
	}

	cfg := new(T)
	if err := viperInstance.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("v.Unmarshal: %w", err)
	}

	return cfg, nil
}
