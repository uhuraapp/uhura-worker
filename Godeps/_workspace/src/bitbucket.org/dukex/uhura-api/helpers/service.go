package helpers

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) (int, error) {
	fromContextUserID, err := c.Get("user_id")
	if err != nil {
		return 0, err
	}

	stringUserID, ok := fromContextUserID.(string)

	if !ok {
		log.Println("stringUserID: ", stringUserID)
		return 0, errors.New("stringUserID inst a string")
	}

	userId, err := strconv.Atoi(stringUserID)
	return userId, err
}
