package parser

import (
	"errors"
	"net/url"
	"time"
)

func URL(url *url.URL) (channels *Channel, body []byte, _error error) {
	// log.Debug("creating fetcher")

	c := make(chan *Channel)
	err := make(chan error)
	bodyC := make(chan []byte)

	go RunFetcher(url.String(), c, err, bodyC)

	return <-c, <-bodyC, <-err
}

func Body(body []byte, url string) (channels *Channel, _error error) {
	// log.Debug("runing parser")

	c := make(chan *Channel)
	err := make(chan error)

	go StartParser(body, url, c, err)

	go func() {
		<-time.After(3 * time.Second)
		err <- errors.New("Timeout")
		c <- nil
	}()

	return <-c, <-err
}
