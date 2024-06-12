.PHONY:docker
docker:
	@rm -f webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f lee/webook:v0.01
	@docker build -t lee/webook:v0.01 .