# MyKinoList API
MyKinoList — RESTful API service with JWT authorization for maintaining the list of watched movies. You can perform CRUD operations over the list, viz:
1) Add movies to the list.
2) Fetch all movies from the list.
3) Update information about the added movie (e.g., add the movie to Favorites, change the rating of the movie, or mark it as 'Plan to Watch').
4) Remove movies from the list.

A third-party, [unofficial Kinopoisk API](https://kinopoisk.dev/) is used to retrieve information about movies. 
The service runs in the Docker container, and there is also documentation in Swagger. The project was designed in accordance with Clean Architecture.

## Tech Stack
1. Router: [gorilla/mux](https://github.com/gorilla/mux).
2. DB and stuff: PostgreSQL, [migrate cli util](https://github.com/golang-migrate/migrate), [database/sql golang package](https://pkg.go.dev/database/sql) and [pq driver](https://github.com/lib/pq).
3. Test: [golang test package](https://pkg.go.dev/testing) and [testify lib](https://github.com/stretchr/testify) for unit-testing, [gomock lib] for mocks.
4. Security: [bcrypt lib](https://pkg.go.dev/golang.org/x/crypto/bcrypt) for hashing passwords and [jwt lib](https://github.com/golang-jwt/jwt) to generate JSONWebTokens.
5. Configuration: [viper lib](https://github.com/spf13/viper) and [gotenv](https://github.com/subosito/gotenv).

## TODO
The project is not perfect and in a good way needs improvement. Here are a few directions on where to go next:
1. SQL transactions. Should to change the architecture of the Repository layer using the unit of work pattern to be able to wrap the function calls of this layer in a transaction.
2. Use the Singleton pattern for the config object.
3. Add the ability to filter the list of movies by various criteria.
