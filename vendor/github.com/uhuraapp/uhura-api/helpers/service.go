package helpers

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) (int, error) {
	fromContextUserID, ok := c.Get("user_id")
	if !ok {
		return 0, errors.New("user_id not found")
	}

	stringUserID, ok := fromContextUserID.(string)

	if !ok {
		log.Println("stringUserID: ", stringUserID)
		return 0, errors.New("stringUserID inst a string")
	}

	userId, err := strconv.Atoi(stringUserID)
	return userId, err
}

func CacheHeader(c *gin.Context, lastModified time.Time) bool {
	lastModifiedAt := lastModified.Format(time.RFC1123)

	c.Writer.Header().Add("Cache-Control", "public, max-age=31536000")
	c.Writer.Header().Add("Last-Modified", lastModifiedAt)

	if ifModifiedSince := c.Request.Header.Get("If-Modified-Since"); ifModifiedSince != "" {
		ifModifiedSinceTime, err := time.Parse(time.RFC1123, ifModifiedSince)
		if err != nil {
			return false
		}

		updatedAt, err := time.Parse(time.RFC1123, lastModifiedAt)
		if err != nil {
			return false
		}

		if ifModifiedSinceTime.Sub(updatedAt) < 1 {
			return true
		}
	}

	return false
}
