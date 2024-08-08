package service

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/article"
	artrepomocks "gitee.com/geekbang/basic-go/webook/internal/repository/mocks/article"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_articleService_Publish(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository)
		art     domain.Article
		wantErr error
		wantId  int64
	}{
		{
			name: "新建发布成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := artrepomocks.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "我的发布",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Title:   "我的发布",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
					Id: 1,
				}).Return(int64(1), nil)
				return author, reader
			},
			art: domain.Article{
				Title:   "我的发布",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: nil,
			wantId:  1,
		},
		{
			name: "修改并发布成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := artrepomocks.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Title:   "我的发布",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
					Id: 2,
				}).Return(nil)
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Title:   "我的发布",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
					Id: 2,
				}).Return(int64(2), nil)
				return author, reader
			},
			art: domain.Article{
				Title:   "我的发布",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
				Id: 2,
			},
			wantId: 2,
		},
		{
			name: "保存到制作库失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository,
				article.ArticleReaderRepository) {
				author := artrepomocks.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("mock db error"))
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)
				return author, reader
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			//wantId:  int64(0),
			wantErr: errors.New("mock db error"),
		},
		{
			name: "保存到制作库成功，但是线上库失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := artrepomocks.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Title:   "我的发布",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
					Id: 2,
				}).Return(nil)
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Title:   "我的发布",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
					Id: 2,
				}).Times(3).Return(int64(0), errors.New("mock db error"))
				return author, reader
			},
			art: domain.Article{
				Title:   "我的发布",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
				Id: 2,
			},
			wantId:  0,
			wantErr: errors.New("mock db error"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			author, reader := tc.mock(ctrl)
			svc := NewArticleServiceV1(author, reader, &logger.NopLogger{})
			id, err := svc.PublishV1(context.Background(), tc.art)
			assert.Equal(t, tc.wantId, id)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
