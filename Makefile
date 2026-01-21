
mod:
	go mod tidy
	go mod verify
	go mod download
	go mod vendor
	go mod verify


test: clean-cache-test ## run all tests
#	go test ./mylist3 -race -v -coverprofile coverage.out
#	go tool cover -html=coverage.out -o coverage.html
#	rm coverage.out
	go test -race ./
	#go test ./

clean-cache-test: ## clean cache
	@echo "Cleaning test cache..."
	go clean -testcache
