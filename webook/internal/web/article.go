package web

import (
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	g.POST("/publish", h.Publish)
	g.GET("/detail/:id", h.Detail)
	//g.POST("/withdraw", h.Withdraw)

	//pub.GET("/:id", ginx.WrapClaims(h.PubDetail))
	//pub.POST("/like", ginx.WrapClaimsAndReq[LikeReq](h.Like))
	//pub.POST("/collect", ginx.WrapClaimsAndReq[CollectReq](h.Collect))

	//g.POST("/list", ginx.WrapClaimsAndReq(h.List))

}

func (h *ArticleHandler) Publish(ctx *gin.Context) {

	req, claims, done := h.saveOrPublish(ctx)
	if done {
		return
	}
	//调用service
	id, err := h.svc.Publish(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("发表失败！", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "ok",
		Data: id,
	})
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	req, claims, done := h.saveOrPublish(ctx)
	if done {
		return
	}
	//调用service
	id, err := h.svc.Save(ctx, req.toDomain(claims.Uid))
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

func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
			//Name: claims.Issuer,
		},
	}
}

func (h *ArticleHandler) saveOrPublish(ctx *gin.Context) (ArticleReq, *ijwt.UserClaims, bool) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return ArticleReq{}, nil, true
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
		return ArticleReq{}, nil, true
	}
	return req, claims, false
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		h.l.Error("前端输入的ID不对", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("user").(ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("获取用户信息失败", logger.Error(err))
		return
	}
	resp, err := h.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("获得文章信息失败", logger.Error(err))
		return
	}
	// 这是不借助数据库查询来判定的方法
	if resp.Author.Id != usr.Id {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			// 也不需要告诉前端究竟发生了什么
			Msg: "输入有误",
		})
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。
		h.l.Error("非法访问文章，创作者 ID 不匹配",
			logger.Int64("uid", usr.Id))
		return
	}
	ctx.JSON(http.StatusOK, Result{})
}
