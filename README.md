![example workflow name](https://github.com/Arthurbp/fizz-buzz-api/workflows/CICD/badge.svg)

# fizz-buzz-api

This api allows you to get fizz-buzz response for parameters of your choice and retrieve the most used request.

## Tech/framework used

Written in Golang, this project uses open source libraries and frameworks such as for example:

- Gorilla/mux router
- MongoDB as a database with official go-driver: https://github.com/mongodb/mongo-go-driver
- Docker
- testcontainers-go for easy integration tests with docker: https://github.com/testcontainers/testcontainers-go

## Installation

### Steps:

```bash
git clone https://github.com/Arthurbp/fizz-buzz-api
```

```bash
docker-compose up --build
```

### Shut down the application

```bash
docker-compose down
```

## REST API

### Get fizz-buzz strings list

#### Request

`GET /fizzbuzz?str1={string1}&str2={string2}&int1={integer1}&int2={integer2}&limit={limit}`

    curl http://localhost:8080/fizzbuzz?str1=fizz&str2=buzz&int1=3&int2=5&limit=15

#### Response

```json
[
  "1",
  "2",
  "fizz",
  "4",
  "buzz",
  "fizz",
  "7",
  "8",
  "fizz",
  "buzz",
  "11",
  "fizz",
  "13",
  "14",
  "fizzbuzz"
]
```

### Get most frequent parameters

#### Request

`GET /fizzbuzz/stats`

    curl http://localhost:8080/fizzbuzz/stats

#### Response

```json
{
  "str1": "fizz",
  "str2": "buzz",
  "int1": 3,
  "int2": 5,
  "limit": 15,
  "nbHits": 22
}
```

## Tests

The tests can be run by using the command: `go test ./...`

Warning: due to issue with testcontainers-go and go1.15 (https://github.com/testcontainers/testcontainers-go/issues/232), go version should be < 1.15
