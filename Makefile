init: down up
up:
	docker-compose up -d --build
down:
	docker-compose down --remove-orphans

build:
	docker build -t minima:latest .
run:
	docker run -d --restart on-failure -p 9004:9004 -p 9002:9002 -p 9003:9003 -v paintmeyellow:/root/.minima minima:latest
