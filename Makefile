all: server client
  
server: cmd/server/main.go
	go build github.com/ensn1to/tcp_server_demo/cmd/server
client: cmd/client/main.go
	go build github.com/ensn1to/tcp_server_demo/cmd/client

clean:
	rm -fr ./server
	rm -fr ./client