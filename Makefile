SHELL := /bin/bash

test:
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v -coverprofile=coverage.txt -covermode=atomic
	go tool cover -html=coverage.txt -o coverage.html

demo:
	go run -mod=mod entgo.io/ent/cmd/ent new User
	go run -mod=mod entgo.io/ent/cmd/ent new Car Group

gogen:
	go generate ./ent