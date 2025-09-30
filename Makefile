# Defines the build command for GreeterFunction for `sam build`
build-GreeterFunction:
	GOARCH=amd64 GOOS=linux go build -o ./functions/greeter/bootstrap ./functions/greeter/
	cp ./functions/greeter/bootstrap $(ARTIFACTS_DIR)/.
