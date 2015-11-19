package parser

import (
	"errors"
	"regexp"

	charset "code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	rss "github.com/jteeuwen/go-pkg-rss"
)

type Fetcher struct {
	channel chan<- *Channel
	err     chan<- error
}

func RunFetcher(url string, c chan<- *Channel, err chan<- error) {
	fetcher := &Fetcher{c, err}
	fetcher.process(url)
	fetcher.end(nil, errors.New("not found channel"))
}

func (f *Fetcher) process(url string) {
	body, err := RequestURL(url)

	if err != nil {
		log.Debug("a new error %s", err.Error())
		f.end(nil, err)
		return
	}

	if f.isHTML(body) {
		log.Debug("url is html")
		links := FindLinks(body)
		log.Debug("links %s", links)
		if len(links) > 0 {
			f.process(links[0])
			return
		}

		err = errors.New("URL is a HTML page, and a LINK to XML Feed URL was not found")
		f.end(nil, err)
		return
	}

	rss.New(0, true, f._c, f.episodeHandler(url)).
		FetchBytes(url, body, charset.NewReader)
}

func (f *Fetcher) episodeHandler(url string) func(*rss.Feed, *rss.Channel, []*rss.Item) {
	return func(feed *rss.Feed, rssChannel *rss.Channel, episodes []*rss.Item) {
		channel := Channel{Feed: rssChannel, URL: url}
		if channel.HasNewURL() && channel.NewURL() != url {
			log.Debug("has new URL: %s != %s", channel.NewURL(), feed.Url)
			f.process(channel.NewURL())
			return
		}

		channel = build(addEpisodes(channel, episodes))
		f.end(&channel, nil)
	}
}

func addEpisodes(c Channel, e []*rss.Item) Channel {
	for k, _ := range e {
		log.Debug("Adding episode (%s) to channel (%s)", e[k].Title, c.Title)
		c.Episodes = append(c.Episodes, &Episode{Feed: e[k]})
	}
	return c
}

func build(c Channel) Channel {
	c.Build()
	episodes := make([]*Episode, 0)
	for _, e := range c.Episodes {
		if e.Build() {
			episodes = append(episodes, e)
		}
	}
	c.Episodes = episodes
	return c
}

func (f *Fetcher) end(channel *Channel, err error) {
	f.channel <- channel
	f.err <- err
}

var hasHTML = regexp.MustCompile(`<\/?html>`)

func (f Fetcher) isHTML(body []byte) bool {
	return hasHTML.Match(body)
}

func (f *Fetcher) _c(feed *rss.Feed, channels []*rss.Channel) {}
