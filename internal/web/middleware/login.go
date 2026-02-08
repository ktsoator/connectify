package middleware

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"time"

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
		// 1. Whitelist Check: Skip authentication for registered paths (e.g., signup, login)
		if slices.Contains(l.paths, c.Request.RequestURI) {
			c.Next()
			return
		}

		// 2. Session Validation: Check if the user is authenticated
		session := sessions.Default(c)
		if session.Get("userId") == nil {
			// No userId in session means the user either isn't logged in or the session has expired
			c.AbortWithStatusJSON(http.StatusOK, resp.Result{
				Code: resp.CodeInvalidCreds,
				Msg:  "please login first",
			})
			return
		}

		// 3. Smart Session Renewal (Throttling Strategy)
		// We avoid calling session.Save() on every request to reduce Redis load/network overhead.
		// Instead, we only refresh the session if a certain amount of time has passed.
		updateTime := session.Get("update_time")
		now := time.Now().UnixMilli()

		// First request after login or session doesn't have an update timestamp yet
		if updateTime == nil {
			fmt.Println("First refresh session time")
			session.Set("update_time", now)
			session.Save() // Sync to Redis and update Cookie expiration
			c.Next()
			return
		}

		// Verify the timestamp format (defensive check)
		updateTimeValue, ok := updateTime.(int64)
		if !ok {
			log.Println("Session time format error")
			c.AbortWithStatusJSON(http.StatusOK, resp.Result{
				Code: resp.CodeServerBusy,
				Msg:  "system error",
			})
			return
		}

		// Periodic Renewal (e.g., every 60 seconds)
		// If more than 60 seconds have passed since the last update, refresh the session.
		// This extends the session life in Redis and the browser for another full term (e.g., 30 mins).
		if now-updateTimeValue > 60*1000 {
			fmt.Println("Refresh session time")
			session.Set("update_time", now)
			// Calling Save() will apply the global store options (MaxAge) to Redis/Cookie
			session.Save()
		}

		c.Next()
	}
}
