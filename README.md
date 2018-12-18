# Back2School REST API [![Go Report Card](https://goreportcard.com/badge/github.com/middleware2018-PSS/Back2School)](https://goreportcard.com/report/github.com/middleware2018-PSS/Back2School)

RESTful API for a high school system written in [Go](https://golang.org/) using the [Buffalo](https://gobuffalo.io/) framework following the [JSON:API](https://jsonapi.org/) spec

## Starting the application

#### Development mode
```
make rundb
make initdb
buffalo dev
```

#### Production mode
```
make up
```

## Tools and libraries used
- [Go](https://golang.org/)
- [Buffalo](https://gobuffalo.io/)
- [Pop](https://github.com/gobuffalo/pop)
- [PostgreSQL](https://www.postgresql.org/)
- [JSON:API](https://github.com/google/jsonapi)
- [Casbin](https://casbin.org/)
- [JWT](https://jwt.io/)
- [Docker](https://www.docker.com/)
