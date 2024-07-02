package article

import (
	"context"
)

type ReaderDAO interface {
	Upsert(ctx context.Context, art Article) error
}

// 线上表
type ArticleReaderDAO struct {
	Article
}
