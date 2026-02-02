package config

import (
	"fmt"
	"path/filepath"
	"reflect"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	zlog "github.com/rs/zerolog/log"
)

type ConfigParam struct {
	FilePath       string
	Dir            string // directory to load from (e.g. "./configs")
	BaseEnvFile    string // e.g. ".env"
	EnvVarName     string // e.g. "APP_ENV"
	DisableEnvFile bool   // allow env-only config
}

func withDefaults(p ConfigParam) ConfigParam {
	if p.Dir == "" {
		p.Dir = "."
	}
	if p.BaseEnvFile == "" {
		p.BaseEnvFile = ".env"
	}
	if p.EnvVarName == "" {
		p.EnvVarName = rootEnvName
	}
	return p
}
func LoadConfigTParam[T any](param ConfigParam) *T {
	param = withDefaults(param)
	// loads values from .env into config
	config := koanf.New(".") // TODO: Add Default Configs here
	// 1. Load OS environment variables first
	if err := config.Load(env.Provider("", ".", nil), nil); err != nil {
		zlog.Fatal().Err(err).Msg("Can not load Environment variables")
	}

	// 2. Load base env file (e.g. ./configs/.env)
	if !param.DisableEnvFile {
		basePath := filepath.Join(param.Dir, param.BaseEnvFile)
		if err := config.Load(file.Provider(basePath), dotenv.Parser()); err != nil {
			zlog.Warn().
				Err(err).
				Msgf("could not load base env file: %s", basePath)
		}
	}
	// 3. Load environment-specific env file
	environment := config.String(param.EnvVarName)
	if environment == "" {
		zlog.Warn().Msgf("%s is empty or not defined: Not Loading Additonal .env files", environment)
	} else if environment == "skip" {
		zlog.Info().Msgf("%s=skip: skipping loading additional .env files", param.EnvVarName)
	} else {
		envFile := fmt.Sprintf("%s.%s", param.BaseEnvFile, environment)
		envPath := filepath.Join(param.Dir, envFile)

		zlog.Info().
			Msgf("%s=%s → loading %s", param.EnvVarName, environment, envPath)

		if err := config.Load(file.Provider(envPath), dotenv.Parser()); err != nil {
			zlog.Error().
				Err(err).
				Msgf("error loading env file: %s", envPath)
		}
	}
	envConfig := new(T)
	err := config.UnmarshalWithConf("", envConfig, koanf.UnmarshalConf{
		DecoderConfig: &mapstructure.DecoderConfig{
			ErrorUnset: true,
			Result:     envConfig,
			Squash:     true,
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToTimeHookFunc(time.RFC3339),
				mapstructure.StringToIntHookFunc(),
				mapstructure.OrComposeDecodeHookFunc(mapstructure.StringToBoolHookFunc(),
					func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
						if f.Kind() != reflect.String || t.Kind() != reflect.Bool {
							return data, nil
						}
						if data.(string) == "" {
							return false, nil
						}
						return true, nil
					}),
			),
		},
	})
	if err != nil {
		zlog.Fatal().Err(err).Msg("Can not load decode Configuration")
	}
	return envConfig
}
