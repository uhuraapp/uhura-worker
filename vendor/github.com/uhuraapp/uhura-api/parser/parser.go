package parser

import (
	"errors"
	"regexp"

	rss "github.com/jteeuwen/go-pkg-rss"
	"golang.org/x/net/html/charset"
)

const (
	ITUNES_EXT = "http://www.itunes.com/dtds/podcast-1.0.dtd"
)

var (
	hasHTML = regexp.MustCompile(`<\/?html>`)

	// log = logging.MustGetLogger("example")

	// format = logging.MustStringFormatter(
	// 	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
	// )
)

type Parser struct {
	channel chan<- *Channel
	err     chan<- error
}

func StartParser(body []byte, url string, c chan<- *Channel, err chan<- error) {
	// log.Debug("Starting Parser")
	p := &Parser{c, err}
	p.process(body, url)
}

func (p *Parser) process(body []byte, url string) {
	// log.Debug("Parser started, processing..")

	if p.isHTML(body) {
		// log.Debug("url is html")
		links := FindLinks(body)
		// log.Debug("links %s", links)
		if len(links) > 0 {
			RunFetcher(links[0], p.channel, p.err, make(chan []byte))
			return
		}

		err := errors.New("URL is a HTML page, and a LINK to XML Feed URL was not found")
		p.end(nil, err)
		return
	}

	// log.Debug("Body is not an html")

	rss.New(0, true, p._c, p.episodeHandler(url, body)).
		FetchBytes(url, body, charset.NewReaderLabel)
}

func (p Parser) isHTML(body []byte) bool {
	return hasHTML.Match(body)
}

func (p *Parser) episodeHandler(url string, body []byte) func(*rss.Feed, *rss.Channel, []*rss.Item) {
	return func(feed *rss.Feed, rssChannel *rss.Channel, episodes []*rss.Item) {
		// log.Debug("XML parsed")

		channel := Channel{Feed: rssChannel, URL: url, Body: body}
		if channel.HasNewURL() && channel.NewURL() != url {
			// log.Debug("has new URL: %s != %s", channel.NewURL(), feed.Url)
			RunFetcher(channel.NewURL(), p.channel, p.err, make(chan []byte))
			return
		}

		channel = build(addEpisodes(channel, episodes))
		p.end(&channel, nil)
	}
}

func addEpisodes(c Channel, e []*rss.Item) Channel {
	// log.Debug("Adding episodes to channel")

	for k := range e {
		c.Episodes = append(c.Episodes, &Episode{Feed: e[k]})
	}
	return c
}

func build(c Channel) Channel {
	// log.Debug("Bulding channel and episodes")
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

func (p *Parser) end(channel *Channel, err error) {
	// log.Debug("end")
	p.channel <- channel
	p.err <- err
}

func (p *Parser) _c(feed *rss.Feed, channels []*rss.Channel) {}
