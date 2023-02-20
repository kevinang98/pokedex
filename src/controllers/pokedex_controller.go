package controllers

import (
	"encoding/json"
	"net/http"
	"pokedex/src/middleware"
	"pokedex/src/models"
	Service "pokedex/src/services"
	"strings"

	"github.com/labstack/echo/v4"
)

type PokedexController interface {
	GetPokemon(c echo.Context) error
	AddPokemon(c echo.Context) error
	UpdatePokemon(c echo.Context) error
	DeletePokemon(c echo.Context) error
}

type PokedexControllerImpl struct {
	ps Service.PokedexService
	m  middleware.Middleware
}

func NewPokedexController(ps Service.PokedexService, m middleware.Middleware) PokedexController {
	return &PokedexControllerImpl{
		ps: ps,
		m:  m,
	}
}

func (p *PokedexControllerImpl) GetPokemon(c echo.Context) error {
	var (
		params models.PokedexRequest
		types  []string
		orders int64
		err    error
	)

	id, _, role := p.m.ParseJwtToken(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Status: "fail",
			Data:   models.ErrorNoAccess,
		})
	}

	if err = c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorInvalidQueryParam,
		})
	}

	if params.Offset < 0 || params.Limit <= 0 {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorOffsetLimitParam,
		})
	}

	if params.Sort != "name" && params.Sort != "pid" && params.Sort != "" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorSortParam,
		})
	}

	if params.Order != "asc" && params.Order != "desc" && params.Order != "" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorOrderParam,
		})
	} else {
		if params.Order == "asc" {
			orders = 1
		} else {
			orders = -1
		}
	}

	if params.PokeType != "" {
		params.PokeType = strings.ReplaceAll(params.PokeType, " ", "")
		types = strings.Split(params.PokeType, ",")
	}

	pokeList, err := p.ps.GetPokemon(id, params.Name, params.Sort, params.Option, types, params.Offset, params.Limit, orders)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "error",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   pokeList,
	})
}

func (p *PokedexControllerImpl) AddPokemon(c echo.Context) error {
	var (
		req     models.Request
		pokedex []models.Pokedex
		err     error
	)

	_, _, role := p.m.ParseJwtToken(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Status: "fail",
			Data:   models.ErrorNoAccess,
		})
	}

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

	err = json.Unmarshal(byteData, &pokedex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorInvalidBodyReq,
		})
	}

	err = p.m.ValidatePokedexBody(pokedex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   err.Error(),
		})
	}

	err = p.ps.AddPokemon(pokedex)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "error",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   "success add pokemons",
	})
}

func (p *PokedexControllerImpl) UpdatePokemon(c echo.Context) error {
	var (
		req     models.Request
		pokedex []models.Pokedex
		err     error
	)

	_, _, role := p.m.ParseJwtToken(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Status: "fail",
			Data:   models.ErrorNoAccess,
		})
	}

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

	err = json.Unmarshal(byteData, &pokedex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorInvalidBodyReq,
		})
	}

	err = p.m.ValidatePokedexBody(pokedex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   err.Error(),
		})
	}

	listID, err := p.ps.UpdatePokemon(pokedex)
	if err != nil {
		errorDesc := ""
		if listID == "" {
			errorDesc = err.Error()
		} else {
			errorDesc = "Error : " + err.Error() + ", Updated ID : " + listID
		}
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "error",
			Data:   errorDesc,
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   "success update pokemons",
	})
}

func (p *PokedexControllerImpl) DeletePokemon(c echo.Context) error {
	var (
		req     models.Request
		pokedex models.DeletePokedex
		err     error
	)

	_, _, role := p.m.ParseJwtToken(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Status: "fail",
			Data:   models.ErrorNoAccess,
		})
	}

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

	err = json.Unmarshal(byteData, &pokedex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   models.ErrorInvalidBodyReq,
		})
	}

	err = p.m.ValidateIntPokedexBody(pokedex.PID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: "fail",
			Data:   err.Error(),
		})
	}

	listID, err := p.ps.DeletePokemon(pokedex.PID)
	if err != nil {
		errorDesc := ""
		if listID == "" {
			errorDesc = err.Error()
		} else {
			errorDesc = "Error : " + err.Error() + ", Deleted ID : " + listID
		}
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: "error",
			Data:   errorDesc,
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   "success delete pokemons",
	})
}
