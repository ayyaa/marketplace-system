# Marketplace System Project

## Requirements

To run this project you need to have the following installed:

1. [Go](https://golang.org/doc/install) version 1.22
2. [Docker](https://docs.docker.com/get-docker/) version 20
3. [Docker Compose](https://docs.docker.com/compose/install/) version 1.29

## Initiate The Project

To start working, execute

```
go mod tidy
```

## Running

To run the project, run the following command:

```
docker-compose up --build
```

You should be able to access the API at http://localhost:8080

If you change `database.sql` file, you need to reinitate the database by running:

```
docker-compose down --volumes
```

# Documentation Api Spec

To see the api spec documentation, run the following command:

```
swagger serve ./docs/swagger.yaml
```

if there changers anotation for api spec, by running:
```
swag init -g main.go --output docs
```

