package pkg

import (
	"crypto/subtle"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Route(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAcceptEncoding, echo.HeaderXCSRFToken, echo.HeaderContentLength, echo.HeaderAuthorization, "Cache-Control"},
		AllowCredentials: true,
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			ignorePaths := []string{"/", "/health", "/metrics"}
			path := c.Path()
			for _, v := range ignorePaths {
				if path == v {
					return true
				}
			}

			return false
		},
		Format: `{"level":"info","time":"${time_custom}","status":${status},"method":"${method}","uri":"${uri}",` +
			`"id":"${id}","remote_ip":"${remote_ip}","message":"${error}","host":"${host}",` +
			`"latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02T15:04:05.999999Z07:00",
	}))

	e.GET("/", Index)
	e.GET("/health", Index)

	basicAuth := middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("")) == 1 {
			return false, nil
		}
		c.Set("username", []byte(username))
		c.Set("password", []byte(password))
		return true, nil
	})
	e.POST("/login", UserLogin, basicAuth)

	e.GET("/users/all", UserIndex)
	e.GET("/users/id/:userID", UserShowByID)
	e.GET("/users/username/:username", UserShowByUsername)
	e.GET("/users/identityid/:identityID", UserShowByIdentityId)

	e.GET("/users/unclaimed/", GetUnclaimedUsers)

	e.GET("/users/random/", GetRandomUser)

	e.PUT("/users/id/:userID/claim", ClaimUser)
	e.OPTIONS("/users/id/:userID/claim", ClaimUser)

	e.POST("/users", UserCreate)
	e.OPTIONS("/users", UserCreate)

	e.PUT("/users/id/:userID", UserUpdate)
	e.OPTIONS("/users/id/:userID", UserUpdate)
}
