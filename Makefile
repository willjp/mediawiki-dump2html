
build:
	go build

test:
	go test ./...

coverage:
	test -d .tmp || mkdir -p .tmp
	go test -coverprofile=.tmp/out.cov ./... || /bin/true

coverage-report: coverage
	go tool cover -html=.tmp/out.cov

