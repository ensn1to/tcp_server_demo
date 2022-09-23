all: server client
  
server: cmd/server/main.go
	go build github.com/ensn1to/tcp_server_demo/cmd/server
client: cmd/client/main.go
	go build github.com/ensn1to/tcp_server_demo/cmd/client

clean:
	rm -fr ./server
	rm -fr ./client	
	docker stop server
	docker rm server
	docker rmi server

devdownv:
	docker-compose -f ./deployments/docker-compose.yml  down -v

dev:
	docker-compose -f ./deployments/docker-compose.yml -f ./deployments/dev.docker-compose.yml up

prod:
	docker-compose -f ./deployments/docker-compose.yml -f ./deployments/prod.docker-compose.yml up -d


logs:
	docker-compose  -f ./deployments/docker-compose.yml logs -f

ps:
	docker-compose ps