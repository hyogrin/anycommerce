package pkg

import (
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

	e.GET("/carts", CartIndex)
	e.GET("/carts/:cartID", CartShowByID)
	e.POST("/carts", CartCreate)
	e.OPTIONS("/carts", CartCreate)
	e.PUT("/carts/:cartID", CartUpdate)
	e.OPTIONS("/carts/:cartID", CartUpdate)

	e.POST("/sign", SignAmazonPayPayload)
	e.OPTIONS("/sign", SignAmazonPayPayload)
}
