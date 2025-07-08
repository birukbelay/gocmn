package config

import (
	"reflect"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	zlog "github.com/rs/zerolog/log"
)

const rootEnvName = "ENVIRONMENT"

func LoadConfigT[T any]() *T {
	// loads values from .env into config
	config := koanf.New(".") // TODO: Add Default Configs here
	if err := config.Load(env.Provider("", ".", nil), nil); err != nil {
		zlog.Fatal().Err(err).Msg("Can not load Environment variables")
	}
	if err := config.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		zlog.Warn().Err(err).Msg("couldn't load '.env' file")
	}
	environment := config.String(rootEnvName)
	if environment == "" {
		zlog.Warn().Msgf("%s is empty or not defined: Not Loading Additonal .env files", rootEnvName)
	} else if environment == "prod" {

	} else {
		zlog.Info().Msgf("%s variable is: %s", rootEnvName, environment)
		addEnv := ".env." + environment
		if err := config.Load(file.Provider(addEnv), dotenv.Parser()); err != nil {
			zlog.Error().Err(err).Msgf("Error loading env file '%s'", addEnv)
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
