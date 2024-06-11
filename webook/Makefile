.PHONY:docker
docker:
	@rm -f webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f lee/webook:v0.01
	@docker build -t lee/webook:v0.01 .