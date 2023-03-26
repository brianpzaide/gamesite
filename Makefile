BINARY_NAME=gamesite
STUFFED_BINARY_NAME=stuffedgamesite

## build: Build binary
build:
	@echo "Building..."
	@env CGO_ENABLED=0  go build -ldflags="-s -w" -o ${BINARY_NAME} ./cmd/web
	@/go/bin/stuffbin -a stuff -in ${BINARY_NAME} -out ${STUFFED_BINARY_NAME} ui/static:/ ui/html:/ ui/templates/wshome.tmpl:/templates/wshome.tmpl
	@echo "Built!"

## run: builds and runs the application
run: clean build 
	@echo "Starting..."
	./${STUFFED_BINARY_NAME} &
	@echo "Started!"

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME} ${STUFFED_BINARY_NAME}
	@echo "Cleaned!"

## start: an alias to run
start: run

## stop: stops the running application
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./${STUFFED_BINARY_NAME}"
	@echo "Stopped!"

## restart: stops and starts the application
restart: stop start

## test: runs all tests
test:
	go test -v ./...