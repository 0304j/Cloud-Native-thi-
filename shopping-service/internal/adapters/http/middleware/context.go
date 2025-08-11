package middleware

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (string, bool) {
	id, ok := c.Get(ContextUserIDKey)
	if !ok {
		return "", false
	}
	s, ok := id.(string)
	return s, ok
}
