package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

var MAX_RANDOM_USER_COUNT_PER_REQEUST int = 20

// Index Handler
func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the Users Web Service")
}

// UserIndex Handler
func UserIndex(c echo.Context) error {
	var offset = 0
	var count = 20

	var offsetParam = c.QueryParam("offset")
	if len(offsetParam) > 0 {
		i, err := strconv.Atoi(offsetParam)
		if err != nil {
			panic(err)
		}

		if i < 0 {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "Offset must be >= 0")
		}
		offset = i
	}

	var countParam = c.QueryParam("count")
	if len(countParam) > 0 {
		i, err := strconv.Atoi(countParam)
		if err != nil {
			panic(err)
		}

		if i < 1 {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "Count must be > 0")
		}

		if i > 10000 {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "Count exceeds maximum value; please use paging by offset")
		}

		count = i
	}

	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")

	ret := make([]User, 0, count)

	idx := offset
	for len(ret) < count && idx < len(users) {
		// Do NOT return any users with an associated identity ID.
		if len(users[idx].IdentityId) == 0 {
			ret = append(ret, users[idx])
		}
		idx++
	}

	return c.JSON(http.StatusOK, ret)
}

// UserShowByID Handler
func UserShowByID(c echo.Context) error {
	userID := c.Param("userID")

	return c.JSON(http.StatusOK, RepoFindUserByID(userID))
}

// UserShowByUsername Handler
func UserShowByUsername(c echo.Context) error {
	username := c.Param("username")

	return c.JSON(http.StatusOK, RepoFindUserByUsername(username))
}

// UserLogin Handler
func UserLogin(c echo.Context) error {
	username := string(c.Get("username").([]byte))
	password := string(c.Get("password").([]byte))

	if val, ok := passwords[username]; ok {
		if val != password {
			return echo.NewHTTPError(http.StatusUnauthorized, "비밀번호가 다릅니다.")
		}
	}
	user := RepoFindUserByUsername(username)
	if user.ID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "유저를 찾을 수 없습니다.")
	}
	return c.JSON(http.StatusOK, user)
}

// UserShowByIdentityId handler
func UserShowByIdentityId(c echo.Context) error {
	identityID := c.Param("identityID")

	return c.JSON(http.StatusOK, RepoFindUserByIdentityID(identityID))
}

// ClaimUser handler
func ClaimUser(c echo.Context) error {
	if c.Request().Method == "OPTIONS" {
		return c.NoContent(http.StatusNoContent)
	}
	var userId int
	userIdVar := c.Param("userID")
	userId, err := strconv.Atoi(userIdVar)
	if err != nil {
		panic(err)
	}
	return c.JSON(http.StatusOK, RepoClaimUser(userId))
}

// GetRandomUser handler
func GetRandomUser(c echo.Context) error {
	var count = 1
	var countParam = c.QueryParam("count")
	if len(countParam) > 0 {
		i, err := strconv.Atoi(countParam)
		if err != nil {
			panic(err)
		}
		if i <= 0 || i > MAX_RANDOM_USER_COUNT_PER_REQEUST {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Sprintf("Count must be grater than 0 and less than %d", MAX_RANDOM_USER_COUNT_PER_REQEUST))
		}
		count = i
	}
	return c.JSON(http.StatusOK, RepoFindRandomUser(count))
}

// GetFilteredUser handler
func GetUnclaimedUsers(c echo.Context) error {
	segment := c.QueryParam("segment")
	var count = 1
	var countParam = c.QueryParam("count")
	if len(countParam) > 0 {
		i, err := strconv.Atoi(countParam)
		if err != nil {
			panic(err)
		}
		if i <= 0 || i > MAX_RANDOM_USER_COUNT_PER_REQEUST {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Sprintf("Count must be greater than 0 and less than %d", MAX_RANDOM_USER_COUNT_PER_REQEUST))
		}
		count = i
	}

	return c.JSON(http.StatusOK, RepoFindRandomUsersBySegment(segment, count))
}

//UserUpdate Func
func UserUpdate(c echo.Context) error {
	if c.Request().Method == "OPTIONS" {
		return c.NoContent(http.StatusNoContent)
	}

	var user User
	body, err := ioutil.ReadAll(io.LimitReader(c.Request().Body, 1048576))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := c.Request().Body.Close(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := json.Unmarshal(body, &user); err != nil {
		c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
		c.Response().WriteHeader(422) // unprocessable entity
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	t := RepoUpdateUser(user)
	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
	return c.JSON(http.StatusCreated, t)
}

//UserCreate Func
func UserCreate(c echo.Context) error {
	if c.Request().Method == "OPTIONS" {
		return c.NoContent(http.StatusNoContent)
	}
	var cu CreateUser
	body, err := ioutil.ReadAll(io.LimitReader(c.Request().Body, 1048576))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := c.Request().Body.Close(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := json.Unmarshal(body, &cu); err != nil {
		c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
		c.Response().WriteHeader(422) // unprocessable entity
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	user := cu.User

	json_data, err := json.Marshal(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp, err := http.Post("http://"+viper.GetString("recommendation.service.host")+":"+viper.GetString("recommendation.service.port")+"/create_user", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if resp.StatusCode != 200 {
		return echo.NewHTTPError(resp.StatusCode, "Unknown error.")
	}
	defer resp.Body.Close()

	t, err := RepoCreateUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	passwords[user.Username] = cu.Password

	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
	return c.JSON(http.StatusCreated, t)
}
