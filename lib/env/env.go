package env

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

func Load[T any]() (T, error) {
	_ = godotenv.Load()

	var obj T

	if err := env.Parse(&obj); err != nil {
		return obj, err
	}

	return obj, nil
}
