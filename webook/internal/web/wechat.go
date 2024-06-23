package web

import (
	"errors"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/service/oauth2/wechat"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
)

type OAuth2WeChatHandler struct {
	svc wechat.Service
	ijwt.Handler
	userSvc  service.UserService
	stateKey []byte
	cfg      WechatHandlerConfig
}

type WechatHandlerConfig struct {
	Secure bool
}

func NewOAuth2WeChatHandler(svc wechat.Service, userService service.UserService, cfg WechatHandlerConfig, jwtHdl ijwt.Handler) *OAuth2WeChatHandler {
	return &OAuth2WeChatHandler{
		svc:      svc,
		userSvc:  userService,
		stateKey: []byte("fb0e22c79ac75679e9781e6ba183b354"),
		cfg:      cfg,
		Handler:  jwtHdl,
	}
}

func (h *OAuth2WeChatHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/oauth2/wechat")
	g.GET("/authrul", h.OAuth2URL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WeChatHandler) OAuth2URL(ctx *gin.Context) {
	state := uuid.New()
	authUrl, err := h.svc.AuthURL(ctx.Request.Context(), state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  err.Error(),
		})
		return
	}
	err = h.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: authUrl,
	})
}

func (h *OAuth2WeChatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	ctx.SetCookie("jwt-state", tokenStr, 60*10, "/oauth2/wechat/callback", "", h.cfg.Secure, true)
	if err != nil {
		return err
	}
	return nil
}

func (h *OAuth2WeChatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	er := h.verifyState(ctx)
	if er != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登陆失败",
		})
		return
	}
	info, err := h.svc.VerifyCode(ctx, code)
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
	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{})
}

func (h *OAuth2WeChatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		//有人监控
		return fmt.Errorf("拿不到state的cookie，%w", err)
	}
	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("token已经过期了，%w", err)
	}

	if sc.State != state {
		return errors.New("state 不一致")
	}
	return nil
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
