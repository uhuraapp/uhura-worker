package main

import (
	"os"

	"github.com/uhuraapp/uhura-worker/setup"
)

func main() {
	setup.Worker(os.Getenv("REDIS_URL"))
}
