build:
	go build -ldflags="-s -w"
	upx deladupe

test:
	go test -coverprofile=coverage.out
	go tool cover -html coverage.out -o coverage.html
