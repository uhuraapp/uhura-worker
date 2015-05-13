package helpers

import (
	"os"

	"github.com/dukex/mixpanel"
)

var (
	clientMixpanel *mixpanel.Mixpanel
)

func init() {
	clientMixpanel = mixpanel.NewMixpanel(os.Getenv("MIXPANEL_TOKEN"))
}

func MixpanelPerson(userId string) *mixpanel.People {
	return clientMixpanel.Identify(userId)
}

func NewEvent(userId string, trackName string, data map[string]interface{}) {
	person := clientMixpanel.Identify(userId)
	person.Track(trackName, data)
}
