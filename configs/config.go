package configs

import (
	"github.com/spf13/viper"
)

type Conf struct {
	InputApiHttpPort                      string `mapstructure:"INPUT_API_HTTP_PORT"`
	InputApiOtelServiceName               string `mapstructure:"INPUT_API_OTEL_SERVICE_NAME"`
	OrchestratorApiPort                   string `mapstructure:"ORCHESTRATOR_API_PORT"`
	OrchestratorApiHost                   string `mapstructure:"ORCHESTRATOR_API_HOST"`
	OpenWeathermapApiKey                  string `mapstructure:"OPEN_WEATHERMAP_API_KEY"`
	OrchestratorApiServiceName            string `mapstructure:"ORCHESTRATOR_API_SERVICE_NAME"`
	OpenTelemetryCollectorExporerEndpoint string `mapstructure:"OPEN_TELEMETRY_COLLECTOR_EXPORTER_ENDPOINT"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
