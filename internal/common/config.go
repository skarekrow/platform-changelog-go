package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

type Config struct {
	PublicPort string
	MetricsPort string
	MetricsPath string
	LogLevel string
	Hostname string
	CloudwatchConfig CloudwatchCfg
	DatabaseConfig DatabaseCfg
}

type DatabaseCfg struct {
	DBUser string
	DBPassword string
	DBName string
	DBHost string
	DBPort string
	RDSCa string
}

type CloudwatchCfg struct {
	CWLogGroup string
	CWRegion string
	CWAccessKey string
	CWSecretKey string
}

func Get() *Config {
	options := viper.New()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// global logging
	options.SetDefault("logLevel", "INFO")
	options.SetDefault("Hostname", hostname)

	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig

		// ports
		options.SetDefault("publicPort", cfg.PublicPort)
		options.SetDefault("metricsPort", cfg.MetricsPort)
		options.SetDefault("metricsPath", cfg.MetricsPath)
		// database
		options.SetDefault("db.user", cfg.Database.Username)
		options.SetDefault("db.password", cfg.Database.Password)
		options.SetDefault("db.name", cfg.Database.Name)
		options.SetDefault("db.host", cfg.Database.Hostname)
		options.SetDefault("db.port", cfg.Database.Port)
		options.SetDefault("rdsCa", cfg.Database.RdsCa)
		// cloudwatch
		options.SetDefault("logGroup", cfg.Logging.Cloudwatch.LogGroup)
		options.SetDefault("cwRegion", cfg.Logging.Cloudwatch.Region)
		options.SetDefault("cwAccessKey", cfg.Logging.Cloudwatch.AccessKeyId)
		options.SetDefault("cwSecretKey", cfg.Logging.Cloudwatch.SecretAccessKey)
	} else {
		// ports
		options.SetDefault("publicPort", "8000")
		options.SetDefault("metricsPort", "9001")
		options.SetDefault("metricsPath", "/metrics")
		// database
		options.SetDefault("db.user", "crc")
		options.SetDefault("db.password", "crc")
		options.SetDefault("db.name", "gumbaroo")
		options.SetDefault("db.host", "0.0.0.0")
		options.SetDefault("db.port", "5432")
		// cloudwatch
		options.SetDefault("logGroup", "platform-dev")
		options.SetDefault("cwRegion", "us-east-1")
		options.SetDefault("cwAccessKey", os.Getenv("CW_AWS_ACCESS_KEY_ID"))
		options.SetDefault("cwSecretKey", os.Getenv("CW_AWS_SECRET_ACCESS_KEY"))
	}

	options.AutomaticEnv()
	options.SetEnvKeyReplacer(strings.NewReplacer(".","_"))

	config := &Config{
		Hostname: options.GetString("Hostname"),
		LogLevel: options.GetString("logLevel"),
		PublicPort: options.GetString("publicPort"),
		MetricsPort: options.GetString("metricsPort"),
		MetricsPath: options.GetString("metricsPath"),
		DatabaseConfig: DatabaseCfg{
			DBUser: options.GetString("db.user"),
			DBPassword: options.GetString("db.password"),
			DBName: options.GetString("db.name"),
			DBHost: options.GetString("db.host"),
			DBPort: options.GetString("db.port"),
		},
		CloudwatchConfig: CloudwatchCfg{
			CWLogGroup: options.GetString("logGroup"),
			CWRegion: options.GetString("cwRegion"),
			CWAccessKey: options.GetString("cwAccessKey"),
			CWSecretKey: options.GetString("cwSecretKey"),
		},
	}

	if clowder.IsClowderEnabled() {
		
		// write the RDS CA using the app-common-go package
		if clowder.LoadedConfig.Database.RdsCa != nil {
			rdsCAPath, err := clowder.LoadedConfig.RdsCa()

			if err != nil {
				panic("RDS CA Failed to Write")
			}

			config.DatabaseConfig.RDSCa = rdsCAPath
		}
	}

	return config
}
