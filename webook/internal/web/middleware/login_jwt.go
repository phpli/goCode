package middleware

import (
	"encoding/gob"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
	cmd   redis.Cmdable
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHdl,
	}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(paths string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, paths)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	//go的方式编码二进制
	gob.Register(time.Now())
	return func(c *gin.Context) {
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		tokenStr := l.ExtractToken(c)
		claims := &ijwt.UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return ijwt.AtKey, nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != c.Request.UserAgent() {
			//严重问题
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = l.CheckSession(c, claims.Ssid)
		if err != nil {
			//要么redis有问题，要不已经退出登陆
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//已经用了长短token
		//now := time.Now()
		//if claims.ExpiresAt.Sub(now) < time.Second*50 {
		//	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//	tokenStr, err = token.SignedString([]byte("fb0e22c79ac75679e9881e6ba183b354"))
		//	c.Header("x-jwt-token", tokenStr)
		//	if err != nil {
		//		// 这边不要中断，因为仅仅是过期时间没有刷新，但是用户是登录了的
		//		log.Println(err)
		//	}
		//}

		//c.Set("userId", claims.Uid)
		c.Set("claims", claims)
	}
}
