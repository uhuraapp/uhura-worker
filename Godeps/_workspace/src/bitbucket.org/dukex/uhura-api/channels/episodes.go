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
	client := &http.Client{}
	request, err := http.NewRequest("GET", sourceURL, nil)
	request.Close = true
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(sourceURL + " status " + response.Status + "")
	}

	e := new(EpisodeAudioData)
	e.ContentLength = response.ContentLength
	e.ContentType = response.Header.Get("Content-Type")

	return e, nil
}
