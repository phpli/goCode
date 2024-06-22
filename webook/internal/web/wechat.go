package web

import (
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WeChatHandler struct {
	svc wechat.Service
	jwtHandler
	userSvc service.UserService
}

func NewOAuth2WeChatHandler(svc wechat.Service, userService service.UserService) *OAuth2WeChatHandler {
	return &OAuth2WeChatHandler{
		svc:     svc,
		userSvc: userService,
	}
}

func (h *OAuth2WeChatHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/oauth2/wechat")
	g.GET("/authrul", h.OAuth2URL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WeChatHandler) OAuth2URL(ctx *gin.Context) {
	authUrl, err := h.svc.AuthURL(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: authUrl,
	})
}

func (h *OAuth2WeChatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	info, err := h.svc.VerifyCode(ctx, code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  err.Error(),
		})
	}
	err = h.setJWTToken(ctx, u.Id)

	ctx.JSON(http.StatusOK, Result{})
}

//type OAuth2Handler struct {
//	wecahtOAuth2Service
//}
//
//func (h *OAuth2Handler) RegisterRoutes(s *gin.Engine) {
//	g := s.Group("/oauth2")
//	g.GET("/:platform/authurl", h.AuthURL)
//	g.Any("/:platform/callback", h.Callback)
//}
//
//func (h *OAuth2Handler) AuthURL(ctx *gin.Context) {
//	platform := ctx.Param("platform")
//	switch platform {
//	case "qq":
//
//	}
//}
//
//func (h *OAuth2Handler) Callback(ctx *gin.Context) {
//
//}
