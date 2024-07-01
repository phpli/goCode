.PHONY:mock
mock:
	@mockgen -source=webook/internal/service/user.go -package=svcmocks -destination=webook/internal/service/mocks/user.moc.go
	@mockgen -source=webook/internal/service/article.go -package=svcmocks -destination=webook/internal/service/mocks/article.moc.go
	@mockgen -source=webook/internal/service/code.go -package=svcmocks -destination=webook/internal/service/mocks/code.moc.go
	@mockgen -source=webook/internal/repository/user.go -package=repomocks -destination=webook/internal/repository/mocks/user.moc.go
	@mockgen -source=webook/internal/repository/code.go -package=repomocks -destination=webook/internal/repository/mocks/code.moc.go
	@mockgen -source=webook/internal/repository/dao/user.go -package=daomocks -destination=webook/internal/repository/dao/mocks/user.moc.go
	@mockgen -source=webook/internal/repository/cache/user.go -package=cachemocks -destination=webook/internal/repository/cache/mocks/user.moc.go
	@mockgen -source=webook/pkg/ratelimit/types.go -package=limitmocks -destination=webook/pkg/ratelimit/mocks/ratelimit.mock.go
	@mockgen -package=redismocks -destination=webook/internal/repository/cache/redismocks/cmdable.moc.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy

