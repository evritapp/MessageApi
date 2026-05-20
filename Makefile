run:
	go run main.go

install_staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Installing staticcheck...'
	go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...