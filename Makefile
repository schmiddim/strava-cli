install:
	go install
test:
	go test ./...
coverage:
	go test --cover ./...
coverage-html:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html
completion-mac:
	strava-cli completion bash > /opt/homebrew/etc/bash_completion.d/strava-cli
completion-linux:
	strava-cli completion bash | sudo tee /etc/bash_completion.d/strava-cli
