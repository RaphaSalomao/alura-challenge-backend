# Alura-Challenge-Backend-2nd-Edition

This project is a challenge proposed by [Alura](https://www.alura.com.br/), which is the largest Brazilian platform for technology courses, where I've created a budget controll REST API using Golang.

This API allows you to create, read, update and delete yours month receipts and expenses. You can also retrieve a month balance of your budget.

All requests that deal with expenses and receipts need authentication, so, you need to create a user to consume those endpoints. After the authentication, every single receipt and expenses are linked to the logged user

The project was built following instructions which was given along four weeks through the following trello charts:

[1st week](https://trello.com/b/EdShXSLz/challenge-backend-1st-week)

[2nd week](https://trello.com/b/mDOu1l92/challenge-backend-2nd-week)

[3rd and 4th week](https://trello.com/b/NImixLgR/challenge-backend-3rd-week)
## Run Locally

What you will need: Golang v1.13 or greater, Docker, Linux operating system and [Swagger](https://github.com/swaggo/swag#getting-started)

Clone the project

```bash
git clone git@github.com:RaphaSalomao/alura-challenge-backend.git
```

Change the `@host` at application.go
```Golang
// @host      localhost
```

Run postgres docker container

```bash
docker-compose up -d
```

Start the server
```bash
go run application.go
```

Access the documentation at http://localhost:5000/swagger/doc.json

To see the swagger-ui, access http://localhost:5000/swagger/index.html, insert `http://localhost:5000/swagger/doc.json` at the dialog box and click on `explore`
## Running Tests

To run tests, run the following commands

Run postgres test database docker container 
```bash
docker-compose -f docker-compose.test.yml -p alura-challenge-backend_test up -d
```

Run tests
```bash
go test ./test/... -v
```