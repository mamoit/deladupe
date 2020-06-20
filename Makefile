build:
	go build -ldflags="-s -w"
	upx deladupe