build-MLGMRFunction:
	GOARCH=amd64 GOOS=linux go build -o ./bootstrap main.go
	cp ./bootstrap $(ARTIFACTS_DIR)/.
