package refersrest

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length,  Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		if c.Request.Header["Content-Type"] == nil || len(c.Request.Header["Content-Type"]) == 0 {
			c.AbortWithStatus(204)
			return
		}

		content_type := c.Request.Header["Content-Type"]
		if content_type[0] != "application/json" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (m *Rest) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := struct {
			Authorization string `header:"Authorization"`
		}{}
		if err := c.BindHeader(&header); err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		var token string

		if header.Authorization == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		} else {
			token = header.Authorization
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//проверка валидности токена
		user, err := m.DB.IsAuch(ctx, token)
		if err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
