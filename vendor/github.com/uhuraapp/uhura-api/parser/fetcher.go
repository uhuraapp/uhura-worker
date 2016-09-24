package parser

type Fetcher struct {
	c chan<- *Channel
	e chan<- error
	b chan<- []byte
}

func RunFetcher(url string, c chan<- *Channel, err chan<- error, body chan<- []byte) {
	fetcher := &Fetcher{c, err, body}
	fetcher.process(url)
}

func (f *Fetcher) process(url string) {
	body, err := RequestURL(url)

	if err != nil {
		// log.Debug("a new error %s", err.Error())
		f.end(nil, err)
		return
	}

	go func() {
		f.b <- body
	}()

	StartParser(body, url, f.c, f.e)
}

func (f *Fetcher) end(channel *Channel, err error) {
	f.c <- channel
	f.e <- err
}
