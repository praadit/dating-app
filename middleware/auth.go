package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/praadit/dating-apps/constant"
	"github.com/praadit/dating-apps/utils"
)

func AuthenticatedOnly(c *gin.Context) {
	accessToken := checkAuthCookie(c)

	if accessToken == "" {
		token, err := checkAuthHeader(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		accessToken = token
	}

	claim, err := utils.ParseJWT(accessToken)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(c, constant.CTX_CLAIM, claim)
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

func checkAuthCookie(c *gin.Context) string {
	cookie, _ := c.Request.Cookie("access-token")
	if cookie == nil {
		return ""
	}

	accToken := strings.TrimPrefix(cookie.Value, " ")
	accToken = strings.TrimPrefix(accToken, "Bearer ")
	accToken = strings.TrimPrefix(accToken, " ")

	return accToken
}

func checkAuthHeader(c *gin.Context) (string, error) {
	rawToken := c.GetHeader("Authorization")
	if rawToken == "" {
		return "", errors.New("Unathorized")
	}

	accToken := strings.TrimPrefix(rawToken, " ")
	accToken = strings.TrimPrefix(accToken, "Bearer")
	accToken = strings.TrimPrefix(accToken, " ")

	return accToken, nil
}
