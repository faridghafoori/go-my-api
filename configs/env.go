package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func ENV_MONGO_URI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("MONGO_URI")
}

func ENV_MONGO_URI_LOCAL() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("MONGO_URI_LOCAL")
}

func ENV_MONGO_DB() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("MONGO_DB")
}

func ENV_REDIS_DSN() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("REDIS_DSN")
}

func ENV_JWT_ACCESS_SECRET() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("JWT_ACCESS_SECRET")
}

func ENV_JWT_TOTP_SECRET() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("JWT_TOTP_SECRET")
}

func ENV_JWT_REFRESH_SECRET() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("JWT_REFRESH_SECRET")
}

func ENV_RUNABLE_PROJECT_URI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("RUNABLE_PROJECT_URI")
}

func ENV_LAUNCH_PROJECT_URI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("LAUNCH_PROJECT_URI")
}

func ENV_MINIO_ENDPOINT() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("MINIO_ENDPOINT")
}

func ENV_MINIO_ACCESS_KEY() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("MINIO_ACCESS_KEY")
}

func ENV_MINIO_SECRET_KEY() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("MINIO_SECRET_KEY")
}
