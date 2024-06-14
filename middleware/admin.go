package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		auth, err := c.Cookie("authenticated")

		if err != nil || auth != "true" {

			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return

		}
		c.Next()

	}

}
