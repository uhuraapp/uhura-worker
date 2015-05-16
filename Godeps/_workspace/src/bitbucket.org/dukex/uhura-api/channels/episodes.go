package channels

import (
	"errors"
	"net/http"
)

type EpisodeAudioData struct {
	ContentType   string
	ContentLength int64
}

func GetEpisodeAudioData(sourceURL string) (*EpisodeAudioData, error) {
	response, err := http.Get(sourceURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusFound {
		return nil, errors.New(sourceURL + " status " + response.Status)
	}

	e := new(EpisodeAudioData)
	e.ContentLength = response.ContentLength
	e.ContentType = response.Header.Get("Content-Type")

	return e, nil
}
