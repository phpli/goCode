package article

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	//Update(ctx context.Context, article domain.Article) error
}
