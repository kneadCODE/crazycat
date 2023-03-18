# base
COMPOSE := docker-compose -f build/docker-compose.yaml -p crazycat

teardown:
	${COMPOSE} down -v

test-go:
	${COMPOSE} run --name crazycat-${APP_NAME}-go --rm ${APP_NAME}-go sh -c "go test -coverprofile=c.out -failfast -timeout 5m -vet '' ./..."

# golib
test-golib: APP_NAME=golib
test-golib: test-go
