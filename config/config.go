package config

import (
	"log"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

var AppConfig Config
var GoogleOAuthConfig *oauth2.Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	AppConfig = Config{
		DBHost:             viper.GetString("DB_HOST"),
		DBPort:             viper.GetString("DB_PORT"),
		DBUser:             viper.GetString("DB_USER"),
		DBPassword:         viper.GetString("DB_PASSWORD"),
		DBName:             viper.GetString("DB_NAME"),
		GoogleClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  viper.GetString("GOOGLE_REDIRECT_URL"),
	}
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     AppConfig.GoogleClientID,
		ClientSecret: AppConfig.GoogleClientSecret,
		RedirectURL:  AppConfig.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

}
