package web

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
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

	// Session storage initialization.
	// We have two options for session storage:

	// Option 1: Cookie-based session storage.
	// All session data is stored directly in the cookie on the client side.
	// "Ktsoator" is the authentication key (for signing).
	// "np6p_m!qY8G@Z-7*fR2&jS9#vT5%kL8B" is the encryption key (for encrypting).
	// In production, these should be loaded from environment variables.
	// store := cookie.NewStore([]byte("Ktsoator"), []byte("np6p_m!qY8G@Z-7*fR2&jS9#vT5%kL8B"))

	// Option 2: Redis-based session storage.
	// Only the session ID is stored in the cookie; the actual data resides in Redis.
	// Parameters:
	// 1. 10: Maximum number of idle connections in the pool.
	// 2. "tcp": Network type.
	// 3. "localhost:16379": Redis server address (mapped in docker-compose).
	// 4. "": Username (empty for default Redis setup).
	// 5. "": Password (empty as per ALLOW_EMPTY_PASSWORD=yes in docker-compose).
	// 6. []byte(...): Authentication key for signing session cookies.
	// 7. []byte(...): Encryption key for encrypting session data (AES).
	store, err := redis.NewStore(10, "tcp", "localhost:16379", "", "",
		[]byte("Ktsoator"), []byte("np6p_m!qY8G@Z-7*fR2&jS9#vT5%kL8B"))
	if err != nil {
		fmt.Println("Failed to initialize Redis session store:", err)
		panic(err)
	}

	// Option B: Set global default session options at the store level.
	// This ensures all sessions created via this store share the same secure defaults.
	store.Options(sessions.Options{
		// Path: The path where the cookie is valid. "/" means the entire site.
		Path: "/",
		// MaxAge: Default session expiration time (30 minutes).
		MaxAge: 30 * 60,
		// HttpOnly: Prevents client-side scripts from accessing the cookie.
		HttpOnly: true,
		// Secure: Set to false for local HTTP development.
		Secure: false,
	})

	// Register global session middleware.
	// "connectify" is the name (key) of the cookie in the browser.
	// When the browser stores the cookie, it will show Name="connectify".
	// 'store' is the storage engine created above, determining where session data is actually stored (here, in the cookie).
	server.Use(sessions.Sessions("connectify", store))

	// server.Use(middleware.NewLoginMiddlewareBuilder().
	// 	IgnorePath("/user/login").
	// 	IgnorePath("/user/signup").
	// 	Build())

	// JWT login middleware
	// Ignore authentication for the following paths
	server.Use(middleware.NewLoginJwtMiddlewareBuilder().
		IgnorePath("/user/login_jwt").
		IgnorePath("/user/signup").
		Build())

	return server
}
