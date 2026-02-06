package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ktsoator/connectify/internal/web/resp"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePath(paths string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, paths)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		if slices.Contains(l.paths, c.Request.RequestURI) {
			c.Next()
			return
		}

		session := sessions.Default(c)
		if session.Get("userId") == nil {
			// If the session does not have userId, it means the user is not logged in
			c.AbortWithStatusJSON(http.StatusOK, resp.Result{
				Code: resp.CodeInvalidCreds,
				Msg:  "please login first",
			})
			return
		}
		c.Next()

	}
}
