package parser

import (
	"errors"
	"regexp"
	"sync"

	charset "code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	rss "github.com/jteeuwen/go-pkg-rss"
)

type Fetcher struct {
	channel chan<- *Channel
	err     chan<- error
	finish  chan<- bool
	urls    []string
}

func NewFetcher(urls []string, c chan<- *Channel, f chan<- bool, err chan<- error) *Fetcher {
	return &Fetcher{c, err, f, urls}
}

func (f *Fetcher) run() {
	var wg sync.WaitGroup

	for _, url := range f.urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			f.process(url)
		}(url)
	}
	wg.Wait()

	f.finish <- true
}

func (f *Fetcher) process(url string) {
	body, err := f.request(url)

	if err != nil {
		log.Debug("a new error %s", err.Error())
		f.err <- err
		return
	}

	log.Debug("url is html: %s", f.isHTML(body))

	if f.isHTML(body) {
		links := FindLinks(body)
		log.Debug("links %s", links)
		if len(links) > 0 {
			for _, link := range links {
				f.process(link)
			}
		} else {
			err = errors.New("URL is a HTML page, and a LINK to XML Feed URL was not found")
			f.err <- err
			return
		}
	}

	rss.New(0, true, f._c, f.episodeHandler(url)).
		FetchBytes(url, body, charset.NewReader)
}

func (f Fetcher) request(url string) ([]byte, error) {
	return RequestURL(url)
}

func (f *Fetcher) episodeHandler(url string) func(*rss.Feed, *rss.Channel, []*rss.Item) {
	return func(feed *rss.Feed, rssChannel *rss.Channel, episodes []*rss.Item) {
		channel := &Channel{Feed: rssChannel}
		if channel.HasNewURL() && channel.NewURL() != url {
			log.Debug("has new URL: %s != %s", channel.NewURL(), feed.Url)
			f.rerun(channel.NewURL())
		} else {
			f.appendEpisodes(channel, episodes)
			f.send(channel, url)
		}
	}
}

func (f *Fetcher) rerun(newURL string) {
	f.process(newURL)
}

func (f *Fetcher) send(channel *Channel, url string) {
	channel.URL = url

	f.buildRecords(channel)

	f.channel <- channel
	close(f.channel)
}

func (f *Fetcher) appendEpisodes(c *Channel, e []*rss.Item) {
	for k, _ := range e {
		c.Episodes = append(c.Episodes, &Episode{Feed: e[k]})
	}
	return
}

func (f *Fetcher) buildRecords(c *Channel) {
	c.Build()
	episodes := make([]*Episode, 0)
	for _, e := range c.Episodes {
		if e.Build() {
			episodes = append(episodes, e)
		}
	}
	c.Episodes = episodes
}

var hasHTML = regexp.MustCompile(`<\/?html>`)

func (f Fetcher) isHTML(body []byte) bool {
	return hasHTML.Match(body)
}

func (f *Fetcher) _c(feed *rss.Feed, channels []*rss.Channel) {}
