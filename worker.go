package main

import (
	"os"

	"bitbucket.org/dukex/uhura-worker/setup"
)

func main() {
	setup.Worker(os.Getenv("REDIS_URL"), true)
}
