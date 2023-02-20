## Pokedex

Why use mysql or mongodb?

```For user, I prefer to use mysql because maybe in the future, there are several case that need relationship between table, for example transaction table for purchase in-game data and the structure table should be rarely to change. For pokedex, I prefer to use mongodb because in game case, there will be update that maybe need to add or remove data from pokemon attribute. Therefore nosql should be cover that and pokemon data has a lot of object data that if implemented in sql database, will have many tables.```

How to run:
   - Copy Paste .env-example file and rename it to .env
   - Change mysql database env in .env and port if you want (make sure not to change env variable name) and make new database with name `pokemon`
   - Change mongodb database env in .env with mongodb URI and make database with name `pokemon`. Heres the query for the collection :
   ```
   db.createCollection('pokedex')
   db.poketest.createIndex({ "pid": 1 }, { unique: true })
   db.poketest.createIndex({ "name": 1 })
   db.poketest.createIndex({ "type": 1 })
   ```
   - Run command `make run` to run program within project or `make build` to build program and run it (or run command from makefile manually if you prefer)
   - Register your admin or user using `/user/register` endpoint, then login to get token with 15 minutes expire time using `/user/login`
   - Documentation postman attached in this project and API documentation link : https://docs.google.com/document/d/1CJCBdYn7yLarBrNg6EYQaYkA8KTJhx8MWbBvBpTrBAM/edit?usp=sharing
   - If want to use auto migration db, make sure env variable `GOOSE` is `ENABLE`, otherwise you can run query table ddl manually in `/src/databases/ddl` and set env variable `GOOSE` value to `DISABLE`.
   - Run command `make mock` to mock function in several package and `make test` to check unit testing and code coverage in this repository.