package settings

import (
	"os"
)

type EnvList struct {
	Port      string
	BasePath  string
	AwsRegion string
}

func GetEnv() *EnvList {
	return &EnvList{
		Port:      osEnv("PORT", "8080"),
		BasePath:  osEnv("BASE_PATH", "/"),
		AwsRegion: osEnv("AWS_REGION", "us-east-1"),
	}
}

func osEnv(env string, envDefault string) string {
	if err := os.Getenv(env); err != "" {
		return err
	}
	return envDefault
}
