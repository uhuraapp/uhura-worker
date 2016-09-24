package helpers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muesli/cache2go"
)

func UseHTTPCache(cacheKey string, store *cache2go.CacheTable, c *gin.Context) bool {
	t := time.Now()

	cacheRes, err := store.Value(cacheKey)
	if err == nil {
		t = *cacheRes.Data().(*time.Time)
	}

	if CacheHeader(c, t) {
		c.AbortWithStatus(http.StatusNotModified)
		return true
	}

	store.Add(cacheKey, ((30 * 24) * time.Hour), &t)
	return false
}
