package parser

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/uhuraapp/uhura-api/helpers"
)

func RequestURL(url string) ([]byte, error) {
	_, checkURL := helpers.ParseURL(url)

	if checkURL != nil {
		return nil, checkURL
	}

	response, err := http.Get(url)

	// log.Debug("requested url: %s - %s", url, err)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("Status is " + response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)

	// log.Debug("request finished")
	return body, err
}
