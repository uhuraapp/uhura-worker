package parser

import (
	"net/url"

	_ "code.google.com/p/go-charset/data"
)

func URL(url *url.URL) ([]*Channel, []error) {
	log.Debug("creating fetcher")

	_errors := make([]error, 0)
	channels := make([]*Channel, 0)

	c := make(chan *Channel)
	finish := make(chan bool)
	err := make(chan error)

	fetcher := NewFetcher([]string{url.String()}, c, finish, err)

	go fetcher.run()

	go func() {
		e := <-err
		_errors = append(_errors, e)
		log.Error("%s", e)
	}()

	go func() {
		channel, ok := <-c
		if ok {
			channel.requestedURL = url.String()
			log.Debug("finishing with %s", channel.Title)
			channels = append(channels, channel)
		} else {
			log.Debug("finishing with error %s")
		}
	}()

	<-finish

	return channels, _errors
}

// func ProcessChannel(channel *rss.Channel) {
// 	CacheImage(channel)
// 	SaveOnDatabase(channel)
// }
//
// func CacheImage(channel *rss.Channel) {
// 	imageURL := GetChannelImageURL(channel)
// 	originalImage, err := RequestURL(imageURL)
//
// 	log.Println(imageURL)
// 	if err != nil {
// 		//	notifyWrongImage
// 	}
// 	image, err := imageproxy.Transform(originalImage, resizeOptions())
// 	if err != nil {
// 		log.Println(err)
// 	}
//
// 	str := base64.StdEncoding.EncodeToString(image)
// 	fmt.Println(str)
// }
//

// func resizeOptions() imageproxy.Options {
// 	return imageproxy.ParseOptions("250x")
// }
//
// func SaveOnDatabase(channel *rss.Channel) {
// }

////

// func GetChannelImageURL(channel *rss.Channel) string {
// 	var imageURL string
// 	// 	if itunesImage := channel.Extensions[ITUNES_EXT]["image"]; itunesImage != nil {
// 	// 		imageURL = itunesImage[0].Attrs["href"]
// 	// 	} else {
// 	imageURL = channel.Image.Url
// 	// 	}
// 	return imageURL
// }

////////

// func GetDelayBetweenEpisodes() {
//
// }
//
// func SetNewRun() {
// }

// If isFeed
//    RequestURL -> ParserXML
//    FindChannel or CreateChannel
//    UpdateChannel
//    CacheImage
//    FindEpisodes or CreateEpisodes
//      SetContentType
//      SetContentLengh
//    GetDelayBetweenEpisodes
//    SetNewRun
// else
//    FindChannelOnWebsite
