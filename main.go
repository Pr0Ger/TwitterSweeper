package main

import (
	"log"

	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetConfigFile("config.toml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	anaconda.SetConsumerKey(viper.GetString("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(viper.GetString("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(viper.GetString("TWITTER_ACCESS_TOKEN"), viper.GetString("TWITTER_ACCESS_TOKEN_SECRET"))

	user, err := api.GetSelf(nil)
	if err != nil {
		log.Fatalf("Unable to get current user: %v", err)
	}

	log.Printf("Authenticated as %v", user.ScreenName)
}
