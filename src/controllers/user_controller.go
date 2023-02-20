package controllers

import (
	"encoding/json"
	"net/http"
	"pokedex/src/middleware"
	"pokedex/src/models"
	Service "pokedex/src/services"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserController interface {
	RegisterUser(c echo.Context) error
	LoginUser(c echo.Context) error
	UpdateCapturedPokemonUser(c echo.Context) error
}

type UserControllerImpl struct {
	us Service.UserService
	m  middleware.Middleware
}

func NewUserController(us Service.UserService, m middleware.Middleware) UserController {
	return &UserControllerImpl{
		us: us,
		m:  m,
	}
}

func (u *UserControllerImpl) RegisterUser(c echo.Context) error {
	var (
		err error
	)

	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")

	if username == "" || password == "" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorUserPassEmpty,
		})
	} else if len(username) > 30 || len(password) > 30 {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorUserPassTooLong,
		})
	}

	if role != "admin" && role != "user" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorInvalidRole,
		})
	}

	if err = u.us.RegisterUser(username, password, role); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "error",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   "username " + username + " success registered",
	})
}

func (u *UserControllerImpl) LoginUser(c echo.Context) error {
	var (
		err error
	)

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorUserPassEmpty,
		})
	} else if len(username) > 30 || len(password) > 30 {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorUserPassTooLong,
		})
	}

	user, err := u.us.LoginUser(username, password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "error",
			Data:   err.Error(),
		})
	}

	token, err := u.m.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "error",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   token,
	})
}

func (u *UserControllerImpl) UpdateCapturedPokemonUser(c echo.Context) error {
	var (
		req     models.Request
		cap     models.PokemonCaptured
		pCap    string
		pokeCap []string
		err     error
	)

	id, _, _ := u.m.ParseJwtToken(c)

	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorFailGetBody,
		})
	}

	byteData, err := json.Marshal(req.Data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "fail",
			Data:   models.ErrorFailMarshalBody,
		})
	}

	err = json.Unmarshal(byteData, &cap)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorInvalidBodyReq,
		})
	}

	err = u.m.ValidateIntPokedexBody(cap.PID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   err.Error(),
		})
	}

	for _, pid := range cap.PID {
		pCap = ""

		pCap = strconv.Itoa(pid)
		pokeCap = append(pokeCap, pCap)
	}

	listPID := strings.Join(pokeCap, ",")
	listPID = strings.ReplaceAll(listPID, " ", "")

	if listPID == "" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorPIDEmpty,
		})
	}

	if err = u.us.UpdateCapturedPokemonUser(id, listPID); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "error",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   "pokemon captured has updated",
	})
}
