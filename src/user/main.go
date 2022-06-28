package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"

	"github.com/anycommerce/user/pkg"
)

func init() {
	// log
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999999Z07:00"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.With().Logger()

	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	pkg.InitUser()
}

func main() {
	// echo
	e := echo.New()
	e.HideBanner = true
	pkg.Route(e)

	// echo start server
	go func() {
		if err := e.Start(":" + viper.GetString("server.port")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
