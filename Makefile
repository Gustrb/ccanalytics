clean:
	rm -rf dist
	rm -rf **/*.db

build:
	go build -o dist/ccanalytics cmd/api/api.go
	go build -o dist/checksign cmd/checksign/checksign.go
	go build -o dist/signer cmd/signer/signer.go
	go build -o dist/migrate cmd/migrate/migrate.go
