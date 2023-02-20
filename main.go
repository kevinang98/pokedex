package main

import (
	"fmt"
	"os"
	"pokedex/src/configs"

	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	port, found := os.LookupEnv("PORT")
	if !found {
		port = "8000"
	}

	config := configs.Config{}
	config.Init(app)

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", port)))
}
