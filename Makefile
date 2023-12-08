.PHONY: all
all: server client

.PHONY: server
server:
	go build -o bin/znfsd cmd/server/main.go

.PHONY: client
client:
	go build -o bin/znfs cmd/client/main.go
