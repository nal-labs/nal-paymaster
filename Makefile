install-dev:
	go install github.com/air-verse/air@latest
	go install github.com/ethereum/go-ethereum/cmd/abigen@latest
	go mod tidy
	cp .env.template .env

dev-run:
	air -c .air.server.toml

generate-contract-pkg:
	abigen --abi=./abi/SingletonPaymasterV6.json --pkg=contract --out=./pkg/contract/bindings.go
	# abigen --abi=./abi/verifyingPaymaster.json --pkg=contract --out=./pkg/contract/bindings.go