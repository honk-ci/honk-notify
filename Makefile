.PHONY: check run

check: 
	golangci-lint run -c .golang-ci.yml ./... 

run:
	$(source .env)
	go run main.go