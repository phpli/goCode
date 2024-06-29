package web

import (
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.LoggerV1
}

func NewArticleHandler(svc service.ArticleService, l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", h.Edit)
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `form:"title"`
		Content string `form:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		//ctx.JSON(http.StatusOK, Result{
		//	Code: 5,
		//	Msg:  err.Error(),
		//})
		return
	}
	//c, _ := ctx.Get("claims") 取值加类型断言
	//claims, ok := c.(*ijwt.UserClaims)
	claims, ok := ctx.MustGet("claims").(*ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("没有用户信息")
		return
	}
	//调用service
	id, err := h.svc.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.Uid,
			//Name: claims.Issuer,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("保存失败！", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "ok",
		Data: id,
	})
}
