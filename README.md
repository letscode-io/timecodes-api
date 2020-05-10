# Backend for Timecodes Chrome extension

# Development

```bash
$ docker-compose up --build app_dev
```

# Run the whole test suite

```bash
$ docker-compose build app_test
$ docker-compose run app_test go test ./... -cover
```

# Run a single test

```bash
$ docker-compose build app_test
$ docker-compose run app_test go test -run Test%TEST_FUNCTION_NAME%
```

# Build for using in production

```bash
$ docker-compose build app
$ docker-compose up --no-deps -d app
```
