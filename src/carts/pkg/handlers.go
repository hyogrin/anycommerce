package pkg

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/labstack/echo/v4"
)

// Index Handler
func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the Carts Web Service")
}

// CartIndex Handler
func CartIndex(c echo.Context) error {
	var values []Cart
	for _, value := range carts {
		values = append(values, value)
	}

	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
	return c.JSON(http.StatusOK, values)
}

// CartShowByID Handler
func CartShowByID(c echo.Context) error {
	cartID := c.Param("cartID")

	return c.JSON(http.StatusOK, RepoFindCartByID(cartID))
}

//CartUpdate Func
func CartUpdate(c echo.Context) error {
	if c.Request().Method == "OPTIONS" {
		return c.NoContent(http.StatusOK)
	}

	var cart Cart
	body, err := ioutil.ReadAll(io.LimitReader(c.Request().Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := c.Request().Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &cart); err != nil {
		c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
		c.Response().WriteHeader(422) // unprocessable entity
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	cartID := c.Param("cartID")

	t := RepoUpdateCart(cartID, cart)

	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
	return c.JSON(http.StatusOK, t)
}

//CartCreate Func
func CartCreate(c echo.Context) error {
	if c.Request().Method == "OPTIONS" {
		return c.NoContent(http.StatusOK)
	}

	var cart Cart
	body, err := ioutil.ReadAll(io.LimitReader(c.Request().Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := c.Request().Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &cart); err != nil {
		c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
		c.Response().WriteHeader(422) // unprocessable entity
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())

	}

	t := RepoCreateCart(cart)

	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
	return c.JSON(http.StatusOK, t)
}

//Sign a payload for Amazon Pay - delegates to a Lambda function for doing this.
func SignAmazonPayPayload(c echo.Context) error {
	if c.Request().Method == "OPTIONS" {
		return c.NoContent(http.StatusOK)
	}

	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(awsSession)

	body, err := ioutil.ReadAll(io.LimitReader(c.Request().Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := c.Request().Body.Close(); err != nil {
		panic(err)
	}

	var requestBody map[string]interface{}
	json.Unmarshal(body, &requestBody)

	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("AmazonPaySigningLambda"), Payload: body})
	if err != nil {
		panic(err)
	}

	var responsePayload map[string]interface{}
	json.Unmarshal(result.Payload, &responsePayload)

	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
	return c.JSON(http.StatusOK, responsePayload)
}
