package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	Controller "pokedex/src/controllers"
	Database "pokedex/src/databases"
	"pokedex/src/middleware"
	Repository "pokedex/src/repositories"
	Service "pokedex/src/services"
)

type Config struct {
}

func (c *Config) Init(app *echo.Echo) {
	m := middleware.NewMiddleware()

	c.GetEnv()

	mysqlDB, err := Database.NewMySQLConn()
	if err != nil {
		log.Fatal("Failed connect to mysql : " + err.Error())
	}

	mongoDB, err := Database.NewMongoConn()
	if err != nil {
		log.Fatal("Failed connect to mongodb : " + err.Error())
	}

	repoPokedex := Repository.NewPokedexRepository(mongoDB)
	repoUser := Repository.NewUserRepository(mysqlDB)

	svcPokedex := Service.NewPokedexService(repoPokedex, repoUser)
	svcUser := Service.NewUserService(repoUser)

	ctlPokedex := Controller.NewPokedexController(svcPokedex, m)
	ctlUser := Controller.NewUserController(svcUser, m)

	r := NewRoute(app, ctlPokedex, ctlUser)
	r.Init()
}

func (c *Config) GetEnv() {
	godotenv.Load(".env")

	listEnv := []string{
		"DB_HOST",
		"DB_PASS",
		"DB_USER",
		"DB_PORT",
		"DB_NAME",
		"MONGO_URI",
	}

	exist := ValidateEnv(listEnv)
	if !exist {
		log.Printf("Environment variable not defined")
	}
}

func ValidateEnv(env []string) bool {
	for _, e := range env {
		v := os.Getenv(e)

		if v == "" {
			return false
		}
	}
	return true
}
