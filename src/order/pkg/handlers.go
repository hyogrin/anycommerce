package pkg

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Index Handler
func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the Orders Web Service")
}

// OrderIndex Handler
func OrderIndex(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")

	return c.JSON(http.StatusOK, orders)
}

// OrderIndexByUsername Handler
func OrderIndexByUsername(c echo.Context) error {
	username := c.Param("username")

	return c.JSON(http.StatusOK, RepoFindOrdersByUsername(username))
}

// OrderShowByID Handler
func OrderShowByID(c echo.Context) error {
	orderID := c.Param("orderID")

	return c.JSON(http.StatusOK, RepoFindOrderByID(orderID))
}

//OrderUpdate Func
func OrderUpdate(c echo.Context) error {
	if c.Request().Method == "OPTIONS" {
		return c.NoContent(http.StatusOK)
	}

	var order Order
	body, err := ioutil.ReadAll(io.LimitReader(c.Request().Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := c.Request().Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &order); err != nil {
		c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	t := RepoUpdateOrder(order)
	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
	return c.JSON(http.StatusCreated, t)
}

//OrderCreate Func
func OrderCreate(c echo.Context) error {
	if c.Request().Method == "OPTIONS" {
		return c.NoContent(http.StatusOK)
	}

	var order Order
	body, err := ioutil.ReadAll(io.LimitReader(c.Request().Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := c.Request().Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &order); err != nil {
		c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())

	}

	log.Info().Str("Personalize", "true").Str("user_username", order.Username).Interface("product", order.Items).Msg("user purchased order.")

	t := RepoCreateOrder(order)
	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
	return c.JSON(http.StatusCreated, t)
}
