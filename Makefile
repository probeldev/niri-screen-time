test-macos:
	go run . -autostart disable 
	go run . -autostart enable

lint:
	golangci-lint run ./...
