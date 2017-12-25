package main

import (
	"log"
	"net/url"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/viper"
)

const MAX_TWEETS_PER_PAGE = 200

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

	v := url.Values{}
	v.Set("count", string(MAX_TWEETS_PER_PAGE))

	var allTweets []anaconda.Tweet

	for {
		timeline, err := api.GetUserTimeline(v)
		if err != nil {
			log.Fatalf("Unable to fetch user timeline %v", err)
		}
		if len(timeline) == 0 {
			break
		}

		allTweets = append(allTweets, timeline...)
		v.Set("max_id", strconv.FormatInt(timeline[len(timeline)-1].Id-1, 10))

		log.Printf("Downloaded %v/%v tweets", len(allTweets), user.StatusesCount)
	}
}
