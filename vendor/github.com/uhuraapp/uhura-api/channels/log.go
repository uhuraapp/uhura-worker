package channels

import "log"

const PACKAGE_NAME = "channels"

func info(message string) {
	log.Println("["+PACKAGE_NAME+"]", message)
}

func errr(err error, message string) {
	log.Println("["+PACKAGE_NAME+"][Error]", err.Error()+" - "+message)
}
