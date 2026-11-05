package config

import (
	"fmt"
	"log/slog"
	"os"
)

// @TODO 여기도 수정이 좀 필요

const (
	_devEnv  = "dev"
	_prodEnv = "prod"
)

const (
	_appEnvKey      = "APP_ENV"
	_appPortKey     = "APP_PORT"
	_appTimezoneKey = "APP_TIMEZONE"
)

type Config struct {
	Env      string
	Port     string
	Timezone string
	LogLevel slog.Level
}

func Load() (*Config, error) {
	env, exists := os.LookupEnv(_appEnvKey)
	if !exists || env == "" {
		return nil, fmt.Errorf("%s 가 설정되지 않았습니다.", _appEnvKey)
	}

	appPort, exists := os.LookupEnv(_appPortKey)
	if !exists || appPort == "" {
		return nil, fmt.Errorf("%s 가 설정되지 않았습니다.", _appPortKey)
	}

	timezone, exists := os.LookupEnv(_appTimezoneKey)
	if !exists || timezone == "" {
		return nil, fmt.Errorf("%s 가 설정되지 않았습니다.", _appTimezoneKey)
	}

	var level slog.Level
	switch env {
	case _devEnv:
		level = slog.LevelDebug
	case _prodEnv:
		level = slog.LevelInfo
	default:
		return nil, fmt.Errorf("지원되지 않는 환경입니다: %s", env)
	}

	return &Config{
		Env:      env,
		Port:     appPort,
		Timezone: timezone,
		LogLevel: level,
	}, nil
}
