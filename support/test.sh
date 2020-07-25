while ! nc -z db 5432; do sleep 1; done;

go test ./... -covermode=count -coverprofile=tmp/coverage.out
