package main

import (
	"log"
	"net/url"
	"strconv"
	"time"

	"regexp"

	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/viper"
)

const MAX_TWEETS_PER_PAGE = 200

func main() {
	viper.SetDefault("OLDS", 365)
	viper.SetDefault("FAVS", 5)
	viper.SetDefault("RT", 5)

	viper.AutomaticEnv()
	viper.SetConfigFile("config.toml")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()

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

	var tweetsForRemoving []anaconda.Tweet
	var skipRemovingBecauseWeRepliedToIt []int64
	var totalTweets, skippedTweets int

	oldestTimestamp := viper.GetInt64("OLDEST_TIMESTAMP")
	olds := time.Duration(viper.GetInt64("OLDS"))

	now := time.Now()
	oldestByDays := now.Add(-olds * time.Hour * 24)

	if oldestTimestamp < oldestByDays.Unix() && olds != 0 {
		oldestTimestamp = oldestByDays.Unix()
	}

	for {
		timeline, err := api.GetUserTimeline(v)
		if err != nil {
			log.Fatalf("Unable to fetch user timeline %v", err)
		}
		if len(timeline) == 0 {
			break
		}
		totalTweets += len(timeline)

		v.Set("max_id", strconv.FormatInt(timeline[len(timeline)-1].Id-1, 10))

		log.Printf("Downloaded %v/%v tweets", totalTweets, user.StatusesCount)

		for _, tweet := range timeline {
			tweetTime, _ := tweet.CreatedAtTime()
			if oldestTimestamp != 0 && tweetTime.Unix() < oldestTimestamp {
				for _, id := range skipRemovingBecauseWeRepliedToIt {
					if tweet.Id == id {
						log.Printf("Skipping tweet because we replied to it (%v): %v", tweet.Id, tweet.FullText)
						continue
					}
				}
				if (tweet.FavoriteCount >= viper.GetInt("FAVS") || tweet.RetweetCount >= viper.GetInt("RT")) && !tweet.Retweeted {
					log.Printf("Skipping tweet because it's popular (%v): %v", tweet.Id, tweet.FullText)
					skippedTweets++
					continue
				}
				tweetsForRemoving = append(tweetsForRemoving, tweet)
			} else {
				if tweet.InReplyToUserID == user.Id {
					skipRemovingBecauseWeRepliedToIt = append(skipRemovingBecauseWeRepliedToIt, tweet.InReplyToStatusID)
				}
			}
		}
	}

	keybase := regexp.MustCompile(`^Verifying myself: I am .* on Keybase\.io\.`)
	for _, tweet := range tweetsForRemoving {
		if keybase.MatchString(tweet.FullText) {
			log.Print("Skip Keybase.io verification tweet")
			continue
		}

		log.Printf("Removing tweet (%v): %v", tweet.Id, tweet.FullText)
		_, err := api.DeleteTweet(tweet.Id, true)
		if err != nil {
			log.Printf("\tUnable to delete tweet: %v", err)
		}
	}

	log.Printf("Total tweets deleted: %v; skipped: %v", len(tweetsForRemoving), skippedTweets)
}
