build: 
	@go build -o ./bin/app .

run: build
	@./bin/app	

docker:
	echo "building docker file"
	@docker build -t order_loader .
	echo "start runnign order_loader in Docker container"
	@docker run --name order_loader order_loader