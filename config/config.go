package config

import (
	"os"
)

const (
	apiGithubAccessToken = "SECRET_GITHUB_ACCESS_TOKEN"
	LogLevel                = "info"
)

var githubAccessToken = os.Getenv(apiGithubAccessToken)

//func init() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//	githubAccessToken = os.Getenv("TOKEN_GITHUB")
//}

func GetGithubAccessToken() string {
	return githubAccessToken
}

//
//func IsProduction() bool {
//	return os.Getenv("GO_ENVIRONMENT") == "production"
//}