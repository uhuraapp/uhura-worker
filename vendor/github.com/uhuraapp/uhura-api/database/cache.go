package database

import (
  "github.com/muesli/cache2go"
)

var CACHE *cache2go.CacheTable

func init() {
  CACHE = cache2go.Cache("cache")
}
