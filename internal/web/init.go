package web

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/ktsoator/connectify/internal/web/middleware"
)

func InitRouter() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	// Initialize the session storage engine.
	// "Ktsoator" is the secret key used to encrypt/sign session data.
	// It prevents users from tampering with cookie content. In production, this should be more complex and kept secret.
	// "Ktsoator" is the authentication key (for signing).
	// "np6p_m!qY8G@Z-7*fR2&jS9#vT5%kL8B" is the encryption key (for encrypting).
	// In production, these should be loaded from environment variables.
	store := cookie.NewStore([]byte("Ktsoator"), []byte("np6p_m!qY8G@Z-7*fR2&jS9#vT5%kL8B"))

	// Register global session middleware.
	// "connectify" is the name (key) of the cookie in the browser.
	// When the browser stores the cookie, it will show Name="connectify".
	// 'store' is the storage engine created above, determining where session data is actually stored (here, in the cookie).
	server.Use(sessions.Sessions("connectify", store))

	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePath("/user/login").
		IgnorePath("/user/signup").
		Build())

	return server
}
