# TwitterSweeper
TwitterSweeper is a small application to delete your old, unpopular tweets.

## Features
- Delete, unretweet old tweets
- Keep tweets based on age, retweet or favourite count
- Keep tweets if we replied to them
- Keep keybase.io verification tweet no matter how old it is

## Usage
To install from sources:
```bash
brew install dep ## for macOS users; in other cases follow instructions here https://github.com/golang/dep
go get -d github.com/Pr0Ger/TwitterSweeper
cd $GOPATH/src/github.com/Pr0Ger/TwitterSweeper
dep ensure
go build 
```

Get the Twitter API variables from https://apps.twitter.com and add the following variables to a config file (example in `config.toml.example`):
```bash
TWITTER_CONSUMER_KEY=...
TWITTER_CONSUMER_SECRET=...
TWITTER_ACCESS_TOKEN=...
TWITTER_ACCESS_TOKEN_SECRET=...
```

Then just run TwitterSweeper
```bash
$GOPATH/bin/TwitterSweeper
```

## Contact
[email](mailto:me@pr0ger.org) [twitter (ironic, huh?)](http://twitter.com/Pr0Ger)

## License
[MIT](LICENSE)
